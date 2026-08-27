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
