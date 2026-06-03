package api

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/core/media"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type UploadResult struct {
	Key         string  `json:"key"`
	ContentType string  `json:"contentType"`
	Bucket      string  `json:"bucket"`
	Size        int64   `json:"size,omitempty"`
	Width       int     `json:"width,omitempty"`
	Height      int     `json:"height,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
	Format      string  `json:"format,omitempty"`
}

func InspectUpload(pathOrURL string, autoMeta bool) (UploadResult, string, error) {
	contentType := DetectContentType(pathOrURL)
	buffer, fileName, size, err := readUploadContent(pathOrURL)
	if err != nil {
		return UploadResult{}, "", err
	}
	result, err := buildUploadMetadata(pathOrURL, buffer, fileName, size, contentType, autoMeta)
	if err != nil {
		return UploadResult{}, "", err
	}
	return result, fileName, nil
}

func (c *Client) Upload(pathOrURL, bucket string, autoMeta bool) (UploadResult, error) {
	if bucket == "" {
		bucket = "cloud-publish"
	}
	result, fileName, err := InspectUpload(pathOrURL, autoMeta)
	if err != nil {
		return UploadResult{}, err
	}

	params := map[string]string{
		"fileKey":     fileName,
		"contentType": result.ContentType,
	}
	if result.Size > 0 {
		params["size"] = fmt.Sprint(result.Size)
	}

	var uploadInfo map[string]interface{}
	if err := c.Get(Query("/storages/"+bucket+"/upload-url", params), &uploadInfo); err != nil {
		return UploadResult{}, err
	}
	data, _ := DataOrSelf(uploadInfo).(map[string]interface{})
	serviceURL, _ := data["serviceUrl"].(string)
	key, _ := data["key"].(string)
	if serviceURL == "" || key == "" {
		return UploadResult{}, yxerrors.Remote("invalid upload info response", uploadInfo).
			WithCategory("remote_response")
	}

	buffer, _, _, err := readUploadContent(pathOrURL)
	if err != nil {
		return UploadResult{}, err
	}
	req, err := http.NewRequest(http.MethodPut, serviceURL, bytes.NewReader(buffer))
	if err != nil {
		return UploadResult{}, err
	}
	req.Header.Set("Content-Type", result.ContentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UploadResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return UploadResult{}, yxerrors.Remote("failed to upload to OSS", string(raw)).
			WithCategory("remote_upload")
	}

	result.Key = key
	result.Bucket = bucket
	return result, nil
}

func probeVideoMetadata(pathOrURL string, raw []byte, fileName string) (media.VideoMetadata, error) {
	lower := strings.ToLower(pathOrURL)
	if !strings.HasPrefix(lower, "http://") && !strings.HasPrefix(lower, "https://") {
		return media.ProbeVideo(pathOrURL)
	}

	tmpFile, err := os.CreateTemp("", "yxer-upload-*"+filepath.Ext(fileName))
	if err != nil {
		return media.VideoMetadata{}, err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	if _, err := tmpFile.Write(raw); err != nil {
		_ = tmpFile.Close()
		return media.VideoMetadata{}, err
	}
	if err := tmpFile.Close(); err != nil {
		return media.VideoMetadata{}, err
	}
	return media.ProbeVideo(tmpPath)
}

func buildUploadMetadata(pathOrURL string, buffer []byte, fileName string, size int64, contentType string, autoMeta bool) (UploadResult, error) {
	width, height := imageDimensions(pathOrURL, buffer, contentType)
	duration := float64(0)
	if strings.HasPrefix(contentType, "video/") && autoMeta {
		videoMeta, err := probeVideoMetadata(pathOrURL, buffer, fileName)
		if err != nil {
			return UploadResult{}, err
		}
		width = videoMeta.Width
		height = videoMeta.Height
		duration = videoMeta.Duration
	}
	return UploadResult{
		ContentType: contentType,
		Size:        size,
		Width:       width,
		Height:      height,
		Duration:    duration,
		Format:      strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), "."),
	}, nil
}

func DetectContentType(pathOrURL string) string {
	ext := strings.ToLower(filepath.Ext(pathOrURL))
	if guessed := mime.TypeByExtension(ext); guessed != "" {
		if strings.Contains(guessed, ";") {
			return strings.Split(guessed, ";")[0]
		}
		return guessed
	}
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}

func readUploadContent(pathOrURL string) ([]byte, string, int64, error) {
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return nil, "", 0, err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, "", 0, yxerrors.Remote("HTTP error downloading file during sync upload", map[string]interface{}{
				"statusCode": resp.StatusCode,
				"url":        pathOrURL,
			}).WithCategory("remote_download")
		}
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", 0, err
		}
		parsed, _ := url.Parse(pathOrURL)
		fileName := filepath.Base(parsed.Path)
		if fileName == "." || fileName == "/" || fileName == "" {
			fileName = "file.jpg"
		}
		if filepath.Ext(fileName) == "" {
			fileName += ".jpg"
		}
		return raw, fileName, int64(len(raw)), nil
	}

	abs, err := filepath.Abs(pathOrURL)
	if err != nil {
		return nil, "", 0, err
	}
	raw, err := os.ReadFile(abs)
	if err != nil {
		return nil, "", 0, err
	}
	stat, err := os.Stat(abs)
	if err != nil {
		return nil, "", 0, err
	}
	return raw, filepath.Base(abs), stat.Size(), nil
}
