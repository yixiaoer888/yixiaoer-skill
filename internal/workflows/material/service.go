package material

import (
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type Service struct{}

type AddInput struct {
	FilePath  string
	ThumbPath string
	Type      string
}

func NewService() Service {
	return Service{}
}

func (Service) Create(payload map[string]interface{}) (map[string]interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
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
	return client.New(cfg).Material(body)
}

func (Service) Add(input AddInput) (map[string]interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	apiClient := client.New(cfg)
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
