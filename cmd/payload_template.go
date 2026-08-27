package cmd

import (
	"sort"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
)

func buildPayloadTemplate(doc schema.Document) map[string]interface{} {
	contentProperties := doc.Properties
	if doc.Type == "article" && !isWeixinAccountArticleDoc(doc) {
		contentProperties = clonePropertyViewsWithoutKeys(doc.Properties, "content")
	}
	contentProperties = clonePropertyViewsWithoutKeys(contentProperties, accountLevelResourceKeys(doc.Type)...)
	contentPublishForm := buildTemplateObject(contentProperties)
	if isDuoduoshipinPlatform(doc.Platform) {
		if _, ok := doc.Properties["shopping_cart"]; ok {
			// The product ID is a user business input. Keep the platform source
			// fixed in the form skeleton instead of asking agents to query goods
			// and accidentally copy yixiaoerId.
			contentPublishForm["shopping_cart"] = map[string]interface{}{"source": "pdd"}
		}
	}
	accountForm := map[string]interface{}{
		"platformAccountId":  "<platformAccountId>",
		"contentPublishForm": contentPublishForm,
	}
	addRequiredAccountResources(accountForm, doc)
	publishArgs := map[string]interface{}{
		"accountForms": []interface{}{
			accountForm,
		},
	}
	if isWeixinAccountArticleDoc(doc) {
		publishArgs["platformForms"] = map[string]interface{}{
			"微信公众号": buildTemplateObject(doc.Properties),
		}
	} else if doc.Type == "article" {
		publishArgs["content"] = "<content>"
	}
	template := map[string]interface{}{
		"action":         "publish",
		"publishType":    doc.Type,
		"platforms":      []interface{}{platformutil.ChineseName(doc.Platform)},
		"publishChannel": "cloud",
		"publishArgs":    publishArgs,
	}
	if doc.Type == "article" {
		template["desc"] = "<desc>"
	}
	return template
}

// buildMinimalPayloadTemplate 构建最小可用模板（仅包含必填字段）
func buildMinimalPayloadTemplate(doc schema.Document) map[string]interface{} {
	// 只提取必填字段
	requiredFields := make(map[string]interface{})
	for key, prop := range clonePropertyViewsWithoutKeys(doc.Properties, append(articleContentTemplateExclusion(doc.Type), accountLevelResourceKeys(doc.Type)...)...) {
		if prop.Required {
			value, ok := buildTemplateValue(key, prop)
			if ok {
				requiredFields[key] = value
			}
		}
	}
	if isWeixinAccountImageTextDoc(doc) {
		for _, key := range []string{"notifySubscribers", "sex", "needOpenComment", "statement", "disableRecommend", "pubType"} {
			if _, exists := requiredFields[key]; exists {
				continue
			}
			if prop, ok := doc.Properties[key]; ok {
				if value, ok := buildTemplateValue(key, prop); ok {
					requiredFields[key] = value
				}
			}
		}
	}

	accountForm := map[string]interface{}{
		"platformAccountId": "<从 accounts list 获取>",
	}
	if isWeixinAccountArticleDoc(doc) {
		accountForm["platformName"] = "微信公众号"
	} else {
		accountForm["contentPublishForm"] = requiredFields
	}
	addRequiredAccountResources(accountForm, doc)

	template := map[string]interface{}{
		"action":      "publish",
		"publishType": doc.Type,
		"platforms":   []interface{}{platformutil.ChineseName(doc.Platform)},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				accountForm,
			},
		},
	}
	if isWeixinAccountArticleDoc(doc) {
		template["publishArgs"].(map[string]interface{})["platformForms"] = map[string]interface{}{
			"微信公众号": requiredFields,
		}
		template["desc"] = "<任务描述>"
	} else if doc.Type == "article" {
		template["publishArgs"].(map[string]interface{})["content"] = "<从正文生成>"
		template["desc"] = "<任务描述>"
	}

	return template
}

func requiresPublishCoverResource(doc schema.Document) bool {
	if doc.Type == "video" {
		return true
	}
	return doc.Type == "imageText" && !platformutil.ImageTextUsesFirstImageAsCover(doc.Platform)
}

func articleContentTemplateExclusion(publishType string) []string {
	if publishType != "article" {
		return nil
	}
	return []string{"content"}
}

func clonePropertyViewsWithoutKeys(src map[string]schema.PropertyView, keys ...string) map[string]schema.PropertyView {
	if len(src) == 0 {
		return nil
	}
	omit := map[string]bool{}
	for _, key := range keys {
		omit[key] = true
	}
	dst := make(map[string]schema.PropertyView, len(src))
	for key, value := range src {
		if omit[key] {
			continue
		}
		dst[key] = value
	}
	return dst
}

func buildTemplateObject(properties map[string]schema.PropertyView) map[string]interface{} {
	if len(properties) == 0 {
		return map[string]interface{}{}
	}
	keys := make([]string, 0, len(properties))
	for key := range properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		view := properties[key]
		value, ok := buildTemplateValue(key, view)
		if !ok {
			continue
		}
		result[key] = value
	}
	return result
}

func buildTemplateValue(name string, view schema.PropertyView) (interface{}, bool) {
	if view.Const != nil {
		return view.Const, true
	}
	if view.Default != nil {
		return view.Default, true
	}
	if len(view.Enum) > 0 {
		return view.Enum[0], true
	}
	if !view.Required {
		return nil, false
	}
	switch view.Type {
	case "object":
		return buildTemplateObject(view.Properties), true
	case "array":
		if view.Items == nil {
			return []interface{}{}, true
		}
		item, ok := buildTemplateValue(name, *view.Items)
		if !ok {
			return []interface{}{}, true
		}
		return []interface{}{item}, true
	case "boolean":
		return false, true
	case "number", "integer":
		return 0, true
	case "string":
		return "<" + name + ">", true
	default:
		return "<" + name + ">", true
	}
}

func accountLevelResourceKeys(publishType string) []string {
	switch publishType {
	case "video":
		return []string{"video", "cover", "coverKey"}
	case "imageText":
		return []string{"images", "cover", "coverKey"}
	default:
		return nil
	}
}

func isWeixinAccountImageTextDoc(doc schema.Document) bool {
	return doc.Type == "imageText" && platformutil.CanonicalKey(doc.Platform) == "weixin.account"
}

func addRequiredAccountResources(accountForm map[string]interface{}, doc schema.Document) {
	switch doc.Type {
	case "video":
		accountForm["video"] = videoResourcePlaceholder()
		accountForm["cover"] = imageResourcePlaceholder()
		accountForm["coverKey"] = "<与 cover.key 一致>"
	case "imageText":
		accountForm["images"] = []interface{}{imageResourcePlaceholder()}
		if requiresPublishCoverResource(doc) {
			accountForm["cover"] = imageResourcePlaceholder()
			accountForm["coverKey"] = "<与 cover.key 一致>"
		}
	}
}

func imageResourcePlaceholder() map[string]interface{} {
	return map[string]interface{}{
		"key":    "<从 yxer upload 获取>",
		"size":   0,
		"width":  0,
		"height": 0,
		"format": "<png|jpg|jpeg|webp>",
	}
}

func videoResourcePlaceholder() map[string]interface{} {
	resource := imageResourcePlaceholder()
	resource["format"] = "<mp4|mov>"
	resource["duration"] = 0
	return resource
}
