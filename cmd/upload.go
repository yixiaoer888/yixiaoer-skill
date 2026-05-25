package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
)

var uploadBucket string

func init() {
	uploadCmd.Flags().StringVar(&uploadBucket, "bucket", "cloud-publish", "upload bucket")
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload <file_path_or_url>",
	Short: "上传资源",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		result, err := api.NewClient(cfg).Upload(args[0], uploadBucket)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "upload", result)
	},
}
