package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
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
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "创建账号分组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := buildCreateAccountGroupBody(args[0])
			if err != nil {
				return err
			}
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "account-group.create.dry-run", map[string]interface{}{
					"dryRun":  true,
					"request": body,
				})
			}

			rt, err := app.Load()
			if err != nil {
				return err
			}
			result, err := queryflow.NewService(rt).CreateAccountGroup(body)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "account-group.create", result)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview create account group request without performing the write")
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

func buildCreateAccountGroupBody(name string) (map[string]interface{}, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, yxerrors.Usage("account group name must not be empty", nil).
			WithHint("请传入非空分组名称，例如 yxer account-group create 核心账号组 --dry-run。")
	}
	return map[string]interface{}{
		"name": trimmed,
	}, nil
}

func newAccountGroupUpdateCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "update <group_id> <name>",
		Short: "更新账号分组",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID := strings.TrimSpace(args[0])
			if groupID == "" {
				return yxerrors.Usage("account group id must not be empty", nil).
					WithHint("请传入有效的分组 ID，例如 yxer account-group update group_1 核心账号组 --dry-run。")
			}
			body, err := buildCreateAccountGroupBody(args[1])
			if err != nil {
				return err
			}
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "account-group.update.dry-run", map[string]interface{}{
					"dryRun":  true,
					"groupId": groupID,
					"request": body,
				})
			}

			rt, err := app.Load()
			if err != nil {
				return err
			}
			result, err := queryflow.NewService(rt).UpdateAccountGroup(groupID, body)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "account-group.update", result)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview update account group request without performing the write")
	return cmd
}

func newAccountGroupDeleteCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "delete <group_id>",
		Short: "删除账号分组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID := strings.TrimSpace(args[0])
			if groupID == "" {
				return yxerrors.Usage("account group id must not be empty", nil).
					WithHint("请传入有效的分组 ID，例如 yxer account-group delete group_1 --dry-run。")
			}
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "account-group.delete.dry-run", map[string]interface{}{
					"dryRun":  true,
					"groupId": groupID,
				})
			}

			rt, err := app.Load()
			if err != nil {
				return err
			}
			result, err := queryflow.NewService(rt).DeleteAccountGroup(groupID)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "account-group.delete", result)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview delete account group request without performing the write")
	return cmd
}
