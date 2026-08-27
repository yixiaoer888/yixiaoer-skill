package cmd

import (
	"strings"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
)

type dynamicFieldExample struct {
	Field        string      `json:"field"`
	Path         string      `json:"path"`
	Source       string      `json:"source"`
	QueryCommand string      `json:"queryCommand,omitempty"`
	Note         string      `json:"note"`
	Value        interface{} `json:"value"`
}

func buildDynamicFieldExamples(doc schema.Document) map[string]dynamicFieldExample {
	examples := map[string]dynamicFieldExample{}
	addDynamicObjectExample(examples, doc, "category", "yxer query categories <account_id> --type <video|article> --json", "分类对象必须完整来自查询结果；有层级分类时保留完整 child 路径。")
	addDynamicObjectExample(examples, doc, "location", "yxer query locations <account_id> [--query 关键词] --json", "位置对象必须完整来自查询结果，并按 schema 显示的前端表单结构回填。")
	addDynamicObjectExample(examples, doc, "music", "yxer query music <account_id> [--query 关键词] --json", "音乐对象必须完整来自查询结果，保留播放地址、时长等查询返回字段以及 raw。")
	addDynamicObjectExample(examples, doc, "collection", "yxer query collections <account_id> --type <video|article> --json", "合集对象必须完整来自查询结果，并保留 raw。")
	addDynamicObjectExample(examples, doc, "sub_collection", "yxer query collections <account_id> --type <video|article> --json", "子合集对象必须完整来自查询结果，并保留 raw。")
	addDramaExample(examples, doc)
	addDynamicObjectExample(examples, doc, "challenge", "yxer query challenges <account_id> [--query 关键词] --type video --json", "话题对象必须完整来自查询结果，并保留 raw。")
	addDynamicObjectExample(examples, doc, "hot_event", "yxer query hot-events <account_id> --type <video|article> --json", "热点对象必须完整来自查询结果，并保留 raw。")
	addDynamicObjectExample(examples, doc, "mini_app", "yxer query miniapps <account_id> [--query 关键词] --json", "小程序对象必须完整来自查询结果，并保留 raw。")
	addDynamicObjectExample(examples, doc, "game", "yxer query games <account_id> [--query 关键词] --json", "游戏对象必须完整来自查询结果，并保留 raw。")
	addDynamicArrayObjectExample(examples, doc, "sync_apps", "yxer query syncapps <account_id> --json", "同步发布应用数组中的每个对象都必须来自查询结果，并保留 raw。")
	addDynamicObjectExample(examples, doc, "activity", "yxer query activities <account_id> --type <video|article> [--category-id ID] [--query 关键词] --json", "活动对象必须完整来自查询结果，并保留 raw。")
	addTagsExample(examples, doc)
	addShoppingCartExample(examples, doc)
	if len(examples) == 0 {
		return nil
	}
	return examples
}

func addDramaExample(examples map[string]dynamicFieldExample, doc schema.Document) {
	if _, ok := doc.Properties["drama"]; !ok {
		return
	}
	examples["drama"] = dynamicFieldExample{
		Field:        "drama",
		Path:         "publishArgs.accountForms[].contentPublishForm.drama",
		Source:       "query",
		QueryCommand: "yxer query drama-tasks <account_id> [--query 关键词] --json",
		Note:         "视频号剧集对象必须完整来自 yxer query drama-tasks 查询结果，并且只保留 yixiaoerId、yixiaoerImageUrl、yixiaoerName；剧集对象不使用 raw。",
		Value: map[string]interface{}{
			"yixiaoerId":       "<from query>",
			"yixiaoerImageUrl": "<from query>",
			"yixiaoerName":     "<from query>",
		},
	}
}

func addDynamicArrayObjectExample(examples map[string]dynamicFieldExample, doc schema.Document, field, command, note string) {
	if _, ok := doc.Properties[field]; !ok {
		return
	}
	item := map[string]interface{}{
		"yixiaoerId":   "<from query>",
		"yixiaoerName": "<from query>",
		"raw":          map[string]interface{}{"...": "copy complete raw object from query result"},
	}
	examples[field] = dynamicFieldExample{
		Field:        field,
		Path:         "publishArgs.accountForms[].contentPublishForm." + field,
		Source:       "query",
		QueryCommand: command,
		Note:         note,
		Value:        []interface{}{item},
	}
}

func addDynamicObjectExample(examples map[string]dynamicFieldExample, doc schema.Document, field, command, note string) {
	if _, ok := doc.Properties[field]; !ok {
		return
	}
	view := doc.Properties[field]
	queryObject := map[string]interface{}{
		"yixiaoerId":   "<from query>",
		"yixiaoerName": "<from query>",
		"raw": map[string]interface{}{
			"...": "copy complete raw object from query result",
		},
	}
	value := queryObject
	if _, ok := view.Properties["isScp"]; ok {
		value = map[string]interface{}{
			"isScp": false,
			"data":  queryObject,
		}
	} else if _, ok := view.Properties["id"]; ok {
		value = map[string]interface{}{
			"id":   "<from query yixiaoerId>",
			"text": "<from query yixiaoerName>",
			"raw":  queryObject,
		}
	}
	if field == "music" {
		value["duration"] = 0
		value["playUrl"] = "<from query>"
	}
	examples[field] = dynamicFieldExample{
		Field:        field,
		Path:         "publishArgs.accountForms[].contentPublishForm." + field,
		Source:       "query",
		QueryCommand: command,
		Note:         note,
		Value:        value,
	}
}

func addTagsExample(examples map[string]dynamicFieldExample, doc schema.Document) {
	if _, ok := doc.Properties["tags"]; !ok {
		return
	}
	examples["tags"] = dynamicFieldExample{
		Field:  "tags",
		Path:   "publishArgs.accountForms[].contentPublishForm.tags",
		Source: "schema",
		Note:   "话题标签使用字符串数组；CLI 不会改变 tags 字段结构。description 中的 #话题 会按描述规则归一化。",
		Value:  []interface{}{"话题1", "话题2"},
	}
}

func addShoppingCartExample(examples map[string]dynamicFieldExample, doc schema.Document) {
	field, view, ok := shoppingCartField(doc)
	if !ok {
		return
	}
	if isDuoduoshipinPlatform(doc.Platform) {
		// Duoduoshipin receives a user-entered business goods_id rather than a
		// query goods object. Keep it out of dynamicFieldExamples so
		// publish form choose cannot silently use yixiaoerId.
		return
	}
	item := map[string]interface{}{
		"yixiaoerId":   "<from query>",
		"yixiaoerName": "<from query>",
		"raw": map[string]interface{}{
			"...": "copy complete goods raw object from query result",
		},
	}
	value := []interface{}{item}
	note := "购物车商品必须来自 yxer query goods 返回的完整对象。"
	if shoppingCartUsesNestedData(view) {
		value = []interface{}{
			map[string]interface{}{
				"sale_title": "点击购买",
				"images":     []interface{}{"<from query images[0]>"},
				"data":       item,
			},
		}
		note = "购物车使用顶层 sale_title/images + 内层 data；直接把商品字段扁平放在根节点会在 dry-run 中被归一化，但新 payload 应直接使用该结构。"
	}
	examples[field] = dynamicFieldExample{
		Field:        field,
		Path:         "publishArgs.accountForms[].contentPublishForm." + field,
		Source:       "query",
		QueryCommand: "yxer query goods <account_id> [--query 关键词] --json",
		Note:         note,
		Value:        value,
	}
}

func isDuoduoshipinPlatform(platform string) bool {
	return platformutil.CanonicalKey(platform) == "duoduoshipin"
}

func shoppingCartField(doc schema.Document) (string, schema.PropertyView, bool) {
	for _, key := range []string{"shopping_cart", "group_shopping", "shoppingCart", "groupShopping"} {
		if view, ok := doc.Properties[key]; ok {
			return key, view, true
		}
	}
	return "", schema.PropertyView{}, false
}

func shoppingCartUsesNestedData(view schema.PropertyView) bool {
	if view.Items != nil {
		_, ok := view.Items.Properties["data"]
		return ok
	}
	_, ok := view.Properties["data"]
	if ok {
		return true
	}
	for key := range view.Properties {
		if strings.EqualFold(key, "data") {
			return true
		}
	}
	return false
}
