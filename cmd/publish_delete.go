package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/cmdflow"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
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
								"path":   api.PublishedTaskDeleteEndpoint(taskID),
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
	cmd.AddCommand(newPublishDeletePreviewCmd())
	cmd.AddCommand(newPublishDeleteFromRecordCmd())
	return cmd
}

func newPublishDeletePreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview <task_set_id>",
		Short: "预览发布记录中的作品",
		Long:  "展示发布记录中的作品信息和可选序号，不显示任务 ID。使用 from-record 和 --index 选择要删除的作品。",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskSetID := strings.TrimSpace(args[0])
			if taskSetID == "" {
				return yxerrors.Usage("task set id must not be empty", nil).
					WithHint("请传入 yxer query records 返回的 id，例如 yxer publish delete preview task_set_1。")
			}
			preview, err := loadPublishedWorkPreview(taskSetID)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "publish.delete.preview", publicPublishedWorkPreview(preview))
		},
	}
	return cmd
}

func newPublishDeleteFromRecordCmd() *cobra.Command {
	var index int
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "from-record <task_set_id>",
		Short: "按发布记录中的作品序号删除",
		Long:  "先用 preview 查看作品列表，再用 --index 选择作品。序号从 1 开始；CLI 会在内部解析任务 ID。",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var taskSetID string
			var selected map[string]interface{}

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					taskSetID = strings.TrimSpace(args[0])
					if taskSetID == "" {
						return yxerrors.Usage("task set id must not be empty", nil).
							WithHint("请先运行 yxer publish delete preview <task_set_id> 查看作品。")
					}
					if index < 1 {
						return yxerrors.Usage("published work index must be at least 1", map[string]interface{}{"index": index}).
							WithHint("请先运行 yxer publish delete preview <task_set_id>，再传入其中的序号，例如 --index 1。")
					}
					preview, err := loadPublishedWorkPreview(taskSetID)
					if err != nil {
						return err
					}
					selected, err = selectPublishedWork(preview, index)
					if err != nil {
						return err
					}
					if selected["status"] == "已删除" {
						return yxerrors.Usage("selected published work is already deleted", map[string]interface{}{"work": publicPublishedWork(selected)}).
							WithHint("请选择状态不是“已删除”的作品。")
					}
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "publish.delete.from-record.dry-run",
						Data: map[string]interface{}{
							"dryRun":    true,
							"taskSetId": taskSetID,
							"work":      publicPublishedWork(selected),
							"request": map[string]string{
								"method": "DELETE",
								"path":   api.PublishedTaskDeleteEndpoint(firstStringField(selected, "taskId")),
							},
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := publishflow.NewService(rt).DeletePublishedTask(firstStringField(selected, "taskId"))
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{
						Action: "publish.delete.from-record",
						Data: map[string]interface{}{
							"taskSetId": taskSetID,
							"work":      publicPublishedWork(selected),
							"response":  result,
						},
					}, nil
				},
			})
		},
	}
	cmd.Flags().IntVar(&index, "index", 0, "1-based work index from publish delete preview")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview delete published work request without performing the write")
	return cmd
}

func loadPublishedWorkPreview(taskSetID string) (map[string]interface{}, error) {
	rt, err := app.Load()
	if err != nil {
		return nil, err
	}
	details, err := queryflow.NewService(rt).Details(taskSetID)
	if err != nil {
		return nil, err
	}
	return publishedWorkPreview(taskSetID, details)
}

func publishedWorkPreview(taskSetID string, details interface{}) (map[string]interface{}, error) {
	payload, ok := details.(map[string]interface{})
	if !ok {
		return nil, yxerrors.Remote("published work details response was not an object", details).
			WithHint("请重新执行 yxer query details <task_set_id> 检查发布记录。")
	}
	rawTasks, ok := payload["tasks"].([]interface{})
	if !ok || len(rawTasks) == 0 {
		return nil, yxerrors.Usage("published work record has no tasks", map[string]interface{}{"taskSetId": taskSetID}).
			WithHint("请使用 yxer query records 查询其他发布记录。")
	}
	works := make([]interface{}, 0, len(rawTasks))
	for position, raw := range rawTasks {
		task, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		taskID := firstStringField(task, "id")
		if taskID == "" {
			continue
		}
		works = append(works, publishedWorkDisplay(task, position+1))
	}
	if len(works) == 0 {
		return nil, yxerrors.Remote("published work record did not contain selectable tasks", payload).
			WithHint("请重新执行 yxer query details <task_set_id> 检查发布记录。")
	}
	return map[string]interface{}{
		"taskSetId": taskSetID,
		"works":     works,
	}, nil
}

func publishedWorkDisplay(task map[string]interface{}, index int) map[string]interface{} {
	work := map[string]interface{}{
		"index":       index,
		"platform":    firstStringField(task, "platformName"),
		"accountName": firstStringField(task, "platformAccountName"),
		"publishType": firstStringField(task, "publishType"),
		"title":       firstStringField(task, "title"),
		"coverUrl":    firstStringField(task, "coverUrl"),
		"openUrl":     firstStringField(task, "openUrl"),
		"status":      publishedWorkStatus(task),
		"taskId":      firstStringField(task, "id"),
	}
	if work["title"] == "" {
		work["title"] = "未提供标题"
	}
	if errorMessage := firstStringField(task, "errorMessage"); errorMessage != "" {
		work["statusMessage"] = errorMessage
	}
	return work
}

func publishedWorkStatus(task map[string]interface{}) string {
	if strings.EqualFold(firstStringField(task, "taskStatus"), "deleted") {
		return "已删除"
	}
	switch strings.ToLower(firstStringField(task, "stageStatus")) {
	case "success":
		return "发布成功"
	case "doing", "pending":
		return "处理中"
	case "fail", "failed":
		return "发布失败"
	default:
		return "状态未知"
	}
}

func selectPublishedWork(preview map[string]interface{}, index int) (map[string]interface{}, error) {
	works, _ := preview["works"].([]interface{})
	if index > len(works) {
		return nil, yxerrors.Usage("published work index is out of range", map[string]interface{}{
			"index":          index,
			"availableWorks": publicPublishedWorks(works),
		}).WithHint("请使用 publish delete preview 输出中的有效序号。")
	}
	work, _ := works[index-1].(map[string]interface{})
	if firstStringField(work, "taskId") == "" {
		return nil, yxerrors.Remote("selected published work did not contain a task id", fmt.Sprintf("index=%d", index))
	}
	return work, nil
}

func publicPublishedWork(work map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{}, len(work))
	for key, value := range work {
		if key != "taskId" {
			copy[key] = value
		}
	}
	return copy
}

func publicPublishedWorks(works []interface{}) []interface{} {
	public := make([]interface{}, 0, len(works))
	for _, raw := range works {
		if work, ok := raw.(map[string]interface{}); ok {
			public = append(public, publicPublishedWork(work))
		}
	}
	return public
}

func publicPublishedWorkPreview(preview map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"taskSetId": preview["taskSetId"],
		"works":     publicPublishedWorks(preview["works"].([]interface{})),
	}
}
