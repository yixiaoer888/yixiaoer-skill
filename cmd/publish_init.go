package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/core/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var publishInitOutput string

func init() {
	publishInitCmd.Flags().StringVarP(&publishInitOutput, "output", "o", "payload.json", "output payload file")
	publishCmd.AddCommand(publishInitCmd)
}

var publishInitCmd = &cobra.Command{
	Use:   "init <中文平台名|platform-key> <type>",
	Short: "生成可编辑的发布 payload 模板",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		schemaDoc, err := schema.NewValidator(cfg.SchemaDir).Schema(args[0], args[1])
		if err != nil {
			return yxerrors.Usage("schema not found", map[string]interface{}{
				"platform": args[0],
				"type":     args[1],
			}).WithHint("请先用 yxer schema list 确认平台和类型。").
				WithNextCommand("yxer schema list")
		}

		template := buildPayloadTemplate(schemaDoc)
		raw, err := json.MarshalIndent(template, "", "  ")
		if err != nil {
			return err
		}
		outputPath := strings.TrimSpace(publishInitOutput)
		if outputPath == "" {
			outputPath = "payload.json"
		}
		absPath, err := filepath.Abs(outputPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(absPath, append(raw, '\n'), 0o644); err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "publish.init", map[string]interface{}{
			"platform": args[0],
			"type":     args[1],
			"file":     filepath.ToSlash(absPath),
			"payload":  template,
			"next": []string{
				"填写 platformAccountId 和 contentPublishForm 业务字段",
				"如有视频/封面，先执行 yxer upload 获取真实资源元数据",
				"执行 yxer validate <platform> <type> <payload.json>",
			},
			"platformName": platformutil.ChineseName(schemaDoc.Platform),
		})
	},
}
