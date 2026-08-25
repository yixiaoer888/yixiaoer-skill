package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/cmdflow"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func newPublishDeleteCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "delete <task_id>",
		Short: "删除已发布作品",
		Long:  "删除指定平台任务对应的已发布作品。task_id 使用 yxer query details <task_set_id> 返回任务项的 id；先执行 --dry-run 确认目标。",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var taskID string

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					taskID = strings.TrimSpace(args[0])
					if taskID == "" {
						return yxerrors.Usage("published task id must not be empty", nil).
							WithHint("请传入 yxer query details 返回任务项的 id，例如 yxer publish delete task_1 --dry-run。")
					}
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "publish.delete.dry-run",
						Data: map[string]interface{}{
							"dryRun": true,
							"taskId": taskID,
							"request": map[string]string{
								"method": "DELETE",
								"path":   "/tasks/" + taskID + "/publish",
							},
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := publishflow.NewService(rt).DeletePublishedTask(taskID)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{
						Action: "publish.delete",
						Data: map[string]interface{}{
							"taskId":   taskID,
							"response": result,
						},
					}, nil
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview delete published work request without performing the write")
	return cmd
}
