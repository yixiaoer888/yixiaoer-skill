package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/cmdflow"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(newPublishCmd())
}

type publishOptions struct {
	Channel                     string
	ClientID                    string
	DryRun                      bool
	AutoFallbackLocal           bool
	ContinueOnContentImageError bool
}

func newPublishCmd() *cobra.Command {
	opts := publishOptions{}
	cmd := &cobra.Command{
		Use:   "publish <type> <中文平台名|platform-key> <payload.json>",
		Short: "发布内容（单平台原子发布）",
		Long:  "仅支持标准 payload.json。默认云发布；本机发布请显式传 --publish-channel local，并通过 --client-id 或 config set-local-client-id 提供客户端标识。第四个位置参数仅为旧版兼容，不再推荐使用。发布前请先执行 validate 和 publish --dry-run。",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPublish(cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.Channel, "publish-channel", "", `publish channel: "cloud" or "local"`)
	cmd.Flags().StringVar(&opts.ClientID, "client-id", "", "client ID for local publish")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview the publish request without performing the write")
	cmd.Flags().BoolVar(&opts.AutoFallbackLocal, "auto-fallback-local", false, "advanced: explicitly authorize retrying with local publish when cloud publish fails due to proxy availability")
	cmd.Flags().BoolVar(&opts.ContinueOnContentImageError, "continue-on-content-image-error", false, "advanced: continue publishing article content even when some img src URLs cannot be materialized")
	cmd.AddCommand(newPublishInitCmd())
	cmd.AddCommand(newPublishFormCmd())
	return cmd
}

func runPublish(cmd *cobra.Command, args []string, opts publishOptions) error {
	var input publishflow.ExecuteInput
	var service publishflow.Service

	return cmdflow.Run(cmd, opts.DryRun, cmdflow.Flow{
		Validate: func() error {
			rt, err := app.Load()
			if err != nil {
				return err
			}
			if err := detectSwappedPublishArgs(args[0], args[1], "publish <type> <platform> <payload.json>"); err != nil {
				return err
			}
			if !looksLikePayloadArg(args[2]) {
				return yxerrors.Usage("publish requires a payload.json file", []string{
					`Run "yxer prepare <platform> <type>" to inspect platform-specific form fields and preflight data.`,
					`Run "yxer schema fields <platform> <type>" to inspect the compact field list before filling the JSON file.`,
					`Run "yxer schema get <platform> <type>" only when you need the full payload skeleton.`,
				}).WithHint("发布命令已移除内容 flags 模式，请先准备 payload.json，再执行 validate / publish。")
			}
			payload, err := readPayload(args[2])
			if err != nil {
				return err
			}
			positionalClientID := ""
			if len(args) == 4 {
				positionalClientID = args[3]
			}
			input = publishflow.ExecuteInput{
				PublishType:                 args[0],
				PlatformInput:               args[1],
				PayloadPath:                 args[2],
				Payload:                     payload,
				PositionalClientID:          positionalClientID,
				FlagChannel:                 opts.Channel,
				FlagClientID:                opts.ClientID,
				AutoFallbackLocal:           opts.AutoFallbackLocal,
				ContinueOnContentImageError: opts.ContinueOnContentImageError,
			}
			service = publishflow.NewService(rt)
			return nil
		},
		DryRun: func() (cmdflow.Result, error) {
			result, err := service.DryRunEnvelope(input)
			if err != nil {
				return cmdflow.Result{}, err
			}
			return cmdflow.Result{Action: result.Action, Data: result.Data}, nil
		},
		Execute: func() (cmdflow.Result, error) {
			result, err := service.ExecuteEnvelope(input)
			if err != nil {
				return cmdflow.Result{}, err
			}
			return cmdflow.Result{Action: result.Action, Data: result.Data}, nil
		},
	})
}
