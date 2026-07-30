package cmd

import (
	"testing"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
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
		{"快手-Open", "video", "visibleType", []float64{0, 1}},
		{"快手-Open", "video", "pubType", []float64{0, 1}},
		{"小红书", "imageText", "visibleType", []float64{0, 1, 3}},
		{"小红书", "imageText", "declaration", []float64{0, 1, 2}},
		{"新浪微博", "imageText", "visibleType", []float64{0, 1}},
		{"新浪微博", "imageText", "declaration", []float64{0, 1, 2, 3, 4}},
		{"视频号", "video", "declaration", []float64{0, 1, 2, 3, 7, 8}},
		{"哔哩哔哩-Open", "video", "allowReprint", []float64{0, 1}},
		{"哔哩哔哩-Open", "video", "createType", []float64{1, 2}},
		{"哔哩哔哩-Open", "video", "type", []float64{1, 2}},
		{"哔哩哔哩-Open", "video", "pubType", []float64{0, 1}},
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

func TestOpenPlatformVideoSchemasUseWebPlatformNames(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	validator := schema.NewValidator("schemas")
	tests := []struct {
		inputPlatform string
		wantKey       string
		wantName      string
		requiredField string
	}{
		{"BiLiBiLi-Open", "bilibili-open/video", "哔哩哔哩-Open", "title"},
		{"哔哩哔哩-Open", "bilibili-open/video", "哔哩哔哩-Open", "tags"},
		{"KuaiShou-Open", "kuaishou-open/video", "快手-Open", "description"},
		{"快手-Open", "kuaishou-open/video", "快手-Open", "description"},
	}
	for _, tc := range tests {
		doc, err := validator.Schema(tc.inputPlatform, "video")
		if err != nil {
			t.Fatalf("%s/video: %v", tc.inputPlatform, err)
		}
		if doc.Key != tc.wantKey {
			t.Fatalf("%s/video key = %q, want %q", tc.inputPlatform, doc.Key, tc.wantKey)
		}
		if got := platformutil.ChineseName(doc.Platform); got != tc.wantName {
			t.Fatalf("%s/video chinese name = %q, want %q", tc.inputPlatform, got, tc.wantName)
		}
		field := doc.Properties[tc.requiredField]
		if !field.Required {
			t.Fatalf("%s/video field %s should be required, got %#v", tc.inputPlatform, tc.requiredField, field)
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
