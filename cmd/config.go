package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetLocalClientIDCmd)
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
