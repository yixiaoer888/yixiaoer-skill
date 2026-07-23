package publish

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type PreflightResult struct {
	AccountIDs []string
	Errors     []string
}

type NormalizationEvent struct {
	Field        string `json:"field"`
	Path         string `json:"path"`
	Action       string `json:"action"`
	Message      string `json:"message"`
	QueryCommand string `json:"queryCommand,omitempty"`
}

type TopicHTMLFields struct {
	HasTopics      bool
	HasDescription bool
	HasContent     bool
}

type TopicHTMLPolicy map[string]TopicHTMLFields

var externalURLPattern = regexp.MustCompile(`(?i)^https?://`)
var placeholderPattern = regexp.MustCompile(`^<[^<>]+>$`)
var hashtagPattern = regexp.MustCompile(`#([^\s#<]+)`)

const shipinhaoImageMaxBytes = 512 * 1024

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
	return PreflightWithTopicHTMLPolicy(publishType, platforms, payload, nil)
}

func PreflightWithTopicHTMLPolicy(publishType string, platforms []string, payload map[string]interface{}, topicPolicy TopicHTMLPolicy) PreflightResult {
	return PreflightWithTopicHTMLPolicyAndTrace(publishType, platforms, payload, topicPolicy, nil)
}

func PreflightWithTopicHTMLPolicyAndTrace(publishType string, platforms []string, payload map[string]interface{}, topicPolicy TopicHTMLPolicy, normalizations *[]NormalizationEvent) PreflightResult {
	var result PreflightResult
	publishType = NormalizePublishType(publishType)
	if err := RequireStandardPayload(payload); err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result
	}
	resolveStandardPayloadResourceMetadata(payload, &result.Errors)
	payload = ValidateAndExtractPublishArgs(publishType, platforms, payload, &result.Errors)
	NormalizeStandardPublishArgs(payload, publishType, platforms)
	normalizePlatformSpecificFields(publishType, platforms, payload, topicPolicy, true, normalizations)
	NormalizeScheduledTimes(payload, &result.Errors)
	rejectTemplatePlaceholders(payload, &result.Errors)
	if publishType != "video" && publishType != "imageText" && publishType != "article" {
		result.Errors = append(result.Errors, fmt.Sprintf("publish type %q is not supported; expected video, imageText, or article", publishType))
	}
	if len(platforms) == 0 {
		result.Errors = append(result.Errors, "at least one target platform is required")
	}
	weixinAccountArticle := isWeixinAccountArticlePublish(platforms, publishType)
	weixinPlatformForm := weixinAccountArticleForm(payload)
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		result.Errors = append(result.Errors, "payload.accountForms must be a non-empty array")
		return result
	}

	if weixinAccountArticle {
		preflightWeixinAccountArticle(payload, weixinPlatformForm, &result.Errors)
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
		if cpf == nil && !weixinAccountArticle {
			result.Errors = append(result.Errors, formPath+": missing contentPublishForm")
		}

		switch publishType {
		case "video":
			if _, exists := form["horizontalCover"]; exists {
				result.Errors = append(result.Errors, formPath+".horizontalCover: unexpected field; use contentPublishForm.horizontalCover or publishArgs.horizontalCover")
			}
			video := objectField(form, "video")
			if video == nil && cpf != nil {
				video = objectField(cpf, "video")
			}
			requireUploadedVideoResource(video, formPath+".video", &result.Errors)
			cover := objectField(form, "cover")
			if cover == nil && cpf != nil {
				cover = objectField(cpf, "cover")
			}
			requireUploadedResource(cover, formPath+".cover", &result.Errors)
			horizontalCover := objectField(cpf, "horizontalCover")
			if horizontalCover != nil {
				requireUploadedResource(horizontalCover, formPath+".contentPublishForm.horizontalCover", &result.Errors)
			}
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
			requirePlatformImageTextConstraints(platforms, images, cover, formPath, &result.Errors)
		case "article":
			if !weixinAccountArticle && stringField(payload, "content") == "" {
				result.Errors = append(result.Errors, "publishArgs.content: article publish requires content")
			}
			cover := objectField(form, "cover")
			if cover == nil && cpf != nil {
				cover = objectField(cpf, "cover")
			}
			if weixinAccountArticle {
				if cover != nil {
					requireUploadedResource(cover, formPath+".cover", &result.Errors)
					requireCoverKey(form, cpf, cover, formPath, &result.Errors)
				}
			} else if articleRequiresCover(platforms) {
				requireUploadedResource(cover, formPath+".cover", &result.Errors)
				requireCoverKey(form, cpf, cover, formPath, &result.Errors)
			}
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

func isWeixinAccountArticlePublish(platforms []string, publishType string) bool {
	if NormalizePublishType(publishType) != "article" {
		return false
	}
	for _, platform := range platforms {
		if platformutil.CanonicalKey(platform) == "weixin.account" {
			return true
		}
	}
	return false
}

func weixinAccountArticleForm(publishArgs map[string]interface{}) map[string]interface{} {
	if publishArgs == nil {
		return nil
	}
	platformForms, _ := publishArgs["platformForms"].(map[string]interface{})
	if platformForms == nil {
		return nil
	}
	for _, key := range []string{"微信公众号", "weixin.account"} {
		form, _ := platformForms[key].(map[string]interface{})
		if form != nil {
			return form
		}
	}
	return nil
}

func preflightWeixinAccountArticle(payload, platformForm map[string]interface{}, errors *[]string) {
	if platformForm == nil {
		*errors = append(*errors, `publishArgs.platformForms["微信公众号"]: missing required platform form`)
		return
	}
	articles, _ := platformForm["articles"].([]interface{})
	if len(articles) == 0 {
		*errors = append(*errors, `publishArgs.platformForms["微信公众号"].articles: must contain 1-8 articles`)
		return
	}
	for i, item := range articles {
		article, _ := item.(map[string]interface{})
		if article == nil {
			*errors = append(*errors, fmt.Sprintf(`publishArgs.platformForms["微信公众号"].articles[%d]: must be an object`, i))
			continue
		}
		articlePath := fmt.Sprintf(`publishArgs.platformForms["微信公众号"].articles[%d]`, i)
		if strings.TrimSpace(stringField(article, "title")) == "" {
			*errors = append(*errors, articlePath+".title: missing required field")
		}
		if strings.TrimSpace(stringField(article, "content")) == "" {
			*errors = append(*errors, articlePath+".content: missing required field")
		}
		requireUploadedResource(objectField(article, "cover"), articlePath+".cover", errors)
		if categories, ok := article["categories"]; ok {
			assertRawObject(categories, articlePath+".categories", errors)
		}
	}
}

func shouldIgnoreExternalURLPath(path string) bool {
	if strings.Contains(path, ".raw.") || strings.HasSuffix(path, ".raw") {
		return true
	}
	if isMusicMetadataURLPath(path) {
		return true
	}
	if strings.Contains(path, ".shopping_cart[") && strings.Contains(path, ".images[") {
		return true
	}
	if strings.Contains(path, ".shoppingCart[") && strings.Contains(path, ".images[") {
		return true
	}
	if strings.Contains(path, ".shopping_cart[") && strings.HasSuffix(path, ".yixiaoerImageUrl") {
		return true
	}
	if strings.Contains(path, ".shoppingCart[") && strings.HasSuffix(path, ".yixiaoerImageUrl") {
		return true
	}
	if strings.HasSuffix(path, ".yixiaoerImageUrl") || strings.HasSuffix(path, ".yixiaoerImage") {
		return true
	}
	return false
}

func isMusicMetadataURLPath(path string) bool {
	if !strings.HasSuffix(path, ".url") && !strings.HasSuffix(path, ".playUrl") {
		return false
	}
	return strings.Contains(path, ".music.") || strings.Contains(path, ".music[")
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

// NormalizeStandardPayload applies the content/resource and platform-specific
// normalization that publish performs before schema validation, so validate,
// dry-run, and publish all evaluate the identical normalized payload. It mutates
// payload in place (pass a clone to preserve the original) and returns the
// resolved publishArgs. platforms must use canonical Chinese platform names so
// shared rules (e.g. description topic HTML) trigger consistently.
func NormalizeStandardPayload(publishType string, platforms []string, payload map[string]interface{}) map[string]interface{} {
	return NormalizeStandardPayloadWithTopicHTMLPolicy(publishType, platforms, payload, nil)
}

func NormalizeStandardPayloadWithTopicHTMLPolicy(publishType string, platforms []string, payload map[string]interface{}, topicPolicy TopicHTMLPolicy) map[string]interface{} {
	return normalizeStandardPayloadInternal(publishType, platforms, payload, topicPolicy, true, false, nil)
}

func NormalizeStandardPayloadForSchemaValidation(publishType string, platforms []string, payload map[string]interface{}) map[string]interface{} {
	return NormalizeStandardPayloadForSchemaValidationWithTrace(publishType, platforms, payload, nil)
}

func NormalizeStandardPayloadForSchemaValidationWithTrace(publishType string, platforms []string, payload map[string]interface{}, normalizations *[]NormalizationEvent) map[string]interface{} {
	return normalizeStandardPayloadInternal(publishType, platforms, payload, nil, false, true, normalizations)
}

func normalizeStandardPayloadInternal(publishType string, platforms []string, payload map[string]interface{}, topicPolicy TopicHTMLPolicy, normalizeTopics bool, copyArticleContentToForm bool, normalizations *[]NormalizationEvent) map[string]interface{} {
	publishArgs := ExtractPublishArgs(payload)
	if publishArgs == nil {
		return nil
	}
	NormalizeStandardPublishArgs(publishArgs, publishType, platforms)
	if copyArticleContentToForm && NormalizePublishType(publishType) == "article" {
		copyArticleContentIntoForms(publishArgs, platforms)
	}
	resolveStandardPayloadResourceMetadata(payload, nil)
	normalizePlatformSpecificFields(publishType, platforms, publishArgs, topicPolicy, normalizeTopics, normalizations)
	return publishArgs
}

func NormalizeStandardPublishArgs(payload map[string]interface{}, publishType string, platforms []string) {
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		return
	}
	allowArticleCovers := NormalizePublishType(publishType) != "article" || articleAllowsContentCovers(platforms)
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
		if NormalizePublishType(publishType) == "video" {
			copyIfMissing(cpf, payload, "horizontalCover")
		}
		if allowArticleCovers {
			copyIfMissing(cpf, payload, "covers")
		}
		copyIfMissing(form, cpf, "images")
		copyIfMissing(form, cpf, "cover")
		copyIfMissing(form, cpf, "coverKey")
	}
}

func copyArticleContentIntoForms(payload map[string]interface{}, platforms []string) {
	content := stringField(payload, "content")
	sharedCovers := articleCoversForPayload(payload)
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		return
	}
	for _, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}
		if content != "" {
			copyIfMissing(cpf, payload, "content")
		}
		if articleAllowsContentCovers(platforms) && len(sharedCovers) > 0 {
			if _, exists := cpf["covers"]; !exists {
				cpf["covers"] = sharedCovers
			}
		} else if articleAllowsContentCovers(platforms) {
			if cover := objectField(form, "cover"); cover != nil {
				if _, exists := cpf["covers"]; !exists {
					cpf["covers"] = []interface{}{cover}
				}
			}
		} else {
			delete(cpf, "covers")
		}
	}
}

func articleRequiresCover(platforms []string) bool {
	for _, platform := range platforms {
		switch platformutil.CanonicalKey(platform) {
		case "douban":
			return false
		}
	}
	return true
}

func articleAllowsContentCovers(platforms []string) bool {
	for _, platform := range platforms {
		switch platformutil.CanonicalKey(platform) {
		case "douban":
			return false
		}
	}
	return true
}

func articleCoversForPayload(payload map[string]interface{}) []interface{} {
	if payload == nil {
		return nil
	}
	if covers, ok := payload["covers"].([]interface{}); ok && len(covers) > 0 {
		return covers
	}
	if cover := objectField(payload, "cover"); cover != nil {
		return []interface{}{cover}
	}
	return nil
}

func NormalizePlatformSpecificFields(publishType string, platforms []string, payload map[string]interface{}) {
	NormalizePlatformSpecificFieldsWithTopicHTMLPolicy(publishType, platforms, payload, nil)
}

func NormalizePlatformSpecificFieldsWithTopicHTMLPolicy(publishType string, platforms []string, payload map[string]interface{}, topicPolicy TopicHTMLPolicy) {
	normalizePlatformSpecificFields(publishType, platforms, payload, topicPolicy, true, nil)
}

func normalizePlatformSpecificFields(publishType string, platforms []string, payload map[string]interface{}, topicPolicy TopicHTMLPolicy, normalizeTopics bool, normalizations *[]NormalizationEvent) {
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		return
	}

	publishType = NormalizePublishType(publishType)
	platformSet := map[string]bool{}
	for _, platform := range platforms {
		platform = strings.TrimSpace(platform)
		platformSet[platform] = true
		if canonical := platformutil.CanonicalKey(platform); canonical != "" {
			platformSet[canonical] = true
		}
	}

	topicTarget := topicHTMLTargetField(platforms, topicPolicy)
	if !normalizeTopics {
		topicTarget = ""
	}
	for formIndex, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}

		formPath := fmt.Sprintf("accountForms[%d].contentPublishForm", formIndex)
		normalizeDynamicObjectFields(cpf, formPath, publishType, platformSet, normalizations)
		if topicTarget != "" {
			normalizeTopicHTML(cpf, topicTarget, formPath, normalizations)
		}
		if (publishType == "video" || publishType == "imageText") && isDouyinPlatformSet(platformSet) {
			normalizeDouyinShoppingCart(cpf, formPath, normalizations)
			normalizeDouyinGroupShopping(cpf, formPath, normalizations)
		} else {
			normalizeFlatShoppingCart(cpf, formPath, normalizations)
		}
	}
}

func TopicHTMLPolicyFromSchema(platform string, properties map[string]schema.PropertyView) TopicHTMLPolicy {
	return TopicHTMLPolicy{
		strings.TrimSpace(platform): TopicHTMLFields{
			HasTopics:      properties["topics"].Type != "",
			HasDescription: properties["description"].Type != "",
			HasContent:     properties["content"].Type != "",
		},
	}
}

func topicHTMLTargetField(platforms []string, topicPolicy TopicHTMLPolicy) string {
	hasPolicy := len(topicPolicy) > 0
	hasDescription := !hasPolicy
	for _, platform := range platforms {
		fields, ok := topicPolicy[strings.TrimSpace(platform)]
		if hasPolicy && !ok {
			continue
		}
		hasDescription = hasDescription || fields.HasDescription
	}
	if hasDescription {
		return "description"
	}
	return ""
}

func normalizeTopicHTML(cpf map[string]interface{}, targetField, formPath string, normalizations *[]NormalizationEvent) {
	description := strings.TrimSpace(stringField(cpf, "description"))
	if description == "" {
		return
	}

	finalHTML := firstNonEmptyTopicHTML(description)
	if finalHTML == "" {
		finalHTML = buildTopicHTMLFromDescription(description)
	}
	if finalHTML == "" {
		return
	}
	changed := cpf[targetField] != finalHTML
	cpf[targetField] = finalHTML
	if changed {
		appendNormalization(normalizations, NormalizationEvent{
			Field:   "description",
			Path:    formPath + "." + targetField,
			Action:  "description_topic_html",
			Message: "Normalized description hashtags into platform topic HTML.",
		})
	}
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

func buildTopicHTMLFromDescription(description string) string {
	description = strings.TrimSpace(description)
	matches := hashtagPattern.FindAllStringSubmatch(description, -1)
	if len(matches) == 0 {
		return ""
	}
	var tags []interface{}
	for _, match := range matches {
		if len(match) > 1 && strings.TrimSpace(match[1]) != "" {
			tags = append(tags, "#"+strings.TrimSpace(match[1]))
		}
	}
	baseText := strings.TrimSpace(hashtagPattern.ReplaceAllString(description, ""))
	baseText = strings.Join(strings.Fields(baseText), " ")
	if len(tags) == 0 {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(baseText), "<p") {
		topicOnly := buildTopicHTML("", tags)
		if topicOnly == "" {
			return ""
		}
		return baseText + topicOnly
	}
	return buildTopicHTML(baseText, tags)
}

func normalizeDynamicObjectFields(cpf map[string]interface{}, formPath, publishType string, platformSet map[string]bool, normalizations *[]NormalizationEvent) {
	if cpf == nil {
		return
	}
	for _, field := range []string{
		"category", "location", "music", "collection", "collections", "sub_collection",
		"challenge", "challenges", "mini_app", "miniapp", "miniapps", "sync_apps",
		"game", "hot_event", "group", "groups", "activity",
	} {
		value, exists := cpf[field]
		if !exists {
			continue
		}
		fieldPath := fmt.Sprintf("%s.%s", formPath, field)
		var normalized interface{}
		var changed bool
		switch {
		case field == "location":
			normalized, changed = normalizeLocationValue(value, fieldPath, publishType, platformSet, normalizations)
		case field == "category":
			normalized, changed = normalizePlatformDataValue(value, fieldPath, field, normalizations)
		default:
			normalized, changed = normalizeDynamicObjectValue(value, fieldPath, field, normalizations)
		}
		if changed {
			cpf[field] = normalized
		}
	}
}

func isDouyinPlatformSet(platformSet map[string]bool) bool {
	return platformSet["抖音"] || platformSet["douyin"]
}

func normalizeLocationValue(value interface{}, path, publishType string, platformSet map[string]bool, normalizations *[]NormalizationEvent) (interface{}, bool) {
	if isDouyinPlatformSet(platformSet) {
		return normalizeDouyinLocationValue(value, path, publishType, normalizations)
	}
	return normalizePlatformDataValue(value, path, "location", normalizations)
}

func normalizeDouyinLocationValue(value interface{}, path, publishType string, normalizations *[]NormalizationEvent) (interface{}, bool) {
	obj, ok := value.(map[string]interface{})
	if !ok || obj == nil {
		return value, false
	}
	if data, _ := obj["data"].(map[string]interface{}); isDynamicQueryObject(data) {
		changed := false
		if _, exists := obj["isScp"]; !exists {
			obj["isScp"] = inferDouyinLocationIsScp(data, publishType)
			changed = true
		}
		if changed {
			appendNormalization(normalizations, NormalizationEvent{
				Field:        "location",
				Path:         path,
				Action:       "complete_frontend_shape",
				Message:      `Completed Douyin frontend location shape with "isScp".`,
				QueryCommand: dynamicObjectQueryCommand("location"),
			})
		}
		return obj, changed
	}
	if isDynamicQueryObject(obj) {
		normalized := map[string]interface{}{
			"isScp": inferDouyinLocationIsScp(obj, publishType),
			"data":  obj,
		}
		appendNormalization(normalizations, NormalizationEvent{
			Field:        "location",
			Path:         path,
			Action:       "wrap_frontend_shape",
			Message:      `Wrapped location query object into Douyin frontend location shape.`,
			QueryCommand: dynamicObjectQueryCommand("location"),
		})
		return normalized, true
	}
	return value, false
}

func inferDouyinLocationIsScp(location map[string]interface{}, publishType string) bool {
	if NormalizePublishType(publishType) != "imageText" {
		return false
	}
	if n, ok := numericValue(location["cpsProductCount"]); ok {
		return n > 0
	}
	if raw, _ := location["raw"].(map[string]interface{}); raw != nil {
		if n, ok := numericValue(raw["cpsProductCount"]); ok {
			return n > 0
		}
	}
	return false
}

func numericValue(value interface{}) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	default:
		return 0, false
	}
}

func normalizePlatformDataValue(value interface{}, path, field string, normalizations *[]NormalizationEvent) (interface{}, bool) {
	if items, ok := value.([]interface{}); ok {
		changed := false
		for i, item := range items {
			normalized, itemChanged := normalizePlatformDataValue(item, fmt.Sprintf("%s[%d]", path, i), field, normalizations)
			if itemChanged {
				items[i] = normalized
				changed = true
			}
		}
		return items, changed
	}
	obj, ok := value.(map[string]interface{})
	if !ok || obj == nil {
		return value, false
	}
	if isPlatformDataObject(obj) {
		return value, false
	}
	if data, _ := obj["data"].(map[string]interface{}); data != nil {
		if normalized, ok := dynamicQueryObjectToPlatformData(data); ok {
			appendNormalization(normalizations, NormalizationEvent{
				Field:        field,
				Path:         path,
				Action:       "unwrap_data_to_frontend_shape",
				Message:      `Unwrapped query data into frontend id/text/raw shape.`,
				QueryCommand: dynamicObjectQueryCommand(field),
			})
			return normalized, true
		}
	}
	if normalized, ok := dynamicQueryObjectToPlatformData(obj); ok {
		appendNormalization(normalizations, NormalizationEvent{
			Field:        field,
			Path:         path,
			Action:       "map_frontend_shape",
			Message:      `Mapped query object into frontend id/text/raw shape.`,
			QueryCommand: dynamicObjectQueryCommand(field),
		})
		return normalized, true
	}
	return value, false
}

func isPlatformDataObject(obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}
	raw, _ := obj["raw"].(map[string]interface{})
	return stringField(obj, "id") != "" && stringField(obj, "text") != "" && raw != nil
}

func dynamicQueryObjectToPlatformData(obj map[string]interface{}) (map[string]interface{}, bool) {
	if !isDynamicQueryObject(obj) {
		return nil, false
	}
	id := firstNonEmptyString(stringField(obj, "id"), stringField(obj, "yixiaoerId"), stringField(obj, "value"))
	text := firstNonEmptyString(stringField(obj, "text"), stringField(obj, "yixiaoerName"), stringField(obj, "name"), stringField(obj, "label"), stringField(obj, "title"))
	if id == "" || text == "" {
		return nil, false
	}
	normalized := map[string]interface{}{
		"id":   id,
		"text": text,
		"raw":  obj,
	}
	if child, _ := obj["child"].([]interface{}); len(child) > 0 {
		normalized["children"] = normalizePlatformDataChildren(child)
	}
	if children, _ := obj["children"].([]interface{}); len(children) > 0 {
		normalized["children"] = normalizePlatformDataChildren(children)
	}
	return normalized, true
}

func normalizePlatformDataChildren(items []interface{}) []interface{} {
	var result []interface{}
	for _, item := range items {
		obj, _ := item.(map[string]interface{})
		if normalized, ok := dynamicQueryObjectToPlatformData(obj); ok {
			result = append(result, normalized)
		}
	}
	return result
}

func normalizeDynamicObjectValue(value interface{}, path, field string, normalizations *[]NormalizationEvent) (interface{}, bool) {
	if items, ok := value.([]interface{}); ok {
		changed := false
		for i, item := range items {
			normalized, itemChanged := normalizeDynamicObjectValue(item, fmt.Sprintf("%s[%d]", path, i), field, normalizations)
			if itemChanged {
				items[i] = normalized
				changed = true
			}
		}
		return items, changed
	}
	obj, ok := value.(map[string]interface{})
	if !ok || obj == nil {
		return value, false
	}
	if data, _ := obj["data"].(map[string]interface{}); isDynamicQueryObject(data) {
		appendNormalization(normalizations, NormalizationEvent{
			Field:        field,
			Path:         path,
			Action:       "unwrap_data",
			Message:      `Unwrapped data envelope into the frontend form shape.`,
			QueryCommand: dynamicObjectQueryCommand(field),
		})
		normalizeDynamicIdentityAliases(data)
		return data, true
	}
	if normalizeDynamicIdentityAliases(obj) {
		appendNormalization(normalizations, NormalizationEvent{
			Field:        field,
			Path:         path,
			Action:       "map_identity_aliases",
			Message:      `Mapped id/text aliases into query object identity fields.`,
			QueryCommand: dynamicObjectQueryCommand(field),
		})
		return obj, true
	}
	return value, false
}

func isDynamicQueryObject(obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}
	raw, _ := obj["raw"].(map[string]interface{})
	if raw == nil {
		return false
	}
	return obj["yixiaoerId"] != nil || obj["yixiaoerName"] != nil || obj["id"] != nil || obj["name"] != nil || obj["text"] != nil
}

func normalizeDynamicIdentityAliases(obj map[string]interface{}) bool {
	if !isDynamicQueryObject(obj) {
		return false
	}
	changed := false
	if empty(obj["yixiaoerId"]) {
		if id := firstNonEmptyString(stringField(obj, "id"), stringField(obj, "value")); id != "" {
			obj["yixiaoerId"] = id
			changed = true
		}
	}
	if empty(obj["yixiaoerName"]) {
		if name := firstNonEmptyString(stringField(obj, "text"), stringField(obj, "name"), stringField(obj, "label"), stringField(obj, "title")); name != "" {
			obj["yixiaoerName"] = name
			changed = true
		}
	}
	return changed
}

func dynamicObjectQueryCommand(field string) string {
	switch field {
	case "location":
		return "yxer query locations <account_id> [--query 关键词] --json"
	case "music":
		return "yxer query music <account_id> [--query 关键词] --json"
	case "category":
		return "yxer query categories <account_id> [--type video|article] --json"
	case "collection", "collections", "sub_collection":
		return "yxer query collections <account_id> [--type video|article] --json"
	case "challenge", "challenges":
		return "yxer query challenges <account_id> [--query 关键词] [--type video] --json"
	case "mini_app", "miniapp", "miniapps":
		return "yxer query miniapps <account_id> [--query 关键词] --json"
	case "sync_apps":
		return "yxer query syncapps <account_id> --json"
	case "game":
		return "yxer query games <account_id> [--query 关键词] --json"
	case "hot_event":
		return "yxer query hot-events <account_id> [--type video|article] --json"
	case "group", "groups":
		return "yxer query groups <account_id> --json"
	case "activity":
		return "yxer query activities <account_id> [--type video|article] [--query 关键词] --json"
	case "shopping_cart", "shoppingCart":
		return "yxer query goods <account_id> [--query 关键词] --json"
	default:
		return ""
	}
}

func normalizeFlatShoppingCart(cpf map[string]interface{}, formPath string, normalizations *[]NormalizationEvent) {
	if cpf == nil {
		return
	}
	value, exists := cpf["shopping_cart"]
	if !exists {
		return
	}
	normalized, changed := normalizeDynamicObjectValue(value, formPath+".shopping_cart", "shopping_cart", normalizations)
	if changed {
		cpf["shopping_cart"] = normalized
	}
}

func normalizeDouyinShoppingCart(cpf map[string]interface{}, formPath string, normalizations *[]NormalizationEvent) {
	if cpf == nil {
		return
	}
	if value, ok := cpf["shoppingCart"]; ok {
		if _, exists := cpf["shopping_cart"]; !exists {
			cpf["shopping_cart"] = value
			appendNormalization(normalizations, NormalizationEvent{
				Field:        "shopping_cart",
				Path:         formPath + ".shoppingCart",
				Action:       "rename_field",
				Message:      `Renamed legacy "shoppingCart" to "shopping_cart".`,
				QueryCommand: "yxer query goods <account_id> [--query 关键词]",
			})
		}
		delete(cpf, "shoppingCart")
	}
	items, ok := cpf["shopping_cart"].([]interface{})
	if !ok {
		return
	}
	for itemIndex, item := range items {
		cart, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		itemPath := fmt.Sprintf("%s.shopping_cart[%d]", formPath, itemIndex)
		if cart["data"] == nil {
			if isDynamicQueryObject(cart) {
				data := cloneObjectExcluding(cart, "sale_title", "images", "data")
				saleTitle := firstNonEmptyString(stringField(cart, "sale_title"), stringField(data, "yixiaoerName"), "点击购买")
				images := interfaceSliceField(cart, "images")
				if len(images) == 0 {
					images = extractShoppingCartImagesFromGoods(data)
				}
				clearObject(cart)
				cart["sale_title"] = truncateRunes(saleTitle, 10)
				if len(images) > 0 {
					cart["images"] = images
					appendNormalization(normalizations, NormalizationEvent{
						Field:        "shopping_cart",
						Path:         itemPath + ".images",
						Action:       "derive_images",
						Message:      "Derived top-level shopping_cart images from goods raw data.",
						QueryCommand: "yxer query goods <account_id> [--query 关键词]",
					})
				}
				cart["data"] = data
				appendNormalization(normalizations, NormalizationEvent{
					Field:        "shopping_cart",
					Path:         itemPath,
					Action:       "wrap_data",
					Message:      `Wrapped complete goods query object into "data".`,
					QueryCommand: "yxer query goods <account_id> [--query 关键词]",
				})
			} else {
				data := map[string]interface{}{}
				for _, key := range []string{"yixiaoerId", "yixiaoerName", "raw"} {
					if value, exists := cart[key]; exists {
						data[key] = value
						delete(cart, key)
					}
				}
				if len(data) > 0 {
					cart["data"] = data
					appendNormalization(normalizations, NormalizationEvent{
						Field:        "shopping_cart",
						Path:         itemPath,
						Action:       "wrap_data",
						Message:      `Moved flat yixiaoerId/yixiaoerName/raw into "data".`,
						QueryCommand: "yxer query goods <account_id> [--query 关键词]",
					})
				}
			}
		}
		if empty(cart["sale_title"]) {
			if data, _ := cart["data"].(map[string]interface{}); data != nil {
				cart["sale_title"] = truncateRunes(firstNonEmptyString(stringField(data, "yixiaoerName"), "点击购买"), 10)
				appendNormalization(normalizations, NormalizationEvent{
					Field:        "shopping_cart",
					Path:         itemPath + ".sale_title",
					Action:       "derive_sale_title",
					Message:      "Derived shopping_cart sale_title from goods data.",
					QueryCommand: "yxer query goods <account_id> [--query 关键词]",
				})
			}
		}
		if cart["images"] == nil {
			if data, _ := cart["data"].(map[string]interface{}); data != nil {
				if images := extractShoppingCartImagesFromGoods(data); len(images) > 0 {
					cart["images"] = images
					appendNormalization(normalizations, NormalizationEvent{
						Field:        "shopping_cart",
						Path:         itemPath + ".images",
						Action:       "derive_images",
						Message:      "Derived top-level shopping_cart images from goods raw data.",
						QueryCommand: "yxer query goods <account_id> [--query 关键词]",
					})
				}
			}
		}
	}
}

func normalizeDouyinGroupShopping(cpf map[string]interface{}, formPath string, normalizations *[]NormalizationEvent) {
	if cpf == nil {
		return
	}
	if value, ok := cpf["groupShopping"]; ok {
		if _, exists := cpf["group_shopping"]; !exists {
			cpf["group_shopping"] = value
			appendNormalization(normalizations, NormalizationEvent{
				Field:        "group_shopping",
				Path:         formPath + ".groupShopping",
				Action:       "rename_field",
				Message:      `Renamed legacy "groupShopping" to frontend "group_shopping".`,
				QueryCommand: "yxer query goods <account_id> [--query 关键词]",
			})
		}
		delete(cpf, "groupShopping")
	}
	group, _ := cpf["group_shopping"].(map[string]interface{})
	if group == nil {
		return
	}
	groupPath := formPath + ".group_shopping"
	if data, _ := group["data"].(map[string]interface{}); data != nil {
		changed := false
		if empty(group["sale_title"]) {
			group["sale_title"] = truncateRunes(firstNonEmptyString(stringField(data, "yixiaoerName"), "点击购买"), 10)
			changed = true
		}
		if _, exists := group["brand_switch_value"]; !exists {
			group["brand_switch_value"] = float64(0)
			changed = true
		}
		if changed {
			appendNormalization(normalizations, NormalizationEvent{
				Field:        "group_shopping",
				Path:         groupPath,
				Action:       "complete_frontend_shape",
				Message:      "Completed Douyin frontend group_shopping fields from goods data.",
				QueryCommand: "yxer query goods <account_id> [--query 关键词]",
			})
		}
		return
	}
	if isDynamicQueryObject(group) {
		data := cloneObjectExcluding(group, "sale_title", "brand_switch_value", "data")
		saleTitle := truncateRunes(firstNonEmptyString(stringField(group, "sale_title"), stringField(data, "yixiaoerName"), "点击购买"), 10)
		brandSwitch := group["brand_switch_value"]
		if brandSwitch == nil {
			brandSwitch = float64(0)
		}
		clearObject(group)
		group["brand_switch_value"] = brandSwitch
		group["sale_title"] = saleTitle
		group["data"] = data
		appendNormalization(normalizations, NormalizationEvent{
			Field:        "group_shopping",
			Path:         groupPath,
			Action:       "wrap_frontend_shape",
			Message:      `Wrapped goods query object into Douyin frontend group_shopping shape.`,
			QueryCommand: "yxer query goods <account_id> [--query 关键词]",
		})
	}
}

func appendNormalization(normalizations *[]NormalizationEvent, event NormalizationEvent) {
	if normalizations == nil {
		return
	}
	*normalizations = append(*normalizations, event)
}

func cloneObjectExcluding(obj map[string]interface{}, excluded ...string) map[string]interface{} {
	if obj == nil {
		return nil
	}
	exclude := map[string]bool{}
	for _, key := range excluded {
		exclude[key] = true
	}
	cloned := map[string]interface{}{}
	for key, value := range obj {
		if !exclude[key] {
			cloned[key] = value
		}
	}
	return cloned
}

func clearObject(obj map[string]interface{}) {
	for key := range obj {
		delete(obj, key)
	}
}

func interfaceSliceField(obj map[string]interface{}, key string) []interface{} {
	items, _ := obj[key].([]interface{})
	return items
}

func truncateRunes(value string, max int) string {
	value = strings.TrimSpace(value)
	if max <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}

func extractShoppingCartImagesFromGoods(goods map[string]interface{}) []interface{} {
	if goods == nil {
		return nil
	}
	if raw, _ := goods["raw"].(map[string]interface{}); raw != nil {
		if images := extractShoppingCartImages(raw); len(images) > 0 {
			return images
		}
	}
	if images := extractShoppingCartImages(goods); len(images) > 0 {
		return images
	}
	if imageURL := strings.TrimSpace(stringField(goods, "yixiaoerImageUrl")); imageURL != "" {
		return []interface{}{imageURL}
	}
	return nil
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
		case "视频号", "微信视频号", "shipinhao":
			requireShipinhaoCoverSize(cover, formPath, errors)
		}
	}
}

func requirePlatformImageTextConstraints(platforms []string, images []interface{}, cover map[string]interface{}, formPath string, errors *[]string) {
	for _, platform := range platforms {
		switch strings.TrimSpace(platform) {
		case "视频号", "微信视频号", "shipinhao":
			requireShipinhaoImageSizes(images, formPath, errors)
			requireShipinhaoCoverSize(cover, formPath, errors)
		}
	}
}

func requireShipinhaoImageSizes(images []interface{}, formPath string, errors *[]string) {
	for i, item := range images {
		imageObj, _ := item.(map[string]interface{})
		if imageObj == nil {
			continue
		}
		size, ok := integerField(imageObj, "size")
		if !ok {
			continue
		}
		if size > shipinhaoImageMaxBytes {
			*errors = append(*errors, fmt.Sprintf("%s.images[%d].size: 视频号图片不能超过 512KB，当前为 %d bytes", formPath, i, size))
		}
	}
}

func requireShipinhaoCoverSize(cover map[string]interface{}, formPath string, errors *[]string) {
	if cover == nil {
		return
	}
	size, ok := integerField(cover, "size")
	if !ok {
		return
	}
	if size > shipinhaoImageMaxBytes {
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

func requireUploadedVideoResource(resource map[string]interface{}, pathLabel string, errors *[]string) {
	requireUploadedResource(resource, pathLabel, errors)
	if resource == nil {
		return
	}
	duration, ok := numericField(resource, "duration")
	if !ok || duration <= 0 {
		*errors = append(*errors, fmt.Sprintf("%s: missing uploaded video field %q; run \"yxer upload <file_path_or_url>\" with --auto-meta enabled and keep the returned duration", pathLabel, "duration"))
	}
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
		*errors = append(*errors, pathLabel+dynamicRawErrorSuffix(pathLabel))
	}
}

func dynamicRawErrorSuffix(pathLabel string) string {
	field, command := dynamicObjectRepairCommand(pathLabel)
	if command == "" {
		return `: dynamic platform object must include complete "raw" data from a yxer query command`
	}
	return fmt.Sprintf(`: dynamic platform object must include complete "raw" data from a yxer query command; %s must be copied from %s`, field, command)
}

func dynamicObjectRepairCommand(pathLabel string) (string, string) {
	switch {
	case strings.Contains(pathLabel, "shopping_cart") || strings.Contains(pathLabel, "group_shopping") || strings.Contains(pathLabel, "shoppingCart") || strings.Contains(pathLabel, "groupShopping"):
		return "shopping cart goods", "yxer query goods <account_id> [--query 关键词] --json"
	case strings.Contains(pathLabel, "location"):
		return "location", "yxer query locations <account_id> [--query 关键词] --json"
	case strings.Contains(pathLabel, "music"):
		return "music", "yxer query music <account_id> [--query 关键词] --json"
	case strings.Contains(pathLabel, "challenge") || strings.Contains(pathLabel, "topics"):
		return "topic/challenge", "yxer query challenges <account_id> [--query 关键词] [--type video] --json"
	default:
		return "", ""
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
	return float64(value), ""
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

func ResolveStandardPayloadResourceMetadata(payload map[string]interface{}) error {
	var errors []string
	resolveStandardPayloadResourceMetadata(payload, &errors)
	if len(errors) == 0 {
		return nil
	}
	return yxerrors.Usage("Publish resource metadata extraction failed", errors).
		WithHint("请在资源对象中保留已上传的 key，并使用 source/path/localPath/filePath 指向本地文件以自动提取元数据；或直接使用 yxer upload 返回的完整对象。")
}

func resolveStandardPayloadResourceMetadata(payload map[string]interface{}, errors *[]string) {
	if payload == nil {
		return
	}
	enrichResourceObjectMetadata(objectField(payload, "cover"), "cover", errors)
	publishArgs := ExtractPublishArgs(payload)
	if publishArgs == nil {
		return
	}
	enrichResourceContainerMetadata(publishArgs, "publishArgs", errors)
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	for i, item := range accountForms {
		form, _ := item.(map[string]interface{})
		if form == nil {
			continue
		}
		formPath := fmt.Sprintf("publishArgs.accountForms[%d]", i)
		enrichResourceContainerMetadata(form, formPath, errors)
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf != nil {
			enrichResourceContainerMetadata(cpf, formPath+".contentPublishForm", errors)
		}
	}
}

func enrichResourceContainerMetadata(container map[string]interface{}, path string, errors *[]string) {
	if container == nil {
		return
	}
	enrichResourceObjectMetadata(objectField(container, "video"), path+".video", errors)
	enrichResourceObjectMetadata(objectField(container, "cover"), path+".cover", errors)
	enrichResourceObjectMetadata(objectField(container, "horizontalCover"), path+".horizontalCover", errors)
	if items, _ := container["images"].([]interface{}); len(items) > 0 {
		for i, item := range items {
			resource, _ := item.(map[string]interface{})
			if resource == nil {
				continue
			}
			enrichResourceObjectMetadata(resource, fmt.Sprintf("%s.images[%d]", path, i), errors)
		}
	}
}

func enrichResourceObjectMetadata(resource map[string]interface{}, path string, errors *[]string) {
	if resource == nil {
		return
	}
	source := consumeResourceMetadataSource(resource)
	if source == "" {
		return
	}
	meta, _, err := api.InspectUpload(source, true)
	if err != nil {
		appendResourceMetadataError(errors, fmt.Sprintf("%s: failed to inspect media metadata from %q: %v", path, source, err))
		return
	}
	if empty(resource["size"]) && meta.Size > 0 {
		resource["size"] = float64(meta.Size)
	}
	if numericFieldMissingOrZero(resource, "width") && meta.Width > 0 {
		resource["width"] = float64(meta.Width)
	}
	if numericFieldMissingOrZero(resource, "height") && meta.Height > 0 {
		resource["height"] = float64(meta.Height)
	}
	if numericFieldMissingOrZero(resource, "duration") && meta.Duration > 0 {
		resource["duration"] = meta.Duration
	}
	if empty(resource["format"]) && meta.Format != "" {
		resource["format"] = meta.Format
	}
}

func consumeResourceMetadataSource(resource map[string]interface{}) string {
	if resource == nil {
		return ""
	}
	keys := []string{"source", "path", "localPath", "filePath"}
	source := ""
	for _, key := range keys {
		if source == "" {
			source = strings.TrimSpace(stringField(resource, key))
		}
		delete(resource, key)
	}
	return source
}

func appendResourceMetadataError(errors *[]string, message string) {
	if errors == nil || strings.TrimSpace(message) == "" {
		return
	}
	*errors = append(*errors, message)
}

func numericFieldMissingOrZero(resource map[string]interface{}, key string) bool {
	value, exists := resource[key]
	if !exists || value == nil {
		return true
	}
	switch typed := value.(type) {
	case float64:
		return typed == 0
	case int:
		return typed == 0
	case int64:
		return typed == 0
	default:
		return empty(value)
	}
}

func numericField(resource map[string]interface{}, key string) (float64, bool) {
	if resource == nil {
		return 0, false
	}
	value, exists := resource[key]
	if !exists || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	default:
		return 0, false
	}
}
