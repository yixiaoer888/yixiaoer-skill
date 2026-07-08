package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/cmdflow"
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
			var payload map[string]interface{}

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					loaded, err := readPayload(args[0])
					if err != nil {
						return err
					}
					payload = loaded
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "draft.save.dry-run",
						Data: map[string]interface{}{
							"dryRun":  true,
							"request": draftflow.PreviewSave(payload),
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := draftflow.NewService(rt).Save(payload)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "draft.save", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the draft payload without performing the write")
	return cmd
}
