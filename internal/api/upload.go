package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/media"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const maxRemoteUploadDownloadSize int64 = 5 * 1024 * 1024 * 1024

var uploadHTTPClient = &http.Client{}

type UploadResult struct {
	Key          string  `json:"key"`
	URL          string  `json:"url,omitempty"`
	ContentType  string  `json:"contentType"`
	Bucket       string  `json:"bucket"`
	Size         int64   `json:"size,omitempty"`
	Width        int     `json:"width,omitempty"`
	Height       int     `json:"height,omitempty"`
	Duration     float64 `json:"duration,omitempty"`
	Format       string  `json:"format,omitempty"`
	Compressed   bool    `json:"compressed,omitempty"`
	OriginalSize int64   `json:"originalSize,omitempty"`
}

type UploadOptions struct {
	MaxImageBytes int64
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
	return c.UploadWithOptions(pathOrURL, bucket, autoMeta, UploadOptions{})
}

func (c *Client) UploadWithOptions(pathOrURL, bucket string, autoMeta bool, opts UploadOptions) (UploadResult, error) {
	if bucket == "" {
		bucket = "cloud-publish"
	}
	result, fileName, buffer, err := c.inspectUpload(pathOrURL, autoMeta)
	if err != nil {
		return UploadResult{}, err
	}
	result, fileName, buffer, err = applyUploadOptions(result, fileName, buffer, opts)
	if err != nil {
		return UploadResult{}, err
	}

	params := map[string]string{
		"fileKey":     uploadObjectName(fileName),
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

	req, err := http.NewRequest(http.MethodPut, serviceURL, bytes.NewReader(buffer))
	if err != nil {
		return UploadResult{}, err
	}
	req.Header.Set("Content-Type", result.ContentType)
	resp, err := uploadHTTPClient.Do(req)
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

func applyUploadOptions(result UploadResult, fileName string, buffer []byte, opts UploadOptions) (UploadResult, string, []byte, error) {
	if opts.MaxImageBytes <= 0 || !strings.HasPrefix(result.ContentType, "image/") || result.Size <= opts.MaxImageBytes {
		return result, fileName, buffer, nil
	}
	compressed, ok, err := media.CompressImageToMaxBytes(buffer, opts.MaxImageBytes)
	if err != nil {
		return UploadResult{}, "", nil, yxerrors.Usage("failed to compress image for upload", map[string]interface{}{
			"fileName": fileName,
			"limit":    opts.MaxImageBytes,
			"cause":    err.Error(),
		}).WithCategory("media_compress").
			WithHint("请确认封面是可解码的图片文件，支持 jpg/png 等常见格式。")
	}
	if !ok {
		return UploadResult{}, "", nil, yxerrors.Usage("image cannot be compressed below platform limit", map[string]interface{}{
			"fileName":     fileName,
			"limitBytes":   opts.MaxImageBytes,
			"originalSize": result.Size,
		}).WithCategory("media_compress").
			WithHint("请改用尺寸更小或内容更简单的封面图后再上传。")
	}
	result.OriginalSize = result.Size
	result.Size = int64(len(compressed.Data))
	result.Width = compressed.Width
	result.Height = compressed.Height
	result.ContentType = compressed.ContentType
	result.Format = compressed.Format
	result.Compressed = true
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".jpg"
	return result, fileName, compressed.Data, nil
}

func (c *Client) inspectUpload(pathOrURL string, autoMeta bool) (UploadResult, string, []byte, error) {
	contentType := DetectContentType(pathOrURL)
	buffer, fileName, size, err := c.readUploadContent(pathOrURL)
	if err != nil {
		return UploadResult{}, "", nil, err
	}
	result, err := buildUploadMetadata(pathOrURL, buffer, fileName, size, contentType, autoMeta)
	if err != nil {
		return UploadResult{}, "", nil, err
	}
	return result, fileName, buffer, nil
}

func (c *Client) readUploadContent(pathOrURL string) ([]byte, string, int64, error) {
	return readUploadContentWithProxy(pathOrURL, baseURL(c.cfg))
}

func (c *Client) StableURL(bucket, key string) (string, error) {
	if bucket == "" {
		bucket = "cloud-publish"
	}
	var result interface{}
	if err := c.Get(Query("/storages/"+bucket+"/stable-url", map[string]string{"fileKey": key}), &result); err != nil {
		return "", err
	}
	if url, ok := stableURLFromValue(result); ok {
		return url, nil
	}
	return "", yxerrors.Remote("invalid stable url response", result).
		WithCategory("remote_response")
}

func stableURLFromValue(value interface{}) (string, bool) {
	switch typed := value.(type) {
	case string:
		url := strings.TrimSpace(typed)
		return url, url != ""
	case map[string]interface{}:
		if data, exists := typed["data"]; exists {
			return stableURLFromValue(data)
		}
		for _, key := range []string{"url", "stableUrl", "hostUrl"} {
			if value, exists := typed[key]; exists {
				if url, ok := stableURLFromValue(value); ok {
					return url, true
				}
			}
		}
	}
	return "", false
}

func uploadObjectName(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return "file"
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	safeBase := asciiSafeFileStem(base)
	if safeBase == "" {
		safeBase = "file"
	}
	if ext == "" {
		return safeBase
	}
	return safeBase + ext
}

func asciiSafeFileStem(name string) string {
	var builder strings.Builder
	lastUnderscore := false
	hasNonASCII := false
	for _, r := range name {
		switch {
		case r > 127:
			hasNonASCII = true
			if !lastUnderscore && builder.Len() > 0 {
				builder.WriteByte('_')
				lastUnderscore = true
			}
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			builder.WriteRune(r)
			lastUnderscore = false
		case r == '-' || r == '_':
			if !lastUnderscore && builder.Len() > 0 {
				builder.WriteRune(r)
				lastUnderscore = true
			}
		default:
			if !lastUnderscore && builder.Len() > 0 {
				builder.WriteByte('_')
				lastUnderscore = true
			}
		}
	}
	safe := strings.Trim(builder.String(), "_-")
	if hasNonASCII {
		sum := sha1.Sum([]byte(name))
		suffix := hex.EncodeToString(sum[:])[:8]
		if safe == "" {
			return "file_" + suffix
		}
		return safe + "_" + suffix
	}
	return safe
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
	return readUploadContentWithProxy(pathOrURL, "")
}

func readUploadContentWithProxy(pathOrURL, proxyBaseURL string) ([]byte, string, int64, error) {
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		raw, size, err := downloadRemoteUploadContent(pathOrURL, pathOrURL)
		if err != nil && shouldRetryRemoteUploadViaProxy(pathOrURL, proxyBaseURL) {
			proxyURL := Query(strings.TrimRight(proxyBaseURL, "/")+"/storages/proxy-url", map[string]string{"url": pathOrURL})
			raw, size, err = downloadRemoteUploadContent(proxyURL, pathOrURL)
		}
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
		return raw, fileName, size, nil
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

func downloadRemoteUploadContent(fetchURL, sourceURL string) ([]byte, int64, error) {
	resp, err := uploadHTTPClient.Get(fetchURL)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, 0, yxerrors.Remote("HTTP error downloading file during sync upload", map[string]interface{}{
			"statusCode": resp.StatusCode,
			"url":        sourceURL,
			"fetchUrl":   fetchURL,
		}).WithCategory("remote_download")
	}
	if resp.ContentLength > maxRemoteUploadDownloadSize {
		return nil, 0, yxerrors.Usage("remote file exceeds size limit", map[string]interface{}{
			"limitBytes":   maxRemoteUploadDownloadSize,
			"contentBytes": resp.ContentLength,
			"url":          sourceURL,
		}).WithHint("请改用更小的素材文件，或先下载到本地压缩后再上传。")
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxRemoteUploadDownloadSize+1))
	if err != nil {
		return nil, 0, err
	}
	if int64(len(raw)) > maxRemoteUploadDownloadSize {
		return nil, 0, yxerrors.Usage("remote file exceeds size limit", map[string]interface{}{
			"limitBytes": maxRemoteUploadDownloadSize,
			"url":        sourceURL,
		}).WithHint("请改用更小的素材文件，或先下载到本地压缩后再上传。")
	}
	return raw, int64(len(raw)), nil
}

func shouldRetryRemoteUploadViaProxy(pathOrURL, proxyBaseURL string) bool {
	if strings.TrimSpace(proxyBaseURL) == "" {
		return false
	}
	lower := strings.ToLower(pathOrURL)
	if strings.Contains(lower, "/storages/proxy-url?url=") || strings.HasPrefix(lower, "https://view.yixiaoer.cn/proxy-url?url=") {
		return false
	}
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}
