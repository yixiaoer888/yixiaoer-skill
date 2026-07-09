package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/cmdflow"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const (
	accountGroupVisibleScopeAll      = "all"
	accountGroupVisibleScopeSpecific = "specific"
)

func init() {
	rootCmd.AddCommand(newAccountGroupCmd())
}

func newAccountGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-group",
		Short: "管理账号分组",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newAccountGroupListCmd())
	cmd.AddCommand(newAccountGroupCreateCmd())
	cmd.AddCommand(newAccountGroupUpdateCmd())
	cmd.AddCommand(newAccountGroupDeleteCmd())
	return cmd
}

func newAccountGroupCreateCmd() *cobra.Command {
	var dryRun bool
	var opts accountGroupUpdateOptions
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "创建账号分组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var body map[string]interface{}

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					built, err := buildAccountGroupBody(args[0], opts)
					if err != nil {
						return err
					}
					body = built
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "account-group.create.dry-run",
						Data: map[string]interface{}{
							"dryRun":  true,
							"request": body,
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := queryflow.NewService(rt).CreateAccountGroup(body)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "account-group.create", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview create account group request without performing the write")
	cmd.Flags().StringVar(&opts.VisibleScope, "visible-scope", "", "visible scope: all or specific")
	cmd.Flags().StringSliceVar(&opts.VisibleUsers, "visible-user", nil, "visible user id; repeat or comma-separate for multiple when visible-scope is specific")
	return cmd
}

func newAccountGroupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "列出账号分组",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "account-group.list", func(service queryflow.Service) (interface{}, error) {
				return service.AccountGroups()
			})
		},
	}
}

func buildAccountGroupBody(name string, opts accountGroupUpdateOptions) (map[string]interface{}, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, yxerrors.Usage("account group name must not be empty", nil).
			WithHint("请传入非空分组名称，例如 yxer account-group create 核心账号组 --dry-run。")
	}
	body := map[string]interface{}{
		"name": trimmed,
	}

	scope := strings.ToLower(strings.TrimSpace(opts.VisibleScope))
	if scope == "" {
		return body, nil
	}
	if scope != accountGroupVisibleScopeAll && scope != accountGroupVisibleScopeSpecific {
		return nil, yxerrors.Usage("account group visible scope must be all or specific", map[string]interface{}{
			"visibleScope": opts.VisibleScope,
		}).WithHint("请传入 --visible-scope all 或 --visible-scope specific。")
	}

	body["visibleScope"] = scope
	if scope == accountGroupVisibleScopeAll {
		return body, nil
	}

	users := make([]string, 0, len(opts.VisibleUsers))
	for _, user := range opts.VisibleUsers {
		trimmed := strings.TrimSpace(user)
		if trimmed != "" {
			users = append(users, trimmed)
		}
	}
	if len(users) == 0 {
		return nil, yxerrors.Usage("account group visibleUsers must not be empty when visibleScope is specific", map[string]interface{}{
			"visibleScope": scope,
		}).WithHint("当使用 --visible-scope specific 时，请至少传入一个 --visible-user <userId>。")
	}
	body["visibleUsers"] = users
	return body, nil
}

type accountGroupUpdateOptions struct {
	VisibleScope string
	VisibleUsers []string
}

func newAccountGroupUpdateCmd() *cobra.Command {
	var dryRun bool
	var opts accountGroupUpdateOptions
	cmd := &cobra.Command{
		Use:   "update <group_id> <name>",
		Short: "更新账号分组",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var groupID string
			var body map[string]interface{}

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					groupID = strings.TrimSpace(args[0])
					if groupID == "" {
						return yxerrors.Usage("account group id must not be empty", nil).
							WithHint("请传入有效的分组 ID，例如 yxer account-group update group_1 核心账号组 --dry-run。")
					}
					built, err := buildAccountGroupBody(args[1], opts)
					if err != nil {
						return err
					}
					body = built
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "account-group.update.dry-run",
						Data: map[string]interface{}{
							"dryRun":  true,
							"groupId": groupID,
							"request": body,
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := queryflow.NewService(rt).UpdateAccountGroup(groupID, body)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "account-group.update", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview update account group request without performing the write")
	cmd.Flags().StringVar(&opts.VisibleScope, "visible-scope", "", "visible scope: all or specific")
	cmd.Flags().StringSliceVar(&opts.VisibleUsers, "visible-user", nil, "visible user id; repeat or comma-separate for multiple when visible-scope is specific")
	return cmd
}

func newAccountGroupDeleteCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "delete <group_id>",
		Short: "删除账号分组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var groupID string

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					groupID = strings.TrimSpace(args[0])
					if groupID == "" {
						return yxerrors.Usage("account group id must not be empty", nil).
							WithHint("请传入有效的分组 ID，例如 yxer account-group delete group_1 --dry-run。")
					}
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "account-group.delete.dry-run",
						Data: map[string]interface{}{
							"dryRun":  true,
							"groupId": groupID,
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := queryflow.NewService(rt).DeleteAccountGroup(groupID)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "account-group.delete", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview delete account group request without performing the write")
	return cmd
}
