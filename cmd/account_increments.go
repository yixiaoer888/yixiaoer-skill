package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type accountIncrementsOptions struct {
	StartDate string
	EndDate   string
	GroupID   string
	Platform  string
	Name      string
}

func newAccountIncrementsCmd() *cobra.Command {
	opts := accountIncrementsOptions{}
	cmd := &cobra.Command{
		Use:   "account-increments",
		Short: "查询账号增量数据",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(opts.StartDate) == "" || strings.TrimSpace(opts.EndDate) == "" {
				return yxerrors.Usage("account increments start-date and end-date are required", nil).
					WithHint("请同时传入 --start-date 和 --end-date，格式为 YYYY-MM-DD。")
			}
			startTime, endTime, err := api.ShanghaiDateRange(opts.StartDate, opts.EndDate)
			if err != nil {
				return err
			}
			return runQuery(cmd, "account-increments", func(service queryflow.Service) (interface{}, error) {
				return service.AccountIncrements(api.AccountIncrementOptions{
					StartTime: startTime,
					EndTime:   endTime,
					GroupID:   opts.GroupID,
					Platform:  opts.Platform,
					Name:      opts.Name,
				})
			})
		},
	}
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "inclusive start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "inclusive end date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.GroupID, "group-id", "", "account group ID")
	cmd.Flags().StringVar(&opts.Platform, "platform", "", "platform name, for example xhs or 小红书")
	cmd.Flags().StringVar(&opts.Name, "account-name", "", "account name keyword")
	return cmd
}
