package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
	accountsflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/accounts"
)

var accountsName string
var accountsStatus int
var accountsPage int
var accountsSize int
var accountsAll bool

func init() {
	accountsCmd.PersistentFlags().StringVar(&accountsName, "name", "", "filter by name")
	accountsCmd.PersistentFlags().IntVar(&accountsStatus, "status", -1, "filter by status")
	accountsCmd.PersistentFlags().IntVar(&accountsPage, "page", 1, "page number")
	accountsCmd.PersistentFlags().IntVar(&accountsSize, "size", 20, "page size")
	accountsCmd.PersistentFlags().BoolVar(&accountsAll, "all", false, "fetch all pages when remote pagination metadata allows it")
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
	if accountsPage <= 0 {
		return yxerrors.Usage("accounts page must be greater than 0", map[string]interface{}{"page": accountsPage}).
			WithCategory("invalid_input")
	}
	if accountsSize <= 0 {
		return yxerrors.Usage("accounts size must be greater than 0", map[string]interface{}{"size": accountsSize}).
			WithCategory("invalid_input")
	}

	platform := ""
	if len(args) > 0 {
		platform = args[0]
	}
	rt, err := app.Load()
	if err != nil {
		return err
	}
	filtered, err := accountsflow.NewService(rt).ListWithOptions(platform, accountsName, accountsStatus, accountsflow.ListOptions{
		Page: accountsPage,
		Size: accountsSize,
		All:  accountsAll,
	})
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), "accounts.list", filtered)
}

func filterAccounts(accounts []map[string]interface{}, name string, status int) []map[string]interface{} {
	return accountsflow.FilterAccounts(accounts, name, status)
}

func accountName(account map[string]interface{}) string {
	return accountsflow.AccountName(account)
}
