package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/core/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
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
	Short: "只返回指定平台和发布类型的字段视图",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSchemaFields(cmd, args[0], args[1])
	},
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
	return output.Success(cmd.OutOrStdout(), "schema.get", map[string]interface{}{
		"key":             schemaDoc.Key,
		"platform":        schemaDoc.Platform,
		"type":            schemaDoc.Type,
		"file":            filepath.ToSlash(schemaDoc.File),
		"rootSchema":      schemaDoc.RootSchema,
		"document":        envelopeSchema,
		"schema":          envelopeSchema,
		"businessSchema":  schemaDoc,
		"payloadTemplate": buildPayloadTemplate(schemaDoc),
	})
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
	return output.Success(cmd.OutOrStdout(), "schema.fields", map[string]interface{}{
		"platform": doc.Platform,
		"type":     doc.Type,
		"key":      doc.Key,
		"fields":   envelopeFields,
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
			Properties: map[string]schema.PropertyView{
				"content": {
					Type: "string",
				},
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
								Properties: businessFields,
							},
						},
					},
				},
			},
		},
	}
}

func buildPayloadTemplate(doc schema.Document) map[string]interface{} {
	return map[string]interface{}{
		"action":         "publish",
		"publishType":    doc.Type,
		"platforms":      []interface{}{platformutil.ChineseName(doc.Platform)},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "<platformAccountId>",
					"contentPublishForm": map[string]interface{}{
						"formType": "<see schema required fields>",
					},
				},
			},
		},
	}
}

func intPtr(value int) *int {
	return &value
}
