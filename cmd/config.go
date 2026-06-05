package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configSetAPIKeyCmd)
	configCmd.AddCommand(configSetLocalClientIDCmd)
	configInitCmd.Flags().StringVar(&configInitAPIKey, "api-key", "", "api key for yxer cli init")
	configInitCmd.Flags().StringVar(&configInitLocalClientID, "local-client-id", "", "default local publish client id")
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理 CLI 本地配置",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查看 CLI 配置",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "config.get", map[string]interface{}{
			"configPath":           cfg.ConfigPath,
			"localPublishClientId": cfg.LocalClientID,
			"apiUrl":               cfg.APIURL,
			"apiKeyPresent":        cfg.APIKey != "",
		})
	},
}

var configInitAPIKey string
var configInitLocalClientID string

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化 CLI 配置",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configInitAPIKey == "" {
			return yxerrors.Usage("apiKey must not be empty", nil).
				WithHint("请传入 --api-key 完成 yxer CLI 初始化。")
		}
		configPath, err := config.SaveAPIKey(configInitAPIKey)
		if err != nil {
			return err
		}

		if configInitLocalClientID != "" {
			configPath, err = config.SaveLocalClientID(configInitLocalClientID)
			if err != nil {
				return err
			}
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "config.init", map[string]interface{}{
			"configPath":           configPath,
			"apiKeyPresent":        true,
			"localPublishClientId": cfg.LocalClientID,
		})
	},
}

var configSetLocalClientIDCmd = &cobra.Command{
	Use:   "set-local-client-id <clientId>",
	Short: "设置本机发布默认 clientId",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "" {
			return yxerrors.Usage("clientId must not be empty", nil).
				WithHint("请传入有效的本机发布 clientId。")
		}
		configPath, err := config.SaveLocalClientID(args[0])
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "config.set-local-client-id", map[string]interface{}{
			"configPath":           configPath,
			"localPublishClientId": args[0],
		})
	},
}

var configSetAPIKeyCmd = &cobra.Command{
	Use:   "set-api-key <apiKey>",
	Short: "设置 CLI 默认 apiKey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "" {
			return yxerrors.Usage("apiKey must not be empty", nil).
				WithHint("请传入有效的 apiKey。")
		}
		configPath, err := config.SaveAPIKey(args[0])
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "config.set-api-key", map[string]interface{}{
			"configPath":    configPath,
			"apiKeyPresent": true,
		})
	},
}
