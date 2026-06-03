package cmd

import (
	"strings"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestDetectSwappedPublishArgsAllowsCorrectOrder(t *testing.T) {
	// publish form: <type> <platform>
	if err := detectSwappedPublishArgs("video", "抖音", "publish <type> <platform> <payload.json>"); err != nil {
		t.Fatalf("expected no error for correct order, got %v", err)
	}
	if err := detectSwappedPublishArgs("imageText", "douyin", "publish <type> <platform> <payload.json>"); err != nil {
		t.Fatalf("expected no error for english platform alias, got %v", err)
	}
}

func TestDetectSwappedPublishArgsFlagsReversedOrder(t *testing.T) {
	// User typed publish 抖音 video ... (platform and type transposed)
	err := detectSwappedPublishArgs("抖音", "video", "publish <type> <platform> <payload.json>")
	if err == nil {
		t.Fatal("expected swapped-order error")
	}
	if !strings.Contains(err.Error(), "wrong order") {
		t.Fatalf("expected wrong-order message, got %v", err)
	}
	typed, ok := err.(*yxerrors.Error)
	if !ok {
		t.Fatalf("expected structured usage error, got %T", err)
	}
	if !strings.Contains(typed.Hint, "publish <type> <platform>") {
		t.Fatalf("expected hint to name canonical order, got %q", typed.Hint)
	}
}

func TestDetectSwappedPublishArgsIgnoresUnknownTokens(t *testing.T) {
	// An unrecognized type with an unrecognized platform should not be flagged as
	// "swapped" — let the downstream schema lookup produce its own error.
	if err := detectSwappedPublishArgs("video", "unknown-platform", "publish <type> <platform> <payload.json>"); err != nil {
		t.Fatalf("expected no swap error for unknown platform, got %v", err)
	}
	if err := detectSwappedPublishArgs("unknown-type", "unknown-platform", "publish <type> <platform> <payload.json>"); err != nil {
		t.Fatalf("expected no swap error for fully unknown args, got %v", err)
	}
}
