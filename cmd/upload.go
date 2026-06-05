package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	uploadflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/upload"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var uploadBucket string
var uploadFile string
var uploadURL string
var uploadDryRun bool
var uploadAutoMeta bool

func init() {
	uploadCmd.Flags().StringVar(&uploadBucket, "bucket", "cloud-publish", "upload bucket")
	uploadCmd.Flags().StringVar(&uploadFile, "file", "", "local file path to upload")
	uploadCmd.Flags().StringVar(&uploadURL, "url", "", "remote URL to upload")
	uploadCmd.Flags().BoolVar(&uploadDryRun, "dry-run", false, "preview upload request without performing the write")
	uploadCmd.Flags().BoolVar(&uploadAutoMeta, "auto-meta", true, "extract media metadata automatically for uploaded assets")
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload [file_path_or_url]",
	Short: "上传资源",
	Long:  "上传本地文件或 URL 资源。\n默认上传到 cloud-publish；素材库资源建议使用 --bucket material-library。\n推荐使用 --file 或 --url，兼容旧的位置参数模式。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source, err := resolveUploadSource(args)
		if err != nil {
			return err
		}
		if uploadDryRun {
			result, err := uploadflow.Preview(source, uploadBucket, uploadAutoMeta)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "upload.dry-run", map[string]interface{}{
				"dryRun":  true,
				"request": result,
			})
		}
		rt, err := app.Load()
		if err != nil {
			return err
		}
		result, err := uploadflow.NewService(rt).Upload(source, uploadBucket, uploadAutoMeta)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "upload", result)
	},
}

func resolveUploadSource(args []string) (string, error) {
	sources := make([]string, 0, 3)
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		sources = append(sources, strings.TrimSpace(args[0]))
	}
	if strings.TrimSpace(uploadFile) != "" {
		sources = append(sources, strings.TrimSpace(uploadFile))
	}
	if strings.TrimSpace(uploadURL) != "" {
		sources = append(sources, strings.TrimSpace(uploadURL))
	}
	if len(sources) == 0 {
		return "", yxerrors.Usage("upload requires a file path or URL", nil).
			WithHint("请传入位置参数，或使用 --file / --url。").
			WithNextCommand("yxer upload --file ./cover.jpg")
	}
	if len(sources) > 1 {
		return "", yxerrors.Usage("upload accepts exactly one source", sources).
			WithHint("请在位置参数、--file、--url 三者中只保留一个输入来源。")
	}
	return sources[0], nil
}
