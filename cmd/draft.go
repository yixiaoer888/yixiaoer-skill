package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	draftflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/draft"
)

var draftDryRun bool

func init() {
	draftCmd.AddCommand(draftSaveCmd)
	draftSaveCmd.Flags().BoolVar(&draftDryRun, "dry-run", false, "preview the draft payload without performing the write")
	rootCmd.AddCommand(draftCmd)
}

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "管理蚁小二草稿",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var draftSaveCmd = &cobra.Command{
	Use:   "save <payload.json>",
	Short: "保存为蚁小二草稿",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, err := readPayload(args[0])
		if err != nil {
			return err
		}
		if draftDryRun {
			return output.Success(cmd.OutOrStdout(), "draft.save.dry-run", map[string]interface{}{
				"dryRun":  true,
				"request": draftflow.PreviewSave(payload),
			})
		}
		rt, err := app.Load()
		if err != nil {
			return err
		}
		result, err := draftflow.NewService(rt).Save(payload)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "draft.save", result)
	},
}
