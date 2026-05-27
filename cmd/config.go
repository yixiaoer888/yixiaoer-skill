package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
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
