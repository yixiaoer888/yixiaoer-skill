package publish

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreflightAcceptsValidVideo(t *testing.T) {
	inner := readPayload(t, "douyin-video-valid.json")
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId":  "acc_001",
				"cover":              uploadedResourceWithKey("cover-key"),
				"coverKey":           "cover-key",
				"video":              inner["video"],
				"contentPublishForm": inner,
			},
		},
	}
	result := Preflight("video", []string{"douyin"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected valid preflight, got %v", result.Errors)
	}
	if len(result.AccountIDs) != 1 {
		t.Fatalf("expected one account id, got %v", result.AccountIDs)
	}
}

func TestPreflightRejectsExternalURL(t *testing.T) {
	payload := readPayload(t, "douyin-video-url-invalid.json")
	result := Preflight("video", []string{"douyin"}, payload)
	if len(result.Errors) == 0 {
		t.Fatal("expected external URL preflight error")
	}
}

func TestPreflightRejectsMissingRaw(t *testing.T) {
	payload := readPayload(t, "douyin-video-location-missing-raw.json")
	result := Preflight("video", []string{"douyin"}, payload)
	if len(result.Errors) == 0 {
		t.Fatal("expected missing raw preflight error")
	}
}

func TestPreflightRejectsUnsupportedTypeAndMissingPlatforms(t *testing.T) {
	payload := validVideoPayload()
	result := Preflight("audio", nil, payload)
	assertHasError(t, result.Errors, `publish type "audio" is not supported`)
	assertHasError(t, result.Errors, "at least one target platform is required")
}

func TestPreflightRejectsMissingAccountForms(t *testing.T) {
	result := Preflight("video", []string{"douyin"}, map[string]interface{}{})
	assertHasError(t, result.Errors, "payload.accountForms must be a non-empty array")
}

func TestPreflightValidatesFullPublishRequestFields(t *testing.T) {
	payload := map[string]interface{}{
		"action":         "save",
		"publishType":    "imageText",
		"platforms":      []interface{}{},
		"publishChannel": "local",
		"publishArgs":    validVideoPayload(),
	}

	result := Preflight("video", []string{"douyin"}, payload)
	assertHasError(t, result.Errors, `action: must equal "publish"`)
	assertHasError(t, result.Errors, `publishType: got "imageText", expected "video"`)
	assertHasError(t, result.Errors, "platforms: must be a non-empty array")
	assertHasError(t, result.Errors, "clientId: required when publishChannel is local")
}

func TestPreflightValidatesTopLevelCoverInFullPublishRequest(t *testing.T) {
	payload := map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
		"coverKey":       "outer-cover-key",
		"cover":          uploadedResourceWithKey("wrong-cover-key"),
		"publishChannel": "cloud",
		"publishArgs":    validVideoPayload(),
	}

	result := Preflight("video", []string{"douyin"}, payload)
	assertHasError(t, result.Errors, "coverKey: must match cover.key")
}

func TestPreflightRejectsMissingAccountIDAndContentForm(t *testing.T) {
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"video": uploadedResource(),
			},
		},
	}
	result := Preflight("video", []string{"douyin"}, payload)
	assertHasError(t, result.Errors, "accountForms[0]: missing platformAccountId")
	assertHasError(t, result.Errors, "accountForms[0]: missing contentPublishForm")
}

func TestPreflightAcceptsAccountIDAlias(t *testing.T) {
	payload := validVideoPayload()
	form := payload["accountForms"].([]interface{})[0].(map[string]interface{})
	delete(form, "platformAccountId")
	form["account_id"] = "acc_alias"

	result := Preflight("video", []string{"douyin"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected account_id alias to pass, got %v", result.Errors)
	}
	if got := result.AccountIDs[0]; got != "acc_alias" {
		t.Fatalf("unexpected account id: %s", got)
	}
}

func TestPreflightRejectsVideoMissingResourceKey(t *testing.T) {
	payload := validVideoPayload()
	form := payload["accountForms"].([]interface{})[0].(map[string]interface{})
	form["video"] = map[string]interface{}{"size": float64(1024), "width": float64(1080)}

	result := Preflight("video", []string{"douyin"}, payload)
	assertHasError(t, result.Errors, `accountForms[0].video: missing uploaded resource field "key"`)
}

func TestPreflightAcceptsStandardPublishArgsSharedResources(t *testing.T) {
	payload := map[string]interface{}{
		"video":    uploadedResource(),
		"images":   []interface{}{uploadedResourceWithKey("image-key")},
		"cover":    uploadedResourceWithKey("cover-key"),
		"coverKey": "cover-key",
		"content":  "文章正文",
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "视频",
					"description": "描述",
				},
			},
		},
	}

	result := Preflight("video", []string{"douyin"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected shared standard resources to be normalized, got %v", result.Errors)
	}

	form := payload["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["video"] == nil || form["cover"] == nil || form["coverKey"] != "cover-key" {
		t.Fatalf("expected shared resource fields to copy into account form, got %+v", form)
	}
}

func TestPreflightAcceptsStandardPublishArgsArticleContent(t *testing.T) {
	payload := map[string]interface{}{
		"cover":    uploadedResourceWithKey("cover-key"),
		"coverKey": "cover-key",
		"content":  "文章正文",
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"contentPublishForm": map[string]interface{}{
					"formType": "task",
					"title":    "文章标题",
				},
			},
		},
	}

	result := Preflight("article", []string{"zhihu"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected top-level content to normalize into contentPublishForm, got %v", result.Errors)
	}
}

func TestPreflightAcceptsImageTextImagesInContentPublishForm(t *testing.T) {
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"cover":             uploadedResourceWithKey("cover-key"),
				"coverKey":          "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "图文",
					"description": "正文",
					"images":      []interface{}{uploadedResource()},
				},
			},
		},
	}
	result := Preflight("image-text", []string{"douyin"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected image-text preflight to pass, got %v", result.Errors)
	}
}

func TestPreflightRejectsImageTextMissingImageKey(t *testing.T) {
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId":  "acc_001",
				"cover":              uploadedResourceWithKey("cover-key"),
				"coverKey":           "cover-key",
				"images":             []interface{}{map[string]interface{}{"size": float64(10)}},
				"contentPublishForm": map[string]interface{}{"formType": "task", "title": "图文", "description": "正文"},
			},
		},
	}
	result := Preflight("image-text", []string{"douyin"}, payload)
	assertHasError(t, result.Errors, `accountForms[0].images[0]: missing uploaded resource field "key"`)
}

func TestPreflightRejectsArticleMissingContent(t *testing.T) {
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId":  "acc_001",
				"cover":              uploadedResourceWithKey("cover-key"),
				"coverKey":           "cover-key",
				"contentPublishForm": map[string]interface{}{"formType": "task", "title": "文章"},
			},
		},
	}
	result := Preflight("article", []string{"zhihu"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].contentPublishForm.content: article publish requires content")
}

func TestPreflightAcceptsArticleContent(t *testing.T) {
	payload := map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId":  "acc_001",
				"cover":              uploadedResourceWithKey("cover-key"),
				"coverKey":           "cover-key",
				"contentPublishForm": map[string]interface{}{"formType": "task", "title": "文章", "content": "正文"},
			},
		},
	}
	result := Preflight("article", []string{"zhihu"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected article preflight to pass, got %v", result.Errors)
	}
}

func TestPreflightAcceptsDynamicObjectWithRaw(t *testing.T) {
	payload := validVideoPayload()
	form := payload["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["location"] = map[string]interface{}{
		"yixiaoerId":   "poi_1",
		"yixiaoerName": "位置",
		"raw":          map[string]interface{}{"id": "poi_1"},
	}

	result := Preflight("video", []string{"douyin"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected dynamic raw object to pass, got %v", result.Errors)
	}
}

func TestPreflightNormalizesScheduledTimeFromMilliseconds(t *testing.T) {
	payload := validVideoPayload()
	cpf := payload["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	cpf["scheduledTime"] = float64(1760000000000)

	result := Preflight("video", []string{"douyin"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected scheduledTime normalization to pass, got %v", result.Errors)
	}
	if got := cpf["scheduledTime"]; got != float64(1760000000) {
		t.Fatalf("expected scheduledTime to normalize to seconds, got %#v", got)
	}
}

func TestPreflightRejectsSecondBasedScheduledTime(t *testing.T) {
	payload := validVideoPayload()
	cpf := payload["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	cpf["scheduledTime"] = float64(1760000000)

	result := Preflight("video", []string{"douyin"}, payload)
	assertHasError(t, result.Errors, "scheduledTime: must be a 13-digit Unix timestamp in milliseconds")
}

func TestPreflightRejectsMiniappsMissingRaw(t *testing.T) {
	payload := validVideoPayload()
	form := payload["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["miniapps"] = []interface{}{
		map[string]interface{}{"id": "mini_1", "name": "小程序"},
	}

	result := Preflight("video", []string{"douyin"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].contentPublishForm.miniapps[0]: dynamic platform object must include complete \"raw\" data")
}

func readPayload(t *testing.T, name string) map[string]interface{} {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "tests", "fixtures", "payloads", name))
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimPrefix(string(raw), "\uFEFF")), &payload); err != nil {
		t.Fatal(err)
	}
	return payload
}

func validVideoPayload() map[string]interface{} {
	return map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"cover":             uploadedResourceWithKey("cover-key"),
				"coverKey":          "cover-key",
				"video":             uploadedResource(),
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "视频",
					"description": "描述",
				},
			},
		},
	}
}

func uploadedResource() map[string]interface{} {
	return uploadedResourceWithKey("resource-key")
}

func uploadedResourceWithKey(key string) map[string]interface{} {
	return map[string]interface{}{
		"key":    key,
		"size":   float64(1024),
		"width":  float64(1080),
		"height": float64(1920),
	}
}

func assertHasError(t *testing.T, errors []string, want string) {
	t.Helper()
	for _, err := range errors {
		if strings.Contains(err, want) {
			return
		}
	}
	t.Fatalf("expected error containing %q, got %v", want, errors)
}
