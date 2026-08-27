package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(newValidateCmd())
}

type validateOptions struct {
	Channel     string
	ClientID    string
	ContentFile string
}

func newValidateCmd() *cobra.Command {
	opts := validateOptions{}
	cmd := &cobra.Command{
		Use:   "validate <中文平台名> <type> <payload.json>",
		Short: "校验发布 Payload",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.Channel, "publish-channel", "", `publish channel: "cloud" or "local"`)
	cmd.Flags().StringVar(&opts.ClientID, "client-id", "", "client ID for local publish")
	cmd.Flags().StringVar(&opts.ContentFile, "content-file", "", "Markdown file used as article content source")
	return cmd
}

func runValidate(cmd *cobra.Command, args []string, opts validateOptions) error {
	platform, publishType, payloadPath := args[0], args[1], args[2]
	if err := detectSwappedPublishArgs(publishType, platform, "validate <platform> <type> <payload.json>"); err != nil {
		return err
	}
	rt, err := app.Load()
	if err != nil {
		return err
	}
	payload, err := readPayload(payloadPath)
	if err != nil {
		return err
	}
	publishType = publishmod.NormalizePublishType(publishType)
	service := publishflow.NewService(rt)
	prepared, err := service.Prepare(publishflow.ExecuteInput{
		PublishType:   publishType,
		PlatformInput: platform,
		Payload:       payload,
		ContentFile:   opts.ContentFile,
		FlagChannel:   opts.Channel,
		FlagClientID:  opts.ClientID,
	}, publishflow.PrepareOptions{TraceNormalizations: true, RemoteChecks: publishflow.RemoteChecksCloudWithKey})
	if err != nil {
		return err
	}
	if err := publishflow.AssertShoppingCartEntitlements(rt.Client, prepared.Payload, prepared.Platforms...); err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), "validate", map[string]interface{}{
		"platform":     platform,
		"type":         publishType,
		"valid":        true,
		"prepared":     true,
		"request":      prepared.PublishBody,
		"remoteChecks": prepared.RemoteChecked,
		"nextStep":     publishNextCommand(publishType, platform, payloadPath, prepared.PublishMode, prepared.ClientID, true, opts.ContentFile),
	})
}

func publishNextCommand(publishType, platform, payloadPath, channel, clientID string, dryRun bool, contentFile ...string) string {
	parts := []string{"yxer", "publish", publishType, quoteCommandArg(platform), quoteCommandArg(payloadPath)}
	if len(contentFile) > 0 && strings.TrimSpace(contentFile[0]) != "" {
		parts = append(parts, "--content-file", quoteCommandArg(contentFile[0]))
	}
	if strings.TrimSpace(channel) == "local" {
		parts = append(parts, "--publish-channel", "local")
		if strings.TrimSpace(clientID) != "" {
			parts = append(parts, "--client-id", quoteCommandArg(clientID))
		}
	}
	if dryRun {
		parts = append(parts, "--dry-run")
	}
	return strings.Join(parts, " ")
}

func validateNextCommand(platform, publishType, payloadPath, channel, clientID string) string {
	parts := []string{"yxer", "validate", quoteCommandArg(platform), publishType, quoteCommandArg(payloadPath)}
	if strings.TrimSpace(channel) == "local" {
		parts = append(parts, "--publish-channel", "local")
		if strings.TrimSpace(clientID) != "" {
			parts = append(parts, "--client-id", quoteCommandArg(clientID))
		}
	}
	return strings.Join(parts, " ")
}

func quoteCommandArg(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	if strings.ContainsAny(value, " \t\"'") {
		return strconv.Quote(value)
	}
	return value
}

func readPayload(path string) (map[string]interface{}, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, yxerrors.Usage("Payload file not found", map[string]interface{}{
				"path": abs,
			}).
				WithCategory("file_not_found").
				WithHint("请检查 payload 文件路径；相对路径按当前工作目录解析，或先用对应的 init/form export 命令生成文件。")
		}
		if os.IsPermission(err) {
			return nil, yxerrors.Usage("Payload file is not readable", map[string]interface{}{
				"path": abs,
			}).
				WithCategory("file_not_readable").
				WithHint("请检查文件权限，确认当前用户可以读取该 payload 文件。")
		}
		return nil, err
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimPrefix(string(raw), "\uFEFF")), &payload); err != nil {
		return nil, yxerrors.Usage("Invalid JSON payload", err.Error()).
			WithHint("请检查 JSON 文件格式，确认没有多余逗号、注释或截断内容。")
	}
	return payload, nil
}
