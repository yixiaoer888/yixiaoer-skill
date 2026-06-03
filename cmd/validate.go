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
			// 增强错误提示
			suggestions := analyzeValidationErrors(result.Errors, platform, publishType)

			return yxerrors.Usage("Schema validation failed", map[string]interface{}{
				"errors":      result.Errors,
				"suggestions": suggestions,
				"checklist": []string{
					"✓ 已执行 yxer schema fields " + platform + " " + publishType + "?",
					"✓ payload 顶层包含 publishArgs?",
					"✓ 业务字段在 publishArgs.accountForms[].contentPublishForm?",
					"✓ 资源已通过 yxer upload 上传并使用返回的完整对象?",
					"✓ 复杂对象（location/music等）已通过查询命令获取?",
				},
			}).
				WithHint("根据上方 suggestions 修正字段，或查看 checklist 确认流程是否正确。").
				WithNextCommand(fmt.Sprintf("yxer schema fields %s %s", platform, publishType))
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
			"nextStep": fmt.Sprintf("yxer publish %s %s %s --dry-run", publishType, platform, payloadPath),
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
