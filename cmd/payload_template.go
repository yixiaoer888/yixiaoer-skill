package cmd

import (
	"sort"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/core/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/schema"
)

func buildPayloadTemplate(doc schema.Document) map[string]interface{} {
	template := map[string]interface{}{
		"action":         "publish",
		"publishType":    doc.Type,
		"platforms":      []interface{}{platformutil.ChineseName(doc.Platform)},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId":  "<platformAccountId>",
					"contentPublishForm": buildTemplateObject(doc.Properties),
				},
			},
		},
	}
	return template
}

// buildMinimalPayloadTemplate 构建最小可用模板（仅包含必填字段）
func buildMinimalPayloadTemplate(doc schema.Document) map[string]interface{} {
	// 只提取必填字段
	requiredFields := make(map[string]interface{})
	for key, prop := range doc.Properties {
		if prop.Required {
			value, ok := buildTemplateValue(key, prop)
			if ok {
				requiredFields[key] = value
			}
		}
	}

	template := map[string]interface{}{
		"action":      "publish",
		"publishType": doc.Type,
		"platforms":   []interface{}{platformutil.ChineseName(doc.Platform)},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId":  "<从 accounts list 获取>",
					"contentPublishForm": requiredFields,
				},
			},
		},
	}

	return template
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
