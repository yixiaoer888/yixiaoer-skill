package publish

import (
	"fmt"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const shipinhaoCoverAutoCompressMaxBytes int64 = 512 * 1024

type CoverCompressionEvent struct {
	Path         string `json:"path"`
	Source       string `json:"source"`
	FromKey      string `json:"fromKey,omitempty"`
	ToKey        string `json:"toKey,omitempty"`
	OriginalSize int64  `json:"originalSize"`
	Size         int64  `json:"size"`
	Status       string `json:"status"`
}

func materializeShipinhaoCoverCompression(apiClient *api.Client, input ExecuteInput) (map[string]interface{}, []CoverCompressionEvent, error) {
	if apiClient == nil || publishmod.NormalizePublishType(input.PublishType) != "video" {
		return input.Payload, nil, nil
	}
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil || platformutil.CanonicalKey(platform) != "shipinhao" {
		return input.Payload, nil, nil
	}
	payload := cloneMap(input.Payload)
	publishArgs := publishmod.ExtractPublishArgs(payload)
	if publishArgs == nil {
		return payload, nil, nil
	}

	var events []CoverCompressionEvent
	uploadedBySource := map[string]api.UploadResult{}
	materialize := func(resource map[string]interface{}, path string, owner map[string]interface{}) error {
		if resource == nil {
			return nil
		}
		source := resourceSource(resource)
		if source == "" {
			return nil
		}
		meta, _, err := api.InspectUpload(source, true)
		if err != nil {
			return yxerrors.Usage("failed to inspect shipinhao cover before publish", map[string]interface{}{
				"path":   path,
				"source": source,
				"cause":  err.Error(),
			}).WithCategory("media_compress").
				WithHint("请确认视频号封面 source/path/localPath/filePath 指向可读取的图片文件或可下载 URL。")
		}
		if meta.Size <= shipinhaoCoverAutoCompressMaxBytes {
			return nil
		}
		uploaded, ok := uploadedBySource[source]
		if !ok {
			uploaded, err = apiClient.UploadWithOptions(source, "cloud-publish", true, api.UploadOptions{MaxImageBytes: shipinhaoCoverAutoCompressMaxBytes})
			if err != nil {
				return err
			}
			uploadedBySource[source] = uploaded
		}
		fromKey := stringField(resource, "key")
		applyUploadResultToResource(resource, uploaded)
		if owner != nil {
			owner["coverKey"] = uploaded.Key
		}
		events = append(events, CoverCompressionEvent{
			Path:         path,
			Source:       source,
			FromKey:      fromKey,
			ToKey:        uploaded.Key,
			OriginalSize: firstPositiveInt64(uploaded.OriginalSize, meta.Size),
			Size:         uploaded.Size,
			Status:       "compressed_uploaded",
		})
		return nil
	}

	if err := materialize(objectField(publishArgs, "cover"), "publishArgs.cover", publishArgs); err != nil {
		return payload, events, err
	}
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	for i, item := range accountForms {
		form, _ := item.(map[string]interface{})
		if form == nil {
			continue
		}
		formPath := fmt.Sprintf("publishArgs.accountForms[%d]", i)
		if err := materialize(objectField(form, "cover"), formPath+".cover", form); err != nil {
			return payload, events, err
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf != nil {
			if err := materialize(objectField(cpf, "cover"), formPath+".contentPublishForm.cover", cpf); err != nil {
				return payload, events, err
			}
		}
	}
	return payload, events, nil
}

func resourceSource(resource map[string]interface{}) string {
	for _, key := range []string{"source", "path", "localPath", "filePath"} {
		if source := strings.TrimSpace(stringField(resource, key)); source != "" {
			return source
		}
	}
	return ""
}

func applyUploadResultToResource(resource map[string]interface{}, uploaded api.UploadResult) {
	resource["key"] = uploaded.Key
	resource["bucket"] = uploaded.Bucket
	resource["contentType"] = uploaded.ContentType
	resource["size"] = uploaded.Size
	if uploaded.Width > 0 {
		resource["width"] = uploaded.Width
	}
	if uploaded.Height > 0 {
		resource["height"] = uploaded.Height
	}
	if uploaded.Format != "" {
		resource["format"] = uploaded.Format
	}
	if uploaded.Compressed {
		resource["compressed"] = true
	}
	if uploaded.OriginalSize > 0 {
		resource["originalSize"] = uploaded.OriginalSize
	}
	for _, key := range []string{"source", "path", "localPath", "filePath"} {
		delete(resource, key)
	}
}

func firstPositiveInt64(values ...int64) int64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}
