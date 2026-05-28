package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/schema"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

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
		result := schema.NewValidator(cfg.SchemaDir).Validate(platform, publishType, payload)
		if !result.Valid {
			return yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint("请根据 schema 要求修正 payload 字段、类型和必填项。").
				WithNextCommand("yxer schema get <platform> <type>")
		}
		if _, hasAccountForms := payload["accountForms"]; hasAccountForms || payload["publishArgs"] != nil {
			preflight := publishmod.Preflight(publishType, []string{platform}, payload)
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
