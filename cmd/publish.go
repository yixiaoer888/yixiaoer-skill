package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
)

var (
	publishChannelFlag string
	publishClientID    string
	publishAccount     string
	publishTitle       string
	publishDescription string
	publishContent     string
	publishImages      []string
	publishVideoPath   string
	publishVideoKey    string
	publishCoverPath   string
	publishVisibleType int
	publishDryRun      bool
)

func init() {
	publishCmd.Flags().StringVar(&publishChannelFlag, "publish-channel", "", `publish channel: "cloud" or "local"`)
	publishCmd.Flags().StringVar(&publishClientID, "client-id", "", "client ID for local publish")
	publishCmd.Flags().StringVar(&publishAccount, "account", "", "account name or ID for flags mode")
	publishCmd.Flags().StringVar(&publishTitle, "title", "", "title for flags mode")
	publishCmd.Flags().StringVar(&publishDescription, "description", "", "description for flags mode")
	publishCmd.Flags().StringVar(&publishContent, "content", "", "content for article flags mode; prefix with @ to read a file")
	publishCmd.Flags().StringSliceVar(&publishImages, "image", nil, "image path or URL for imageText flags mode; repeatable")
	publishCmd.Flags().StringVar(&publishVideoPath, "video", "", "local video path or URL for video flags mode")
	publishCmd.Flags().StringVar(&publishVideoKey, "video-key", "", "uploaded video key for video flags mode")
	publishCmd.Flags().StringVar(&publishCoverPath, "cover", "", "cover image path or URL for flags mode")
	publishCmd.Flags().IntVar(&publishVisibleType, "visible-type", -1, "visible type for supported platforms in flags mode")
	publishCmd.Flags().BoolVar(&publishDryRun, "dry-run", false, "preview the publish request without performing the write")
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish <type> <中文平台名|platform-key> [payload.json] [clientId]",
	Short: "发布内容（单平台原子发布）",
	Long: strings.TrimSpace(`
支持两种用法：

1. 兼容旧模式：直接传 payload.json
   yxer publish imageText 小红书 ./publish-payload.json

2. 推荐 flags 模式：由 CLI 自动解析账号、上传资源并构造 payload
   yxer publish imageText 小红书 --account "图文账号" --title "标题" --description "正文" --image ./1.jpg
   yxer publish article 知乎 --account "知乎账号" --title "文章标题" --content @./article.html --cover ./cover.png
   yxer publish video 抖音 --account "视频账号" --title "视频标题" --description "视频描述" --video ./clip.mp4 --cover ./cover.png

flags 模式目前最完整支持 imageText 和 article；video 已支持本地视频探测与上传。
`),
	Args: cobra.RangeArgs(2, 4),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			payload            map[string]interface{}
			err                error
			positionalClientID string
		)
		if len(args) >= 3 && looksLikePayloadArg(args[2]) {
			payload, err = readPayload(args[2])
			if err != nil {
				return err
			}
			if len(args) == 4 {
				positionalClientID = args[3]
			}
		} else {
			payload, err = publishflow.NewService().BuildPayload(publishflow.BuildInput{
				PublishType:   args[0],
				PlatformInput: args[1],
				Account:       publishAccount,
				Title:         publishTitle,
				Description:   publishDescription,
				Content:       publishContent,
				Images:        publishImages,
				VideoPath:     publishVideoPath,
				VideoKey:      publishVideoKey,
				CoverPath:     publishCoverPath,
				VisibleType:   publishVisibleType,
			})
			if err != nil {
				return err
			}
		}
		input := publishflow.ExecuteInput{
			PublishType:        args[0],
			PlatformInput:      args[1],
			Payload:            payload,
			PositionalClientID: positionalClientID,
			FlagChannel:        publishChannelFlag,
			FlagClientID:       publishClientID,
		}
		if publishDryRun {
			result, err := publishflow.NewService().DryRun(input)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "publish.dry-run", map[string]interface{}{
				"dryRun":  true,
				"request": result.PublishBody,
				"meta": map[string]interface{}{
					"platform":       result.Platform,
					"publishType":    result.PublishType,
					"publishChannel": result.PublishMode,
					"clientId":       result.ClientID,
					"accountIds":     result.AccountIDs,
					"schemaChecked":  result.SchemaChecked,
				},
			})
		}
		result, err := publishflow.NewService().Execute(input)
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), "publish", result)
	},
}

func looksLikePayloadArg(value string) bool {
	value = strings.TrimSpace(value)
	return strings.HasSuffix(strings.ToLower(value), ".json")
}
