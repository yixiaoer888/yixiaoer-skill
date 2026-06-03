package platform

import "strings"

var chineseNames = map[string]string{
	"douyin":           "抖音",
	"kuaishou":         "快手",
	"xhs":              "小红书",
	"xiaohongshu":      "小红书",
	"xiaohongshushop":  "小红书商家号",
	"shipinghao":       "视频号",
	"weixin.account":   "微信公众号",
	"bilibili":         "哔哩哔哩",
	"baijiahao":        "百家号",
	"toutiaohao":       "头条号",
	"zhihu":            "知乎",
	"qiehao":           "企鹅号",
	"xinlang":          "新浪微博",
	"souhuhao":         "搜狐号",
	"souhushipin":      "搜狐视频",
	"yidianhao":        "一点号",
	"dayuhao":          "大鱼号",
	"wangyihao":        "网易号",
	"aiqiyi":           "爱奇艺",
	"weishi":           "腾讯微视",
	"tengxunshipin":    "腾讯视频",
	"pipixia":          "皮皮虾",
	"duoduoshipin":     "多多视频",
	"meipai":           "美拍",
	"kuaichuanhao":     "快传号",
	"xueqiuhao":        "雪球号",
	"chejiahao":        "车家号",
	"yichehao":         "易车号",
	"fengwang":         "蜂网",
	"douban":           "豆瓣",
	"dewu":             "得物",
	"jianshu":          "简书",
	"meiyou":           "美柚",
}

func ChineseName(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if name, ok := chineseNames[strings.ToLower(trimmed)]; ok {
		return name
	}
	return trimmed
}

// chineseNameSet holds every canonical Chinese platform name for reverse lookup.
var chineseNameSet = func() map[string]bool {
	set := make(map[string]bool, len(chineseNames))
	for _, name := range chineseNames {
		set[name] = true
	}
	return set
}()

// IsKnown reports whether value is a recognized platform, either by English
// alias (e.g. "douyin") or canonical Chinese name (e.g. "抖音"). It is used to
// detect swapped <platform>/<type> command arguments.
func IsKnown(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}
	if _, ok := chineseNames[strings.ToLower(trimmed)]; ok {
		return true
	}
	return chineseNameSet[trimmed]
}
