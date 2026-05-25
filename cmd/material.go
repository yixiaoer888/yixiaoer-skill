package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	materialCmd.AddCommand(materialCreateCmd)
	rootCmd.AddCommand(materialCmd)
}

var materialCmd = &cobra.Command{
	Use:   "material",
	Short: "管理素材库",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var materialCreateCmd = &cobra.Command{
	Use:   "create <payload.json>",
	Short: "将已上传资源登记到素材库",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		payload, err := readPayload(args[0])
		if err != nil {
			return err
		}
		body := map[string]interface{}{}
		for _, field := range []string{"filePath", "fileName", "width", "height", "type", "thumbPath"} {
			if value, ok := payload[field]; ok {
				body[field] = value
			}
		}
		for _, required := range []string{"filePath", "fileName", "width", "height", "type"} {
			if _, ok := body[required]; !ok {
				return yxerrors.Usage("material create requires payload fields", []string{
					"filePath",
					"fileName",
					"width",
					"height",
					"type",
				})
			}
		}
		result, err := api.NewClient(cfg).Material(body)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "material.create", result)
	},
}
