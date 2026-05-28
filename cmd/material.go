package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	materialflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/material"
)

var (
	materialFilePath  string
	materialThumbPath string
	materialType      string
	materialDryRun    bool
	materialCreateDryRun bool
)

func init() {
	materialCmd.AddCommand(materialCreateCmd)
	materialCmd.AddCommand(materialAddCmd)
	rootCmd.AddCommand(materialCmd)
}

var materialCmd = &cobra.Command{
	Use:   "material",
	Short: "管理素材库",
	Long: "支持两种模式：\n1. create <payload.json>：兼容旧模式，直接提交素材登记 payload\n2. add --file ...：推荐模式，自动完成上传并登记到素材库",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var materialCreateCmd = &cobra.Command{
	Use:   "create <payload.json>",
	Short: "将已上传资源登记到素材库",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, err := readPayload(args[0])
		if err != nil {
			return err
		}
		if materialCreateDryRun {
			return output.Success(cmd.OutOrStdout(), "material.create.dry-run", map[string]interface{}{
				"dryRun":  true,
				"request": materialflow.BuildMaterialBody(payload),
			})
		}
		result, err := materialflow.NewService().Create(payload)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "material.create", result)
	},
}

var materialAddCmd = &cobra.Command{
	Use:   "add",
	Short: "上传资源并登记到素材库",
	RunE: func(cmd *cobra.Command, args []string) error {
		input := materialflow.AddInput{
			FilePath:  materialFilePath,
			ThumbPath: materialThumbPath,
			Type:      materialType,
		}
		if materialDryRun {
			result, err := materialflow.PreviewAdd(input)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "material.add.dry-run", map[string]interface{}{
				"dryRun":  true,
				"request": result.Request,
				"upload":  result.Upload,
				"thumb":   result.Thumb,
			})
		}
		result, err := materialflow.NewService().Add(input)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "material.add", result)
	},
}

func init() {
	materialCreateCmd.Flags().BoolVar(&materialCreateDryRun, "dry-run", false, "preview material create request without performing the write")
	materialAddCmd.Flags().StringVar(&materialFilePath, "file", "", "local file path or URL to upload and register")
	materialAddCmd.Flags().StringVar(&materialThumbPath, "thumb", "", "optional thumbnail path or URL")
	materialAddCmd.Flags().StringVar(&materialType, "type", "", "optional material type override: image, video, file")
	materialAddCmd.Flags().BoolVar(&materialDryRun, "dry-run", false, "preview upload and material request without performing the write")
}
