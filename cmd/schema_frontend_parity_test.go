package cmd

import (
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
)

func TestDouyinSchemasExposeCurrentDeclarationOptions(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	validator := schema.NewValidator("schemas")
	for _, publishType := range []string{"video", "imageText"} {
		doc, err := validator.Schema("抖音", publishType)
		if err != nil {
			t.Fatal(err)
		}
		declaration, ok := doc.Properties["declaration"]
		if !ok {
			t.Fatalf("%s schema has no declaration field", publishType)
		}
		if !containsSchemaEnum(declaration.Enum, float64(7)) || !containsSchemaEnum(declaration.Enum, float64(8)) {
			t.Fatalf("%s declaration enum does not include current frontend values: %#v", publishType, declaration.Enum)
		}
	}
}

func TestWebParitySchemasExposeCurrentFieldNamesAndEnums(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	validator := schema.NewValidator("schemas")
	tests := []struct {
		platform string
		kind     string
		field    string
		values   []float64
	}{
		{"快手", "imageText", "visibleType", []float64{0, 1, 2}},
		{"快手", "imageText", "declaration", []float64{0, 1, 2, 3}},
		{"小红书", "imageText", "visibleType", []float64{0, 1, 2}},
		{"小红书", "imageText", "declaration", []float64{0, 1, 2}},
		{"新浪微博", "imageText", "visibleType", []float64{0, 1}},
		{"新浪微博", "imageText", "declaration", []float64{0, 1, 2, 3, 4}},
		{"视频号", "video", "declaration", []float64{0, 1, 2, 3, 7, 8}},
		{"一点号", "video", "createType", []float64{1, 2}},
	}
	for _, tc := range tests {
		doc, err := validator.Schema(tc.platform, tc.kind)
		if err != nil {
			t.Fatalf("%s/%s: %v", tc.platform, tc.kind, err)
		}
		field, ok := doc.Properties[tc.field]
		if !ok {
			t.Fatalf("%s/%s missing Web field %q", tc.platform, tc.kind, tc.field)
		}
		for _, value := range tc.values {
			if !containsSchemaEnum(field.Enum, value) {
				t.Fatalf("%s/%s field %s missing enum %v: %#v", tc.platform, tc.kind, tc.field, value, field.Enum)
			}
		}
	}
}

func containsSchemaEnum(values []interface{}, expected float64) bool {
	for _, value := range values {
		if number, ok := value.(float64); ok && number == expected {
			return true
		}
	}
	return false
}
