package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	materialflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/material"
)

func init() {
	rootCmd.AddCommand(newMaterialCmd())
}

type materialAddOptions struct {
	FilePath  string
	ThumbPath string
	Type      string
	DryRun    bool
}

func newMaterialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "material",
		Short: "管理素材库",
		Long:  "管理素材库资源。\n1. create <payload.json>：提交素材登记 payload\n2. add --file ...：上传资源并登记到素材库",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newMaterialCreateCmd())
	cmd.AddCommand(newMaterialAddCmd())
	return cmd
}

func newMaterialCreateCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "create <payload.json>",
		Short: "将已上传资源登记到素材库",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := readPayload(args[0])
			if err != nil {
				return err
			}
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "material.create.dry-run", map[string]interface{}{
					"dryRun":  true,
					"request": materialflow.BuildMaterialBody(payload),
				})
			}
			rt, err := app.Load()
			if err != nil {
				return err
			}
			result, err := materialflow.NewService(rt).Create(payload)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "material.create", result)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview material create request without performing the write")
	return cmd
}

func newMaterialAddCmd() *cobra.Command {
	opts := materialAddOptions{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "上传资源并登记到素材库",
		RunE: func(cmd *cobra.Command, args []string) error {
			input := materialflow.AddInput{
				FilePath:  opts.FilePath,
				ThumbPath: opts.ThumbPath,
				Type:      opts.Type,
			}
			if opts.DryRun {
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
			rt, err := app.Load()
			if err != nil {
				return err
			}
			result, err := materialflow.NewService(rt).Add(input)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "material.add", result)
		},
	}
	cmd.Flags().StringVar(&opts.FilePath, "file", "", "local file path or URL to upload and register")
	cmd.Flags().StringVar(&opts.ThumbPath, "thumb", "", "optional thumbnail path or URL")
	cmd.Flags().StringVar(&opts.Type, "type", "", "optional material type override: image, video, file")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview upload and material request without performing the write")
	return cmd
}
