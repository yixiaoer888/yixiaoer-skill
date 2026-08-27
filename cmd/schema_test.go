package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestSchemaListCommandOutputsAgentDiscoverableItems(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaListCmd()
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["count"].(float64) == 0 {
		t.Fatal("expected schema list count")
	}
	items := data["items"].([]interface{})
	found := false
	for _, item := range items {
		entry := item.(map[string]interface{})
		if entry["platform"] == "douyin" && entry["type"] == "video" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected douyin video schema in items")
	}
}

func TestSchemaGetCommandOutputsSchemaForChinesePlatformAlias(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaGetCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"抖音", "video"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["platform"] != "douyin" || data["type"] != "video" {
		t.Fatalf("unexpected schema metadata: %#v", data)
	}
	if data["key"] != "douyin/video" {
		t.Fatalf("unexpected schema key: %#v", data["key"])
	}
	// businessFields holds the platform-specific content fields directly.
	businessFields := data["businessFields"].(map[string]interface{})
	title := businessFields["title"].(map[string]interface{})
	if title["type"] != "string" || title["required"] != true {
		t.Fatalf("expected required string title in businessFields, got %#v", title)
	}
	fieldPlacements := data["fieldPlacements"].(map[string]interface{})
	titlePlacement := fieldPlacements["title"].(map[string]interface{})
	titlePaths := titlePlacement["inputPaths"].([]interface{})
	if len(titlePaths) != 1 || titlePaths[0] != "publishArgs.accountForms[].contentPublishForm.title" {
		t.Fatalf("expected title placement under contentPublishForm, got %#v", titlePlacement)
	}
	// minimalTemplate provides a ready-to-edit skeleton using the standard envelope.
	template := data["minimalTemplate"].(map[string]interface{})
	if template["action"] != "publish" {
		t.Fatalf("expected minimalTemplate action=publish, got %#v", template)
	}
	templateArgs := template["publishArgs"].(map[string]interface{})
	templateForms := templateArgs["accountForms"].([]interface{})
	if len(templateForms) != 1 {
		t.Fatalf("expected single template account form, got %#v", templateForms)
	}
	// default (non-verbose) output must omit the debug-only schema views.
	for _, key := range []string{"fullDocument", "accountFormSchema", "contentPublishFormSchema", "businessSchema"} {
		if _, ok := data[key]; ok {
			t.Fatalf("expected default schema.get output to omit %q", key)
		}
	}
	if data["recommendedCommand"] != "yxer schema fields 抖音 video" {
		t.Fatalf("expected recommended schema.fields command, got %#v", data["recommendedCommand"])
	}
	guidance := data["guidance"].([]interface{})
	if len(guidance) < 3 {
		t.Fatalf("expected schema.get guidance, got %#v", guidance)
	}
}

func TestSchemaFieldsShipinhaoExposesDramaQueryExample(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaFieldsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"视频号", "video"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	queryCommands := data["queryCommands"].(map[string]interface{})
	if got := queryCommands["drama"]; got != "yxer query drama-tasks <account_id> [--query 关键词]" {
		t.Fatalf("unexpected drama query command: %#v", got)
	}
	examples := data["dynamicFieldExamples"].(map[string]interface{})
	drama := examples["drama"].(map[string]interface{})
	if drama["path"] != "publishArgs.accountForms[].contentPublishForm.drama" {
		t.Fatalf("unexpected drama example path: %#v", drama)
	}
	if drama["queryCommand"] != "yxer query drama-tasks <account_id> [--query 关键词] --json" {
		t.Fatalf("unexpected drama example command: %#v", drama["queryCommand"])
	}
	value := drama["value"].(map[string]interface{})
	if _, ok := value["raw"]; ok {
		t.Fatalf("drama example must not contain raw: %#v", value)
	}
	for _, field := range []string{"yixiaoerId", "yixiaoerImageUrl", "yixiaoerName"} {
		if _, ok := value[field]; !ok {
			t.Fatalf("drama example missing %q: %#v", field, value)
		}
	}
}

func TestSchemaFieldsDuoduoshipinTreatsShoppingCartAsManualGoodsID(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaFieldsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"多多视频", "video"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if queryCommands := data["queryCommands"].(map[string]interface{}); queryCommands["goods"] != nil {
		t.Fatalf("Duoduoshipin shopping_cart must not advertise goods query, got %#v", queryCommands)
	}
	if examples, ok := data["dynamicFieldExamples"].(map[string]interface{}); ok && examples["shopping_cart"] != nil {
		t.Fatalf("Duoduoshipin shopping_cart must be manually entered, got dynamic example %#v", examples["shopping_cart"])
	}

	placements := data["fieldPlacements"].(map[string]interface{})
	shoppingCart := placements["shopping_cart"].(map[string]interface{})
	if note := shoppingCart["note"].(string); !strings.Contains(note, "goods_id") || !strings.Contains(note, "手工") {
		t.Fatalf("expected manual goods_id placement note, got %#v", shoppingCart)
	}

	foundGoodsID := false
	for _, raw := range data["flatFields"].([]interface{}) {
		field := raw.(map[string]interface{})
		if field["path"] != "publishArgs.accountForms[].contentPublishForm.shopping_cart.goods_id" {
			continue
		}
		foundGoodsID = true
		if field["required"] != true {
			t.Fatalf("expected goods_id to be required when shopping_cart is present, got %#v", field)
		}
	}
	if !foundGoodsID {
		t.Fatal("expected manual shopping_cart.goods_id field in schema fields")
	}
}

func TestSchemaGetShipinhaoExposesStrictDramaSchema(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaGetCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"视频号", "video"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	fields := data["businessFields"].(map[string]interface{})
	drama := fields["drama"].(map[string]interface{})
	if drama["additionalProperties"] != false {
		t.Fatalf("expected drama additionalProperties=false, got %#v", drama)
	}
	properties := drama["properties"].(map[string]interface{})
	if len(properties) != 3 {
		t.Fatalf("expected exactly three drama properties, got %#v", properties)
	}
	if _, hasRaw := properties["raw"]; hasRaw {
		t.Fatalf("drama schema must not expose raw: %#v", properties)
	}
}

func TestSchemaGetCommandExplainsDuplicatedCoverPlacementForImageText(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaGetCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"抖音", "imageText"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	fieldPlacements := data["fieldPlacements"].(map[string]interface{})

	coverPlacement := fieldPlacements["cover"].(map[string]interface{})
	coverPaths := coverPlacement["inputPaths"].([]interface{})
	if len(coverPaths) != 2 ||
		coverPaths[0] != "publishArgs.accountForms[].cover" ||
		coverPaths[1] != "publishArgs.accountForms[].contentPublishForm.cover" {
		t.Fatalf("expected duplicated cover placement, got %#v", coverPlacement)
	}
	if coverPlacement["note"] == nil {
		t.Fatalf("expected cover placement note, got %#v", coverPlacement)
	}

	coverKeyPlacement := fieldPlacements["coverKey"].(map[string]interface{})
	coverKeyPaths := coverKeyPlacement["inputPaths"].([]interface{})
	if len(coverKeyPaths) != 2 ||
		coverKeyPaths[0] != "publishArgs.accountForms[].coverKey" ||
		coverKeyPaths[1] != "publishArgs.accountForms[].contentPublishForm.coverKey" {
		t.Fatalf("expected duplicated coverKey placement, got %#v", coverKeyPlacement)
	}
}

func TestSchemaGetCommandPlacesHorizontalCoverInContentPublishForm(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)

	for _, platform := range []string{"抖音", "视频号", "大鱼号"} {
		t.Run(platform, func(t *testing.T) {
			var out bytes.Buffer
			cmd := newSchemaGetCmd()
			cmd.SetOut(&out)
			cmd.SetArgs([]string{platform, "video"})

			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(out.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			data := response["data"].(map[string]interface{})
			fieldPlacements := data["fieldPlacements"].(map[string]interface{})
			horizontalPlacement := fieldPlacements["horizontalCover"].(map[string]interface{})
			horizontalPaths := horizontalPlacement["inputPaths"].([]interface{})
			if len(horizontalPaths) != 2 ||
				horizontalPaths[0] != "publishArgs.horizontalCover" ||
				horizontalPaths[1] != "publishArgs.accountForms[].contentPublishForm.horizontalCover" {
				t.Fatalf("expected horizontalCover placement with shared and contentPublishForm paths, got %#v", horizontalPlacement)
			}

			template := data["minimalTemplate"].(map[string]interface{})
			form := template["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
			if _, exists := form["horizontalCover"]; exists {
				t.Fatalf("did not expect horizontalCover at account form level in minimalTemplate, got %#v", form)
			}
			cpf := form["contentPublishForm"].(map[string]interface{})
			if _, exists := cpf["horizontalCover"]; exists {
				t.Fatalf("did not expect optional horizontalCover in minimalTemplate, got %#v", cpf)
			}
		})
	}
}

func TestSchemaGetCommandImageTextFirstImageCoverPlatformsOmitExternalCover(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)

	for _, platform := range []string{"新浪微博", "小红书", "视频号", "知乎", "头条号"} {
		t.Run(platform, func(t *testing.T) {
			var out bytes.Buffer
			cmd := newSchemaGetCmd()
			cmd.SetOut(&out)
			cmd.SetArgs([]string{platform, "imageText"})

			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(out.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			data := response["data"].(map[string]interface{})
			template := data["minimalTemplate"].(map[string]interface{})
			form := template["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
			if _, exists := form["cover"]; exists {
				t.Fatalf("did not expect %s imageText minimalTemplate to include cover, got %#v", platform, form)
			}
			if _, exists := form["coverKey"]; exists {
				t.Fatalf("did not expect %s imageText minimalTemplate to include coverKey, got %#v", platform, form)
			}
			if images, _ := form["images"].([]interface{}); len(images) != 1 {
				t.Fatalf("expected %s imageText minimalTemplate to include images, got %#v", platform, form)
			}
		})
	}
}

func TestSchemaGetCommandWeixinAccountImageTextUsesStandardFormAndDefaults(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaGetCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"WeiXinGongZhongHao", "imageText"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["key"] != "weixin.account/imageText" {
		t.Fatalf("unexpected schema key: %#v", data["key"])
	}
	placements := data["fieldPlacements"].(map[string]interface{})
	commentPlacement := placements["needOpenComment"].(map[string]interface{})
	paths := commentPlacement["inputPaths"].([]interface{})
	if len(paths) != 1 || paths[0] != "publishArgs.accountForms[].contentPublishForm.needOpenComment" {
		t.Fatalf("unexpected needOpenComment placement: %#v", commentPlacement)
	}
	template := data["minimalTemplate"].(map[string]interface{})
	form := template["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	if _, exists := form["cover"]; exists {
		t.Fatalf("did not expect a separate cover in WeChat imageText template: %#v", form)
	}
	if _, exists := form["coverKey"]; exists {
		t.Fatalf("did not expect a separate coverKey in WeChat imageText template: %#v", form)
	}
	fields := data["businessFields"].(map[string]interface{})
	for field, want := range map[string]interface{}{"disableRecommend": float64(0), "pubType": float64(1)} {
		property := fields[field].(map[string]interface{})
		if property["default"] != want {
			t.Fatalf("expected %s default in schema, got %#v", field, property)
		}
	}
}

func TestSchemaGetCommandWeixinImageTextKeepsImagesAtAccountLevel(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaGetCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"WeiXinGongZhongHao", "imageText"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	form := data["minimalTemplate"].(map[string]interface{})["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	images, ok := form["images"].([]interface{})
	if !ok || len(images) != 1 {
		t.Fatalf("expected one account-level image placeholder, got %#v", form["images"])
	}
	cpForm := form["contentPublishForm"].(map[string]interface{})
	if _, exists := cpForm["images"]; exists {
		t.Fatalf("did not expect an empty CPF images field when images are account-level: %#v", cpForm)
	}
	if cpForm["statement"] != float64(0) || cpForm["pubType"] != float64(1) {
		t.Fatalf("expected canonical WeChat imageText account form defaults, got %#v", cpForm)
	}
}

func TestSchemaFieldsCommandWeixinImageTextExcludesUnrelatedResources(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaFieldsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"WeiXinGongZhongHao", "imageText"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	paths := map[string]bool{}
	for _, entry := range data["flatFields"].([]interface{}) {
		paths[entry.(map[string]interface{})["path"].(string)] = true
	}
	for _, field := range []string{
		"publishArgs.accountForms[].cover",
		"publishArgs.accountForms[].cover.key",
		"publishArgs.accountForms[].video",
		"publishArgs.accountForms[].video.duration",
		"publishArgs.accountForms[].video.key",
	} {
		if paths[field] {
			t.Fatalf("did not expect unrelated WeChat imageText schema field %q", field)
		}
	}
	if !paths["publishArgs.accountForms[].images[].key"] {
		t.Fatal("expected the required WeChat image resource key field")
	}
}

func TestSchemaGetCommandVerboseOutputsDebugViews(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaGetCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"抖音", "video", "--verbose"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	accountFormSchema := data["accountFormSchema"].(map[string]interface{})
	accountFormProps := accountFormSchema["properties"].(map[string]interface{})
	if accountFormProps["platformAccountId"].(map[string]interface{})["required"] != true {
		t.Fatalf("expected accountFormSchema to require platformAccountId, got %#v", accountFormProps["platformAccountId"])
	}
	contentSchema := data["contentPublishFormSchema"].(map[string]interface{})
	contentProps := contentSchema["properties"].(map[string]interface{})
	if contentProps["title"].(map[string]interface{})["required"] != true {
		t.Fatalf("expected contentPublishFormSchema title to be required, got %#v", contentProps["title"])
	}
}

func TestSchemaCatalogCommandOutputsRootSchemasAndPlatforms(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaCatalogCmd()
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	rootSchemas := data["rootSchemas"].([]interface{})
	if len(rootSchemas) < 2 {
		t.Fatalf("expected root schemas, got %#v", data)
	}
	platforms := data["platforms"].([]interface{})
	if len(platforms) == 0 {
		t.Fatal("expected platform schema entries")
	}
}

func TestSchemaFieldsCommandOutputsHorizontalCoverFieldView(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	for _, platform := range []string{"抖音", "视频号", "大鱼号"} {
		t.Run(platform, func(t *testing.T) {
			var out bytes.Buffer
			cmd := newSchemaFieldsCmd()
			cmd.SetOut(&out)
			cmd.SetArgs([]string{platform, "video"})

			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(out.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			data := response["data"].(map[string]interface{})
			if data["recommendedResponse"] != "required + optional（按需查看 complex）" {
				t.Fatalf("expected grouped recommended response, got %#v", data["recommendedResponse"])
			}
			flatFields := data["flatFields"].([]interface{})
			foundHorizontalCover := false
			for _, entry := range flatFields {
				item := entry.(map[string]interface{})
				if item["path"] == "publishArgs.accountForms[].contentPublishForm.horizontalCover" {
					foundHorizontalCover = true
					if item["required"] == true {
						t.Fatalf("expected horizontalCover to be optional contentPublishForm field, got %#v", item)
					}
				}
				if item["path"] == "publishArgs.accountForms[].horizontalCover" {
					t.Fatalf("did not expect account-level horizontalCover field, got %#v", item)
				}
			}
			if !foundHorizontalCover {
				t.Fatal("expected contentPublishForm.horizontalCover in flatFields")
			}
		})
	}
}

func TestSchemaFieldsCommandExposesDouyinVideoDescriptionLimit(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaFieldsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"抖音", "video"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	fields := data["fields"].(map[string]interface{})
	publishArgs := fields["publishArgs"].(map[string]interface{})
	accountForms := publishArgs["properties"].(map[string]interface{})["accountForms"].(map[string]interface{})
	contentPublishForm := accountForms["items"].(map[string]interface{})["properties"].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	description := contentPublishForm["properties"].(map[string]interface{})["description"].(map[string]interface{})
	if description["required"] != true || description["maxLength"] != float64(1000) {
		t.Fatalf("expected douyin video description to be required with maxLength=1000, got %#v", description)
	}
	if _, ok := contentPublishForm["properties"].(map[string]interface{})["accountForms"]; ok {
		t.Fatalf("did not expect nested accountForms inside contentPublishForm: %#v", contentPublishForm)
	}

	for _, item := range data["flatFields"].([]interface{}) {
		field := item.(map[string]interface{})
		path, _ := field["path"].(string)
		if strings.Contains(path, ".contentPublishForm.accountForms") {
			t.Fatalf("did not expect recursive accountForms path in flatFields: %s", path)
		}
	}

	notes := data["platformNotes"].([]interface{})
	foundLimitNote := false
	for _, note := range notes {
		if note == "标题最大长度为30字符，描述最大长度为1000字符" {
			foundLimitNote = true
		}
	}
	if !foundLimitNote {
		t.Fatalf("expected douyin video platformNotes to mention frontend title/description limits, got %#v", notes)
	}
}

func TestSchemaFieldsCommandOutputsDynamicFieldExamples(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaFieldsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"抖音", "video"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	examples := data["dynamicFieldExamples"].(map[string]interface{})
	for _, key := range []string{"shopping_cart", "location", "music", "tags"} {
		if _, ok := examples[key]; !ok {
			t.Fatalf("expected dynamicFieldExamples.%s, got %#v", key, examples)
		}
	}
	cart := examples["shopping_cart"].(map[string]interface{})
	cartValue := cart["value"].([]interface{})[0].(map[string]interface{})
	if cartValue["data"] == nil || cartValue["images"] == nil {
		t.Fatalf("expected douyin shopping_cart example to use nested data/images, got %#v", cart)
	}
	if cart["queryCommand"] != "yxer query goods <account_id> [--query 关键词] --json" {
		t.Fatalf("expected goods query command in shopping cart example, got %#v", cart)
	}
}

func TestSchemaGetCommandSupportsOverseasVideoPlatforms(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)

	cases := []struct {
		name          string
		platform      string
		wantKey       string
		requiredField string
		enumField     string
	}{
		{
			name:          "tiktok",
			platform:      "TikTok",
			wantKey:       "tiktok/video",
			requiredField: "description",
			enumField:     "visible",
		},
		{
			name:          "youtube",
			platform:      "YouTube",
			wantKey:       "youtube/video",
			requiredField: "title",
			enumField:     "category",
		},
		{
			name:      "facebook",
			platform:  "Facebook",
			wantKey:   "facebook/video",
			enumField: "formType",
		},
		{
			name:      "instagram",
			platform:  "Instagram",
			wantKey:   "instagram/video",
			enumField: "share_to_feed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := newSchemaGetCmd()
			cmd.SetOut(&out)
			cmd.SetArgs([]string{tc.platform, "video"})

			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(out.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			data := response["data"].(map[string]interface{})
			if data["key"] != tc.wantKey {
				t.Fatalf("unexpected schema key: %#v", data["key"])
			}
			businessFields := data["businessFields"].(map[string]interface{})
			if tc.requiredField != "" {
				field := businessFields[tc.requiredField].(map[string]interface{})
				if field["required"] != true {
					t.Fatalf("expected %s to be required, got %#v", tc.requiredField, field)
				}
			}
			if _, ok := businessFields[tc.enumField]; !ok {
				t.Fatalf("expected overseas field %s in businessFields, got %#v", tc.enumField, businessFields)
			}
		})
	}
}

func TestSchemaFieldsCommandPlacesArticleContentUnderPublishArgs(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaFieldsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"知乎", "article"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	flatFields := data["flatFields"].([]interface{})
	foundPublishArgsContent := false
	for _, entry := range flatFields {
		item := entry.(map[string]interface{})
		if item["path"] == "publishArgs.content" {
			foundPublishArgsContent = true
			if item["required"] != true {
				t.Fatalf("expected publishArgs.content required for article, got %#v", item)
			}
		}
		if item["path"] == "publishArgs.accountForms[].contentPublishForm.content" {
			t.Fatalf("did not expect article content in contentPublishForm flatFields, got %#v", item)
		}
	}
	if !foundPublishArgsContent {
		t.Fatal("expected publishArgs.content in article flatFields")
	}
}

func TestSchemaFieldsCommandExposesDouyinArticleDescriptionAlias(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaFieldsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"抖音", "article"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	flatFields := data["flatFields"].([]interface{})
	foundDesc := false
	foundDescription := false
	for _, entry := range flatFields {
		item := entry.(map[string]interface{})
		if item["path"] == "publishArgs.accountForms[].contentPublishForm.desc" {
			foundDesc = true
		}
		if item["path"] == "publishArgs.accountForms[].contentPublishForm.description" {
			foundDescription = true
		}
	}
	if !foundDesc {
		t.Fatal("expected legacy article desc field in contentPublishForm")
	}
	if !foundDescription {
		t.Fatal("expected douyin article description field in contentPublishForm")
	}
}

func withGoBuildCache(t *testing.T) {
	t.Helper()
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOCACHE", filepath.Join(repoRoot, ".gocache"))
	t.Setenv("GOMODCACHE", filepath.Join(repoRoot, ".gomodcache"))
}
