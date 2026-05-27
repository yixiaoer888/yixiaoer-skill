package schema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateWithPlatformSchema(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := readTestPayload(t, "douyin-video-valid.json")
	result := validator.Validate("抖音", "video", payload)
	if !result.Valid {
		t.Fatalf("expected valid payload, got errors: %v", result.Errors)
	}
}

func TestValidateRejectsAdditionalProperties(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := readTestPayload(t, "douyin-extra-field.json")
	result := validator.Validate("抖音", "video", payload)
	if result.Valid {
		t.Fatal("expected schema validation to reject extra field")
	}
	if !containsError(result.Errors, `unexpected field "extraField"`) {
		t.Fatalf("expected extra field error, got %v", result.Errors)
	}
}

func TestTypeKeyMapsImageText(t *testing.T) {
	if got := TypeKey("image-text"); got != "imageText" {
		t.Fatalf("expected imageText, got %s", got)
	}
}

func TestDisplayTypeMapsImageText(t *testing.T) {
	if got := DisplayType("imageText"); got != "image-text" {
		t.Fatalf("expected image-text, got %s", got)
	}
}

func TestSchemaReturnsAliasMatchedSchema(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	schemaDoc, err := validator.Schema("抖音", "video")
	if err != nil {
		t.Fatal(err)
	}
	if schemaDoc.Title == "" {
		t.Fatalf("expected schema title, got %#v", schemaDoc.Title)
	}
	if !strings.HasSuffix(filepath.ToSlash(schemaDoc.File), "schemas/platforms/douyin.video.schema.json") {
		t.Fatalf("expected douyin video schema path, got %s", schemaDoc.File)
	}
	if schemaDoc.RootSchema != "schemas/publish.schema.json" {
		t.Fatalf("expected publish root schema, got %s", schemaDoc.RootSchema)
	}
	if schemaDoc.Key != "douyin/video" {
		t.Fatalf("expected schema key, got %s", schemaDoc.Key)
	}
}

func TestListIncludesImageTextAsDisplayType(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	entries, err := validator.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Platform == "douyin" && entry.Type == "image-text" {
			return
		}
	}
	t.Fatalf("expected douyin image-text schema in list, got %d entries", len(entries))
}

func TestValidateImageTextUsesMappedSchemaFile(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "图文标题",
		"description": "图文内容",
		"images": []interface{}{
			map[string]interface{}{
				"key":    "image-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
			},
		},
		"extraField": true,
	}
	result := validator.Validate("抖音", "image-text", payload)
	if result.Valid {
		t.Fatal("expected mapped imageText schema to reject extra field")
	}
	if !containsError(result.Errors, `unexpected field "extraField"`) {
		t.Fatalf("expected extra field error from imageText schema, got %v", result.Errors)
	}
}

func TestValidateFullPayloadPrefixesAccountFormErrors(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_1",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"description": "缺少标题",
				},
			},
		},
	}
	result := validator.Validate("抖音", "video", payload)
	if result.Valid {
		t.Fatal("expected missing title error")
	}
	if !containsError(result.Errors, `accountForms[0].contentPublishForm: /: missing required field "title"`) {
		t.Fatalf("expected prefixed accountForms error, got %v", result.Errors)
	}
}

func TestValidateAcceptsStandardPublishRequestBusinessFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"video": map[string]interface{}{
			"key":    "video-key",
			"size":   float64(100),
			"width":  float64(10),
			"height": float64(10),
		},
		"images": []interface{}{
			map[string]interface{}{
				"key":    "image-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
			},
		},
		"cover": map[string]interface{}{
			"key":    "cover-key",
			"size":   float64(100),
			"width":  float64(10),
			"height": float64(10),
		},
		"coverKey": "cover-key",
		"content":  "正文",
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_1",
				"mediaId":           "media_1",
				"platformName":      "抖音",
				"publishContentId":  "content_1",
				"fps":               float64(0),
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "标题",
					"description": "描述",
					"scheduledTime": float64(1760000000000),
				},
			},
		},
		"isAppContent": false,
	}

	result := validator.Validate("抖音", "video", payload)
	if !result.Valid {
		t.Fatalf("expected standard business fields to be allowed, got %v", result.Errors)
	}
}

func TestValidateMissingSchemaFallsBackToBasicValidation(t *testing.T) {
	validator := NewValidator(t.TempDir())
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId":  "acc_1",
				"contentPublishForm": map[string]interface{}{"title": "ok"},
			},
		},
	}
	result := validator.Validate("unknown", "video", payload)
	if !result.Valid {
		t.Fatalf("expected fallback basic validation to pass, got %v", result.Errors)
	}

	delete(payload["accountForms"].([]interface{})[0].(map[string]interface{}), "contentPublishForm")
	result = validator.Validate("unknown", "video", payload)
	if result.Valid {
		t.Fatal("expected fallback basic validation to reject missing contentPublishForm")
	}
	if !containsError(result.Errors, "missing contentPublishForm") {
		t.Fatalf("expected missing contentPublishForm error, got %v", result.Errors)
	}
}

func readTestPayload(t *testing.T, name string) map[string]interface{} {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "payloads", name))
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimPrefix(string(raw), "\uFEFF")), &payload); err != nil {
		t.Fatal(err)
	}
	return payload
}

func containsError(errors []string, needle string) bool {
	for _, err := range errors {
		if strings.Contains(err, needle) {
			return true
		}
	}
	return false
}
