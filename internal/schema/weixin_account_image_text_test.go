package schema

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestWeixinAccountImageTextSchemaExposesPublishSettings(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	doc, err := validator.Schema("WeiXinGongZhongHao", "imageText")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Key != "weixin.account/imageText" {
		t.Fatalf("unexpected schema key: %q", doc.Key)
	}

	images := doc.Properties["images"]
	if !images.Required || images.MinItems == nil || *images.MinItems != 1 {
		t.Fatalf("images must be required with minItems=1, got %#v", images)
	}
	title := doc.Properties["title"]
	if title.MaxLength == nil || *title.MaxLength != 20 {
		t.Fatalf("title must be limited to 20 characters, got %#v", title)
	}

	assertSchemaDefault(t, doc, "notifySubscribers", float64(0))
	assertSchemaDefault(t, doc, "sex", float64(0))
	assertSchemaDefault(t, doc, "needOpenComment", float64(0))
	assertSchemaDefault(t, doc, "statement", float64(0))
	assertSchemaDefault(t, doc, "disableRecommend", float64(0))
	assertSchemaDefault(t, doc, "pubType", float64(1))
	assertSchemaEnum(t, doc, "needOpenComment", 0, 1, 2, 3)
	assertSchemaEnum(t, doc, "statement", 0, 1, 2, 3, 4, 5)
	assertSchemaEnum(t, doc, "disableRecommend", 0, 1)
	assertSchemaEnum(t, doc, "sex", 0, 1, 2)
}

func TestWeixinAccountImageTextSchemaRequiresUploadedImages(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	result, err := validator.ValidateStrict("微信公众号", "imageText", map[string]interface{}{
		"formType": "task",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Valid {
		t.Fatalf("expected missing images to fail schema validation")
	}
	if len(result.Errors) == 0 {
		t.Fatalf("expected a schema error for missing images")
	}
}

func TestWeixinAccountImageTextSchemaAcceptsStandardAccountLevelImages(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	result, err := validator.ValidateStrict("微信公众号", "imageText", map[string]interface{}{
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_wx_1",
					"images": []interface{}{
						map[string]interface{}{"key": "uploaded-image"},
					},
					"contentPublishForm": map[string]interface{}{
						"formType": "task",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Valid {
		t.Fatalf("expected account-level images to satisfy imageText schema, got %v", result.Errors)
	}
}

func TestWeixinAccountImageTextSchemaAcceptsRecordedPayloadShape(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	result, err := validator.ValidateStrict("微信公众号", "imageText", map[string]interface{}{
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"images": []interface{}{
						map[string]interface{}{
							"key":    "yfb/local/image-1",
							"width":  float64(768),
							"height": float64(848),
							"size":   float64(131072),
						},
					},
					"platformAccountId": "6a75ae6089ecda991fd0d584",
					"coverKey":          "ml/local/cover-1",
					"contentPublishForm": map[string]interface{}{
						"pubType":           float64(1),
						"notifySubscribers": float64(1),
						"sex":               float64(0),
						"title":             "山东20岁女子失联多日",
						"desc":              "<p>图文正文</p>",
						"images":            []interface{}{},
						"needOpenComment":   float64(0),
						"statement":         float64(0),
						"disableRecommend":  float64(0),
					},
				},
			},
			"platformForms": map[string]interface{}{
				"微信公众号": map[string]interface{}{
					"pubType":           float64(1),
					"notifySubscribers": float64(0),
					"sex":               float64(0),
					"title":             "",
					"desc":              "",
					"images":            []interface{}{},
					"needOpenComment":   float64(0),
					"statement":         float64(0),
					"disableRecommend":  float64(0),
					"formType":          "task",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Valid {
		t.Fatalf("expected the recorded successful payload shape to validate, got %v", result.Errors)
	}
}

func TestWeixinAccountImageTextSchemaRejectsInvalidTitleAndSettings(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	result, err := validator.ValidateStrict("微信公众号", "imageText", map[string]interface{}{
		"formType": "task",
		"title":    "这是一个超过微信公众号图文二十字限制的标题内容",
		"images": []interface{}{
			map[string]interface{}{"key": "uploaded-image"},
		},
		"needOpenComment":  float64(4),
		"disableRecommend": float64(2),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Valid {
		t.Fatal("expected invalid title and settings to fail schema validation")
	}
	for _, want := range []string{"title: must NOT", "needOpenComment", "disableRecommend"} {
		found := false
		for _, item := range result.Errors {
			if strings.Contains(item, want) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected schema error containing %q, got %v", want, result.Errors)
		}
	}
}

func assertSchemaDefault(t *testing.T, doc Document, field string, want interface{}) {
	t.Helper()
	property, ok := doc.Properties[field]
	if !ok {
		t.Fatalf("schema is missing %s", field)
	}
	if property.Default != want {
		t.Fatalf("%s default = %#v, want %#v", field, property.Default, want)
	}
}

func assertSchemaEnum(t *testing.T, doc Document, field string, wants ...float64) {
	t.Helper()
	property, ok := doc.Properties[field]
	if !ok {
		t.Fatalf("schema is missing %s", field)
	}
	for _, want := range wants {
		found := false
		for _, value := range property.Enum {
			if number, ok := value.(float64); ok && number == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("%s enum is missing %v: %#v", field, want, property.Enum)
		}
	}
}
