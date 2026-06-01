package publish

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type PreflightResult struct {
	AccountIDs []string
	Errors     []string
}

var externalURLPattern = regexp.MustCompile(`(?i)^https?://`)

const shipinghaoCoverMaxBytes = 512 * 1024

func RequireStandardPayload(payload map[string]interface{}) error {
	if payload == nil {
		return yxerrors.Usage("Standard publish payload is required", []string{
			"missing payload body",
		}).WithHint("请使用标准请求体：顶层包含 action、publishType、platforms、publishArgs。")
	}
	publishArgs, ok := payload["publishArgs"].(map[string]interface{})
	if !ok || publishArgs == nil {
		return yxerrors.Usage("Standard publish payload is required", []string{
			"missing publishArgs object",
		}).WithHint("请使用标准请求体：顶层包含 action、publishType、platforms、publishArgs，业务字段放在 publishArgs.accountForms[].contentPublishForm。")
	}
	if _, ok := publishArgs["accountForms"].([]interface{}); !ok {
		return yxerrors.Usage("Standard publish payload is required", []string{
			"publishArgs.accountForms must be a non-empty array",
		}).WithHint("请将账号发布数据放到 publishArgs.accountForms[] 下，不再支持顶层 accountForms 或直接内层表单结构。")
	}
	return nil
}

func Preflight(publishType string, platforms []string, payload map[string]interface{}) PreflightResult {
	var result PreflightResult
	publishType = NormalizePublishType(publishType)
	if err := RequireStandardPayload(payload); err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result
	}
	payload = ValidateAndExtractPublishArgs(publishType, platforms, payload, &result.Errors)
	NormalizeStandardPublishArgs(payload)
	NormalizeScheduledTimes(payload, &result.Errors)
	if publishType != "video" && publishType != "imageText" && publishType != "article" {
		result.Errors = append(result.Errors, fmt.Sprintf("publish type %q is not supported; expected video, imageText, or article", publishType))
	}
	if len(platforms) == 0 {
		result.Errors = append(result.Errors, "at least one target platform is required")
	}
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		result.Errors = append(result.Errors, "payload.accountForms must be a non-empty array")
		return result
	}

	for i, item := range accountForms {
		form, ok := item.(map[string]interface{})
		formPath := fmt.Sprintf("accountForms[%d]", i)
		if !ok {
			result.Errors = append(result.Errors, formPath+": must be an object")
			continue
		}
		accountID := stringField(form, "platformAccountId")
		if accountID == "" {
			accountID = stringField(form, "account_id")
		}
		if accountID == "" {
			result.Errors = append(result.Errors, formPath+": missing platformAccountId")
		} else {
			result.AccountIDs = append(result.AccountIDs, accountID)
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			result.Errors = append(result.Errors, formPath+": missing contentPublishForm")
		}

		switch publishType {
		case "video":
			video := objectField(form, "video")
			if video == nil && cpf != nil {
				video = objectField(cpf, "video")
			}
			requireUploadedResource(video, formPath+".video", &result.Errors)
			cover := objectField(form, "cover")
			if cover == nil && cpf != nil {
				cover = objectField(cpf, "cover")
			}
			requireUploadedResource(cover, formPath+".cover", &result.Errors)
			requireCoverKey(form, cpf, cover, formPath, &result.Errors)
			requirePlatformConstraints(platforms, cover, formPath, &result.Errors)
		case "imageText":
			images, _ := form["images"].([]interface{})
			if len(images) == 0 && cpf != nil {
				images, _ = cpf["images"].([]interface{})
			}
			if len(images) == 0 {
				result.Errors = append(result.Errors, formPath+".images: imageText publish requires at least one uploaded image")
			}
			for imageIndex, image := range images {
				imageObj, _ := image.(map[string]interface{})
				requireUploadedResource(imageObj, fmt.Sprintf("%s.images[%d]", formPath, imageIndex), &result.Errors)
			}
			cover := objectField(form, "cover")
			if cover == nil && cpf != nil {
				cover = objectField(cpf, "cover")
			}
			requireUploadedResource(cover, formPath+".cover", &result.Errors)
			requireCoverKey(form, cpf, cover, formPath, &result.Errors)
		case "article":
			if cpf == nil || stringField(cpf, "content") == "" {
				result.Errors = append(result.Errors, formPath+".contentPublishForm.content: article publish requires content")
			}
			cover := objectField(form, "cover")
			if cover == nil && cpf != nil {
				cover = objectField(cpf, "cover")
			}
			requireUploadedResource(cover, formPath+".cover", &result.Errors)
			requireCoverKey(form, cpf, cover, formPath, &result.Errors)
		}

		walk(form, func(value interface{}, path string) {
			if text, ok := value.(string); ok && externalURLPattern.MatchString(text) {
				result.Errors = append(result.Errors, formPath+path[1:]+": external URL is not allowed in publish payload; upload resources first")
			}
		}, "$")

		for _, field := range []string{"location", "music", "collection", "collections", "challenge", "challenges", "goods", "group", "groups", "miniapp", "miniapps", "shopping_cart", "shoppingCart"} {
			if value, ok := form[field]; ok {
				assertRawObject(value, formPath+"."+field, &result.Errors)
			}
			if cpf != nil {
				if value, ok := cpf[field]; ok {
					assertRawObject(value, formPath+".contentPublishForm."+field, &result.Errors)
				}
			}
		}
	}
	return result
}

func ExtractPublishArgs(payload map[string]interface{}) map[string]interface{} {
	if publishArgs, ok := payload["publishArgs"].(map[string]interface{}); ok {
		return publishArgs
	}
	return nil
}

func ValidateAndExtractPublishArgs(publishType string, platforms []string, payload map[string]interface{}, errors *[]string) map[string]interface{} {
	publishArgs, ok := payload["publishArgs"].(map[string]interface{})
	if !ok {
		*errors = append(*errors, "publishArgs: missing required object")
		return nil
	}
	if action := stringField(payload, "action"); action != "publish" {
		*errors = append(*errors, `action: must equal "publish"`)
	}
	if apiType := stringField(payload, "publishType"); apiType == "" {
		*errors = append(*errors, "publishType: missing required field")
	} else if !samePublishType(apiType, publishType) {
		*errors = append(*errors, fmt.Sprintf("publishType: got %q, expected %q", apiType, publishType))
	}
	if rawPlatforms, ok := payload["platforms"].([]interface{}); !ok || len(rawPlatforms) == 0 {
		*errors = append(*errors, "platforms: must be a non-empty array")
	} else {
		for i, item := range rawPlatforms {
			if text, ok := item.(string); !ok || strings.TrimSpace(text) == "" {
				*errors = append(*errors, fmt.Sprintf("platforms[%d]: must be a non-empty string", i))
			}
		}
	}
	if channel := stringField(payload, "publishChannel"); channel != "" && channel != "cloud" && channel != "local" {
		*errors = append(*errors, `publishChannel: must be "cloud" or "local"`)
	}
	if stringField(payload, "publishChannel") == "local" && stringField(payload, "clientId") == "" {
		*errors = append(*errors, "clientId: required when publishChannel is local")
	}
	if _, exists := payload["cover"]; exists {
		cover := objectField(payload, "cover")
		if cover == nil {
			*errors = append(*errors, "cover: expected object")
		} else {
			requireUploadedResource(cover, "cover", errors)
			if coverKey := stringField(payload, "coverKey"); coverKey != "" {
				if key := stringField(cover, "key"); key != "" && key != coverKey {
					*errors = append(*errors, "coverKey: must match cover.key")
				}
			}
		}
	}
	for _, field := range []string{"coverKey", "taskSetId", "desc", "clientId"} {
		if value, ok := payload[field]; ok && !matchesString(value) {
			*errors = append(*errors, fmt.Sprintf("%s: expected string", field))
		}
	}
	if value, ok := payload["isDraft"]; ok {
		if _, ok := value.(bool); !ok {
			*errors = append(*errors, "isDraft: expected boolean")
		}
	}
	return publishArgs
}

func NormalizeStandardPublishArgs(payload map[string]interface{}) {
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		return
	}
	for _, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		copyIfMissing(form, payload, "video")
		copyIfMissing(form, payload, "images")
		copyIfMissing(form, payload, "cover")
		copyIfMissing(form, payload, "coverKey")

		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}
		copyIfMissing(form, cpf, "images")
		copyIfMissing(cpf, payload, "content")
	}
}

func requirePlatformConstraints(platforms []string, cover map[string]interface{}, formPath string, errors *[]string) {
	for _, platform := range platforms {
		switch strings.TrimSpace(platform) {
		case "视频号", "微信视频号", "shipinghao":
			requireShipinghaoCoverSize(cover, formPath, errors)
		}
	}
}

func requireShipinghaoCoverSize(cover map[string]interface{}, formPath string, errors *[]string) {
	if cover == nil {
		return
	}
	size, ok := integerField(cover, "size")
	if !ok {
		return
	}
	if size > shipinghaoCoverMaxBytes {
		*errors = append(*errors, fmt.Sprintf("%s.cover.size: 视频号封面不能超过 512KB，当前为 %d bytes", formPath, size))
	}
}

func copyIfMissing(dst, src map[string]interface{}, key string) {
	if dst == nil || src == nil {
		return
	}
	if _, exists := dst[key]; exists {
		return
	}
	if value, exists := src[key]; exists {
		dst[key] = value
	}
}

func requireUploadedResource(resource map[string]interface{}, pathLabel string, errors *[]string) {
	if resource == nil {
		*errors = append(*errors, pathLabel+": missing uploaded resource object")
		return
	}
	if empty(resource["key"]) {
		*errors = append(*errors, fmt.Sprintf("%s: missing uploaded resource field %q", pathLabel, "key"))
	}
	walk(resource, func(value interface{}, path string) {
		if text, ok := value.(string); ok && externalURLPattern.MatchString(text) {
			*errors = append(*errors, pathLabel+path[1:]+`: external URL is not allowed; run "yxer upload" and use the returned key`)
		}
	}, "$")
}

func requireCoverKey(form, cpf map[string]interface{}, cover map[string]interface{}, formPath string, errors *[]string) {
	coverKey := stringField(form, "coverKey")
	if coverKey == "" && cpf != nil {
		coverKey = stringField(cpf, "coverKey")
	}
	if coverKey == "" {
		*errors = append(*errors, formPath+`: missing coverKey`)
		return
	}
	if cover != nil {
		key := stringField(cover, "key")
		if key != "" && key != coverKey {
			*errors = append(*errors, formPath+`.coverKey: must match cover.key`)
		}
	}
}

func assertRawObject(value interface{}, pathLabel string, errors *[]string) {
	if items, ok := value.([]interface{}); ok {
		for i, item := range items {
			assertRawObject(item, fmt.Sprintf("%s[%d]", pathLabel, i), errors)
		}
		return
	}
	obj, ok := value.(map[string]interface{})
	if !ok {
		return
	}
	if nested, ok := obj["data"].(map[string]interface{}); ok && nested != nil {
		assertRawObject(nested, pathLabel+".data", errors)
		return
	}
	hasIdentity := obj["yixiaoerId"] != nil || obj["yixiaoerName"] != nil || obj["id"] != nil || obj["name"] != nil
	raw, rawOK := obj["raw"].(map[string]interface{})
	if hasIdentity && (!rawOK || raw == nil) {
		*errors = append(*errors, pathLabel+`: dynamic platform object must include complete "raw" data from a yxer query command`)
	}
}

func walk(value interface{}, visit func(interface{}, string), currentPath string) {
	visit(value, currentPath)
	switch typed := value.(type) {
	case []interface{}:
		for i, child := range typed {
			walk(child, visit, fmt.Sprintf("%s[%d]", currentPath, i))
		}
	case map[string]interface{}:
		for key, child := range typed {
			walk(child, visit, currentPath+"."+key)
		}
	}
}

func objectField(obj map[string]interface{}, key string) map[string]interface{} {
	value, _ := obj[key].(map[string]interface{})
	return value
}

func stringField(obj map[string]interface{}, key string) string {
	value := obj[key]
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func matchesString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func samePublishType(left, right string) bool {
	return TypeKey(left) == TypeKey(right)
}

func NormalizePublishType(publishType string) string {
	return strings.TrimSpace(publishType)
}

func TypeKey(publishType string) string {
	publishType = NormalizePublishType(publishType)
	return publishType
}

func empty(value interface{}) bool {
	if value == nil {
		return true
	}
	if text, ok := value.(string); ok {
		return text == ""
	}
	return false
}

func NormalizeScheduledTimes(value interface{}, errors *[]string) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, child := range typed {
			if key == "scheduledTime" {
				normalized, err := normalizeScheduledTime(child)
				if err != "" {
					*errors = append(*errors, "scheduledTime: "+err)
				} else {
					typed[key] = normalized
				}
				continue
			}
			NormalizeScheduledTimes(child, errors)
		}
	case []interface{}:
		for _, child := range typed {
			NormalizeScheduledTimes(child, errors)
		}
	}
}

func normalizeScheduledTime(value interface{}) (interface{}, string) {
	switch typed := value.(type) {
	case float64:
		if typed != math.Trunc(typed) {
			return nil, "must be an integer 13-digit Unix timestamp in milliseconds"
		}
		return normalizeScheduledTimeInt64(int64(typed))
	case int64:
		return normalizeScheduledTimeInt64(typed)
	case int:
		return normalizeScheduledTimeInt64(int64(typed))
	case string:
		text := strings.TrimSpace(typed)
		if len(text) != 13 {
			return nil, "must be a 13-digit Unix timestamp in milliseconds"
		}
		var n int64
		for _, r := range text {
			if r < '0' || r > '9' {
				return nil, "must contain digits only"
			}
			n = n*10 + int64(r-'0')
		}
		return normalizeScheduledTimeInt64(n)
	default:
		return nil, "must be a 13-digit Unix timestamp in milliseconds"
	}
}

func normalizeScheduledTimeInt64(value int64) (float64, string) {
	if value < 1_000_000_000_000 || value > 9_999_999_999_999 {
		return 0, "must be a 13-digit Unix timestamp in milliseconds"
	}
	return float64(value / 1000), ""
}

func integerField(obj map[string]interface{}, key string) (int64, bool) {
	value, ok := obj[key]
	if !ok || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case float64:
		if typed != math.Trunc(typed) {
			return 0, false
		}
		return int64(typed), true
	case int:
		return int64(typed), true
	case int64:
		return typed, true
	default:
		return 0, false
	}
}
