package cmd

import (
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
)

// buildXiaohongshuMentionContract documents the inline description contract
// used by the Yixiaoer client. A mention is not a standalone friends payload
// field: the complete user object returned by query friends is serialized into
// the <friend> raw attribute.
func buildXiaohongshuMentionContract(doc schema.Document) map[string]interface{} {
	if platformutil.CanonicalKey(doc.Platform) != "xhs" || (doc.Type != "imageText" && doc.Type != "video") {
		return nil
	}
	return map[string]interface{}{
		"field":             "description",
		"path":              "publishArgs.accountForms[].contentPublishForm.description",
		"queryCommand":      "yxer query friends <account_id> --red-id <小红书号> --json",
		"keywordFlags":      []string{"--query", "--keyword"},
		"supportsStrangers": true,
		"search": map[string]interface{}{
			"gatewayParameter": "keyWord",
			"clientParameter":  "keyword",
			"exactMatchField":  "raw.red_id",
		},
		"syntax": "<friend raw='完整用户对象 JSON'>@用户昵称</friend>",
		"rules": []string{
			"目标用户可以不在当前账号的好友列表中。",
			"raw 必须使用本次 query friends 返回的完整对象，不能只手写小红书号、昵称或内部 raw。",
			"不要在 payload 中新增未被 schema 声明的 friends 字段。",
		},
	}
}
