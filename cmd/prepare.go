package cmd

import (
	"github.com/spf13/cobra"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
)

func newPrepareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prepare <platform> <type>",
		Short: "获取发布前置数据",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			publishType := "video"
			if len(args) > 1 {
				publishType = args[1]
			}
			return runQuery(cmd, "prepare", func(service queryflow.Service) (interface{}, error) {
				return service.Prepare(args[0], publishType)
			})
		},
	}
}
