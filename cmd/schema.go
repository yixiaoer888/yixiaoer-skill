package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
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
			return yxerrors.Usage("schema requires <platform> and <type>", nil)
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
		})
	}
	return output.Success(cmd.OutOrStdout(), "schema.get", map[string]interface{}{
		"key":        schemaDoc.Key,
		"platform":   schemaDoc.Platform,
		"type":       schemaDoc.Type,
		"file":       filepath.ToSlash(schemaDoc.File),
		"rootSchema": schemaDoc.RootSchema,
		"document":   schemaDoc,
		"schema":     schemaDoc,
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
		})
	}
	fields, err := validator.Fields(platform, publishType)
	if err != nil {
		return yxerrors.Usage("schema not found", map[string]interface{}{
			"platform": platform,
			"type":     publishType,
		})
	}
	return output.Success(cmd.OutOrStdout(), "schema.fields", map[string]interface{}{
		"platform": doc.Platform,
		"type":     doc.Type,
		"key":      doc.Key,
		"fields":   fields,
	})
}
