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
			t.Fatalf("account setting %s default = %#v, want %#v", field, got, want)
		}
	}
	if got := form["coverKey"]; got != "wx-first-image" {
		t.Fatalf("expected first image to become coverKey, got %#v", got)
	}
	cover := form["cover"].(map[string]interface{})
	if cover["key"] != "wx-first-image" {
		t.Fatalf("expected first image to become cover, got %#v", cover)
	}

	if _, exists := publishArgsOf(payload)["platformForms"]; exists {
		t.Fatalf("did not expect a 微信公众号 imageText platform form: %#v", publishArgsOf(payload))
	}
}

func TestNormalizeWeixinAccountImageTextConvertsLegacyStatementObject(t *testing.T) {
	payload := standardPayload("imageText", []string{"微信公众号"}, map[string]interface{}{
		"accountForms": []interface{}{map[string]interface{}{
			"platformAccountId": "acc_wx_1",
			"contentPublishForm": map[string]interface{}{
				"formType":          "task",
				"statement":         map[string]interface{}{"type": float64(4)},
				"notifySubscribers": float64(1),
				"needOpenComment":   float64(3),
				"disableRecommend":  float64(0),
				"pubType":           float64(1),
			},
		}},
	})

	NormalizeStandardPayloadForSchemaValidation("imageText", []string{"微信公众号"}, payload)
	args := publishArgsOf(payload)
	contentForm := args["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	for field, want := range map[string]interface{}{
		"statement":         float64(4),
		"notifySubscribers": float64(1),
		"needOpenComment":   float64(3),
		"disableRecommend":  float64(0),
		"pubType":           float64(1),
	} {
		if got := contentForm[field]; got != want {
			t.Fatalf("content form %s = %#v, want %#v", field, got, want)
		}
	}
}

func TestNormalizeWeixinAccountImageTextMigratesLegacyPlatformFormToAccountForm(t *testing.T) {
	payload := standardPayload("imageText", []string{"微信公众号"}, map[string]interface{}{
		"accountForms": []interface{}{map[string]interface{}{
			"platformAccountId": "acc_wx_1",
			"contentPublishForm": map[string]interface{}{
				"formType":  "task",
				"title":     "公众号图文",
				"statement": float64(4),
			},
		}},
		"platformForms": map[string]interface{}{
			"微信公众号": map[string]interface{}{
				"statement":         float64(4),
				"notifySubscribers": float64(0),
				"sex":               float64(0),
				"needOpenComment":   float64(3),
				"disableRecommend":  float64(0),
				"pubType":           float64(1),
			},
		},
	})
	NormalizeStandardPayloadForSchemaValidation("imageText", []string{"微信公众号"}, payload)

	args := publishArgsOf(payload)
	contentForm := args["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	for field, want := range map[string]interface{}{
		"notifySubscribers": float64(0),
		"sex":               float64(0),
		"needOpenComment":   float64(3),
		"disableRecommend":  float64(0),
		"pubType":           float64(1),
	} {
		if got := contentForm[field]; got != want {
			t.Fatalf("contentPublishForm.%s = %#v, want %#v", field, got, want)
		}
	}
	if got := contentForm["statement"]; got != float64(4) {
		t.Fatalf("contentPublishForm.statement = %#v, want 4", got)
	}
	if _, exists := args["platformForms"]; exists {
		t.Fatalf("legacy platform form must not be sent for imageText: %#v", args)
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
				"notifySubscribers": float64(1),
				"sex":               float64(0),
				"needOpenComment":   float64(0),
				"statement":         float64(0),
				"disableRecommend":  float64(0),
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

func TestPreflightRejectsInvalidWeixinAccountImageTextSettings(t *testing.T) {
	result := Preflight("imageText", []string{"微信公众号"}, standardPayload("imageText", []string{"微信公众号"}, map[string]interface{}{
		"accountForms": []interface{}{map[string]interface{}{
			"platformAccountId": "acc_wx_1",
			"images":            []interface{}{uploadedResourceWithKey("wx-image")},
			"contentPublishForm": map[string]interface{}{
				"formType":        "task",
				"statement":       float64(2),
				"needOpenComment": float64(4),
			},
		}},
	}))
	assertHasError(t, result.Errors, `accountForms[0].contentPublishForm.statement: must be one of [0 1 3 4 5 6]`)
	assertHasError(t, result.Errors, `accountForms[0].contentPublishForm.needOpenComment: must be one of [0 1 2 3]`)
}
