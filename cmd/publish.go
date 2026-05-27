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

var (
	publishChannelFlag string
	publishClientID    string
)

func init() {
	publishCmd.Flags().StringVar(&publishChannelFlag, "publish-channel", "", `publish channel: "cloud" or "local"`)
	publishCmd.Flags().StringVar(&publishClientID, "client-id", "", "client ID for local publish")
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
		resolvedPayload := cloneMap(payload)
		channel, clientID, err := resolvePublishMode(cfg, resolvedPayload, args)
		if err != nil {
			return err
		}
		resolvedPayload["publishChannel"] = channel
		if clientID != "" {
			resolvedPayload["clientId"] = clientID
		} else {
			delete(resolvedPayload, "clientId")
		}

		validator := schema.NewValidator(cfg.SchemaDir)
		for _, platform := range platforms {
			result := validator.Validate(platform, publishType, resolvedPayload)
			if !result.Valid {
				return yxerrors.Usage("Schema validation failed", result.Errors)
			}
		}
		preflight := publishmod.Preflight(publishType, platforms, resolvedPayload)
		if len(preflight.Errors) > 0 {
			return yxerrors.Usage("Publish preflight failed", preflight.Errors)
		}

		client := api.NewClient(cfg)
		if err := assertAccountsOnline(client, platforms, preflight.AccountIDs); err != nil {
			return err
		}

		body := buildPublishBody(resolvedPayload, publishArgs, publishType, platforms)
		body["publishChannel"] = channel
		if clientID != "" {
			body["clientId"] = clientID
		} else {
			delete(body, "clientId")
		}
		result, err := client.Publish(body)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "publish", result)
	},
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return map[string]interface{}{}
	}
	dst := make(map[string]interface{}, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func buildPublishBody(payload, publishArgs map[string]interface{}, publishType string, platforms []string) map[string]interface{} {
	body := map[string]interface{}{
		"action":         "publish",
		"publishType":    publishType,
		"platforms":      platforms,
		"publishArgs":    publishArgs,
		"publishChannel": "cloud",
	}
	if _, ok := payload["publishArgs"].(map[string]interface{}); ok {
		for key, value := range payload {
			if key == "publishArgs" {
				body[key] = publishArgs
				continue
			}
			body[key] = value
		}
		body["action"] = "publish"
		body["publishType"] = publishType
		body["platforms"] = platforms
		if _, ok := body["publishChannel"]; !ok {
			body["publishChannel"] = "cloud"
		}
		return body
	}
	copyOptionalPublishFields(body, payload)
	return body
}

func resolvePublishMode(cfg config.Config, payload map[string]interface{}, args []string) (string, string, error) {
	channel := "cloud"
	clientID := ""
	if value, ok := payload["publishChannel"]; ok {
		channel = strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := payload["clientId"]; ok {
		clientID = strings.TrimSpace(fmt.Sprint(value))
	}
	if len(args) == 4 && strings.TrimSpace(args[3]) != "" {
		channel = "local"
		clientID = strings.TrimSpace(args[3])
	}
	if strings.TrimSpace(publishChannelFlag) != "" {
		channel = strings.TrimSpace(publishChannelFlag)
	}
	if strings.TrimSpace(publishClientID) != "" {
		clientID = strings.TrimSpace(publishClientID)
	}
	if clientID == "" {
		clientID = strings.TrimSpace(cfg.LocalClientID)
	}
	switch channel {
	case "", "cloud":
		channel = "cloud"
	case "local":
		if clientID == "" {
			return "", "", yxerrors.Usage(`clientId is required when publishChannel is "local"`, []string{
				`Run: yxer config set-local-client-id <clientId>`,
				`Or pass a fourth positional argument: yxer publish video <platform> payload.json <clientId>`,
				`Or pass flags: yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>`,
			})
		}
	default:
		return "", "", yxerrors.Usage(`publishChannel must be "cloud" or "local"`, []string{
			fmt.Sprintf("got %q", channel),
		})
	}
	return channel, clientID, nil
}

func copyOptionalPublishFields(dst, src map[string]interface{}) {
	for _, field := range []string{"cover", "coverKey", "taskSetId", "desc", "isDraft", "isAppContent"} {
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
