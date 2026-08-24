package publish

import (
	"fmt"
	"strings"
	"time"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
)

const weixinAccountImageTextScheduleLeadTime = 2 * time.Hour

// normalizeWeixinAccountImageTextFields applies defaults owned by the
// 微信公众号 imageText form. It intentionally does not touch the existing
// 微信公众号 article platform form.
func normalizeWeixinAccountImageTextFields(publishType string, platforms []string, publishArgs map[string]interface{}) {
	if NormalizePublishType(publishType) != "imageText" || !hasWeixinAccountPlatform(platforms) || publishArgs == nil {
		return
	}
	ensureWeixinAccountImageTextPlatformForm(publishArgs)
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	for _, item := range accountForms {
		form, _ := item.(map[string]interface{})
		if form == nil {
			continue
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}
		setWeixinImageTextNumberDefault(cpf, "notifySubscribers", 0)
		setWeixinImageTextNumberDefault(cpf, "sex", 0)
		setWeixinImageTextNumberDefault(cpf, "needOpenComment", 0)
		setWeixinImageTextNumberDefault(cpf, "statement", 0)
		setWeixinImageTextNumberDefault(cpf, "disableRecommend", 0)
		setWeixinImageTextNumberDefault(cpf, "pubType", 1)
	}
}

func ensureWeixinAccountImageTextPlatformForm(publishArgs map[string]interface{}) {
	platformForms, _ := publishArgs["platformForms"].(map[string]interface{})
	if platformForms == nil {
		platformForms = map[string]interface{}{}
		publishArgs["platformForms"] = platformForms
	}

	var form map[string]interface{}
	for _, key := range []string{"微信公众号", "weixin.account"} {
		candidate, _ := platformForms[key].(map[string]interface{})
		if candidate != nil {
			form = candidate
			break
		}
	}
	if form == nil {
		form = map[string]interface{}{}
		platformForms["微信公众号"] = form
	}

	setWeixinImageTextNumberDefault(form, "pubType", 1)
	setWeixinImageTextNumberDefault(form, "notifySubscribers", 0)
	setWeixinImageTextNumberDefault(form, "sex", 0)
	setWeixinImageTextNumberDefault(form, "needOpenComment", 0)
	setWeixinImageTextNumberDefault(form, "statement", 0)
	setWeixinImageTextNumberDefault(form, "disableRecommend", 0)
	if _, exists := form["title"]; !exists {
		form["title"] = ""
	}
	if _, exists := form["desc"]; !exists {
		form["desc"] = ""
	}
	if _, exists := form["images"]; !exists {
		form["images"] = []interface{}{}
	}
	if _, exists := form["formType"]; !exists {
		form["formType"] = "task"
	}
}

func hasWeixinAccountPlatform(platforms []string) bool {
	for _, value := range platforms {
		if platformutil.CanonicalKey(value) == "weixin.account" {
			return true
		}
	}
	return false
}

func setWeixinImageTextNumberDefault(form map[string]interface{}, field string, value float64) {
	if _, exists := form[field]; !exists {
		form[field] = value
	}
}

func validateWeixinAccountImageTextSchedules(publishType string, platforms []string, publishArgs map[string]interface{}, now time.Time, errors *[]string) {
	if NormalizePublishType(publishType) != "imageText" || !hasWeixinAccountPlatform(platforms) || publishArgs == nil {
		return
	}
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	for index, item := range accountForms {
		form, _ := item.(map[string]interface{})
		if form == nil {
			continue
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}
		value, exists := cpf["scheduledTime"]
		if !exists || value == nil {
			continue
		}
		path := fmt.Sprintf("accountForms[%d].contentPublishForm.scheduledTime", index)
		validateWeixinAccountImageTextSchedule(value, now, path, errors)
	}
}

func validateWeixinAccountImageTextSchedule(value interface{}, now time.Time, path string, errors *[]string) {
	normalized, normalizeErr := normalizeScheduledTime(value)
	if normalizeErr != "" {
		// NormalizeScheduledTimes already emits the canonical timestamp error.
		return
	}
	scheduled, ok := normalized.(float64)
	if !ok {
		return
	}
	minimum := now.Add(weixinAccountImageTextScheduleLeadTime).UnixMilli()
	if int64(scheduled) < minimum {
		*errors = append(*errors, strings.TrimSpace(path)+": scheduledTime must be at least 2 hours after the current time (微信公众号图文定时发布时间不得早于当前时间2小时)")
	}
}
