package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
)

var accountsName string
var accountsStatus int

func init() {
	accountsCmd.Flags().StringVar(&accountsName, "name", "", "filter by name")
	accountsCmd.Flags().IntVar(&accountsStatus, "status", -1, "filter by status")
	rootCmd.AddCommand(accountsCmd)
}

var accountsCmd = &cobra.Command{
	Use:   "accounts [中文平台名]",
	Short: "查询账号",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		platform := ""
		if len(args) > 0 {
			platform = args[0]
		}
		accounts, err := api.NewClient(cfg).Accounts(platform)
		if err != nil {
			return err
		}
		filtered := filterAccounts(accounts, accountsName, accountsStatus)
		if wantJSON(cmd) {
			return output.Success(cmd.OutOrStdout(), "accounts", filtered)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "账号列表")
		if platform != "" {
			fmt.Fprintf(cmd.OutOrStdout(), " (%s)", platform)
		}
		fmt.Fprintln(cmd.OutOrStdout(), ":")
		for i, account := range filtered {
			icon := "offline"
			if api.AccountStatus(account) == 1 {
				icon = "online"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s (%s) %s\n", i+1, accountName(account), api.AccountID(account), icon)
		}
		if len(filtered) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  (无在线账号)")
		}
		return nil
	},
}

func filterAccounts(accounts []map[string]interface{}, name string, status int) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, account := range accounts {
		if name != "" && !strings.Contains(accountName(account), name) {
			continue
		}
		if status >= 0 && api.AccountStatus(account) != status {
			continue
		}
		filtered = append(filtered, account)
	}
	return filtered
}

func accountName(account map[string]interface{}) string {
	for _, key := range []string{"platformAccountName", "name", "nickname", "remark", "platformAuthorId"} {
		if value := fmt.Sprint(account[key]); value != "" && value != "<nil>" {
			return value
		}
	}
	return "未命名"
}
