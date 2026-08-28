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

func TestValidateAcceptsYidianhaoVideoUsingCurrentCreateType(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	doc, err := validator.Schema("一点号", "video")
	if err != nil {
		t.Fatal(err)
	}
	for _, requiredField := range doc.Required {
		if requiredField == "type" {
			t.Fatalf("yidianhao video schema must not require removed type field: %#v", doc.Required)
		}
	}
	if _, ok := doc.Properties["type"]; ok {
		t.Fatalf("yidianhao video schema must not expose legacy type field: %#v", doc.Properties)
	}
	createType, ok := doc.Properties["createType"]
	if !ok {
		t.Fatalf("yidianhao video schema must expose createType, got %#v", doc.Properties)
	}
	if createType.Required {
		t.Fatalf("createType should remain optional because it has a default, got %#v", createType)
	}
	assertPropertyEnum(t, "一点号.createType", createType, 1, 2)

	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "一点号视频标题",
		"description": "一点号视频描述",
		"tags":        []interface{}{"生活"},
		"category": []interface{}{
			map[string]interface{}{
				"id":   "6",
				"text": "美食",
				"raw":  map[string]interface{}{"id": "6"},
			},
		},
		"createType": float64(1),
	}

	result, err := validator.ValidateStrict("一点号", "video", payload)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Valid {
		t.Fatalf("expected yidianhao video payload with createType to pass, got %v", result.Errors)
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
	declaration, ok := schemaDoc.Properties["declaration"]
	if !ok {
		t.Fatalf("expected shipinhao video schema to expose declaration, got %+v", schemaDoc.Properties)
	}
	if declaration.Default != float64(0) {
		t.Fatalf("expected declaration default 0, got %#v", declaration.Default)
	}
	drama, ok := schemaDoc.Properties["drama"]
	if !ok {
		t.Fatalf("expected shipinhao video schema to expose drama, got %+v", schemaDoc.Properties)
	}
	if drama.AdditionalProperties == nil || *drama.AdditionalProperties {
		t.Fatalf("expected drama schema to explicitly reject additional properties, got %#v", drama.AdditionalProperties)
	}
	if _, ok := drama.Properties["yixiaoerImageUrl"]; !ok {
		t.Fatalf("expected drama schema to expose yixiaoerImageUrl, got %+v", drama.Properties)
	}
}

func TestValidateAcceptsExactShipinhaoDramaObject(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":   "task",
		"createType": float64(2),
		"pubType":    float64(1),
		"drama": map[string]interface{}{
			"yixiaoerId":       "event/1",
			"yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover",
			"yixiaoerName":     "风浪过后护妻安康",
		},
	}

	result := validator.Validate("视频号", "video", payload)
	if !result.Valid {
		t.Fatalf("expected exact shipinhao drama object to pass, got %v", result.Errors)
	}
}

func TestValidateRejectsShipinhaoDramaRawAndUnknownFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":   "task",
		"createType": float64(2),
		"pubType":    float64(1),
		"drama": map[string]interface{}{
			"yixiaoerId":       "event/1",
			"yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover",
			"yixiaoerName":     "风浪过后护妻安康",
			"raw":              map[string]interface{}{"id": "event/1"},
			"extra":            "must reject",
		},
	}

	result := validator.Validate("视频号", "video", payload)
	if result.Valid {
		t.Fatal("expected drama raw and unknown fields to be rejected")
	}
	if !containsError(result.Errors, `unexpected field "raw"`) || !containsError(result.Errors, `unexpected field "extra"`) {
		t.Fatalf("expected strict drama field errors, got %v", result.Errors)
	}
}

func TestValidateRejectsShipinhaoDramaMissingAndCLICommonFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	newPayload := func() map[string]interface{} {
		return map[string]interface{}{
			"formType":   "task",
			"createType": float64(2),
			"pubType":    float64(1),
			"drama": map[string]interface{}{
				"yixiaoerId":       "event/1",
				"yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover",
				"yixiaoerName":     "风浪过后护妻安康",
			},
		}
	}

	tests := map[string]func(map[string]interface{}){
		"missing image": func(payload map[string]interface{}) {
			delete(payload["drama"].(map[string]interface{}), "yixiaoerImageUrl")
		},
		"content": func(payload map[string]interface{}) {
			payload["drama"].(map[string]interface{})["content"] = "must reject"
		},
		"clientId": func(payload map[string]interface{}) {
			payload["drama"].(map[string]interface{})["clientId"] = "must reject"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			payload := newPayload()
			mutate(payload)
			result := validator.Validate("视频号", "video", payload)
			if result.Valid {
				t.Fatalf("expected strict drama rejection: %#v", payload)
			}
		})
	}
}

func TestValidateAcceptsShipinhaoVideoAiDeclaration(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"createType":  float64(2),
		"declaration": float64(1),
		"pubType":     float64(1),
	}

	result := validator.Validate("视频号", "video", payload)
	if !result.Valid {
		t.Fatalf("expected shipinhao video AI declaration payload to pass, got %v", result.Errors)
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
					"formType":      "task",
					"title":         "标题",
					"description":   "描述",
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

func TestValidateAcceptsDuoduoshipinUserProvidedGoodsID(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType": "task",
		"shopping_cart": map[string]interface{}{
			"goods_id": "998877",
			"source":   "pdd",
		},
	}

	result := validator.Validate("多多视频", "video", payload)
	if !result.Valid {
		t.Fatalf("expected Duoduoshipin user-provided goods_id to pass, got %v", result.Errors)
	}
}

func TestValidateRejectsDuoduoshipinShoppingCartWithoutUserGoodsID(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType": "task",
		"shopping_cart": map[string]interface{}{
			"source": "pdd",
		},
	}

	result := validator.Validate("多多视频", "video", payload)
	if result.Valid {
		t.Fatal("expected Duoduoshipin shopping_cart without goods_id to be rejected")
	}
	if !containsError(result.Errors, `missing required field "goods_id"`) {
		t.Fatalf("expected missing user goods_id error, got %v", result.Errors)
	}
}

func TestValidateRejectsDuoduoshipinQueryGoodsObject(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType": "task",
		"shopping_cart": map[string]interface{}{
			"yixiaoerId":   "goods_001",
			"yixiaoerName": "查询商品",
			"raw":          map[string]interface{}{"id": "goods_001"},
		},
	}

	result := validator.Validate("多多视频", "video", payload)
	if result.Valid {
		t.Fatal("expected Duoduoshipin query goods object to be rejected")
	}
	if !containsError(result.Errors, `unexpected field "yixiaoerId"`) {
		t.Fatalf("expected query object identity field rejection, got %v", result.Errors)
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

func TestValidateAcceptsDouyinFrontendLocationAndGroupShoppingShape(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "带位置视频",
		"description": "视频描述",
		"declaration": float64(0),
		"tagType":     "团购",
		"location": map[string]interface{}{
			"isScp": false,
			"data": map[string]interface{}{
				"yixiaoerId":   "loc_001",
				"yixiaoerName": "上海",
				"raw":          map[string]interface{}{"id": "loc_001"},
			},
		},
		"group_shopping": map[string]interface{}{
			"brand_switch_value": float64(0),
			"sale_title":         "点击购买",
			"data": map[string]interface{}{
				"yixiaoerId":   "goods_001",
				"yixiaoerName": "测试商品",
				"raw":          map[string]interface{}{"id": "goods_001"},
			},
		},
	}

	result := validator.Validate("抖音", "video", payload)
	if !result.Valid {
		t.Fatalf("expected douyin frontend location/group_shopping structure to pass, got %v", result.Errors)
	}
}

func TestValidateAcceptsDouyinVideoDescriptionOverThirtyCharacters(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "视频标题",
		"description": strings.Repeat("字", 31),
	}

	result := validator.Validate("抖音", "video", payload)
	if !result.Valid {
		t.Fatalf("expected douyin video description over 30 characters to pass, got %v", result.Errors)
	}
}

func TestValidateRejectsDouyinVideoDescriptionOverOneThousandCharacters(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"title":       "视频标题",
		"description": strings.Repeat("字", 1001),
	}

	result := validator.Validate("抖音", "video", payload)
	if result.Valid {
		t.Fatal("expected douyin video description over 1000 characters to fail")
	}
	if !containsError(result.Errors, "description: must NOT have more than 1000 characters") {
		t.Fatalf("expected description maxLength error, got %v", result.Errors)
	}
}

func TestValidateAcceptsFrontendPlatformDataLocationShape(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"description": "快手视频描述",
		"location": map[string]interface{}{
			"id":   "loc_001",
			"text": "上海",
			"raw":  map[string]interface{}{"id": "loc_001"},
		},
	}

	result := validator.Validate("快手", "video", payload)
	if !result.Valid {
		t.Fatalf("expected frontend id/text/raw location structure to pass, got %v", result.Errors)
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
		"formType":      "task",
		"title":         "百家号征文文章",
		"content":       "<p>正文</p>",
		"pubType":       float64(1),
		"declaration":   float64(1),
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

func TestValidateAcceptsXhsImageTextScheduledFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":      "task",
		"description":   "<p>小红书图文内容</p>",
		"visibleType":   float64(0),
		"createType":    float64(1),
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
	}

	result := validator.Validate("小红书", "imageText", payload)
	if !result.Valid {
		t.Fatalf("expected xiaohongshu imageText scheduled payload to pass, got %v", result.Errors)
	}
}

func TestValidateRejectsXhsImageTextMusicField(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":    "task",
		"description": "<p>小红书图文内容</p>",
		"visibleType": float64(0),
		"music": map[string]interface{}{
			"yixiaoerId":   "music_1",
			"yixiaoerName": "背景音乐",
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
	if result.Valid {
		t.Fatal("expected xiaohongshu imageText music field to be rejected")
	}
	if !containsError(result.Errors, `unexpected field "music"`) {
		t.Fatalf("expected music additionalProperties error, got %v", result.Errors)
	}
}

func TestSchemaExposesXhsImageTextCreateType(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	schemaDoc, err := validator.Schema("小红书", "imageText")
	if err != nil {
		t.Fatal(err)
	}
	createType, ok := schemaDoc.Properties["createType"]
	if !ok {
		t.Fatalf("expected xiaohongshu imageText schema to expose createType, got %+v", schemaDoc.Properties)
	}
	if createType.Default != float64(0) {
		t.Fatalf("expected createType default 0, got %#v", createType.Default)
	}
	if len(createType.Enum) != 2 || createType.Enum[0] != float64(0) || createType.Enum[1] != float64(1) {
		t.Fatalf("expected createType enum [0 1], got %#v", createType.Enum)
	}
}

func TestValidateAcceptsBaijiahaoImageTextPayload(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	payload := map[string]interface{}{
		"formType":      "task",
		"title":         "百家号图文标题",
		"description":   "<p>百家号图文内容</p>",
		"pubType":       float64(0),
		"declaration":   float64(0),
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
			"id":   "loc_1",
			"text": "北京",
			"raw":  map[string]interface{}{"id": "loc_1"},
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
			"id":   "loc_1",
			"text": "上海",
			"raw":  map[string]interface{}{"id": "loc_1"},
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

func TestValidateAcceptsWebPushedVideoPlatformFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	resource := func(key string) map[string]interface{} {
		return map[string]interface{}{
			"key":      key,
			"size":     float64(100),
			"width":    float64(10),
			"height":   float64(10),
			"duration": float64(10),
		}
	}
	queryItem := func(id string) map[string]interface{} {
		return map[string]interface{}{
			"yixiaoerId":   id,
			"yixiaoerName": id,
			"raw":          map[string]interface{}{"id": id},
		}
	}

	cases := []struct {
		name     string
		platform string
		payload  map[string]interface{}
	}{
		{
			name:     "douyin film and cooperation_info",
			platform: "抖音",
			payload: map[string]interface{}{
				"formType":    "task",
				"title":       "抖音视频",
				"description": "视频描述",
				"film":        queryItem("film_1"),
				"cooperation_info": []interface{}{
					map[string]interface{}{
						"co_type": queryItem("co_type_1"),
						"user_list": []interface{}{
							queryItem("friend_1"),
						},
					},
				},
			},
		},
		{
			name:     "kuaishou author service and pk cover",
			platform: "快手",
			payload: map[string]interface{}{
				"formType":            "task",
				"description":         "快手视频描述",
				"author_service_type": "shopping_cart",
				"pk_cover":            resource("pk-cover-key"),
			},
		},
		{
			name:     "shipinhao shopping_cart",
			platform: "视频号",
			payload: map[string]interface{}{
				"formType":      "task",
				"createType":    float64(2),
				"pubType":       float64(1),
				"shopping_cart": queryItem("goods_1"),
			},
		},
		{
			name:     "baijiahao statement without legacy title tags",
			platform: "百家号",
			payload: map[string]interface{}{
				"formType":    "task",
				"description": "百家号视频描述",
				"pubType":     float64(1),
				"statement": map[string]interface{}{
					"type":    float64(1),
					"subType": float64(2),
				},
			},
		},
		{
			name:     "xinlang category and visibleType",
			platform: "新浪微博",
			payload: map[string]interface{}{
				"formType":    "task",
				"title":       "微博视频",
				"description": "微博视频描述",
				"video":       resource("video-key"),
				"createType":  float64(1),
				"visibleType": float64(0),
				"category": []interface{}{
					queryItem("category_1"),
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.Validate(tt.platform, "video", tt.payload)
			if !result.Valid {
				t.Fatalf("expected %s web-pushed video payload to pass, got %v", tt.platform, result.Errors)
			}
		})
	}
}

func TestValidateAcceptsWebPushedArticlePlatformFields(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	cover := map[string]interface{}{
		"key":    "cover-key",
		"size":   float64(100),
		"width":  float64(10),
		"height": float64(10),
	}

	cases := []struct {
		name     string
		platform string
		payload  map[string]interface{}
	}{
		{
			name:     "douyin description field",
			platform: "抖音",
			payload: map[string]interface{}{
				"formType":    "task",
				"title":       "抖音文章",
				"content":     "<p>正文</p>",
				"covers":      []interface{}{cover},
				"visibleType": float64(0),
				"description": "摘要",
			},
		},
		{
			name:     "baijiahao coverType field",
			platform: "百家号",
			payload: map[string]interface{}{
				"formType":  "task",
				"title":     "百家号文章",
				"content":   "<p>正文</p>",
				"pubType":   float64(1),
				"coverType": "single",
				"covers":    []interface{}{cover},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.Validate(tt.platform, "article", tt.payload)
			if !result.Valid {
				t.Fatalf("expected %s web-pushed article payload to pass, got %v", tt.platform, result.Errors)
			}
		})
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
	raw, err := os.ReadFile(filepath.Join("..", "..", "test", "fixtures", "payloads", name))
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
