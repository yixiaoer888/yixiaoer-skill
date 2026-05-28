package publish

import base "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"

type PreflightResult = base.PreflightResult

func Preflight(publishType string, platforms []string, payload map[string]interface{}) PreflightResult {
	return base.Preflight(publishType, platforms, payload)
}

func ExtractPublishArgs(payload map[string]interface{}) map[string]interface{} {
	return base.ExtractPublishArgs(payload)
}

func ValidateAndExtractPublishArgs(publishType string, platforms []string, payload map[string]interface{}, errors *[]string) map[string]interface{} {
	return base.ValidateAndExtractPublishArgs(publishType, platforms, payload, errors)
}

func NormalizeStandardPublishArgs(payload map[string]interface{}) {
	base.NormalizeStandardPublishArgs(payload)
}

func NormalizeScheduledTimes(value interface{}, errors *[]string) {
	base.NormalizeScheduledTimes(value, errors)
}
