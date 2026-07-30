package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func newPrepareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prepare <platform> <type>",
		Short: "获取发布前置数据",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := app.Load()
			if err != nil {
				return err
			}
			prepared, err := queryflow.NewService(rt).Prepare(args[0], args[1])
			if err != nil {
				return err
			}
			doc, err := schema.NewValidator(rt.Config.SchemaDir).Schema(args[0], args[1])
			if err != nil {
				return yxerrors.Usage("schema not found", map[string]interface{}{"platform": args[0], "type": args[1]}).
					WithHint("未找到对应平台和发布类型的 schema，请先查看支持的平台和类型列表。").
					WithNextCommand("yxer schema list")
			}
			prepared.Form = buildPublishFormContract(doc)
			return output.Success(cmd.OutOrStdout(), "prepare", prepared)
		},
	}
}
