package publish

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreflightRequiresStandardPayload(t *testing.T) {
	result := Preflight("video", []string{"抖音"}, map[string]interface{}{})
	assertHasError(t, result.Errors, "Standard publish payload is required")
}

func TestPreflightRejectsLegacyTopLevelAccountForms(t *testing.T) {
	result := Preflight("video", []string{"抖音"}, map[string]interface{}{
		"action":      "publish",
		"publishType": "video",
		"platforms":   []interface{}{"抖音"},
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
			},
		},
	})
	assertHasError(t, result.Errors, "Legacy publish payload is not supported")
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
		"covers":   []interface{}{uploadedResourceWithKey("cover-key")},
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
	cpf := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if _, ok := cpf["covers"].([]interface{}); !ok {
		t.Fatalf("expected publishArgs.covers to normalize into contentPublishForm.covers, got %+v", cpf)
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

func TestPreflightAcceptsBaijiahaoImageTextDraftFields(t *testing.T) {
	payload := standardPayload("imageText", []string{"百家号"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_bjh_1",
				"cover":             uploadedResourceWithKey("cover-key"),
				"coverKey":          "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":      "task",
					"title":         "百家号图文标题",
					"description":   "<p>百家号图文内容</p>",
					"pubType":       float64(0),
					"declaration":   float64(0),
					"scheduledTime": float64(1760000000000),
					"images":        []interface{}{uploadedResourceWithKey("image-key")},
				},
			},
		},
	})

	result := Preflight("imageText", []string{"百家号"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected baijiahao imageText preflight to pass, got %v", result.Errors)
	}
}

func TestPreflightAcceptsSouhuhaoVideoFields(t *testing.T) {
	payload := standardPayload("video", []string{"搜狐号"}, map[string]interface{}{
		"video": uploadedResource(),
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_sh_1",
				"cover":             uploadedResourceWithKey("cover-key"),
				"coverKey":          "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "搜狐号视频标题示例",
					"description": "这是搜狐号视频描述内容。",
					"tags":        []interface{}{"科技"},
					"declaration": float64(2),
					"pubType":     float64(1),
					"category": []interface{}{
						map[string]interface{}{
							"id":   "1",
							"text": "科技",
							"raw":  map[string]interface{}{"id": "1"},
						},
					},
				},
			},
		},
	})

	result := Preflight("video", []string{"搜狐号"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected souhuhao video preflight to pass, got %v", result.Errors)
	}
}

func TestPreflightAcceptsToutiaohaoArticleExtendedFields(t *testing.T) {
	payload := standardPayload("article", []string{"头条号"}, map[string]interface{}{
		"content": "<p>文章正文</p>",
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_tt_1",
				"cover":             uploadedResourceWithKey("cover-key"),
				"coverKey":          "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":      "task",
					"title":         "头条号文章标题",
					"pubType":       float64(0),
					"isFirst":       true,
					"advertisement": float64(3),
					"declaration":   float64(3),
					"scheduledTime": float64(1760000000000),
					"location": map[string]interface{}{
						"yixiaoerId":   "loc_1",
						"yixiaoerName": "上海",
						"raw":          map[string]interface{}{"id": "loc_1"},
					},
				},
			},
		},
	})

	result := Preflight("article", []string{"头条号"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected toutiaohao article preflight to pass, got %v", result.Errors)
	}
	cpf := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["scheduledTime"] != float64(1760000000) {
		t.Fatalf("expected scheduledTime normalized to seconds, got %#v", cpf["scheduledTime"])
	}
}

func TestResolveStandardPayloadResourceMetadataFillsImageDimensionsFromLocalSource(t *testing.T) {
	imagePath := writePNG(t, 8, 6)
	payload := standardPayload("imageText", []string{"抖音"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"cover": map[string]interface{}{
					"key":    "cover-key",
					"source": imagePath,
				},
				"coverKey": "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "图文",
					"description": "正文",
					"images": []interface{}{
						map[string]interface{}{
							"key":    "image-key",
							"source": imagePath,
						},
					},
				},
			},
		},
	})

	if err := ResolveStandardPayloadResourceMetadata(payload); err != nil {
		t.Fatal(err)
	}

	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	cover := form["cover"].(map[string]interface{})
	if cover["width"] != float64(8) || cover["height"] != float64(6) {
		t.Fatalf("expected cover metadata to be filled, got %+v", cover)
	}
	if _, exists := cover["source"]; exists {
		t.Fatalf("expected source helper field to be removed, got %+v", cover)
	}
	imageItem := form["contentPublishForm"].(map[string]interface{})["images"].([]interface{})[0].(map[string]interface{})
	if imageItem["width"] != float64(8) || imageItem["height"] != float64(6) {
		t.Fatalf("expected image metadata to be filled, got %+v", imageItem)
	}
}

func TestNormalizeTopicHTMLForVideoWithoutTopicsField(t *testing.T) {
	payload := standardPayload("video", []string{"抖音"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "视频",
					"description": "今日穿搭分享",
					"tags":        []interface{}{"穿搭", "#夏日"},
				},
			},
		},
	})

	NormalizeStandardPayloadWithTopicHTMLPolicy("video", []string{"抖音"}, payload, TopicHTMLPolicy{
		"抖音": TopicHTMLFields{HasDescription: true},
	})

	cpf := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	expected := `<p>今日穿搭分享</p><p><topic text="穿搭">#穿搭</topic><topic text="夏日">#夏日</topic></p>`
	if cpf["description"] != expected {
		t.Fatalf("expected video description topic HTML, got %#v", cpf["description"])
	}
	if publishArgsOf(payload)["content"] != expected {
		t.Fatalf("expected publishArgs.content topic HTML, got %#v", publishArgsOf(payload)["content"])
	}
	if cpf["content"] != expected {
		t.Fatalf("expected contentPublishForm.content topic HTML, got %#v", cpf["content"])
	}
}

func TestNormalizeTopicHTMLSkipsPlatformsWithoutDescriptionField(t *testing.T) {
	payload := standardPayload("article", []string{"CSDN"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"contentPublishForm": map[string]interface{}{
					"formType": "task",
					"title":    "文章",
					"content":  "今日穿搭分享",
					"tags":     []interface{}{"穿搭", "#夏日"},
				},
			},
		},
	})

	NormalizeStandardPayloadWithTopicHTMLPolicy("article", []string{"CSDN"}, payload, TopicHTMLPolicy{
		"CSDN": TopicHTMLFields{HasContent: true},
	})

	cpf := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	expected := `<p>今日穿搭分享</p><p><topic text="穿搭">#穿搭</topic><topic text="夏日">#夏日</topic></p>`
	if cpf["content"] != expected {
		t.Fatalf("expected content topic HTML when description is unsupported, got %#v", cpf["content"])
	}
	if publishArgsOf(payload)["content"] != expected {
		t.Fatalf("expected publishArgs.content topic HTML, got %#v", publishArgsOf(payload)["content"])
	}
	if _, exists := cpf["description"]; exists {
		t.Fatalf("expected description to remain absent when description is unsupported, got %+v", cpf)
	}
}

func TestNormalizeTopicHTMLSkipsSchemasWithTopicsField(t *testing.T) {
	payload := standardPayload("article", []string{"知乎"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "文章",
					"description": "今日穿搭分享",
					"tags":        []interface{}{"穿搭"},
				},
			},
		},
	})

	NormalizeStandardPayloadWithTopicHTMLPolicy("article", []string{"知乎"}, payload, TopicHTMLPolicy{
		"知乎": TopicHTMLFields{HasTopics: true, HasDescription: true, HasContent: true},
	})

	cpf := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["description"] != "今日穿搭分享" {
		t.Fatalf("expected description to remain unchanged for topics-capable schema, got %#v", cpf["description"])
	}
	if cpf["content"] != nil || publishArgsOf(payload)["content"] != nil {
		t.Fatalf("expected topics-capable schema not to synthesize content fields, got cpf=%+v publishArgs=%+v", cpf, publishArgsOf(payload))
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

func TestPreflightRejectsShipinhaoImageTextImageOver512KB(t *testing.T) {
	payload := standardPayload("imageText", []string{"视频号"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"cover":             uploadedResourceWithKey("cover-key"),
				"coverKey":          "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "图文",
					"description": "正文",
					"images": []interface{}{
						map[string]interface{}{
							"key":    "image-key",
							"size":   float64(512*1024 + 1),
							"width":  float64(1080),
							"height": float64(1440),
						},
					},
				},
			},
		},
	})

	result := Preflight("imageText", []string{"视频号"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].images[0].size: 视频号图片不能超过 512KB")
}

func TestPreflightRejectsShipinhaoCoverOver512KB(t *testing.T) {
	payload := standardPayload("video", []string{"视频号"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"cover": map[string]interface{}{
					"key":    "cover-key",
					"size":   float64(512*1024 + 1),
					"width":  float64(1080),
					"height": float64(1440),
				},
				"coverKey": "cover-key",
				"video":    uploadedResource(),
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "视频",
					"description": "正文",
				},
			},
		},
	})

	result := Preflight("video", []string{"视频号"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].cover.size: 视频号封面不能超过 512KB")
}

func TestPreflightDoesNotApplyShipinhao512KBRuleToOtherPlatforms(t *testing.T) {
	payload := standardPayload("imageText", []string{"抖音"}, map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"cover": map[string]interface{}{
					"key":    "cover-key",
					"size":   float64(512*1024 + 1),
					"width":  float64(1080),
					"height": float64(1440),
				},
				"coverKey": "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "图文",
					"description": "正文",
					"images": []interface{}{
						map[string]interface{}{
							"key":    "image-key",
							"size":   float64(512*1024 + 1),
							"width":  float64(1080),
							"height": float64(1440),
						},
					},
				},
			},
		},
	})

	result := Preflight("imageText", []string{"抖音"}, payload)
	for _, err := range result.Errors {
		if strings.Contains(err, "512KB") {
			t.Fatalf("expected non-shipinhao platform to skip 512KB limit, got %v", result.Errors)
		}
	}
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
	assertHasError(t, result.Errors, "publishArgs.content: article publish requires content")
}

func TestPreflightRejectsUnresolvedTemplatePlaceholders(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["platformAccountId"] = "<platformAccountId>"
	form["contentPublishForm"].(map[string]interface{})["title"] = "<title>"

	result := Preflight("video", []string{"抖音"}, payload)
	assertHasError(t, result.Errors, `accountForms[0].platformAccountId: unresolved template placeholder "<platformAccountId>"`)
	assertHasError(t, result.Errors, `accountForms[0].contentPublishForm.title: unresolved template placeholder "<title>"`)
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

func TestPreflightAcceptsXiaohongshuFlatShoppingCartRaw(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["shopping_cart"] = []interface{}{
		map[string]interface{}{
			"yixiaoerId":       "goods_001",
			"yixiaoerName":     "测试商品",
			"yixiaoerImageUrl": "https://example.invalid/goods.png",
			"yixiaoerDesc":     "--",
			"price":            float64(19900),
			"raw":              map[string]interface{}{"id": "goods_001"},
		},
	}

	result := Preflight("video", []string{"小红书"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected flat shopping_cart.raw to pass, got %v", result.Errors)
	}
}

func TestPreflightAcceptsCompleteMusicObjectWithPlayURLs(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["music"] = map[string]interface{}{
		"yixiaoerId":   "music_001",
		"yixiaoerName": "稻香",
		"duration":     float64(240),
		"url":          "https://example.invalid/music",
		"playUrl":      "https://example.invalid/preview.mp3",
		"raw": map[string]interface{}{
			"id": "music_001",
		},
	}

	result := Preflight("video", []string{"抖音"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected complete music object with metadata URLs to pass, got %v", result.Errors)
	}
}

func TestPreflightNormalizesLegacyDouyinShoppingCartShape(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["shoppingCart"] = []interface{}{
		map[string]interface{}{
			"sale_title":   "点击购买",
			"yixiaoerId":   "goods_001",
			"yixiaoerName": "测试商品",
			"raw": map[string]interface{}{
				"gid":        "goods_001",
				"goods_imgs": []interface{}{"https://example.invalid/goods.png"},
			},
		},
	}

	result := Preflight("video", []string{"抖音"}, payload)
	if len(result.Errors) > 0 {
		t.Fatalf("expected legacy douyin shoppingCart to normalize, got %v", result.Errors)
	}

	cpf := form["contentPublishForm"].(map[string]interface{})
	if _, exists := cpf["shoppingCart"]; exists {
		t.Fatalf("expected legacy shoppingCart key to be removed, got %+v", cpf)
	}
	items := cpf["shopping_cart"].([]interface{})
	item := items[0].(map[string]interface{})
	if len(item["images"].([]interface{})) != 1 {
		t.Fatalf("expected images to be derived from raw, got %+v", item)
	}
	data := item["data"].(map[string]interface{})
	if data["yixiaoerId"] != "goods_001" || data["yixiaoerName"] != "测试商品" {
		t.Fatalf("expected nested shopping cart data, got %+v", item)
	}
}

func TestPreflightRejectsXiaohongshuShoppingCartMissingRaw(t *testing.T) {
	payload := validVideoPayload()
	form := publishArgsOf(payload)["accountForms"].([]interface{})[0].(map[string]interface{})
	form["contentPublishForm"].(map[string]interface{})["shopping_cart"] = []interface{}{
		map[string]interface{}{
			"yixiaoerId":   "goods_001",
			"yixiaoerName": "测试商品",
		},
	}

	result := Preflight("video", []string{"小红书"}, payload)
	assertHasError(t, result.Errors, "accountForms[0].contentPublishForm.shopping_cart[0]: dynamic platform object must include complete \"raw\" data")
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

func writePNG(t *testing.T, width, height int) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{R: 20, G: 120, B: 200, A: 255})
		}
	}
	path := filepath.Join(t.TempDir(), "resource.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
	return path
}
