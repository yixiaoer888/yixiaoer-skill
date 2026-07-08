package cmd

import (
	"fmt"
	"strings"
)

// analyzeValidationErrors 分析校验错误并提供修复建议
func analyzeValidationErrors(errors []string, platform, publishType string) []map[string]interface{} {
	suggestions := []map[string]interface{}{}

	for _, err := range errors {
		suggestion := map[string]interface{}{
			"error": err,
		}

		// 分析错误类型并给出建议
		if strings.Contains(err, "shopping_cart") || strings.Contains(err, "group_shopping") || strings.Contains(err, "shoppingCart") || strings.Contains(err, "groupShopping") {
			suggestion["reason"] = "购物车商品结构或来源不正确"
			suggestion["fix"] = "使用 dynamicFieldExamples 中的前端表单结构，并从 yxer query goods 返回结果复制完整商品对象；抖音购物车必须包含 sale_title、images、data.raw，团购使用 group_shopping"
			suggestion["reference"] = "yxer query goods <account_id> [--query 关键词] --json"
			suggestion["exampleField"] = "dynamicFieldExamples.shopping_cart"

		} else if strings.Contains(err, "location") && strings.Contains(err, "raw") {
			suggestion["reason"] = "位置对象缺少完整 raw"
			suggestion["fix"] = "通过位置查询命令获取完整对象，并按 schema 前端表单结构回填；抖音位置是 {isScp,data}，其他位置通常是 {id,text,raw}"
			suggestion["reference"] = "yxer query locations <account_id> [--query 关键词] --json"
			suggestion["exampleField"] = "dynamicFieldExamples.location"

		} else if strings.Contains(err, "music") && strings.Contains(err, "raw") {
			suggestion["reason"] = "音乐对象缺少完整 raw"
			suggestion["fix"] = "通过音乐查询命令获取完整对象，保留 playUrl/url/duration/raw 等返回字段"
			suggestion["reference"] = "yxer query music <account_id> [--query 关键词] --json"
			suggestion["exampleField"] = "dynamicFieldExamples.music"

		} else if strings.Contains(err, "tags") || strings.Contains(err, "topic") || strings.Contains(err, "challenge") {
			suggestion["reason"] = "话题标签字段结构不符合平台要求"
			suggestion["fix"] = "普通 tags 使用字符串数组；challenge/topics 等动态对象必须从查询命令复制完整 raw"
			suggestion["reference"] = fmt.Sprintf("yxer schema fields %s %s", platform, publishType)
			suggestion["exampleField"] = "dynamicFieldExamples.tags"

		} else if strings.Contains(err, "publishArgs") && strings.Contains(err, "required") {
			suggestion["reason"] = "payload 结构不符合标准"
			suggestion["fix"] = "确保顶层有 publishArgs 字段"
			suggestion["reference"] = fmt.Sprintf("yxer schema get %s %s", platform, publishType)

		} else if strings.Contains(err, "platformAccountId") {
			suggestion["reason"] = "缺少账号 ID"
			suggestion["fix"] = "从 accounts list 获取 platformAccountId"
			suggestion["reference"] = fmt.Sprintf("yxer accounts list %s --status 1", platform)

		} else if strings.Contains(err, "accountForms") && strings.Contains(err, "required") {
			suggestion["reason"] = "缺少账号表单数组"
			suggestion["fix"] = "在 publishArgs 下添加 accountForms 数组"
			suggestion["reference"] = fmt.Sprintf("yxer schema fields %s %s", platform, publishType)

		} else if strings.Contains(err, "contentPublishForm") && strings.Contains(err, "required") {
			suggestion["reason"] = "缺少业务表单对象"
			suggestion["fix"] = "在 accountForms[].contentPublishForm 中填写业务字段"
			suggestion["reference"] = fmt.Sprintf("yxer schema fields %s %s", platform, publishType)

		} else if strings.Contains(err, "formType") {
			suggestion["reason"] = "缺少表单类型"
			suggestion["fix"] = "在 contentPublishForm 中添加 formType: 'task'"
			suggestion["reference"] = "formType 是固定值"

		} else if strings.Contains(err, "title") || strings.Contains(err, "description") {
			suggestion["reason"] = "缺少必填文本字段"
			suggestion["fix"] = "添加 title 和 description 字段"
			suggestion["reference"] = "title 和 description 通常是必填字段"

		} else if strings.Contains(err, "key") || strings.Contains(err, "size") || strings.Contains(err, "width") || strings.Contains(err, "height") {
			suggestion["reason"] = "资源对象不完整"
			suggestion["fix"] = "使用 yxer upload 返回的完整对象，包含 key/size/width/height 等字段"
			suggestion["reference"] = "yxer upload <file_path>"

		} else if strings.Contains(err, "raw") {
			suggestion["reason"] = "复杂对象缺少 raw 字段"
			suggestion["fix"] = "通过查询命令获取完整对象（包含 raw），不能手动构造"
			suggestion["reference"] = "参考 queryCommands 中的查询命令"

		} else if strings.Contains(err, "must be one of") {
			suggestion["reason"] = "字段值不在枚举范围内"
			suggestion["fix"] = "检查字段的可选值，使用 schema fields 查看 enum 选项"
			suggestion["reference"] = fmt.Sprintf("yxer schema fields %s %s", platform, publishType)

		} else if strings.Contains(err, "expected") && strings.Contains(err, "type") {
			suggestion["reason"] = "字段类型不匹配"
			suggestion["fix"] = "检查字段类型，确保 string/number/object/array 正确"
			suggestion["reference"] = fmt.Sprintf("yxer schema fields %s %s", platform, publishType)

		} else if strings.Contains(err, "must NOT have fewer than") {
			suggestion["reason"] = "字符串长度不足"
			suggestion["fix"] = "增加字段内容长度"
			suggestion["reference"] = "查看 schema 中的 minLength 要求"

		} else if strings.Contains(err, "must NOT have more than") {
			suggestion["reason"] = "字符串长度超限"
			suggestion["fix"] = "减少字段内容长度"
			suggestion["reference"] = "查看 schema 中的 maxLength 要求"

		} else if strings.Contains(err, "must have at least") {
			suggestion["reason"] = "数组元素数量不足"
			suggestion["fix"] = "增加数组元素"
			suggestion["reference"] = "查看 schema 中的 minItems 要求"

		} else if strings.Contains(err, "unexpected field") {
			suggestion["reason"] = "字段不在 schema 中"
			suggestion["fix"] = "移除该字段或检查字段名拼写"
			suggestion["reference"] = fmt.Sprintf("yxer schema fields %s %s", platform, publishType)

		} else {
			suggestion["reason"] = "请查看 schema 定义"
			suggestion["fix"] = "根据错误信息修正字段"
			suggestion["reference"] = fmt.Sprintf("yxer schema fields %s %s", platform, publishType)
		}

		suggestions = append(suggestions, suggestion)
	}

	return suggestions
}
