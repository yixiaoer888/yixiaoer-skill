package platform

import "strings"

var canonicalChineseNames = map[string]string{
	"douyin":          "抖音",
	"kuaishou":        "快手",
	"xhs":             "小红书",
	"xiaohongshushop": "小红书商家号",
	"shipinhao":       "视频号",
	"weixin.account":  "微信公众号",
	"bilibili":        "哔哩哔哩",
	"baijiahao":       "百家号",
	"toutiaohao":      "头条号",
	"zhihu":           "知乎",
	"qiehao":          "企鹅号",
	"xinlang":         "新浪微博",
	"souhuhao":        "搜狐号",
	"souhushipin":     "搜狐视频",
	"yidianhao":       "一点号",
	"dayuhao":         "大鱼号",
	"wangyihao":       "网易号",
	"aiqiyi":          "爱奇艺",
	"weishi":          "腾讯微视",
	"tengxunshipin":   "腾讯视频",
	"pipixia":         "皮皮虾",
	"duoduoshipin":    "多多视频",
	"meipai":          "美拍",
	"kuaichuanhao":    "快传号",
	"xueqiuhao":       "雪球号",
	"chejiahao":       "车家号",
	"yichehao":        "易车号",
	"fengwang":        "蜂网",
	"douban":          "豆瓣",
	"dewu":            "得物",
	"jianshu":         "简书",
	"meiyou":          "美柚",
}

var aliasesToCanonicalKeys = map[string]string{
	"douyin":           "douyin",
	"抖音":               "douyin",
	"kuaishou":         "kuaishou",
	"快手":               "kuaishou",
	"xhs":              "xhs",
	"xiaohongshu":      "xhs",
	"小红书":              "xhs",
	"xiaohongshushop":  "xiaohongshushop",
	"小红书商家号":           "xiaohongshushop",
	"shipinhao":        "shipinhao",
	"shipinghao":       "shipinhao",
	"视频号":              "shipinhao",
	"微信视频号":            "shipinhao",
	"weixin.account":   "weixin.account",
	"微信公众号":            "weixin.account",
	"bilibili":         "bilibili",
	"哔哩哔哩":             "bilibili",
	"baijiahao":        "baijiahao",
	"百家号":              "baijiahao",
	"toutiaohao":       "toutiaohao",
	"头条号":              "toutiaohao",
	"zhihu":            "zhihu",
	"知乎":               "zhihu",
	"qiehao":           "qiehao",
	"企鹅号":              "qiehao",
	"xinlang":          "xinlang",
	"新浪微博":             "xinlang",
	"souhuhao":         "souhuhao",
	"搜狐号":              "souhuhao",
	"souhushipin":      "souhushipin",
	"搜狐视频":             "souhushipin",
	"yidianhao":        "yidianhao",
	"一点号":              "yidianhao",
	"dayuhao":          "dayuhao",
	"大鱼号":              "dayuhao",
	"wangyihao":        "wangyihao",
	"网易号":              "wangyihao",
	"aiqiyi":           "aiqiyi",
	"爱奇艺":              "aiqiyi",
	"weishi":           "weishi",
	"腾讯微视":             "weishi",
	"tengxunshipin":    "tengxunshipin",
	"腾讯视频":             "tengxunshipin",
	"pipixia":          "pipixia",
	"皮皮虾":              "pipixia",
	"duoduoshipin":     "duoduoshipin",
	"多多视频":             "duoduoshipin",
	"meipai":           "meipai",
	"美拍":               "meipai",
	"kuaichuanhao":     "kuaichuanhao",
	"快传号":              "kuaichuanhao",
	"xueqiuhao":        "xueqiuhao",
	"雪球号":              "xueqiuhao",
	"chejiahao":        "chejiahao",
	"车家号":              "chejiahao",
	"yichehao":         "yichehao",
	"易车号":              "yichehao",
	"fengwang":         "fengwang",
	"蜂网":               "fengwang",
	"douban":           "douban",
	"豆瓣":               "douban",
	"dewu":             "dewu",
	"得物":               "dewu",
	"jianshu":          "jianshu",
	"简书":               "jianshu",
	"meiyou":           "meiyou",
	"美柚":               "meiyou",
}

var chineseNames = map[string]string{
	"douyin":           "抖音",
	"kuaishou":         "快手",
	"xhs":              "小红书",
	"xiaohongshu":      "小红书",
	"xiaohongshushop":  "小红书商家号",
	"shipinhao":        "视频号",
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

func CanonicalKey(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if key, ok := aliasesToCanonicalKeys[strings.ToLower(trimmed)]; ok {
		return key
	}
	if key, ok := aliasesToCanonicalKeys[trimmed]; ok {
		return key
	}
	return strings.ToLower(trimmed)
}

func ChineseName(value string) string {
	key := CanonicalKey(value)
	if key == "" {
		return ""
	}
	if name, ok := canonicalChineseNames[key]; ok {
		return name
	}
	trimmed := strings.TrimSpace(value)
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
	if _, ok := aliasesToCanonicalKeys[strings.ToLower(trimmed)]; ok {
		return true
	}
	if _, ok := aliasesToCanonicalKeys[trimmed]; ok {
		return true
	}
	return chineseNameSet[trimmed]
}
