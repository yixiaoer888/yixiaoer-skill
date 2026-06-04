package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	accountsflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/accounts"
)

var accountsName string
var accountsStatus int

func init() {
	accountsCmd.PersistentFlags().StringVar(&accountsName, "name", "", "filter by name")
	accountsCmd.PersistentFlags().IntVar(&accountsStatus, "status", -1, "filter by status")
	accountsCmd.AddCommand(accountsListCmd)
	rootCmd.AddCommand(accountsCmd)
}

var accountsCmd = &cobra.Command{
	Use:   "accounts [中文平台名]",
	Short: "查询账号",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAccountsList(cmd, args)
	},
}

var accountsListCmd = &cobra.Command{
	Use:     "list [中文平台名]",
	Short:   "列出账号",
	Aliases: []string{"ls"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAccountsList(cmd, args)
	},
}

func runAccountsList(cmd *cobra.Command, args []string) error {
	platform := ""
	if len(args) > 0 {
		platform = args[0]
	}
	filtered, err := accountsflow.NewService().List(platform, accountsName, accountsStatus)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return output.Success(cmd.OutOrStdout(), "accounts.list", filtered)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "账号列表")
	if platform != "" {
		fmt.Fprintf(cmd.OutOrStdout(), " (%s)", platform)
	}
	fmt.Fprintln(cmd.OutOrStdout(), ":")
	for i, account := range filtered {
		icon := "offline"
		if client.AccountStatus(account) == 1 {
			icon = "online"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s (%s) %s\n", i+1, accountName(account), client.AccountID(account), icon)
	}
	if len(filtered) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  (无在线账号)")
	}
	return nil
}

func filterAccounts(accounts []map[string]interface{}, name string, status int) []map[string]interface{} {
	return accountsflow.FilterAccounts(accounts, name, status)
}

func accountName(account map[string]interface{}) string {
	return accountsflow.AccountName(account)
}
