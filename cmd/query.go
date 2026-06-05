package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var (
	categoriesType    string
	locationsQuery    string
	locationsKeyword  string
	locationsType     string
	musicQuery        string
	musicKeyword      string
	goodsQuery        string
	goodsKeyword      string
	collectionsType   string
	challengesQuery   string
	challengesKeyword string
	challengesType    string
	recordsPlatform   string
	recordsLimit      string
	recordsStatus     string
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "查询发布前置资源和发布记录",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(newCategoriesCmd())
	rootCmd.AddCommand(newLocationsCmd())
	rootCmd.AddCommand(newMusicCmd())
	rootCmd.AddCommand(newGoodsCmd())
	rootCmd.AddCommand(newCollectionsCmd())
	rootCmd.AddCommand(newChallengesCmd())
	rootCmd.AddCommand(newRecordsCmd())
	queryCmd.AddCommand(newCategoriesCmd())
	queryCmd.AddCommand(newLocationsCmd())
	queryCmd.AddCommand(newMusicCmd())
	queryCmd.AddCommand(newGoodsCmd())
	queryCmd.AddCommand(newCollectionsCmd())
	queryCmd.AddCommand(newChallengesCmd())
	queryCmd.AddCommand(newRecordsCmd())
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(prepareCmd)
}

func newCategoriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "categories <account_id>",
		Short: "查询分类（兼容入口，推荐使用 yxer query categories）",
		Long:  "查询分类。\n\n当前支持平台：百家号、爱奇艺、哔哩哔哩、企鹅号、网易号、一点号、知乎、蜂网、AcFun。",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "categories", func(service queryflow.Service) (interface{}, error) {
				return service.Categories(args[0], categoriesType)
			})
		},
	}
	cmd.Flags().StringVar(&categoriesType, "type", "video", "publish type")
	return cmd
}

func newLocationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locations <account_id>",
		Short: "查询位置（兼容入口，推荐使用 yxer query locations）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "locations", func(service queryflow.Service) (interface{}, error) {
				return service.Locations(args[0], resolveQueryAlias(locationsQuery, locationsKeyword), locationsType)
			})
		},
	}
	cmd.Flags().StringVar(&locationsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&locationsKeyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&locationsType, "type", "1", "location type")
	return cmd
}

func newMusicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "music <account_id>",
		Short: "查询音乐（兼容入口，推荐使用 yxer query music）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "music", func(service queryflow.Service) (interface{}, error) {
				return service.Music(args[0], resolveQueryAlias(musicQuery, musicKeyword))
			})
		},
	}
	cmd.Flags().StringVar(&musicQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&musicKeyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newGoodsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goods <account_id>",
		Short: "查询商品（兼容入口，推荐使用 yxer query goods）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "goods", func(service queryflow.Service) (interface{}, error) {
				return service.Goods(args[0], resolveQueryAlias(goodsQuery, goodsKeyword))
			})
		},
	}
	cmd.Flags().StringVar(&goodsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&goodsKeyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collections <account_id>",
		Short: "查询合集（兼容入口，推荐使用 yxer query collections）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "collections", func(service queryflow.Service) (interface{}, error) {
				return service.Collections(args[0], collectionsType)
			})
		},
	}
	cmd.Flags().StringVar(&collectionsType, "type", "video", "publish type")
	return cmd
}

func newChallengesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "challenges <account_id>",
		Short: "查询话题/挑战（兼容入口，推荐使用 yxer query challenges）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "challenges", func(service queryflow.Service) (interface{}, error) {
				return service.Challenges(args[0], resolveQueryAlias(challengesQuery, challengesKeyword), challengesType)
			})
		},
	}
	cmd.Flags().StringVar(&challengesQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&challengesKeyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&challengesType, "type", "video", "publish type")
	return cmd
}

func newRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records",
		Short: "查询发布记录（兼容入口，推荐使用 yxer query records）",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecordsList(cmd)
		},
	}
	cmd.Flags().StringVar(&recordsPlatform, "platform", "", "filter by platform")
	cmd.Flags().StringVar(&recordsLimit, "limit", "", "result limit (required)")
	cmd.Flags().StringVar(&recordsStatus, "status", "", "filter by status")
	cmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "列出发布记录",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecordsList(cmd)
		},
	})
	return cmd
}

func runRecordsList(cmd *cobra.Command) error {
	if strings.TrimSpace(recordsLimit) == "" {
		return yxerrors.Usage("records limit must not be empty", nil).
			WithHint("请传入有效的 --limit 值，例如 10。")
	}
	return runQuery(cmd, "records.list", func(service queryflow.Service) (interface{}, error) {
		return service.Records(recordsPlatform, recordsLimit, recordsStatus)
	})
}

var prepareCmd = &cobra.Command{
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

func runQuery(cmd *cobra.Command, action string, query func(queryflow.Service) (interface{}, error)) error {
	rt, err := app.Load()
	if err != nil {
		return err
	}
	result, err := query(queryflow.NewService(rt))
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), action, result)
}

func resolveQueryAlias(primary, alias string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return alias
}
