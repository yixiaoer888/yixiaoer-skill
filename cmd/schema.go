package cmd

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(newSchemaCmd())
}

func newSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "查询 Agent 可用的参数 Schema",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			if len(args) != 2 {
				return yxerrors.Usage("schema requires <platform> and <type>", nil).
					WithHint("请同时提供平台和发布类型，例如：yxer schema get 抖音 video。").
					WithNextCommand("yxer schema list")
			}
			return runSchemaGet(cmd, args[0], args[1], false)
		},
	}
	cmd.AddCommand(newSchemaCatalogCmd())
	cmd.AddCommand(newSchemaListCmd())
	cmd.AddCommand(newSchemaGetCmd())
	cmd.AddCommand(newSchemaFieldsCmd())
	return cmd
}

func newSchemaListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有平台和发布类型 Schema",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSchemaList(cmd)
		},
	}
}

func newSchemaCatalogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "catalog",
		Short: "返回 schema 根目录、根 schema 和平台 schema 索引",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := app.Load()
			if err != nil {
				return err
			}
			catalog, err := schema.NewValidator(rt.Config.SchemaDir).Catalog()
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "schema.catalog", catalog)
		},
	}
}

func newSchemaGetCmd() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "get <中文平台名|platform-key> <type>",
		Short: "返回指定平台和发布类型的 JSON Schema",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSchemaGet(cmd, args[0], args[1], verbose)
		},
	}
	cmd.Flags().BoolVar(&verbose, "verbose", false, "include duplicated debug schema views")
	return cmd
}

func newSchemaFieldsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fields <中文平台名|platform-key> <type>",
		Short: "返回指定平台和发布类型的紧凑字段视图",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSchemaFields(cmd, args[0], args[1])
		},
	}
}

type flatFieldView struct {
	Path     string        `json:"path"`
	Type     string        `json:"type,omitempty"`
	Required bool          `json:"required,omitempty"`
	Enum     []interface{} `json:"enum,omitempty"`
	Const    interface{}   `json:"const,omitempty"`
	Default  interface{}   `json:"default,omitempty"`
}

type fieldPlacementView struct {
	SchemaPath string   `json:"schemaPath"`
	InputPaths []string `json:"inputPaths"`
	Note       string   `json:"note,omitempty"`
}

func runSchemaList(cmd *cobra.Command) error {
	rt, err := app.Load()
	if err != nil {
		return err
	}
	entries, err := schema.NewValidator(rt.Config.SchemaDir).List()
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), "schema.list", map[string]interface{}{
		"schemaDir": filepath.ToSlash(rt.Config.SchemaDir),
		"count":     len(entries),
		"items":     entries,
	})
}

func runSchemaGet(cmd *cobra.Command, platform, publishType string, verbose bool) error {
	rt, err := app.Load()
	if err != nil {
		return err
	}
	schemaDoc, err := schema.NewValidator(rt.Config.SchemaDir).Schema(platform, publishType)
	if err != nil {
		return yxerrors.Usage("schema not found", map[string]interface{}{
			"platform": platform,
			"type":     publishType,
		}).
			WithHint("未找到对应平台和发布类型的 schema，请先查看支持的平台和类型列表。").
			WithNextCommand("yxer schema list")
	}

	envelopeSchema := buildStandardPublishSchema(schemaDoc)

	// 基础返回结果（简化版）
	result := map[string]interface{}{
		"key":      schemaDoc.Key,
		"platform": schemaDoc.Platform,
		"type":     schemaDoc.Type,
		"file":     filepath.ToSlash(schemaDoc.File),

		// 只返回业务字段定义（最核心的部分）
		"businessFields":  schemaDoc.Properties,
		"fieldPlacements": buildFieldPlacements(schemaDoc),

		// 标准结构说明（文本形式）
		"standardStructure": map[string]interface{}{
			"description": "所有平台统一使用的标准 payload 结构",
			"envelope": append([]string{
				"action: 'publish' (固定值)",
				"publishType: '" + publishType + "' (固定值)",
				"platforms: ['" + platformutil.ChineseName(schemaDoc.Platform) + "'] (固定值)",
				"publishChannel: 'cloud' | 'local' (默认 cloud)",
				"publishArgs: { ... } (必填，包含 accountForms)",
				"publishArgs.accountForms[]: 账号级表单数组",
				"publishArgs.accountForms[].platformAccountId: 账号ID (必填)",
				"publishArgs.accountForms[].cover / coverKey: 账号层资源字段；仅在该发布类型或平台需要单独封面时填写",
				"publishArgs.accountForms[].contentPublishForm: 业务字段 (必填，见 businessFields)",
			}, platformSpecificEnvelopeNotes(schemaDoc)...),
		},

		// 最小可用模板
		"minimalTemplate": buildMinimalPayloadTemplate(schemaDoc),

		// 动态字段公共示例（仅包含当前 schema 暴露的字段）
		"dynamicFieldExamples": buildDynamicFieldExamples(schemaDoc),

		// 使用指引
		"guidance": []string{
			"1. 优先使用 'yxer schema fields' 查看紧凑字段列表",
			"2. businessFields 只描述平台字段定义；实际填写位置请看 fieldPlacements，不能默认全部写进 contentPublishForm",
			"3. 复杂对象（location/music/challenge等）必须通过查询命令获取完整对象",
			"4. 资源（video/images/cover）必须先通过 'yxer upload' 上传并使用返回的完整对象；图文首图封面平台只需提供 images",
			"5. minimalTemplate 提供最小可用骨架，实际使用时需填入真实值",
		},

		"recommendedCommand": "yxer schema fields " + platform + " " + publishType,
	}

	// verbose 模式返回完整 schema（用于调试）
	if verbose {
		result["fullDocument"] = envelopeSchema
		result["accountFormSchema"] = buildAccountFormSchema(schemaDoc)
		result["contentPublishFormSchema"] = buildContentPublishFormSchema(schemaDoc)
		result["verboseNote"] = "完整 schema 仅用于调试，日常使用建议查看 businessFields 或使用 schema fields 命令"
	}

	return output.Success(cmd.OutOrStdout(), "schema.get", result)
}

func buildFieldPlacements(doc schema.Document) map[string]fieldPlacementView {
	if len(doc.Properties) == 0 {
		return nil
	}
	placements := make(map[string]fieldPlacementView, len(doc.Properties))
	for _, key := range sortedPropertyKeys(doc.Properties) {
		placements[key] = fieldPlacementFor(doc, key)
	}
	return placements
}

func fieldPlacementFor(doc schema.Document, key string) fieldPlacementView {
	if isWeixinAccountArticleDoc(doc) {
		return fieldPlacementView{
			SchemaPath: "businessFields." + key,
			InputPaths: []string{"publishArgs.platformForms.微信公众号." + key},
			Note:       "微信公众号文章平台字段应填写在 publishArgs.platformForms[\"微信公众号\"] 下，不走 accountForms[].contentPublishForm。",
		}
	}
	view := fieldPlacementView{
		SchemaPath: "businessFields." + key,
		InputPaths: []string{"publishArgs.accountForms[].contentPublishForm." + key},
	}
	switch key {
	case "video":
		view.InputPaths = []string{
			"publishArgs.accountForms[].video",
			"publishArgs.accountForms[].contentPublishForm.video",
		}
		view.Note = "视频资源运行时从 accountForms[] 层读取；请使用 yxer upload 返回的完整对象，并保留 duration。"
	case "images":
		view.InputPaths = []string{
			"publishArgs.accountForms[].images",
			"publishArgs.accountForms[].contentPublishForm.images",
		}
		view.Note = "图文图片运行时推荐填写 accountForms[].images；CLI 会兼容 contentPublishForm.images，并把账号级图片用于校验。"
	case "cover":
		view.InputPaths = []string{
			"publishArgs.accountForms[].cover",
			"publishArgs.accountForms[].contentPublishForm.cover",
		}
		view.Note = "平台端可能从 accountForms[] 层读取 cover；若 schema 在 contentPublishForm 暴露该字段，也要同步填写 accountForms[].cover。"
	case "coverKey":
		view.InputPaths = []string{
			"publishArgs.accountForms[].coverKey",
			"publishArgs.accountForms[].contentPublishForm.coverKey",
		}
		view.Note = "coverKey 需要和 accountForms[].cover.key 保持一致；若 contentPublishForm 也有该字段，两个层级都要同步。"
	case "horizontalCover":
		if doc.Type == "video" {
			view.InputPaths = []string{
				"publishArgs.horizontalCover",
				"publishArgs.accountForms[].contentPublishForm.horizontalCover",
			}
			view.Note = "横版封面最终写入 contentPublishForm.horizontalCover；可在 publishArgs.horizontalCover 共享填写，CLI 会自动补齐。"
		}
	case "content":
		if doc.Type == "article" {
			view.InputPaths = []string{"publishArgs.content"}
			view.Note = "文章正文应写在 publishArgs.content；CLI 会在校验阶段补齐内层副本，并在最终发布体移除 contentPublishForm.content。"
		}
	}
	return view
}

func runSchemaFields(cmd *cobra.Command, platform, publishType string) error {
	rt, err := app.Load()
	if err != nil {
		return err
	}
	validator := schema.NewValidator(rt.Config.SchemaDir)
	doc, err := validator.Schema(platform, publishType)
	if err != nil {
		return yxerrors.Usage("schema not found", map[string]interface{}{
			"platform": platform,
			"type":     publishType,
		}).
			WithHint("未找到对应平台和发布类型的 schema，请确认平台别名和类型名称是否正确。").
			WithNextCommand("yxer schema list")
	}
	fields, err := validator.Fields(platform, publishType)
	if err != nil {
		return yxerrors.Usage("schema not found", map[string]interface{}{
			"platform": platform,
			"type":     publishType,
		}).
			WithHint("字段视图生成失败，请先确认 schema 文件存在且格式有效。").
			WithNextCommand("yxer schema get <platform> <type>")
	}
	envelopeFields := buildStandardPublishFieldView(doc, fields)
	flatFields := flattenFieldViews(envelopeFields)

	// 按重要性分组字段
	grouped := groupFieldsByImportance(flatFields, platform, publishType)

	return output.Success(cmd.OutOrStdout(), "schema.fields", map[string]interface{}{
		"platform": doc.Platform,
		"type":     doc.Type,
		"key":      doc.Key,

		// 汇总统计
		"summary": map[string]interface{}{
			"total":         len(flatFields),
			"requiredCount": len(grouped.Required),
			"optionalCount": len(grouped.Optional),
			"complexCount":  len(grouped.Complex),
		},

		// 分组展示（AI 优先查看必填字段）
		"required": grouped.Required,
		"optional": grouped.Optional,
		"complex":  grouped.Complex,

		// 复杂字段的查询命令提示
		"queryCommands": buildQueryCommandHints(grouped.Complex, platform),

		// 动态字段公共示例（仅包含当前 schema 暴露的字段）
		"dynamicFieldExamples": buildDynamicFieldExamples(doc),

		// 平台特定说明
		"platformNotes": getPlatformSpecificNotes(platform, publishType),

		// 保留完整数据（向后兼容）
		"flatFields": flatFields,
		"fields":     envelopeFields,

		// 推荐使用方式
		"recommendedResponse": "required + optional（按需查看 complex）",
	})
}

func buildStandardPublishSchema(doc schema.Document) schema.Document {
	envelope := schema.Document{
		Key:                  doc.Key,
		Platform:             doc.Platform,
		Type:                 doc.Type,
		File:                 doc.File,
		RootSchema:           doc.RootSchema,
		Title:                doc.Title + " Payload",
		Required:             []string{"action", "publishType", "platforms", "publishArgs"},
		AdditionalProperties: true,
		Properties:           buildStandardPublishFieldView(doc, doc.Properties),
	}
	return envelope
}

func buildStandardPublishFieldView(doc schema.Document, businessFields map[string]schema.PropertyView) map[string]schema.PropertyView {
	platformName := platformutil.ChineseName(doc.Platform)
	contentPublishFields := contentPublishFormFieldsForEnvelope(doc)
	accountResourceFields := accountResourceFieldViews(doc)
	publishArgsProperties := map[string]schema.PropertyView{
		"cover": {
			Type: "object",
		},
		"coverKey": {
			Type: "string",
		},
		"accountForms": {
			Type:     "array",
			Required: true,
			MinItems: intPtr(1),
			Items: &schema.PropertyView{
				Type: "object",
				Properties: map[string]schema.PropertyView{
					"platformAccountId": {
						Type:     "string",
						Required: true,
					},
					"account_id": {
						Type: "string",
					},
					"video": {
						Type:       "object",
						Required:   accountResourceFields["video"].Required,
						Properties: resourceFieldProperties(true),
					},
					"images": {
						Type:     "array",
						Required: accountResourceFields["images"].Required,
						MinItems: accountResourceFields["images"].MinItems,
						Items: &schema.PropertyView{
							Type:       "object",
							Properties: resourceFieldProperties(false),
						},
					},
					"cover": {
						Type:       "object",
						Required:   accountResourceFields["cover"].Required,
						Properties: resourceFieldProperties(false),
					},
					"coverKey": {
						Type:     "string",
						Required: accountResourceFields["coverKey"].Required,
					},
					"contentPublishForm": {
						Type:       "object",
						Required:   true,
						Properties: contentPublishFields,
					},
				},
			},
		},
	}
	if supportsHorizontalCover(doc) {
		publishArgsProperties["horizontalCover"] = schema.PropertyView{
			Type:       "object",
			Properties: resourceFieldProperties(false),
		}
	}
	if doc.Type == "article" {
		publishArgsProperties["content"] = schema.PropertyView{
			Type:     "string",
			Required: true,
		}
	} else {
		publishArgsProperties["content"] = schema.PropertyView{
			Type: "string",
		}
	}
	if isWeixinAccountArticleDoc(doc) {
		delete(publishArgsProperties, "content")
		publishArgsProperties["platformForms"] = schema.PropertyView{
			Type: "object",
			Properties: map[string]schema.PropertyView{
				"微信公众号": {
					Type:       "object",
					Required:   true,
					Properties: doc.Properties,
				},
			},
		}
		accountFormItem := publishArgsProperties["accountForms"].Items
		if accountFormItem != nil {
			accountFormItem.Properties["contentPublishForm"] = schema.PropertyView{
				Type: "object",
			}
		}
	}
	return map[string]schema.PropertyView{
		"action": {
			Type:     "string",
			Required: true,
			Const:    "publish",
		},
		"publishType": {
			Type:     "string",
			Required: true,
			Const:    doc.Type,
		},
		"platforms": {
			Type:     "array",
			Required: true,
			MinItems: intPtr(1),
			Items: &schema.PropertyView{
				Type: "string",
				Enum: []interface{}{platformName},
			},
		},
		"publishChannel": {
			Type:    "string",
			Default: "cloud",
			Enum:    []interface{}{"cloud", "local"},
		},
		"clientId": {
			Type: "string",
		},
		"cover": {
			Type: "object",
		},
		"coverKey": {
			Type: "string",
		},
		"desc": {
			Type: "string",
		},
		"isDraft": {
			Type:    "boolean",
			Default: false,
		},
		"isAppContent": {
			Type:    "boolean",
			Default: false,
		},
		"publishArgs": {
			Type:       "object",
			Required:   true,
			Properties: publishArgsProperties,
		},
	}
}

func buildAccountFormSchema(doc schema.Document) schema.PropertyView {
	contentPublishFormSchema := schema.PropertyView{
		Type:       "object",
		Required:   true,
		Properties: contentPublishFormFieldsForEnvelope(doc),
	}
	if isWeixinAccountArticleDoc(doc) {
		contentPublishFormSchema = schema.PropertyView{Type: "object"}
	}
	properties := map[string]schema.PropertyView{
		"platformAccountId": {
			Type:     "string",
			Required: true,
		},
		"account_id": {
			Type: "string",
		},
		"video": {
			Type: "object",
		},
		"images": {
			Type: "array",
		},
		"cover": {
			Type: "object",
		},
		"coverKey": {
			Type: "string",
		},
		"contentPublishForm": {
			Type:       contentPublishFormSchema.Type,
			Required:   contentPublishFormSchema.Required,
			Properties: contentPublishFormSchema.Properties,
		},
	}
	return schema.PropertyView{
		Type:       "object",
		Required:   true,
		Properties: properties,
	}
}

func supportsHorizontalCover(doc schema.Document) bool {
	if doc.Type != "video" {
		return false
	}
	_, ok := doc.Properties["horizontalCover"]
	return ok
}

func buildContentPublishFormSchema(doc schema.Document) schema.Document {
	return schema.Document{
		Key:                  doc.Key + "#/publishArgs/accountForms/contentPublishForm",
		Platform:             doc.Platform,
		Type:                 doc.Type,
		File:                 doc.File,
		RootSchema:           doc.RootSchema,
		Title:                doc.Title,
		Required:             requiredPropertyKeys(contentPublishFormFieldsForEnvelope(doc)),
		AdditionalProperties: doc.AdditionalProperties,
		Properties:           contentPublishFormFieldsForEnvelope(doc),
	}
}

func contentPublishFormFieldsForEnvelope(doc schema.Document) map[string]schema.PropertyView {
	if isWeixinAccountArticleDoc(doc) {
		return nil
	}
	if doc.Type != "article" {
		return clonePropertyViewsWithoutKeys(doc.Properties, contentPublishFormExcludedKeys(doc.Type)...)
	}
	return clonePropertyViewsWithoutKeys(doc.Properties, contentPublishFormExcludedKeys(doc.Type)...)
}

func contentPublishFormExcludedKeys(publishType string) []string {
	keys := []string{"accountForms", "platformForms", "publishArgs", "publishChannel", "clientId", "action", "platforms", "publishType"}
	if publishType == "article" {
		keys = append(keys, "content")
		return keys
	}
	return append(keys, accountLevelResourceKeys(publishType)...)
}

func accountResourceFieldViews(doc schema.Document) map[string]schema.PropertyView {
	fields := map[string]schema.PropertyView{}
	if doc.Type == "video" {
		fields["video"] = schema.PropertyView{Required: true}
		fields["cover"] = schema.PropertyView{Required: true}
		fields["coverKey"] = schema.PropertyView{Required: true}
	}
	if doc.Type == "imageText" {
		fields["images"] = schema.PropertyView{Required: true, MinItems: intPtr(1)}
		if !platformutil.ImageTextUsesFirstImageAsCover(doc.Platform) {
			fields["cover"] = schema.PropertyView{Required: true}
			fields["coverKey"] = schema.PropertyView{Required: true}
		}
	}
	return fields
}

func resourceFieldProperties(includeDuration bool) map[string]schema.PropertyView {
	props := map[string]schema.PropertyView{
		"key":    {Type: "string", Required: true},
		"size":   {Type: "integer"},
		"width":  {Type: "integer"},
		"height": {Type: "integer"},
		"format": {Type: "string"},
	}
	if includeDuration {
		props["duration"] = schema.PropertyView{Type: "number", Required: true}
	}
	return props
}

func isWeixinAccountArticleDoc(doc schema.Document) bool {
	return doc.Type == "article" && platformutil.CanonicalKey(doc.Platform) == "weixin.account"
}

func requiredPropertyKeys(fields map[string]schema.PropertyView) []string {
	keys := make([]string, 0, len(fields))
	for key, prop := range fields {
		if prop.Required {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}

func intPtr(value int) *int {
	return &value
}

func flattenFieldViews(fields map[string]schema.PropertyView) []flatFieldView {
	var items []flatFieldView
	appendFlatFieldViews(&items, fields, "")
	return items
}

func appendFlatFieldViews(items *[]flatFieldView, fields map[string]schema.PropertyView, prefix string) {
	for _, key := range sortedPropertyKeys(fields) {
		view := fields[key]
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		*items = append(*items, flatFieldView{
			Path:     path,
			Type:     view.Type,
			Required: view.Required,
			Enum:     view.Enum,
			Const:    view.Const,
			Default:  view.Default,
		})
		if len(view.Properties) > 0 {
			appendFlatFieldViews(items, view.Properties, path)
		}
		if view.Items != nil {
			itemPath := path + "[]"
			*items = append(*items, flatFieldView{
				Path:     itemPath,
				Type:     view.Items.Type,
				Required: view.Required,
				Enum:     view.Items.Enum,
				Const:    view.Items.Const,
				Default:  view.Items.Default,
			})
			if len(view.Items.Properties) > 0 {
				appendFlatFieldViews(items, view.Items.Properties, itemPath)
			}
		}
	}
}

func sortedPropertyKeys(fields map[string]schema.PropertyView) []string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		left := fields[keys[i]]
		right := fields[keys[j]]
		if left.Required != right.Required {
			return left.Required
		}
		return keys[i] < keys[j]
	})
	return keys
}

// groupedFields 字段分组结果
type groupedFields struct {
	Required []flatFieldView
	Optional []flatFieldView
	Complex  []flatFieldView
}

// groupFieldsByImportance 按重要性分组字段
func groupFieldsByImportance(flatFields []flatFieldView, platform, publishType string) groupedFields {
	var required []flatFieldView
	var optional []flatFieldView
	var complex []flatFieldView

	for _, field := range flatFields {
		// 跳过标准顶层字段（这些字段在文档中已说明）
		if isStandardTopLevelField(field.Path) {
			continue
		}

		// 复杂对象字段（需要查询命令获取）
		if isComplexField(field.Path) {
			complex = append(complex, field)
		} else if field.Required {
			required = append(required, field)
		} else {
			optional = append(optional, field)
		}
	}

	return groupedFields{
		Required: required,
		Optional: optional,
		Complex:  complex,
	}
}

// isStandardTopLevelField 判断是否为标准顶层字段
func isStandardTopLevelField(path string) bool {
	standardFields := []string{
		"action", "publishType", "platforms", "publishChannel", "clientId",
	}
	for _, field := range standardFields {
		if path == field {
			return true
		}
	}
	return false
}

// isComplexField 判断是否为复杂对象字段（需要查询命令）
func isComplexField(path string) bool {
	complexPatterns := []string{
		"location", "music", "challenge", "collection", "sub_collection",
		"category", "goods", "shopping_cart", "group_shopping", "groupShopping",
		"mini_app", "hot_event", "game", "sync_apps",
		"cooperation_info", "friends", "group",
	}
	for _, pattern := range complexPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

// buildQueryCommandHints 构建复杂字段的查询命令提示
func buildQueryCommandHints(complexFields []flatFieldView, platform string) map[string]string {
	hints := make(map[string]string)
	seenTypes := make(map[string]bool)

	for _, field := range complexFields {
		path := field.Path

		// 提取字段类型
		var fieldType string
		if strings.Contains(path, "location") {
			fieldType = "location"
		} else if strings.Contains(path, "music") {
			fieldType = "music"
		} else if strings.Contains(path, "challenge") {
			fieldType = "challenge"
		} else if strings.Contains(path, "collection") || strings.Contains(path, "sub_collection") {
			fieldType = "collection"
		} else if strings.Contains(path, "category") {
			fieldType = "category"
		} else if strings.Contains(path, "goods") || strings.Contains(path, "shopping_cart") || strings.Contains(path, "group_shopping") || strings.Contains(path, "groupShopping") {
			fieldType = "goods"
		} else if strings.Contains(path, "mini_app") {
			fieldType = "mini_app"
		} else if strings.Contains(path, "hot_event") {
			fieldType = "hot_event"
		} else if strings.Contains(path, "game") {
			fieldType = "game"
		} else if strings.Contains(path, "friends") {
			fieldType = "friends"
		} else if strings.Contains(path, "group") {
			fieldType = "group"
		}

		// 避免重复添加
		if fieldType != "" && !seenTypes[fieldType] {
			seenTypes[fieldType] = true
			hints[fieldType] = getQueryCommand(fieldType)
		}
	}

	return hints
}

// getQueryCommand 获取字段类型对应的查询命令
func getQueryCommand(fieldType string) string {
	commands := map[string]string{
		"location":   "yxer query locations <account_id> [--query 关键词]",
		"music":      "yxer query music <account_id> [--query 关键词]",
		"challenge":  "yxer query challenges <account_id> [--query 关键词] [--type video]",
		"collection": "yxer query collections <account_id> [--type video|article]",
		"category":   "yxer query categories <account_id> [--type video|article]",
		"goods":      "yxer query goods <account_id> [--query 关键词]",
		"mini_app":   "yxer query miniapps <account_id> [--query 关键词]",
		"hot_event":  "yxer query hot-events <account_id> [--query 关键词]",
		"game":       "yxer query games <account_id> [--query 关键词]",
		"friends":    "yxer query friends <account_id>",
		"group":      "yxer query groups <account_id>",
	}
	if cmd, ok := commands[fieldType]; ok {
		return cmd
	}
	return ""
}

// getPlatformSpecificNotes 获取平台特定说明
func getPlatformSpecificNotes(platform, publishType string) []string {
	notes := []string{}

	// 标准化平台名
	platform = strings.ToLower(strings.TrimSpace(platform))

	switch platform {
	case "douyin", "抖音":
		if publishType == "video" {
			notes = append(notes, "抖音视频支持挂车(shopping_cart)、话题(challenge)、合集(collection)、热点(hot_event)等高级功能")
			notes = append(notes, "标题最大长度为30字符，描述最大长度为1000字符")
		} else if publishType == "imageText" {
			notes = append(notes, "抖音图文需要1-35张图片")
		}

	case "kuaishou", "快手":
		if publishType == "video" {
			notes = append(notes, "快手视频支持话题(challenge)和位置(location)")
		}

	case "kuaishou-open", "kuaishouopen", "快手-open", "快手-Open":
		if publishType == "video" {
			notes = append(notes, "快手-Open 视频使用开放平台发布，只要求 description；不支持浏览器发布通道")
			notes = append(notes, "平台草稿使用 contentPublishForm.pubType=0；私密发布使用 visibleType=1")
		}

	case "xiaohongshu", "xhs", "小红书":
		if publishType == "imageText" {
			notes = append(notes, "小红书图文需要1-9张图片，支持话题标签")
		} else if publishType == "video" {
			notes = append(notes, "小红书视频支持话题和位置")
		}

	case "weixin", "shipinhao", "视频号", "微信视频号":
		if publishType == "imageText" {
			notes = append(notes, "视频号图文不需要外部传入 cover/coverKey，CLI 会默认使用 images[0] 作为内部封面")
			notes = append(notes, "平台草稿使用 contentPublishForm.pubType=0；这不同于蚁小二草稿 isDraft=true")
		}
		if publishType == "video" {
			notes = append(notes, "视频号支持位置(location)和话题")
		}

	case "bilibili", "哔哩哔哩":
		if publishType == "video" {
			notes = append(notes, "B站视频需要选择分区(category)")
		}

	case "bilibili-open", "bilibiliopen", "哔哩哔哩-open", "哔哩哔哩-Open":
		if publishType == "video" {
			notes = append(notes, "哔哩哔哩-Open 视频需要选择分类(category)，分类对象必须来自 yxer query categories")
			notes = append(notes, "createType 为 web 表单字段；type 为后端 DTO 兼容字段，二者枚举均为 1-原创、2-转载")
		}

	case "weixin.account", "weixingongzhonghao", "WeiXinGongZhongHao", "微信公众号":
		if publishType == "imageText" {
			notes = append(notes, "微信公众号图文使用 accountForms[].contentPublishForm；images 放在 accountForms[].images，必须全部先通过 yxer upload 上传，默认首图为封面")
			notes = append(notes, "微信公众号图文默认不群发、关闭留言、无需声明、允许平台推荐并直接发布；需要保存平台草稿时设置 pubType=0")
			notes = append(notes, "微信公众号图文 scheduledTime 必须不早于当前时间 2 小时")
		}
		if publishType == "article" {
			notes = append(notes, "公众号文章使用 publishArgs.platformForms[\"微信公众号\"].articles[] 传递文章包，而不是通用 contentPublishForm")
			notes = append(notes, "公众号文章必须单平台发布；accountForms 仅用于声明目标账号")
		}
	}

	return notes
}

func platformSpecificEnvelopeNotes(doc schema.Document) []string {
	if doc.Type == "imageText" && platformutil.ImageTextUsesFirstImageAsCover(doc.Platform) {
		if platformutil.CanonicalKey(doc.Platform) == "weixin.account" {
			return []string{
				"微信公众号图文图片放在 accountForms[].images；未提供 coverKey 时 CLI 默认使用首图，若接口已返回 coverKey 可原样保留",
			}
		}
		return []string{
			platformutil.ChineseName(doc.Platform) + "图文不需要外部传入 cover/coverKey；CLI 会默认使用 images[0] 作为内部封面",
		}
	}
	return nil
}
