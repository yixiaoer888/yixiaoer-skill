package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	draftflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/draft"
)

func init() {
	rootCmd.AddCommand(newDraftCmd())
}

func newDraftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draft",
		Short: "管理蚁小二草稿",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newDraftSaveCmd())
	return cmd
}

func newDraftSaveCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "save <payload.json>",
		Short: "保存为蚁小二草稿",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := readPayload(args[0])
			if err != nil {
				return err
			}
			if dryRun {
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
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the draft payload without performing the write")
	return cmd
}
