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
		field    []string
		values   []float64
	}{
		{"快手", "imageText", []string{"visibleType"}, []float64{0, 1, 2}},
		{"快手", "imageText", []string{"declaration"}, []float64{0, 1, 2, 3}},
		{"快手-Open", "video", []string{"visibleType"}, []float64{0, 1}},
		{"快手-Open", "video", []string{"pubType"}, []float64{0, 1}},
		{"小红书", "imageText", []string{"visibleType"}, []float64{0, 1, 3}},
		{"小红书", "imageText", []string{"declaration"}, []float64{0, 1, 2}},
		{"新浪微博", "imageText", []string{"visibleType"}, []float64{0, 1}},
		{"新浪微博", "imageText", []string{"declaration"}, []float64{0, 1, 2, 3, 4}},
		{"视频号", "video", []string{"declaration"}, []float64{0, 1, 2, 3, 7, 8}},
		{"WeiXinGongZhongHao", "imageText", []string{"needOpenComment"}, []float64{0, 1, 2, 3}},
		{"WeiXinGongZhongHao", "imageText", []string{"statement"}, []float64{0, 1, 2, 3, 4, 5}},
		{"WeiXinGongZhongHao", "imageText", []string{"disableRecommend"}, []float64{0, 1}},
		{"哔哩哔哩-Open", "video", []string{"contentPublishForm", "allowReprint"}, []float64{0, 1}},
		{"哔哩哔哩-Open", "video", []string{"contentPublishForm", "createType"}, []float64{1, 2}},
		{"哔哩哔哩-Open", "video", []string{"contentPublishForm", "type"}, []float64{1, 2}},
		{"哔哩哔哩-Open", "video", []string{"contentPublishForm", "pubType"}, []float64{0, 1}},
		{"一点号", "video", []string{"createType"}, []float64{1, 2}},
	}
	for _, tc := range tests {
		doc, err := validator.Schema(tc.platform, tc.kind)
		if err != nil {
			t.Fatalf("%s/%s: %v", tc.platform, tc.kind, err)
		}
		field, ok := schemaField(doc.Properties, tc.field...)
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
		requiredField []string
	}{
		{"BiLiBiLi-Open", "bilibili-open/video", "哔哩哔哩-Open", []string{"contentPublishForm", "title"}},
		{"哔哩哔哩-Open", "bilibili-open/video", "哔哩哔哩-Open", []string{"contentPublishForm", "tags"}},
		{"KuaiShou-Open", "kuaishou-open/video", "快手-Open", []string{"description"}},
		{"快手-Open", "kuaishou-open/video", "快手-Open", []string{"description"}},
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
		field, ok := schemaField(doc.Properties, tc.requiredField...)
		if !ok {
			t.Fatalf("%s/video missing required field %q", tc.inputPlatform, tc.requiredField)
		}
		if !field.Required {
			t.Fatalf("%s/video field %s should be required, got %#v", tc.inputPlatform, tc.requiredField, field)
		}
	}
}

func schemaField(properties map[string]schema.PropertyView, path ...string) (schema.PropertyView, bool) {
	for index, name := range path {
		field, ok := properties[name]
		if !ok {
			return schema.PropertyView{}, false
		}
		if index == len(path)-1 {
			return field, true
		}
		properties = field.Properties
	}
	return schema.PropertyView{}, false
}

func containsSchemaEnum(values []interface{}, expected float64) bool {
	for _, value := range values {
		if number, ok := value.(float64); ok && number == expected {
			return true
		}
	}
	return false
}
