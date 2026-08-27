package publish

import "strings"

func schemaValidationHint(errors []string) string {
	joined := strings.Join(errors, "\n")
	switch {
	case strings.Contains(joined, "/drama") || strings.Contains(joined, `field "drama"`):
		return "视频号剧集对象必须先通过 yxer query drama-tasks 查询，并只保留 yixiaoerId、yixiaoerImageUrl、yixiaoerName 三个字段；剧集对象不使用 raw。"
	case strings.Contains(joined, "shopping_cart") || strings.Contains(joined, "group_shopping") || strings.Contains(joined, "shoppingCart") || strings.Contains(joined, "groupShopping"):
		return "商品相关字段必须使用 schema/dynamicFieldExamples 中的前端表单结构；抖音购物车需要 sale_title、images、data，团购使用 group_shopping，并应从 yxer query goods 返回结果复制。"
	case strings.Contains(joined, "location") && strings.Contains(joined, "raw"):
		return "位置对象必须从 yxer query locations 返回结果复制完整对象，并按前端表单结构回填；抖音位置是 {isScp,data}，其他位置通常是 {id,text,raw}。"
	case strings.Contains(joined, "music") && strings.Contains(joined, "raw"):
		return "音乐对象必须从 yxer query music 返回结果复制完整对象，保留 playUrl/url/duration/raw 等字段。"
	case strings.Contains(joined, "tags") || strings.Contains(joined, "topic") || strings.Contains(joined, "challenge"):
		return "话题标签请先查看 schema fields 的 dynamicFieldExamples.tags；普通 tags 使用字符串数组，challenge/topics 等动态对象必须使用查询返回的完整 raw。"
	default:
		return "请根据对应平台 schema 修正 payload 字段后重试。"
	}
}
