package material

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type Service struct {
	rt *app.Runtime
}

type AddInput struct {
	FilePath  string
	ThumbPath string
	Type      string
}

type MoveInput struct {
	GroupID string
}

type MaterialMatch struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	Type     string `json:"type,omitempty"`
	FilePath string `json:"filePath,omitempty"`
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func (s Service) Create(payload map[string]interface{}) (map[string]interface{}, error) {
	body := BuildMaterialBody(payload)
	for _, required := range []string{"filePath", "fileName", "width", "height", "type"} {
		if _, ok := body[required]; !ok {
			return nil, yxerrors.Usage("material create requires payload fields", []string{
				"filePath",
				"fileName",
				"width",
				"height",
				"type",
			}).
				WithHint("请提供已上传素材的完整登记字段，至少包含 filePath、fileName、width、height、type。")
		}
	}
	return s.rt.Client.Material(body)
}

func (s Service) Move(materialID string, input MoveInput) (map[string]interface{}, error) {
	if err := ValidateMoveInput(input); err != nil {
		return nil, err
	}
	return s.rt.Client.MoveMaterial(materialID, BuildMoveBody(materialID, input))
}

func (s Service) List(opts api.MaterialListOptions) (interface{}, error) {
	return s.rt.Client.Materials(opts)
}

func (s Service) ResolveByFileName(fileName string, opts api.MaterialListOptions) (MaterialMatch, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return MaterialMatch{}, yxerrors.Usage("material file name must not be empty", nil).
			WithHint("请传入素材文件名，例如 yxer material move-by-name demo.png --group-id group_1 --dry-run。")
	}

	opts.FileName = fileName
	result, err := s.List(opts)
	if err != nil {
		return MaterialMatch{}, err
	}
	matches := FindExactFileNameMatches(result, fileName)
	switch len(matches) {
	case 0:
		return MaterialMatch{}, yxerrors.Usage("no material matched the file name", map[string]interface{}{
			"fileName": fileName,
		}).
			WithHint("请先用 yxer material list --name <文件名> 查询素材；文件名需要包含扩展名并完全匹配。").
			WithNextCommand(fmt.Sprintf("yxer material list --name %q", fileName))
	case 1:
		return matches[0], nil
	default:
		candidates := make([]map[string]string, 0, len(matches))
		for _, candidate := range matches {
			candidates = append(candidates, map[string]string{
				"id":       candidate.ID,
				"fileName": candidate.FileName,
				"type":     candidate.Type,
				"filePath": candidate.FilePath,
			})
		}
		return MaterialMatch{}, yxerrors.Usage("multiple materials matched the file name", map[string]interface{}{
			"fileName":   fileName,
			"candidates": candidates,
		}).
			WithHint("文件名命中多条素材。请从 candidates 中选择 id，再执行 material move 并先使用 --dry-run。")
	}
}

func (s Service) Add(input AddInput) (map[string]interface{}, error) {
	apiClient := s.rt.Client
	if strings.TrimSpace(input.FilePath) == "" {
		return nil, yxerrors.Usage("material add requires file", nil).
			WithHint("请传入 --file，本地路径或 URL 均可。")
	}
	uploadResult, err := apiClient.Upload(input.FilePath, "material-library", true)
	if err != nil {
		return nil, err
	}
	fileType := strings.TrimSpace(input.Type)
	if fileType == "" {
		fileType = detectMaterialType(uploadResult.ContentType)
	}
	body := map[string]interface{}{
		"filePath": uploadResult.Key,
		"fileName": filepath.Base(input.FilePath),
		"width":    uploadResult.Width,
		"height":   uploadResult.Height,
		"type":     fileType,
	}
	if strings.TrimSpace(input.ThumbPath) != "" {
		thumbResult, err := apiClient.Upload(input.ThumbPath, "material-library", true)
		if err != nil {
			return nil, err
		}
		body["thumbPath"] = thumbResult.Key
	}
	return apiClient.Material(body)
}

func BuildMaterialBody(payload map[string]interface{}) map[string]interface{} {
	body := map[string]interface{}{}
	for _, field := range []string{"filePath", "fileName", "width", "height", "type", "thumbPath"} {
		if value, ok := payload[field]; ok {
			body[field] = value
		}
	}
	return body
}

func ValidateMoveInput(input MoveInput) error {
	if strings.TrimSpace(input.GroupID) == "" {
		return yxerrors.Usage("material group id must not be empty", nil).
			WithHint("请传入目标分组 ID，例如 --group-id group_1。")
	}
	return nil
}

func BuildMoveBody(materialID string, input MoveInput) map[string]interface{} {
	return map[string]interface{}{
		"materialIds": []string{strings.TrimSpace(materialID)},
		"groupId":     strings.TrimSpace(input.GroupID),
	}
}

// FindExactFileNameMatches extracts material rows from common paginated API
// response shapes and only keeps complete, exact file-name matches.
func FindExactFileNameMatches(result interface{}, fileName string) []MaterialMatch {
	target := strings.TrimSpace(fileName)
	if target == "" {
		return nil
	}

	seen := map[string]bool{}
	matches := []MaterialMatch{}
	visitMaterialRows(result, func(row map[string]interface{}) {
		match := MaterialMatch{
			ID:       firstString(row, "id", "materialId", "yixiaoerId"),
			FileName: firstString(row, "fileName", "name"),
			Type:     firstString(row, "type"),
			FilePath: firstString(row, "filePath", "path"),
		}
		if match.ID == "" || match.FileName != target || seen[match.ID] {
			return
		}
		seen[match.ID] = true
		matches = append(matches, match)
	})
	return matches
}

func visitMaterialRows(value interface{}, visit func(map[string]interface{})) {
	switch typed := value.(type) {
	case map[string]interface{}:
		if _, hasFileName := typed["fileName"]; hasFileName {
			visit(typed)
			return
		}
		for _, key := range []string{"data", "items", "list", "records", "rows", "results"} {
			if nested, ok := typed[key]; ok {
				visitMaterialRows(nested, visit)
			}
		}
	case []interface{}:
		for _, item := range typed {
			visitMaterialRows(item, visit)
		}
	}
}

func firstString(row map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := row[key]; ok {
			if text, ok := value.(string); ok {
				return strings.TrimSpace(text)
			}
		}
	}
	return ""
}

func detectMaterialType(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return "image"
	case strings.HasPrefix(contentType, "video/"):
		return "video"
	default:
		return "file"
	}
}
