package upload

import (
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type DryRunResult struct {
	Bucket      string `json:"bucket"`
	Source      string `json:"source"`
	SourceType  string `json:"sourceType"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Format      string `json:"format,omitempty"`
}

func Preview(pathOrURL, bucket string) (DryRunResult, error) {
	source := strings.TrimSpace(pathOrURL)
	if source == "" {
		return DryRunResult{}, yxerrors.Usage("upload requires a file path or URL", nil).
			WithHint("请传入位置参数，或使用 --file / --url。")
	}
	if bucket == "" {
		bucket = "cloud-publish"
	}
	fileName := sourceFileName(source)
	return DryRunResult{
		Bucket:      bucket,
		Source:      source,
		SourceType:  detectSourceType(source),
		FileName:    fileName,
		ContentType: api.DetectContentType(source),
		Format:      strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), "."),
	}, nil
}

func sourceFileName(source string) string {
	name := filepath.Base(source)
	if name == "." || name == "/" || name == "" {
		return "file"
	}
	return name
}

func detectSourceType(source string) string {
	lower := strings.ToLower(source)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return "url"
	}
	return "file"
}
