package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var (
	publishChannelFlag       string
	publishClientID          string
	publishDryRun            bool
	publishAutoFallbackLocal bool
)

func init() {
	publishCmd.Flags().StringVar(&publishChannelFlag, "publish-channel", "", `publish channel: "cloud" or "local"`)
	publishCmd.Flags().StringVar(&publishClientID, "client-id", "", "client ID for local publish")
	publishCmd.Flags().BoolVar(&publishDryRun, "dry-run", false, "preview the publish request without performing the write")
	publishCmd.Flags().BoolVar(&publishAutoFallbackLocal, "auto-fallback-local", false, "automatically retry with local publish when cloud publish fails due to proxy availability")
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish <type> <中文平台名|platform-key> <payload.json> [clientId]",
	Short: "发布内容（单平台原子发布）",
	Long:  "仅支持 payload.json 模式。发布前请先通过 prepare / schema fields 获取表单字段和前置数据；需要完整骨架时再补 schema get，随后执行 validate 和 publish。",
	Args:  cobra.RangeArgs(3, 4),
	RunE: func(cmd *cobra.Command, args []string) error {
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
		input := publishflow.ExecuteInput{
			PublishType:        args[0],
			PlatformInput:      args[1],
			Payload:            payload,
			PositionalClientID: positionalClientID,
			FlagChannel:        publishChannelFlag,
			FlagClientID:       publishClientID,
			AutoFallbackLocal:  publishAutoFallbackLocal,
		}
		if publishDryRun {
			result, err := publishflow.NewService().DryRun(input)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "publish.dry-run", map[string]interface{}{
				"dryRun":  true,
				"request": result.PublishBody,
				"meta": map[string]interface{}{
					"platform":       result.Platform,
					"publishType":    result.PublishType,
					"publishChannel": result.PublishMode,
					"clientId":       result.ClientID,
					"accountIds":     result.AccountIDs,
					"schemaChecked":  result.SchemaChecked,
				},
			})
		}
		result, err := publishflow.NewService().Execute(input)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "publish", result)
	},
}
