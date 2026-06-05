package cmd

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/core/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var schemaGetVerbose bool

func init() {
	schemaGetCmd.Flags().BoolVar(&schemaGetVerbose, "verbose", false, "include duplicated debug schema views")
	schemaCmd.AddCommand(schemaCatalogCmd)
	schemaCmd.AddCommand(schemaListCmd)
	schemaCmd.AddCommand(schemaGetCmd)
	schemaCmd.AddCommand(schemaFieldsCmd)
	rootCmd.AddCommand(schemaCmd)
}

var schemaCmd = &cobra.Command{
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
		return runSchemaGet(cmd, args[0], args[1])
	},
}

var schemaListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有平台和发布类型 Schema",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSchemaList(cmd)
	},
}

var schemaCatalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "返回 schema 根目录、根 schema 和平台 schema 索引",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		catalog, err := schema.NewValidator(cfg.SchemaDir).Catalog()
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "schema.catalog", catalog)
	},
}

var schemaGetCmd = &cobra.Command{
	Use:   "get <中文平台名|platform-key> <type>",
	Short: "返回指定平台和发布类型的 JSON Schema",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSchemaGet(cmd, args[0], args[1])
	},
}

var schemaFieldsCmd = &cobra.Command{
	Use:   "fields <中文平台名|platform-key> <type>",
	Short: "返回指定平台和发布类型的紧凑字段视图",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSchemaFields(cmd, args[0], args[1])
	},
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
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	entries, err := schema.NewValidator(cfg.SchemaDir).List()
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), "schema.list", map[string]interface{}{
		"schemaDir": filepath.ToSlash(cfg.SchemaDir),
		"count":     len(entries),
		"items":     entries,
	})
}

func runSchemaGet(cmd *cobra.Command, platform, publishType string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	schemaDoc, err := schema.NewValidator(cfg.SchemaDir).Schema(platform, publishType)
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
		"businessFields": schemaDoc.Properties,
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
				"publishArgs.accountForms[].cover / coverKey: 账号层资源字段；若 businessFields 也出现 cover / coverKey，需要同步填写",
				"publishArgs.accountForms[].contentPublishForm: 业务字段 (必填，见 businessFields)",
			}, platformSpecificEnvelopeNotes(schemaDoc)...),
		},

		// 最小可用模板
		"minimalTemplate": buildMinimalPayloadTemplate(schemaDoc),

		// 使用指引
		"guidance": []string{
			"1. 优先使用 'yxer schema fields' 查看紧凑字段列表",
			"2. businessFields 只描述平台字段定义；实际填写位置请看 fieldPlacements，不能默认全部写进 contentPublishForm",
			"3. 复杂对象（location/music/challenge等）必须通过查询命令获取完整对象",
			"4. 资源（video/images/cover）必须先通过 'yxer upload' 上传并使用返回的完整对象",
			"5. minimalTemplate 提供最小可用骨架，实际使用时需填入真实值",
		},

		"recommendedCommand": "yxer schema fields " + platform + " " + publishType,
	}

	// verbose 模式返回完整 schema（用于调试）
	if schemaGetVerbose {
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
	view := fieldPlacementView{
		SchemaPath: "businessFields." + key,
		InputPaths: []string{"publishArgs.accountForms[].contentPublishForm." + key},
	}
	switch key {
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
	case "content":
		if doc.Type == "article" {
			view.InputPaths = []string{"publishArgs.content"}
			view.Note = "文章正文应写在 publishArgs.content；CLI 会在校验阶段补齐内层副本，并在最终发布体移除 contentPublishForm.content。"
		}
	}
	return view
}

func runSchemaFields(cmd *cobra.Command, platform, publishType string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	validator := schema.NewValidator(cfg.SchemaDir)
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
						Type:       "object",
						Required:   true,
						Properties: contentPublishFields,
					},
				},
			},
		},
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
			Type:     "object",
			Required: true,
			Properties: publishArgsProperties,
		},
	}
}

func buildAccountFormSchema(doc schema.Document) schema.PropertyView {
	return schema.PropertyView{
		Type:     "object",
		Required: true,
		Properties: map[string]schema.PropertyView{
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
				Type:       "object",
				Required:   true,
				Properties: contentPublishFormFieldsForEnvelope(doc),
			},
		},
	}
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
	if doc.Type != "article" {
		return doc.Properties
	}
	return clonePropertyViewsWithoutKeys(doc.Properties, "content")
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
		"category", "goods", "shopping_cart", "groupShopping",
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
		} else if strings.Contains(path, "goods") || strings.Contains(path, "shopping_cart") || strings.Contains(path, "groupShopping") {
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
		"location":    "yxer locations <account_id> [--query 关键词]",
		"music":       "yxer music <account_id> [--query 关键词]",
		"challenge":   "yxer challenges <account_id> [--query 关键词] [--type video]",
		"collection":  "yxer collections <account_id> [--type video|article]",
		"category":    "yxer categories <account_id> [--type video|article]",
		"goods":       "yxer goods <account_id> [--query 关键词]",
		"mini_app":    "yxer miniapps <account_id> [--query 关键词]",
		"hot_event":   "yxer hot-events <account_id> [--query 关键词]",
		"game":        "yxer games <account_id> [--query 关键词]",
		"friends":     "yxer friends <account_id>",
		"group":       "yxer groups <account_id>",
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
			notes = append(notes, "标题和描述最大长度均为30字符")
		} else if publishType == "imageText" {
			notes = append(notes, "抖音图文需要1-35张图片")
		}

	case "kuaishou", "快手":
		if publishType == "video" {
			notes = append(notes, "快手视频支持话题(challenge)和位置(location)")
		}

	case "xiaohongshu", "xhs", "小红书":
		if publishType == "imageText" {
			notes = append(notes, "小红书图文需要1-9张图片，支持话题标签")
		} else if publishType == "video" {
			notes = append(notes, "小红书视频支持话题和位置")
		}

	case "weixin", "shipinhao", "视频号", "微信视频号":
		if publishType == "imageText" {
			notes = append(notes, "视频号图文除了 contentPublishForm.images，还需要在 accountForms[] 层同时提供 cover 和 coverKey")
			notes = append(notes, "平台草稿使用 contentPublishForm.pubType=0；这不同于蚁小二草稿 isDraft=true")
		}
		if publishType == "video" {
			notes = append(notes, "视频号支持位置(location)和话题")
		}

	case "bilibili", "哔哩哔哩":
		if publishType == "video" {
			notes = append(notes, "B站视频需要选择分区(category)")
		}

	case "weixin.account", "微信公众号":
		if publishType == "article" {
			notes = append(notes, "公众号文章支持富文本和多媒体，需要使用 content 字段传递文章内容")
		}
	}

	return notes
}

func platformSpecificEnvelopeNotes(doc schema.Document) []string {
	if doc.Platform == "shipinhao" && doc.Type == "imageText" {
		return []string{
			"视频号图文额外要求 publishArgs.accountForms[].cover 和 coverKey；建议与首图保持一致",
		}
	}
	return nil
}
