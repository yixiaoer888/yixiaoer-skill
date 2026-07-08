package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	uploadflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/upload"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(newUploadCmd())
}

type uploadOptions struct {
	Bucket   string
	File     string
	URL      string
	DryRun   bool
	AutoMeta bool
}

func newUploadCmd() *cobra.Command {
	opts := uploadOptions{
		Bucket:   "cloud-publish",
		AutoMeta: true,
	}
	cmd := &cobra.Command{
		Use:   "upload [file_path_or_url]",
		Short: "上传资源",
		Long:  "上传本地文件或 URL 资源。\n默认上传到 cloud-publish；素材库资源使用 --bucket material-library。\n可通过位置参数、--file 或 --url 指定唯一资源来源。",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpload(cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.Bucket, "bucket", "cloud-publish", "upload bucket")
	cmd.Flags().StringVar(&opts.File, "file", "", "local file path to upload")
	cmd.Flags().StringVar(&opts.URL, "url", "", "remote URL to upload")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview upload request without performing the write")
	cmd.Flags().BoolVar(&opts.AutoMeta, "auto-meta", true, "extract media metadata automatically for uploaded assets")
	return cmd
}

func runUpload(cmd *cobra.Command, args []string, opts uploadOptions) error {
	source, err := resolveUploadSource(args, opts)
	if err != nil {
		return err
	}
	if opts.DryRun {
		result, err := uploadflow.Preview(source, opts.Bucket, opts.AutoMeta)
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
	result, err := uploadflow.NewService(rt).Upload(source, opts.Bucket, opts.AutoMeta)
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), "upload", result)
}

func resolveUploadSource(args []string, opts uploadOptions) (string, error) {
	sources := make([]string, 0, 3)
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		sources = append(sources, strings.TrimSpace(args[0]))
	}
	if strings.TrimSpace(opts.File) != "" {
		sources = append(sources, strings.TrimSpace(opts.File))
	}
	if strings.TrimSpace(opts.URL) != "" {
		sources = append(sources, strings.TrimSpace(opts.URL))
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
