package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
	accountsflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/accounts"
)

func init() {
	rootCmd.AddCommand(newAccountsCmd())
}

type accountsListOptions struct {
	Name   string
	Status int
	Page   int
	Size   int
	All    bool
}

func newAccountsCmd() *cobra.Command {
	opts := accountsListOptions{
		Status: -1,
		Page:   1,
		Size:   20,
	}
	cmd := &cobra.Command{
		Use:   "accounts [中文平台名]",
		Short: "查询账号",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAccountsListWithOptions(cmd, args, opts)
		},
	}
	cmd.PersistentFlags().StringVar(&opts.Name, "name", "", "filter by name")
	cmd.PersistentFlags().IntVar(&opts.Status, "status", -1, "filter by status")
	cmd.PersistentFlags().IntVar(&opts.Page, "page", 1, "page number")
	cmd.PersistentFlags().IntVar(&opts.Size, "size", 20, "page size")
	cmd.PersistentFlags().BoolVar(&opts.All, "all", false, "fetch all pages when remote pagination metadata allows it")
	cmd.AddCommand(newAccountsListCmd(&opts))
	cmd.AddCommand(newAccountsUpdateCmd())
	return cmd
}

func newAccountsListCmd(opts *accountsListOptions) *cobra.Command {
	return &cobra.Command{
		Use:     "list [中文平台名]",
		Short:   "列出账号",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAccountsListWithOptions(cmd, args, *opts)
		},
	}
}

func runAccountsListWithOptions(cmd *cobra.Command, args []string, opts accountsListOptions) error {
	if opts.Page <= 0 {
		return yxerrors.Usage("accounts page must be greater than 0", map[string]interface{}{"page": opts.Page}).
			WithCategory("invalid_input")
	}
	if opts.Size <= 0 {
		return yxerrors.Usage("accounts size must be greater than 0", map[string]interface{}{"size": opts.Size}).
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
	filtered, err := accountsflow.NewService(rt).ListWithOptions(platform, opts.Name, opts.Status, accountsflow.ListOptions{
		Page: opts.Page,
		Size: opts.Size,
		All:  opts.All,
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
