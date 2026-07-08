package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func newAccountsUpdateCmd() *cobra.Command {
	return newUpdateAccountCommand("update <account_id>", false)
}

func newUpdateAccountCompatCmd() *cobra.Command {
	return newUpdateAccountCommand("update-account <account_id>", true)
}

func newUpdateAccountCommand(use string, hidden bool) *cobra.Command {
	opts := updateAccountOptions{}
	cmd := &cobra.Command{
		Use:    use,
		Short:  "更新账号代理或备注",
		Hidden: hidden,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := opts.body()
			if len(body) == 0 {
				return yxerrors.Usage("update account request must not be empty", nil).
					WithHint("请至少传入 --proxy-id、--kuaidaili-area、--remark 或 --group。")
			}
			if opts.DryRun {
				return output.Success(cmd.OutOrStdout(), "accounts.update.dry-run", map[string]interface{}{
					"dryRun":  true,
					"account": args[0],
					"request": body,
				})
			}
			rt, err := app.Load()
			if err != nil {
				return err
			}
			result, err := queryflow.NewService(rt).UpdateAccount(args[0], body)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "accounts.update", result)
		},
	}
	cmd.Flags().StringVar(&opts.ProxyID, "proxy-id", "", "team proxy id")
	cmd.Flags().StringVar(&opts.KuaidailiArea, "kuaidaili-area", "", "built-in proxy area code")
	cmd.Flags().StringVar(&opts.Remark, "remark", "", "account remark")
	cmd.Flags().StringSliceVar(&opts.Groups, "group", nil, "group id; repeat or comma-separate for multiple")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview update request without performing the write")
	return cmd
}

type updateAccountOptions struct {
	ProxyID       string
	KuaidailiArea string
	Remark        string
	Groups        []string
	DryRun        bool
}

func (opts updateAccountOptions) body() map[string]interface{} {
	body := map[string]interface{}{}
	if strings.TrimSpace(opts.ProxyID) != "" {
		body["proxyId"] = opts.ProxyID
	}
	if strings.TrimSpace(opts.KuaidailiArea) != "" {
		body["kuaidailiArea"] = opts.KuaidailiArea
	}
	if strings.TrimSpace(opts.Remark) != "" {
		body["remark"] = opts.Remark
	}
	if len(opts.Groups) > 0 {
		groups := make([]string, 0, len(opts.Groups))
		for _, group := range opts.Groups {
			if strings.TrimSpace(group) != "" {
				groups = append(groups, group)
			}
		}
		if len(groups) > 0 {
			body["groups"] = groups
		}
	}
	return body
}
