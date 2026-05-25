package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
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
			return yxerrors.Usage("Schema validation failed", result.Errors)
		}
		if _, hasAccountForms := payload["accountForms"]; hasAccountForms || payload["publishArgs"] != nil {
			preflight := publishmod.Preflight(publishType, []string{platform}, payload)
			if len(preflight.Errors) > 0 {
				return yxerrors.Usage("Publish preflight failed", preflight.Errors)
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
		return nil, yxerrors.Usage("Invalid JSON payload", err.Error())
	}
	return payload, nil
}
