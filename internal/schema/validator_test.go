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
	if got := TypeKey("imageText"); got != "imageText" {
		t.Fatalf("expected imageText, got %s", got)
	}
}

func TestDisplayTypeMapsImageText(t *testing.T) {
	if got := DisplayType("imageText"); got != "imageText" {
		t.Fatalf("expected imageText, got %s", got)
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
		if entry.Platform == "douyin" && entry.Type == "imageText" {
			return
		}
	}
	t.Fatalf("expected douyin imageText schema in list, got %d entries", len(entries))
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
	result := validator.Validate("抖音", "imageText", payload)
	if result.Valid {
		t.Fatal("expected mapped imageText schema to reject extra field")
	}
	if !containsError(result.Errors, `unexpected field "extraField"`) {
		t.Fatalf("expected extra field error from imageText schema, got %v", result.Errors)
	}
}

func TestSchemaResolvesVideoAccountAliasesToCanonicalKeys(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))

	tests := []struct {
		platform string
		path     string
	}{
		{platform: "快手", path: "schemas/platforms/kuaishou.video.schema.json"},
		{platform: "视频号", path: "schemas/platforms/shipinhao.video.schema.json"},
		{platform: "微信视频号", path: "schemas/platforms/shipinhao.video.schema.json"},
		{platform: "shipinghao", path: "schemas/platforms/shipinhao.video.schema.json"},
		{platform: "shipinhao", path: "schemas/platforms/shipinhao.video.schema.json"},
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			schemaDoc, err := validator.Schema(tt.platform, "video")
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasSuffix(filepath.ToSlash(schemaDoc.File), tt.path) {
				t.Fatalf("expected schema path %s, got %s", tt.path, schemaDoc.File)
			}
		})
	}
}

func TestSchemaReturnsValidShipinhaoVideoSchema(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	schemaDoc, err := validator.Schema("视频号", "video")
	if err != nil {
		t.Fatal(err)
	}
	if schemaDoc.Title == "" {
		t.Fatal("expected schema title for shipinhao video")
	}
	if _, ok := schemaDoc.Properties["createType"]; !ok {
		t.Fatalf("expected shipinhao video schema to expose createType, got %+v", schemaDoc.Properties)
	}
	if _, ok := schemaDoc.Properties["pubType"]; !ok {
		t.Fatalf("expected shipinhao video schema to expose pubType, got %+v", schemaDoc.Properties)
	}
}

func TestSchemaResolvesShipinhaoImageTextWithoutLegacyAlias(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	for _, platform := range []string{"视频号", "shipinhao", "shipinghao"} {
		schemaDoc, err := validator.Schema(platform, "imageText")
		if err != nil {
			t.Fatalf("%s: %v", platform, err)
		}
		if !strings.HasSuffix(filepath.ToSlash(schemaDoc.File), "schemas/platforms/shipinhao.imageText.schema.json") {
			t.Fatalf("%s: expected shipinhao imageText schema path, got %s", platform, schemaDoc.File)
		}
	}
}

func TestSchemaResolvesBaijiahaoImageTextSchema(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	schemaDoc, err := validator.Schema("百家号", "imageText")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(filepath.ToSlash(schemaDoc.File), "schemas/platforms/baijiahao.imageText.schema.json") {
		t.Fatalf("expected baijiahao imageText schema path, got %s", schemaDoc.File)
	}
}

func TestSchemaResolvesSouhuhaoVideoSchema(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	schemaDoc, err := validator.Schema("搜狐号", "video")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(filepath.ToSlash(schemaDoc.File), "schemas/platforms/souhuhao.video.schema.json") {
		t.Fatalf("expected souhuhao video schema path, got %s", schemaDoc.File)
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
		"covers": []interface{}{
			map[string]interface{}{
				"key":    "cover-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
			},
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

func TestValidateAcceptsXiaohongshuFlatShoppingCartStructure(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"visibleType": float64(0),
		"shopping_cart": []interface{}{
			map[string]interface{}{
				"yixiaoerId":       "goods_001",
				"yixiaoerName":     "测试商品",
				"yixiaoerImageUrl": "https://example.invalid/goods.png",
				"yixiaoerDesc":     "--",
				"price":            float64(19900),
				"raw":              map[string]interface{}{"id": "goods_001"},
			},
		},
	}

	result := validator.Validate("小红书", "video", payload)
	if !result.Valid {
		t.Fatalf("expected xiaohongshu flat shopping_cart structure to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsDouyinNestedShoppingCartStructure(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "带货视频",
		"description": "视频描述",
		"shopping_cart": []interface{}{
			map[string]interface{}{
				"sale_title": "点击购买",
				"images":     []interface{}{"https://example.invalid/goods.png"},
				"data": map[string]interface{}{
					"yixiaoerId":   "goods_001",
					"yixiaoerName": "测试商品",
					"raw": map[string]interface{}{
						"gid": "goods_001",
					},
				},
			},
		},
	}

	result := validator.Validate("抖音", "video", payload)
	if !result.Valid {
		t.Fatalf("expected douyin nested shopping_cart structure to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsBaijiahaoCategoryPathArray(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType": "task",
		"title":    "百家号文章",
		"content":  "<p>正文</p>",
		"pubType":  float64(1),
		"covers": []interface{}{
			map[string]interface{}{
				"key":    "cover-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
			},
		},
		"category": []interface{}{
			map[string]interface{}{
				"yixiaoerId":   "32",
				"yixiaoerName": "财经",
			},
			map[string]interface{}{
				"yixiaoerId":   "9",
				"yixiaoerName": "财经综合",
			},
		},
	}

	result := validator.Validate("百家号", "article", payload)
	if !result.Valid {
		t.Fatalf("expected baijiahao category path array to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsBaijiahaoArticleDraftPubType(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType": "task",
		"title":    "百家号草稿文章",
		"content":  "<p>正文</p>",
		"pubType":  float64(0),
		"covers": []interface{}{
			map[string]interface{}{
				"key":    "cover-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
			},
		},
	}

	result := validator.Validate("百家号", "article", payload)
	if !result.Valid {
		t.Fatalf("expected baijiahao article draft pubType payload to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsBaijiahaoArticleActivityAndScheduledFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType": "task",
		"title":    "百家号征文文章",
		"content":  "<p>正文</p>",
		"pubType":  float64(1),
		"declaration": float64(1),
		"scheduledTime": float64(1760000000000),
		"activity": map[string]interface{}{
			"yixiaoerId":   "activity_1",
			"yixiaoerName": "征文活动",
			"raw":          map[string]interface{}{"id": "activity_1"},
		},
		"covers": []interface{}{
			map[string]interface{}{
				"key":    "cover-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
			},
		},
	}

	result := validator.Validate("百家号", "article", payload)
	if !result.Valid {
		t.Fatalf("expected baijiahao article activity/scheduled payload to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsArticleContentUnderPublishArgs(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"publishArgs": map[string]interface{}{
			"content": "<p>文章正文</p>",
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_1",
					"contentPublishForm": map[string]interface{}{
						"formType": "task",
						"title":    "知乎文章标题示例一",
						"pubType":  float64(1),
					},
				},
			},
		},
	}

	result := validator.Validate("知乎", "article", payload)
	if !result.Valid {
		t.Fatalf("expected publishArgs.content article payload to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsWeixinAccountArticlePlatformForms(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_weixin_1",
					"platformName":      "微信公众号",
				},
			},
			"platformForms": map[string]interface{}{
				"微信公众号": map[string]interface{}{
					"articles": []interface{}{
						map[string]interface{}{
							"title":   "公众号文章标题",
							"content": "<p>公众号文章正文</p>",
							"type":    float64(1),
							"cover": map[string]interface{}{
								"key": "wx-cover-key",
								"raw": map[string]interface{}{},
							},
						},
					},
					"notifySubscribers": float64(1),
					"pubType":           float64(1),
				},
			},
		},
	}

	result := validator.Validate("微信公众号", "article", payload)
	if !result.Valid {
		t.Fatalf("expected weixin account article payload to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsXhsImageTextMusicAndScheduledFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"description": "<p>小红书图文内容</p>",
		"visibleType": float64(0),
		"scheduledTime": float64(1760000000000),
		"music": map[string]interface{}{
			"yixiaoerId":   "music_1",
			"yixiaoerName": "背景音乐",
			"duration":     float64(30),
			"playUrl":      "https://example.invalid/music.mp3",
			"raw":          map[string]interface{}{"id": "music_1"},
		},
		"images": []interface{}{
			map[string]interface{}{
				"key":    "image-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
				"format": "jpg",
			},
		},
	}

	result := validator.Validate("小红书", "imageText", payload)
	if !result.Valid {
		t.Fatalf("expected xiaohongshu imageText music/scheduled payload to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsBaijiahaoImageTextPayload(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "百家号图文标题",
		"description": "<p>百家号图文内容</p>",
		"pubType":     float64(0),
		"declaration": float64(0),
		"scheduledTime": float64(1760000000000),
		"images": []interface{}{
			map[string]interface{}{
				"key":    "image-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
				"format": "jpg",
			},
		},
		"location": map[string]interface{}{
			"yixiaoerId":   "loc_1",
			"yixiaoerName": "北京",
			"raw":          map[string]interface{}{"id": "loc_1"},
		},
	}

	result := validator.Validate("百家号", "imageText", payload)
	if !result.Valid {
		t.Fatalf("expected baijiahao imageText payload to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsSouhuhaoVideoPayload(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "搜狐号视频标题示例",
		"description": "这是搜狐号视频描述内容。",
		"declaration": float64(2),
		"pubType":     float64(1),
		"tags":        []interface{}{"科技"},
		"category": []interface{}{
			map[string]interface{}{
				"id":   "1",
				"text": "科技",
				"raw":  map[string]interface{}{"id": "1"},
			},
		},
	}

	result := validator.Validate("搜狐号", "video", payload)
	if !result.Valid {
		t.Fatalf("expected souhuhao video payload to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsToutiaohaoArticleExtendedFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":      "task",
		"title":         "头条号文章标题",
		"content":       "<p>文章正文</p>",
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
		"covers": []interface{}{
			map[string]interface{}{
				"key":    "cover-key",
				"size":   float64(100),
				"width":  float64(10),
				"height": float64(10),
			},
		},
	}

	result := validator.Validate("头条号", "article", payload)
	if !result.Valid {
		t.Fatalf("expected toutiaohao article extended fields payload to pass, got %v", result.Errors)
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
