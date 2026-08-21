package publish

import (
	"testing"
	"time"
)

func TestNormalizeWeixinAccountImageTextDefaults(t *testing.T) {
	payload := standardPayload("imageText", []string{"WeiXinGongZhongHao"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_wx_1",
				"contentPublishForm": map[string]interface{}{
					"formType": "task",
					"title":    "公众号图文",
					"images":   []interface{}{uploadedResourceWithKey("wx-first-image")},
				},
			},
		},
	})

	NormalizeStandardPayloadForSchemaValidation("imageText", []string{"WeiXinGongZhongHao"}, payload)

	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	cpf := form["contentPublishForm"].(map[string]interface{})
	for field, want := range map[string]interface{}{
		"notifySubscribers": float64(0),
		"sex":               float64(0),
		"needOpenComment":   float64(0),
		"statement":         float64(0),
		"disableRecommend":  float64(0),
		"pubType":           float64(1),
	} {
		if got := cpf[field]; got != want {
			t.Fatalf("%s default = %#v, want %#v", field, got, want)
		}
	}
	if got := form["coverKey"]; got != "wx-first-image" {
		t.Fatalf("expected first image to become coverKey, got %#v", got)
	}
	cover := form["cover"].(map[string]interface{})
	if cover["key"] != "wx-first-image" {
		t.Fatalf("expected first image to become cover, got %#v", cover)
	}

	platformForms, ok := publishArgsOf(payload)["platformForms"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected WeChat imageText platformForms to be initialized, got %#v", publishArgsOf(payload)["platformForms"])
	}
	platformForm, ok := platformForms["微信公众号"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 微信公众号 platform form, got %#v", platformForms)
	}
	for field, want := range map[string]interface{}{
		"pubType":           float64(1),
		"notifySubscribers": float64(0),
		"sex":               float64(0),
		"title":             "",
		"desc":              "",
		"needOpenComment":   float64(0),
		"statement":         float64(0),
		"disableRecommend":  float64(0),
		"formType":          "task",
	} {
		if got := platformForm[field]; got != want {
			t.Fatalf("platform form %s default = %#v, want %#v", field, got, want)
		}
	}
	if images, ok := platformForm["images"].([]interface{}); !ok || len(images) != 0 {
		t.Fatalf("expected platform form images to default to an empty array, got %#v", platformForm["images"])
	}
}

func TestWeixinAccountImageTextScheduleRequiresTwoHours(t *testing.T) {
	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
	for _, test := range []struct {
		name      string
		value     interface{}
		wantError bool
	}{
		{name: "one hour", value: float64(now.Add(time.Hour).UnixMilli()), wantError: true},
		{name: "exactly two hours", value: float64(now.Add(2 * time.Hour).UnixMilli()), wantError: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			errors := []string{}
			validateWeixinAccountImageTextSchedule(test.value, now, "contentPublishForm.scheduledTime", &errors)
			if test.wantError && len(errors) == 0 {
				t.Fatalf("expected schedule validation error")
			}
			if !test.wantError && len(errors) > 0 {
				t.Fatalf("did not expect schedule validation error: %v", errors)
			}
		})
	}
}

func TestPreflightRequiresUploadedImagesForWeixinAccountImageText(t *testing.T) {
	result := Preflight("imageText", []string{"微信公众号"}, standardPayload("imageText", []string{"微信公众号"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_wx_1",
				"contentPublishForm": map[string]interface{}{
					"formType": "task",
					"title":    "公众号图文",
				},
			},
		},
	}))
	assertHasError(t, result.Errors, "imageText publish requires at least one uploaded image")
}

func TestPreflightAcceptsRecordedWeixinAccountImageTextShape(t *testing.T) {
	result := Preflight("imageText", []string{"微信公众号"}, standardPayload("imageText", []string{"微信公众号"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "6a75ae6089ecda991fd0d584",
				"images": []interface{}{
					uploadedResourceWithKey("yfb/local/image-1"),
				},
				"coverKey": "ml/local/cover-1",
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
	}))
	if len(result.Errors) > 0 {
		t.Fatalf("expected the recorded successful payload shape to pass preflight, got %v", result.Errors)
	}
}

func TestPreflightRejectsWeixinAccountImageTextScheduleUnderTwoHours(t *testing.T) {
	result := Preflight("imageText", []string{"微信公众号"}, standardPayload("imageText", []string{"微信公众号"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_wx_1",
				"contentPublishForm": map[string]interface{}{
					"formType":      "task",
					"scheduledTime": float64(time.Now().Add(time.Hour).UnixMilli()),
					"images":        []interface{}{uploadedResourceWithKey("wx-image")},
				},
			},
		},
	}))
	assertHasError(t, result.Errors, "scheduledTime must be at least 2 hours")
}
