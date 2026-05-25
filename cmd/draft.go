package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
)

func init() {
	draftCmd.AddCommand(draftSaveCmd)
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
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		payload, err := readPayload(args[0])
		if err != nil {
			return err
		}
		delete(payload, "action")
		payload["isDraft"] = true
		result, err := api.NewClient(cfg).SaveDraft(payload)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "draft.save", result)
	},
}
