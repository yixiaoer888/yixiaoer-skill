package publish

import (
	"strings"
	"testing"
)

func TestPreflightRequiresStandardPayload(t *testing.T) {
	result := Preflight("video", []string{"抖音"}, map[string]interface{}{})
	assertHasError(t, result.Errors, "Standard publish payload is required")
}

func TestPreflightAcceptsValidStandardVideoPayload(t *testing.T) {
	payload := validVideoPayload()
	result := Preflight("video", []string{"抖音"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected valid preflight, got %v", result.Errors)
	}
	if len(result.AccountIDs) != 1 || result.AccountIDs[0] != "acc_001" {
		t.Fatalf("unexpected account ids: %v", result.AccountIDs)
	}
}

func TestPreflightValidatesFullPublishRequestFields(t *testing.T) {
	payload := map[string]interface{}{
		"action":         "save",
		"publishType":    "imageText",
		"platforms":      []interface{}{},
		"publishChannel": "local",
		"publishArgs":    validPublishArgs(),
	}

	result := Preflight("video", []string{"抖音"}, payload)
	assertHasError(t, result.Errors, `action: must equal "publish"`)
	assertHasError(t, result.Errors, `publishType: got "imageText", expected "video"`)
	assertHasError(t, result.Errors, "platforms: must be a non-empty array")
	assertHasError(t, result.Errors, "clientId: required when publishChannel is local")
}

func TestPreflightRejectsMissingAccountIDAndContentForm(t *testing.T) {
	payload := standardPayload("video", []string{"抖音"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"video": uploadedResource(),
			},
		},
	})
	result := Preflight("video", []string{"抖音"}, payload)
	assertHasError(t, result.Errors, "accountForms[0]: missing platformAccountId")
	assertHasError(t, result.Errors, "accountForms[0]: missing contentPublishForm")
}

func TestPreflightAcceptsAccountIDAlias(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	delete(form, "platformAccountId")
	form["account_id"] = "acc_alias"

	result := Preflight("video", []string{"抖音"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected account_id alias to pass, got %v", result.Errors)
	}
	if got := result.AccountIDs[0]; got != "acc_alias" {
		t.Fatalf("unexpected account id: %s", got)
	}
}

func TestPreflightAcceptsSharedResourcesUnderPublishArgs(t *testing.T) {
	payload := standardPayload("video", []string{"抖音"}, map[string]interface{}{
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
	})

	result := Preflight("video", []string{"抖音"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected shared resources to normalize, got %v", result.Errors)
	}

	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["video"] == nil || form["cover"] == nil || form["coverKey"] != "cover-key" {
		t.Fatalf("expected shared resource fields on account form, got %+v", form)
	}
}

func TestPreflightAcceptsArticleContentFromPublishArgs(t *testing.T) {
	payload := standardPayload("article", []string{"知乎"}, map[string]interface{}{
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
	})

	result := Preflight("article", []string{"知乎"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected article content normalization to pass, got %v", result.Errors)
	}
}

func TestPreflightAcceptsImageTextImagesInContentPublishForm(t *testing.T) {
	payload := standardPayload("imageText", []string{"抖音"}, map[string]interface{}{
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
	})
	result := Preflight("imageText", []string{"抖音"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected imageText preflight to pass, got %v", result.Errors)
	}

	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	if images, _ := form["images"].([]interface{}); len(images) != 1 {
		t.Fatalf("expected contentPublishForm.images to normalize into account form, got %+v", form)
	}
}

func TestPreflightRejectsImageTextMissingImageKey(t *testing.T) {
	payload := standardPayload("imageText", []string{"抖音"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId":  "acc_001",
				"cover":              uploadedResourceWithKey("cover-key"),
				"coverKey":           "cover-key",
				"images":             []interface{}{map[string]interface{}{"size": float64(10)}},
				"contentPublishForm": map[string]interface{}{"formType": "task", "title": "图文", "description": "正文"},
			},
		},
	})
	result := Preflight("imageText", []string{"抖音"}, payload)
	assertHasError(t, result.Errors, `accountForms[0].images[0]: missing uploaded resource field "key"`)
}

func TestPreflightRejectsArticleMissingContent(t *testing.T) {
	payload := standardPayload("article", []string{"知乎"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId":  "acc_001",
				"cover":              uploadedResourceWithKey("cover-key"),
				"coverKey":           "cover-key",
				"contentPublishForm": map[string]interface{}{"formType": "task", "title": "文章"},
			},
		},
	})
	result := Preflight("article", []string{"知乎"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].contentPublishForm.content: article publish requires content")
}

func TestPreflightNormalizesScheduledTimeFromMilliseconds(t *testing.T) {
	payload := validVideoPayload()
	cpf := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	cpf["scheduledTime"] = float64(1760000000000)

	result := Preflight("video", []string{"抖音"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected scheduledTime normalization to pass, got %v", result.Errors)
	}
	if got := cpf["scheduledTime"]; got != float64(1760000000) {
		t.Fatalf("expected scheduledTime to normalize to seconds, got %#v", got)
	}
}

func TestPreflightRejectsMiniappsMissingRaw(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["miniapps"] = []interface{}{
		map[string]interface{}{"id": "mini_1", "name": "小程序"},
	}

	result := Preflight("video", []string{"抖音"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].contentPublishForm.miniapps[0]: dynamic platform object must include complete \"raw\" data")
}

func TestPreflightAcceptsNestedShoppingCartDataRaw(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["shopping_cart"] = []interface{}{
		map[string]interface{}{
			"sale_title": "同款商品",
			"images":     []interface{}{"goods-cover-key"},
			"data": map[string]interface{}{
				"yixiaoerId":   "goods_001",
				"yixiaoerName": "测试商品",
				"raw":          map[string]interface{}{"id": "goods_001"},
			},
		},
	}

	result := Preflight("video", []string{"小红书"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected nested shopping_cart.data.raw to pass, got %v", result.Errors)
	}
}

func TestPreflightRejectsNestedShoppingCartDataMissingRaw(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["shopping_cart"] = []interface{}{
		map[string]interface{}{
			"sale_title": "同款商品",
			"data": map[string]interface{}{
				"yixiaoerId":   "goods_001",
				"yixiaoerName": "测试商品",
			},
		},
	}

	result := Preflight("video", []string{"小红书"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].contentPublishForm.shopping_cart[0].data: dynamic platform object must include complete \"raw\" data")
}

func validVideoPayload() map[string]interface{} {
	return standardPayload("video", []string{"抖音"}, map[string]interface{}{
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
	})
}

func validPublishArgs() map[string]interface{} {
	return publishArgsOf(validVideoPayload())
}

func publishArgsOf(payload map[string]interface{}) map[string]interface{} {
	return payload["publishArgs"].(map[string]interface{})
}

func standardPayload(publishType string, platforms []string, publishArgs map[string]interface{}) map[string]interface{} {
	rawPlatforms := make([]interface{}, 0, len(platforms))
	for _, platform := range platforms {
		rawPlatforms = append(rawPlatforms, platform)
	}
	return map[string]interface{}{
		"action":         "publish",
		"publishType":    publishType,
		"platforms":      rawPlatforms,
		"publishChannel": "cloud",
		"publishArgs":    publishArgs,
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
