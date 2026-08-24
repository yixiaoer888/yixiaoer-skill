package publish

import (
	"fmt"
	"strings"
	"time"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
)

const weixinAccountImageTextScheduleLeadTime = 2 * time.Hour

// normalizeWeixinAccountImageTextFields applies 微信公众号 imageText settings to
// each account form. The gateway reads contentPublishForm before platformForms
// for this publish type, so the account form is authoritative.
func normalizeWeixinAccountImageTextFields(publishType string, platforms []string, publishArgs map[string]interface{}) {
	if NormalizePublishType(publishType) != "imageText" || !hasWeixinAccountPlatform(platforms) || publishArgs == nil {
		return
	}
	legacyPlatformForm := weixinAccountImageTextLegacyPlatformForm(publishArgs)
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
		for _, key := range weixinAccountImageTextSettingKeys {
			if _, exists := cpf[key]; exists {
				continue
			}
			if value, exists := legacyPlatformForm[key]; exists {
				cpf[key] = normalizeWeixinImageTextStatement(key, value)
			}
		}
		if value, exists := cpf["statement"]; exists {
			cpf["statement"] = normalizeWeixinImageTextStatement("statement", value)
		}
		setWeixinImageTextNumberDefault(cpf, "notifySubscribers", 0)
		setWeixinImageTextNumberDefault(cpf, "sex", 0)
		setWeixinImageTextNumberDefault(cpf, "needOpenComment", 0)
		setWeixinImageTextNumberDefault(cpf, "statement", 0)
		setWeixinImageTextNumberDefault(cpf, "disableRecommend", 0)
		setWeixinImageTextNumberDefault(cpf, "pubType", 1)
	}
	removeWeixinAccountImageTextLegacyPlatformForm(publishArgs)
}

var weixinAccountImageTextSettingKeys = []string{
	"prePubTime",
	"scheduledTime",
	"notifySubscribers",
	"sex",
	"needOpenComment",
	"statement",
	"disableRecommend",
	"pubType",
}

func weixinAccountImageTextLegacyPlatformForm(publishArgs map[string]interface{}) map[string]interface{} {
	platformForms, _ := publishArgs["platformForms"].(map[string]interface{})
	if platformForms == nil {
		return nil
	}
	for _, key := range []string{"微信公众号", "weixin.account"} {
		if form, _ := platformForms[key].(map[string]interface{}); form != nil {
			return form
		}
	}
	return nil
}

func removeWeixinAccountImageTextLegacyPlatformForm(publishArgs map[string]interface{}) {
	platformForms, _ := publishArgs["platformForms"].(map[string]interface{})
	if platformForms == nil {
		return
	}
	delete(platformForms, "微信公众号")
	delete(platformForms, "weixin.account")
	if len(platformForms) == 0 {
		delete(publishArgs, "platformForms")
	}
}

func normalizeWeixinImageTextStatement(key string, value interface{}) interface{} {
	if key != "statement" {
		return value
	}
	if statement, ok := value.(map[string]interface{}); ok {
		if statementType, exists := statement["type"]; exists {
			return statementType
		}
	}
	return value
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
		contentForm, _ := form["contentPublishForm"].(map[string]interface{})
		if contentForm == nil {
			continue
		}
		if value, exists := contentForm["scheduledTime"]; exists && value != nil {
			validateWeixinAccountImageTextSchedule(value, now, fmt.Sprintf("accountForms[%d].contentPublishForm.scheduledTime", index), errors)
		}
	}
}

func validateWeixinAccountImageTextPlatformSettings(publishType string, platforms []string, publishArgs map[string]interface{}, errors *[]string) {
	if NormalizePublishType(publishType) != "imageText" || !hasWeixinAccountPlatform(platforms) || publishArgs == nil {
		return
	}
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	for index, item := range accountForms {
		form, _ := item.(map[string]interface{})
		contentForm, _ := form["contentPublishForm"].(map[string]interface{})
		if contentForm == nil {
			continue
		}
		for field, allowed := range map[string][]int{
			"notifySubscribers": {0, 1},
			"sex":               {0, 1, 2},
			"needOpenComment":   {0, 1, 2, 3},
			"statement":         {0, 1, 3, 4, 5, 6},
			"disableRecommend":  {0, 1},
			"pubType":           {0, 1},
		} {
			if value, exists := contentForm[field]; exists {
				validateWeixinImageTextSettingEnum(value, allowed, fmt.Sprintf("accountForms[%d].contentPublishForm.%s", index, field), errors)
			}
		}
	}
}

func validateWeixinImageTextSettingEnum(value interface{}, allowed []int, path string, errors *[]string) {
	number, ok := weixinImageTextSettingNumber(value)
	if !ok || number != float64(int(number)) {
		*errors = append(*errors, path+": must be a numeric enum value")
		return
	}
	for _, candidate := range allowed {
		if number == float64(candidate) {
			return
		}
	}
	*errors = append(*errors, path+": must be one of "+formatWeixinImageTextSettingEnum(allowed))
}

func weixinImageTextSettingNumber(value interface{}) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case int32:
		return float64(typed), true
	default:
		return 0, false
	}
}

func formatWeixinImageTextSettingEnum(values []int) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, string(rune('0'+value)))
	}
	return "[" + strings.Join(parts, " ") + "]"
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
