package publish

import (
	"strings"
	"testing"
)

func TestSchemaValidationHintExplainsDramaQueryShape(t *testing.T) {
	hint := schemaValidationHint([]string{`/drama: unexpected field "raw"`})
	for _, expected := range []string{"drama-tasks", "yixiaoerImageUrl", "raw"} {
		if !strings.Contains(hint, expected) {
			t.Fatalf("expected drama hint to contain %q, got %q", expected, hint)
		}
	}
}

func TestSchemaValidationHintExplainsDuoduoshipinManualGoodsID(t *testing.T) {
	hint := schemaValidationHintForPlatform("多多视频", []string{`/shopping_cart: missing required field "goods_id"`})
	if !strings.Contains(hint, "用户") || !strings.Contains(hint, "goods_id") || !strings.Contains(hint, "source=pdd") {
		t.Fatalf("expected Duoduoshipin manual goods_id hint, got %q", hint)
	}
	if strings.Contains(hint, "query goods") {
		t.Fatalf("Duoduoshipin hint must not direct users to query goods, got %q", hint)
	}
}
