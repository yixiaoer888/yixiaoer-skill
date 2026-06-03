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
	uploadMeta, _, err := api.InspectUpload(input.FilePath, true)
	if err != nil {
		return DryRunResult{}, err
	}
	contentType := uploadMeta.ContentType
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
			"size":        uploadMeta.Size,
			"width":       uploadMeta.Width,
			"height":      uploadMeta.Height,
			"duration":    uploadMeta.Duration,
			"format":      uploadMeta.Format,
		},
		Request: map[string]interface{}{
			"filePath": "material-library/" + fileName,
			"fileName": fileName,
			"width":    uploadMeta.Width,
			"height":   uploadMeta.Height,
			"type":     fileType,
		},
	}
	if strings.TrimSpace(input.ThumbPath) != "" {
		thumbName := filepath.Base(strings.TrimSpace(input.ThumbPath))
		thumbMeta, _, err := api.InspectUpload(input.ThumbPath, true)
		if err != nil {
			return DryRunResult{}, err
		}
		result.Thumb = map[string]interface{}{
			"bucket":      "material-library",
			"source":      input.ThumbPath,
			"contentType": thumbMeta.ContentType,
			"fileName":    thumbName,
			"size":        thumbMeta.Size,
			"width":       thumbMeta.Width,
			"height":      thumbMeta.Height,
			"duration":    thumbMeta.Duration,
			"format":      thumbMeta.Format,
		}
		result.Request["thumbPath"] = "material-library/" + thumbName
	}
	return result, nil
}
