package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/linkedapp"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var linkedAppAccountID string
var linkedAppAccountName string

func init() {
	linkedAppCmd.AddCommand(linkedAppStatusCmd)
	linkedAppCmd.AddCommand(linkedAppConnectCmd)
	linkedAppCmd.AddCommand(linkedAppDisconnectCmd)
	linkedAppCmd.AddCommand(linkedAppToggleCmd)
	linkedAppConnectCmd.Flags().StringVar(&linkedAppAccountID, "account-id", "", "linked app account id")
	linkedAppConnectCmd.Flags().StringVar(&linkedAppAccountName, "account-name", "", "linked app account name")
	rootCmd.AddCommand(linkedAppCmd)
}

var linkedAppCmd = &cobra.Command{
	Use:   "linked-app",
	Short: "管理 yixiaoer 链接应用状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLinkedAppStatus(cmd)
	},
}

var linkedAppStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看链接应用状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLinkedAppStatus(cmd)
	},
}

var linkedAppConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "标记链接应用为已连接",
	RunE: func(cmd *cobra.Command, args []string) error {
		if linkedAppAccountID == "" && linkedAppAccountName == "" {
			return yxerrors.Usage("account-id or account-name is required", nil).
				WithHint("请至少提供 --account-id 或 --account-name，用于标识当前连接到的蚁小二账号。")
		}
		configPath, err := config.SaveLinkedAppState("yixiaoer", linkedAppAccountID, linkedAppAccountName, true)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "linked-app.connect", map[string]interface{}{
			"configPath": configPath,
			"metadata":   linkedapp.DefaultMetadata(),
			"state": map[string]interface{}{
				"appId":       "yixiaoer",
				"connected":   true,
				"accountId":   linkedAppAccountID,
				"accountName": linkedAppAccountName,
			},
		})
	},
}

var linkedAppDisconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "标记链接应用为未连接",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := config.SaveLinkedAppState("yixiaoer", "", "", false)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "linked-app.disconnect", map[string]interface{}{
			"configPath": configPath,
			"metadata":   linkedapp.DefaultMetadata(),
			"state": map[string]interface{}{
				"appId":     "yixiaoer",
				"connected": false,
			},
		})
	},
}

var linkedAppToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "切换链接应用连接状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.LinkedApp.Connected {
			return linkedAppDisconnectCmd.RunE(cmd, args)
		}
		if linkedAppAccountID == "" && linkedAppAccountName == "" {
			return yxerrors.Usage("linked app is disconnected", nil).
				WithHint("当前是未连接状态；请改用 yxer linked-app connect 并传入 --account-id 或 --account-name 完成连接。")
		}
		return linkedAppConnectCmd.RunE(cmd, args)
	},
}

func runLinkedAppStatus(cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), "linked-app.status", map[string]interface{}{
		"metadata": linkedapp.DefaultMetadata(),
		"state":    cfg.LinkedApp,
	})
}
