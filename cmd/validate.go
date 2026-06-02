package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/schema"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	validateCmd.Flags().StringVar(&validateChannelFlag, "publish-channel", "", `publish channel: "cloud" or "local"`)
	validateCmd.Flags().StringVar(&validateClientID, "client-id", "", "client ID for local publish")
	rootCmd.AddCommand(validateCmd)
}

var (
	validateChannelFlag string
	validateClientID    string
)

var validateCmd = &cobra.Command{
	Use:   "validate <中文平台名> <type> <payload.json>",
	Short: "校验发布 Payload",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		platform, publishType, payloadPath := args[0], args[1], args[2]
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		payload, err := readPayload(payloadPath)
		if err != nil {
			return err
		}
		publishType = publishmod.NormalizePublishType(publishType)
		channel, clientID, err := publishflow.ResolvePublishMode(cfg, payload, "", validateChannelFlag, validateClientID)
		if err != nil {
			return err
		}
		result := schema.NewValidator(cfg.SchemaDir).Validate(platform, publishType, payload)
		if !result.Valid {
			return yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint("请根据 schema 要求修正 payload 字段、类型和必填项。").
				WithNextCommand("yxer schema fields <platform> <type>")
		}
		if err := requireStandardPublishPayload(payload, platform, publishType); err != nil {
			return err
		}
		if _, hasAccountForms := payload["accountForms"]; hasAccountForms || payload["publishArgs"] != nil {
			preflight := publishmod.Preflight(publishType, []string{platform}, payloadWithResolvedPublishMode(payload, channel, clientID))
			if len(preflight.Errors) > 0 {
				return yxerrors.Usage("Publish preflight failed", preflight.Errors).
					WithHint("请先完成资源上传，并确保 payload 中引用的是上传后的 key，而不是外部 URL。").
					WithNextCommand("yxer upload <file_path_or_url>")
			}
		}
		return output.Success(cmd.OutOrStdout(), "validate", map[string]interface{}{
			"platform": platform,
			"type":     publishType,
			"valid":    true,
		})
	},
}

func payloadWithResolvedPublishMode(payload map[string]interface{}, channel, clientID string) map[string]interface{} {
	withMode := make(map[string]interface{}, len(payload)+2)
	for key, value := range payload {
		withMode[key] = value
	}
	withMode["publishChannel"] = channel
	if clientID != "" {
		withMode["clientId"] = clientID
	} else {
		delete(withMode, "clientId")
	}
	return withMode
}

func requireStandardPublishPayload(payload map[string]interface{}, platform, publishType string) error {
	if _, ok := payload["publishArgs"].(map[string]interface{}); ok {
		return nil
	}
	return yxerrors.Usage("Standard publish payload is required", []string{
		fmt.Sprintf("platform=%s", platform),
		fmt.Sprintf("type=%s", publishType),
		"missing publishArgs",
	}).
		WithHint("请使用标准请求体：顶层保留 action/publishType/platforms/publishArgs，实际业务字段放到 publishArgs.accountForms[].contentPublishForm。").
		WithNextCommand(fmt.Sprintf("yxer schema fields %s %s", platform, publishType))
}

func readPayload(path string) (map[string]interface{}, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimPrefix(string(raw), "\uFEFF")), &payload); err != nil {
		return nil, yxerrors.Usage("Invalid JSON payload", err.Error()).
			WithHint("请检查 JSON 文件格式，确认没有多余逗号、注释或截断内容。")
	}
	return payload, nil
}
