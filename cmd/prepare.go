package cmd

import (
	"github.com/spf13/cobra"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
)

func newPrepareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prepare <platform> <type>",
		Short: "获取发布前置数据",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "prepare", func(service queryflow.Service) (interface{}, error) {
				return service.Prepare(args[0], args[1])
			})
		},
	}
}
