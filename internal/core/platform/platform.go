package platform

import base "github.com/yixiaoer/yixiaoer-skill/internal/platform"

func CanonicalKey(value string) string {
	return base.CanonicalKey(value)
}

func ChineseName(value string) string {
	return base.ChineseName(value)
}

// IsKnown reports whether value is a recognized platform alias or Chinese name.
func IsKnown(value string) bool {
	return base.IsKnown(value)
}
