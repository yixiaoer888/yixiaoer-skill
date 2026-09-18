package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/cmdflow"
	draftflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/draft"
)

func init() {
	rootCmd.AddCommand(newDraftCmd())
}

func newDraftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draft",
		Short: "管理蚁小二草稿",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newDraftSaveCmd())
	cmd.AddCommand(newDraftImportCmd())
	return cmd
}

func newDraftImportCmd() *cobra.Command {
	var (
		coverPath      string
		accountID      string
		draftID        string
		title          string
		digest         string
		author         string
		originalAuthor string
		original       bool
		dryRun         bool
	)
	cmd := &cobra.Command{
		Use:   "import <article.docx>",
		Short: "导入 DOCX 为微信公众号文章草稿",
		Long:  "读取 DOCX 正文和内嵌图片，上传封面与正文图片并保存为蚁小二内部草稿。封面必须是横版；默认不声明原创。",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := draftflow.ImportInput{
				DocxPath:       args[0],
				CoverPath:      coverPath,
				AccountID:      accountID,
				DraftID:        draftID,
				Title:          title,
				Digest:         digest,
				Author:         author,
				OriginalAuthor: originalAuthor,
				Original:       original,
			}
			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					_, err := draftflow.PreviewImport(input)
					return err
				},
				DryRun: func() (cmdflow.Result, error) {
					preview, err := draftflow.PreviewImport(input)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "draft.import.dry-run", Data: map[string]interface{}{
						"dryRun":  true,
						"request": preview,
					}}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := draftflow.NewService(rt).ImportDOCX(input)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "draft.import", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().StringVar(&coverPath, "cover", "", "horizontal cover image path (required)")
	cmd.Flags().StringVar(&accountID, "account-id", "", "online WeChat account ID from yxer accounts list (required)")
	cmd.Flags().StringVar(&draftID, "draft-id", "", "existing Yixiaoer draft ID to update instead of creating a new draft")
	cmd.Flags().StringVar(&title, "title", "", "override the article title extracted from DOCX")
	cmd.Flags().StringVar(&digest, "digest", "", "override the article digest")
	cmd.Flags().StringVar(&author, "author", "", "article author name")
	cmd.Flags().StringVar(&originalAuthor, "original-author", "", "original author name; used only with --original")
	cmd.Flags().BoolVar(&original, "original", false, "declare the article as original; omitted means not original")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the import without uploading or saving")
	return cmd
}

func newDraftSaveCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "save <payload.json>",
		Short: "保存为蚁小二草稿",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var payload map[string]interface{}

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					loaded, err := readPayload(args[0])
					if err != nil {
						return err
					}
					payload = loaded
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "draft.save.dry-run",
						Data: map[string]interface{}{
							"dryRun":  true,
							"request": draftflow.PreviewSave(payload),
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := draftflow.NewService(rt).Save(payload)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "draft.save", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the draft payload without performing the write")
	return cmd
}
