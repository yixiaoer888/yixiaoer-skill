package cmd

import (
	"strings"
	"testing"
)

func TestAnalyzeValidationErrorsUsesManualDuoduoshipinGoodsIDGuidance(t *testing.T) {
	suggestions := analyzeValidationErrors([]string{`shopping_cart: missing required field "goods_id"`}, "多多视频", "video")
	if len(suggestions) != 1 {
		t.Fatalf("expected one validation suggestion, got %#v", suggestions)
	}
	suggestion := suggestions[0]
	fix := suggestion["fix"].(string)
	if !strings.Contains(fix, "goods_id") || !strings.Contains(fix, "source") || !strings.Contains(fix, "用户") {
		t.Fatalf("expected manual goods_id guidance, got %#v", suggestion)
	}
	if suggestion["exampleField"] != "publishArgs.accountForms[].contentPublishForm.shopping_cart.goods_id" {
		t.Fatalf("expected manual goods_id field path, got %#v", suggestion["exampleField"])
	}
	if suggestion["reference"] == "yxer query goods <account_id> [--query 关键词] --json" {
		t.Fatalf("Duoduoshipin guidance must not use goods query as the reference, got %#v", suggestion)
	}
}
