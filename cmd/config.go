package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(newConfigCmd())
}

type configInitOptions struct {
	APIKey        string
	LocalClientID string
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "管理 CLI 本地配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigInitCmd())
	cmd.AddCommand(newConfigSetAPIKeyCmd())
	cmd.AddCommand(newConfigSetLocalClientIDCmd())
	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
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
}

func newConfigInitCmd() *cobra.Command {
	opts := configInitOptions{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "初始化 CLI 配置",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.APIKey == "" {
				return yxerrors.Usage("apiKey must not be empty", nil).
					WithHint("请传入 --api-key 完成 yxer CLI 初始化。")
			}
			configPath, err := config.SaveAPIKey(opts.APIKey)
			if err != nil {
				return err
			}

			if opts.LocalClientID != "" {
				configPath, err = config.SaveLocalClientID(opts.LocalClientID)
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
	cmd.Flags().StringVar(&opts.APIKey, "api-key", "", "api key for yxer cli init")
	cmd.Flags().StringVar(&opts.LocalClientID, "local-client-id", "", "default local publish client id")
	return cmd
}

func newConfigSetLocalClientIDCmd() *cobra.Command {
	return &cobra.Command{
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
}

func newConfigSetAPIKeyCmd() *cobra.Command {
	return &cobra.Command{
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
}
