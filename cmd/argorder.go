package cmd

import (
	"fmt"
	"strings"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/core/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

// knownPublishType reports whether value is a recognized publish type. Matching
// is case-insensitive so "ImageText" and "imagetext" both resolve.
func knownPublishType(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "video", "imagetext", "article":
		return true
	default:
		return false
	}
}

// detectSwappedPublishArgs guards against the most common command-line mistake:
// publish expects "<type> <platform>" while validate expects "<platform> <type>",
// so users routinely transpose the two. When the supplied (typeArg, platformArg)
// pair only makes sense in the opposite order, return a usage error that names
// the correct order instead of letting a confusing schema-not-found error
// surface downstream. usageOrder is the canonical argument order for the command
// (e.g. "publish <type> <platform>").
func detectSwappedPublishArgs(typeArg, platformArg, usageOrder string) error {
	typeOK := knownPublishType(typeArg)
	platformOK := platformutil.IsKnown(platformArg)
	if typeOK && platformOK {
		return nil
	}
	// Only flag when the reversed reading is unambiguously correct: the supposed
	// type is actually a platform AND the supposed platform is actually a type.
	if platformutil.IsKnown(typeArg) && knownPublishType(platformArg) {
		return yxerrors.Usage("publish arguments appear to be in the wrong order", map[string]interface{}{
			"received":     fmt.Sprintf("%s %s", typeArg, platformArg),
			"expectedForm": usageOrder,
		}).
			WithHint(fmt.Sprintf("参数顺序应为 %s；注意 publish 用 <type> <platform>，validate 用 <platform> <type>，二者相反。", usageOrder)).
			WithNextCommand("yxer " + usageOrder)
	}
	return nil
}
