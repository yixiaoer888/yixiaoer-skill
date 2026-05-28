package material

import (
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type DryRunResult struct {
	Upload  map[string]interface{} `json:"upload,omitempty"`
	Thumb   map[string]interface{} `json:"thumbUpload,omitempty"`
	Request map[string]interface{} `json:"request"`
}

func PreviewAdd(input AddInput) (DryRunResult, error) {
	if strings.TrimSpace(input.FilePath) == "" {
		return DryRunResult{}, yxerrors.Usage("material add requires file", nil).
			WithHint("请传入 --file，本地路径或 URL 均可。")
	}
	fileName := filepath.Base(strings.TrimSpace(input.FilePath))
	contentType := api.DetectContentType(input.FilePath)
	fileType := strings.TrimSpace(input.Type)
	if fileType == "" {
		fileType = detectMaterialType(contentType)
	}

	result := DryRunResult{
		Upload: map[string]interface{}{
			"bucket":      "material-library",
			"source":      input.FilePath,
			"contentType": contentType,
			"fileName":    fileName,
		},
		Request: map[string]interface{}{
			"filePath": "material-library/" + fileName,
			"fileName": fileName,
			"width":    0,
			"height":   0,
			"type":     fileType,
		},
	}
	if strings.TrimSpace(input.ThumbPath) != "" {
		thumbName := filepath.Base(strings.TrimSpace(input.ThumbPath))
		result.Thumb = map[string]interface{}{
			"bucket":      "material-library",
			"source":      input.ThumbPath,
			"contentType": api.DetectContentType(input.ThumbPath),
			"fileName":    thumbName,
		}
		result.Request["thumbPath"] = "material-library/" + thumbName
	}
	return result, nil
}
