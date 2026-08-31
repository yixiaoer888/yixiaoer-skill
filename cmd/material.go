package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/cmdflow"
	materialflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/material"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(newMaterialCmd())
}

type materialAddOptions struct {
	FilePath  string
	ThumbPath string
	Type      string
	DryRun    bool
}

func newMaterialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "material",
		Short: "管理素材库",
		Long:  "管理素材库资源。\n1. create <payload.json>：提交素材登记 payload\n2. add --file ...：上传资源并登记到素材库\n3. list --name ...：按文件名查询素材及其 ID\n4. move <material_id> --group-id ...：移动素材到指定分组\n5. move-by-name <file_name> --group-id ...：按文件名匹配后移动素材",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newMaterialCreateCmd())
	cmd.AddCommand(newMaterialAddCmd())
	cmd.AddCommand(newMaterialListCmd())
	cmd.AddCommand(newMaterialMoveCmd())
	cmd.AddCommand(newMaterialMoveByNameCmd())
	cmd.AddCommand(newMaterialGroupsCmd())
	return cmd
}

func newMaterialListCmd() *cobra.Command {
	var name, materialType, groupID string
	var page, size int
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "search"},
		Short:   "查询素材列表及素材 ID",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := api.MaterialListOptions{
				Page:     page,
				Size:     size,
				FileName: strings.TrimSpace(name),
				Type:     strings.TrimSpace(materialType),
				GroupID:  strings.TrimSpace(groupID),
			}
			return runQuery(cmd, "material.list", func(service queryflow.Service) (interface{}, error) {
				return service.Materials(opts)
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "exact material file name, including extension")
	cmd.Flags().StringVar(&materialType, "type", "", "material type: image, video, file")
	cmd.Flags().StringVar(&groupID, "group-id", "", "filter by material group id")
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&size, "size", 100, "page size")
	return cmd
}

func newMaterialGroupsCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{
		Use:     "groups",
		Aliases: []string{"group-list"},
		Short:   "查询素材分组列表",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "material.groups", func(service queryflow.Service) (interface{}, error) {
				return service.MaterialGroups(api.MaterialGroupOptions{Page: page, Size: size})
			})
		},
	}
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&size, "size", 50, "page size")
	return cmd
}

func newMaterialMoveCmd() *cobra.Command {
	var groupID string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "move <material_id>",
		Short: "移动素材到指定分组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var materialID string
			var input materialflow.MoveInput

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					materialID = strings.TrimSpace(args[0])
					if materialID == "" {
						return yxerrors.Usage("material id must not be empty", nil).
							WithHint("请传入有效素材 ID，例如 yxer material move material_1 --group-id group_1 --dry-run。")
					}
					input = materialflow.MoveInput{GroupID: strings.TrimSpace(groupID)}
					return materialflow.ValidateMoveInput(input)
				},
				DryRun: func() (cmdflow.Result, error) {
					return cmdflow.Result{
						Action: "material.move.dry-run",
						Data: map[string]interface{}{
							"dryRun":     true,
							"materialId": materialID,
							"groupId":    input.GroupID,
							"request":    materialflow.BuildMoveBody(materialID, input),
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := materialflow.NewService(rt).Move(materialID, input)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "material.move", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().StringVar(&groupID, "group-id", "", "destination material group id")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview material move request without performing the write")
	return cmd
}

func newMaterialMoveByNameCmd() *cobra.Command {
	var groupID, materialType string
	var page, size int
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "move-by-name <file_name>",
		Short: "按文件名精确匹配素材后移动到指定分组",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var fileName string
			var input materialflow.MoveInput
			opts := api.MaterialListOptions{}

			resolve := func() (materialflow.Service, materialflow.MaterialMatch, error) {
				rt, err := app.Load()
				if err != nil {
					return materialflow.Service{}, materialflow.MaterialMatch{}, err
				}
				service := materialflow.NewService(rt)
				match, err := service.ResolveByFileName(fileName, opts)
				return service, match, err
			}

			return cmdflow.Run(cmd, dryRun, cmdflow.Flow{
				Validate: func() error {
					fileName = strings.TrimSpace(args[0])
					input = materialflow.MoveInput{GroupID: strings.TrimSpace(groupID)}
					opts = api.MaterialListOptions{
						Page: page,
						Size: size,
						Type: strings.TrimSpace(materialType),
					}
					if err := materialflow.ValidateMoveInput(input); err != nil {
						return err
					}
					if fileName == "" {
						return yxerrors.Usage("material file name must not be empty", nil).
							WithHint("请传入素材文件名，例如 yxer material move-by-name demo.png --group-id group_1 --dry-run。")
					}
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					_, match, err := resolve()
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{
						Action: "material.move-by-name.dry-run",
						Data: map[string]interface{}{
							"dryRun":     true,
							"fileName":   fileName,
							"material":   match,
							"materialId": match.ID,
							"groupId":    input.GroupID,
							"request":    materialflow.BuildMoveBody(match.ID, input),
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					service, match, err := resolve()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := service.Move(match.ID, input)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "material.move-by-name", Data: map[string]interface{}{
						"fileName":   fileName,
						"material":   match,
						"materialId": match.ID,
						"result":     result,
					}}, nil
				},
			})
		},
	}
	cmd.Flags().StringVar(&groupID, "group-id", "", "destination material group id")
	cmd.Flags().StringVar(&materialType, "type", "image", "material type used while matching")
	cmd.Flags().IntVar(&page, "page", 1, "material page to search")
	cmd.Flags().IntVar(&size, "size", 100, "material page size to search")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "resolve the material id and preview the move request without performing the write")
	return cmd
}

func newMaterialCreateCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "create <payload.json>",
		Short: "将已上传资源登记到素材库",
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
						Action: "material.create.dry-run",
						Data: map[string]interface{}{
							"dryRun":  true,
							"request": materialflow.BuildMaterialBody(payload),
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := materialflow.NewService(rt).Create(payload)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "material.create", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview material create request without performing the write")
	return cmd
}

func newMaterialAddCmd() *cobra.Command {
	opts := materialAddOptions{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "上传资源并登记到素材库",
		RunE: func(cmd *cobra.Command, args []string) error {
			var input materialflow.AddInput

			return cmdflow.Run(cmd, opts.DryRun, cmdflow.Flow{
				Validate: func() error {
					input = materialflow.AddInput{
						FilePath:  opts.FilePath,
						ThumbPath: opts.ThumbPath,
						Type:      opts.Type,
					}
					return nil
				},
				DryRun: func() (cmdflow.Result, error) {
					result, err := materialflow.PreviewAdd(input)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{
						Action: "material.add.dry-run",
						Data: map[string]interface{}{
							"dryRun":  true,
							"request": result.Request,
							"upload":  result.Upload,
							"thumb":   result.Thumb,
						},
					}, nil
				},
				Execute: func() (cmdflow.Result, error) {
					rt, err := app.Load()
					if err != nil {
						return cmdflow.Result{}, err
					}
					result, err := materialflow.NewService(rt).Add(input)
					if err != nil {
						return cmdflow.Result{}, err
					}
					return cmdflow.Result{Action: "material.add", Data: result}, nil
				},
			})
		},
	}
	cmd.Flags().StringVar(&opts.FilePath, "file", "", "local file path or URL to upload and register")
	cmd.Flags().StringVar(&opts.ThumbPath, "thumb", "", "optional thumbnail path or URL")
	cmd.Flags().StringVar(&opts.Type, "type", "", "optional material type override: image, video, file")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview upload and material request without performing the write")
	return cmd
}
