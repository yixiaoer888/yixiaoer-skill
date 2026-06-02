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
var placeholderPattern = regexp.MustCompile(`^<[^<>]+>$`)

const shipinghaoCoverMaxBytes = 512 * 1024

func RequireStandardPayload(payload map[string]interface{}) error {
	if payload == nil {
		return yxerrors.Usage("Standard publish payload is required", []string{
			"missing payload body",
		}).WithHint("请使用标准请求体：顶层包含 action、publishType、platforms、publishArgs。")
	}
	if _, exists := payload["accountForms"]; exists {
		return yxerrors.Usage("Legacy publish payload is not supported", []string{
			"top-level accountForms is deprecated",
		}).WithHint("请改用标准请求体：顶层保留 action、publishType、platforms、publishArgs，账号数据放到 publishArgs.accountForms[]。")
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
	NormalizePlatformSpecificFields(publishType, platforms, payload)
	NormalizeScheduledTimes(payload, &result.Errors)
	rejectTemplatePlaceholders(payload, &result.Errors)
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
			if shouldIgnoreExternalURLPath(path) {
				return
			}
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

func shouldIgnoreExternalURLPath(path string) bool {
	if strings.Contains(path, ".raw.") || strings.HasSuffix(path, ".raw") {
		return true
	}
	if strings.Contains(path, ".shopping_cart[") && strings.Contains(path, ".images[") {
		return true
	}
	if strings.Contains(path, ".shoppingCart[") && strings.Contains(path, ".images[") {
		return true
	}
	return false
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

func NormalizePlatformSpecificFields(publishType string, platforms []string, payload map[string]interface{}) {
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		return
	}

	publishType = NormalizePublishType(publishType)
	platformSet := map[string]bool{}
	for _, platform := range platforms {
		platformSet[strings.TrimSpace(platform)] = true
	}

	content, _ := payload["content"].(string)
	for _, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}

		if publishType == "imageText" && (platformSet["抖音"] || platformSet["小红书"]) {
			normalizeTopicHTML(payload, cpf, content)
		}
		if publishType == "video" && platformSet["抖音"] {
			normalizeDouyinShoppingCart(cpf)
		}
	}
}

func normalizeTopicHTML(publishArgs, cpf map[string]interface{}, publishArgsContent string) {
	tags, ok := cpf["tags"].([]interface{})
	if !ok || len(tags) == 0 {
		return
	}
	description := strings.TrimSpace(stringField(cpf, "description"))
	content := strings.TrimSpace(publishArgsContent)
	if content == "" {
		content = description
	}
	if description == "" && content == "" {
		return
	}

	finalHTML := firstNonEmptyTopicHTML(description, content)
	if finalHTML == "" {
		baseText := firstNonEmptyString(content, description)
		finalHTML = buildTopicHTML(baseText, tags)
	}
	if finalHTML == "" {
		return
	}
	cpf["description"] = finalHTML
	publishArgs["content"] = finalHTML
}

func firstNonEmptyTopicHTML(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && strings.Contains(strings.ToLower(value), "<topic") {
			return value
		}
	}
	return ""
}

func buildTopicHTML(description string, tags []interface{}) string {
	var topicParts []string
	for _, item := range tags {
		tag := strings.TrimSpace(fmt.Sprint(item))
		if tag == "" {
			continue
		}
		if !strings.HasPrefix(tag, "#") {
			tag = "#" + tag
		}
		text := strings.TrimPrefix(tag, "#")
		if text == "" {
			continue
		}
		topicParts = append(topicParts, fmt.Sprintf(`<topic text="%s">%s</topic>`, text, tag))
	}
	if len(topicParts) == 0 {
		return strings.TrimSpace(description)
	}
	descHTML := strings.TrimSpace(description)
	if descHTML == "" {
		return "<p>" + strings.Join(topicParts, "") + "</p>"
	}
	return "<p>" + descHTML + "</p><p>" + strings.Join(topicParts, "") + "</p>"
}

func normalizeDouyinShoppingCart(cpf map[string]interface{}) {
	if cpf == nil {
		return
	}
	if value, ok := cpf["shoppingCart"]; ok {
		if _, exists := cpf["shopping_cart"]; !exists {
			cpf["shopping_cart"] = value
		}
		delete(cpf, "shoppingCart")
	}
	items, ok := cpf["shopping_cart"].([]interface{})
	if !ok {
		return
	}
	for _, item := range items {
		cart, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if cart["data"] == nil {
			data := map[string]interface{}{}
			for _, key := range []string{"yixiaoerId", "yixiaoerName", "raw"} {
				if value, exists := cart[key]; exists {
					data[key] = value
					delete(cart, key)
				}
			}
			if len(data) > 0 {
				cart["data"] = data
			}
		}
		if cart["images"] == nil {
			if data, _ := cart["data"].(map[string]interface{}); data != nil {
				if raw, _ := data["raw"].(map[string]interface{}); raw != nil {
					if images := extractShoppingCartImages(raw); len(images) > 0 {
						cart["images"] = images
					}
				}
			}
		}
	}
}

func extractShoppingCartImages(raw map[string]interface{}) []interface{} {
	candidates := [][]interface{}{}
	for _, key := range []string{"images", "imgs", "goods_imgs"} {
		if items, ok := raw[key].([]interface{}); ok && len(items) > 0 {
			candidates = append(candidates, items)
		}
	}
	for _, items := range candidates {
		var urls []interface{}
		for _, item := range items {
			switch typed := item.(type) {
			case string:
				if strings.TrimSpace(typed) != "" {
					urls = append(urls, typed)
				}
			case map[string]interface{}:
				for _, key := range []string{"url", "src"} {
					if value := strings.TrimSpace(fmt.Sprint(typed[key])); value != "" && value != "<nil>" {
						urls = append(urls, value)
						break
					}
				}
			}
		}
		if len(urls) > 0 {
			return urls
		}
	}
	return nil
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
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

func rejectTemplatePlaceholders(value interface{}, errors *[]string) {
	walk(value, func(current interface{}, path string) {
		text, ok := current.(string)
		if !ok {
			return
		}
		text = strings.TrimSpace(text)
		if !placeholderPattern.MatchString(text) {
			return
		}
		*errors = append(*errors, fmt.Sprintf("%s: unresolved template placeholder %q; run prepare/schema fields (and schema get if needed) and replace template values before validate/publish", strings.TrimPrefix(path, "$."), text))
	}, "$")
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
