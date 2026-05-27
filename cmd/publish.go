package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish <type> <中文平台名|platform-key> <payload.json> [clientId]",
	Short: "发布内容（单平台原子发布）",
	Args:  cobra.RangeArgs(3, 4),
	RunE: func(cmd *cobra.Command, args []string) error {
		publishType := args[0]
		platform, err := singlePlatform(args[1])
		if err != nil {
			return err
		}
		platforms := []string{platform}
		payload, err := readPayload(args[2])
		if err != nil {
			return err
		}
		publishArgs := publishmod.ExtractPublishArgs(payload)
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		validator := schema.NewValidator(cfg.SchemaDir)
		for _, platform := range platforms {
			result := validator.Validate(platform, publishType, payload)
			if !result.Valid {
				return yxerrors.Usage("Schema validation failed", result.Errors)
			}
		}
		preflight := publishmod.Preflight(publishType, platforms, payload)
		if len(preflight.Errors) > 0 {
			return yxerrors.Usage("Publish preflight failed", preflight.Errors)
		}

		client := api.NewClient(cfg)
		if err := assertAccountsOnline(client, platforms, preflight.AccountIDs); err != nil {
			return err
		}

		body := map[string]interface{}{
			"action":         "publish",
			"publishType":    publishType,
			"platforms":      platforms,
			"publishArgs":    publishArgs,
			"publishChannel": "cloud",
		}
		copyOptionalPublishFields(body, payload)
		if len(args) == 4 && args[3] != "" {
			body["publishChannel"] = "local"
			body["clientId"] = args[3]
		}
		result, err := client.Publish(body)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "publish", result)
	},
}

func copyOptionalPublishFields(dst, src map[string]interface{}) {
	for _, field := range []string{"cover", "coverKey", "taskSetId", "desc", "isDraft"} {
		if value, ok := src[field]; ok {
			dst[field] = value
		}
	}
	if channel, ok := src["publishChannel"]; ok {
		dst["publishChannel"] = channel
	}
	if clientID, ok := src["clientId"]; ok {
		dst["clientId"] = clientID
	}
}

func splitPlatforms(value string) []string {
	var platforms []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			platforms = append(platforms, platformutil.ChineseName(item))
		}
	}
	return platforms
}

func singlePlatform(value string) (string, error) {
	platforms := splitPlatforms(value)
	if len(platforms) != 1 {
		return "", yxerrors.Usage("publish supports exactly one platform per command", []string{
			`Use Agent orchestration for multi-platform publishing: run "yxer accounts", "yxer schema get", "yxer validate", and "yxer publish" once per platform.`,
			`Example: yxer publish image-text xhs xhs-payload.json; then yxer publish image-text zhihu zhihu-payload.json`,
		})
	}
	return platforms[0], nil
}

func assertAccountsOnline(client *api.Client, platforms []string, accountIDs []string) error {
	wanted := map[string]bool{}
	for _, id := range accountIDs {
		wanted[id] = true
	}
	found := map[string]map[string]interface{}{}
	for _, platform := range platforms {
		accounts, err := client.Accounts(platform)
		if err != nil {
			return err
		}
		for _, account := range accounts {
			id := api.AccountID(account)
			if wanted[id] {
				found[id] = account
			}
		}
	}
	var errors []string
	for id := range wanted {
		account, ok := found[id]
		if !ok {
			errors = append(errors, "account "+id+": not found in target platform account list")
			continue
		}
		if status := api.AccountStatus(account); status != 1 {
			errors = append(errors, fmt.Sprintf("account %s: status=%d; publish requires status=1", id, status))
		}
	}
	if len(errors) > 0 {
		return yxerrors.Usage("Account preflight failed", errors)
	}
	return nil
}
