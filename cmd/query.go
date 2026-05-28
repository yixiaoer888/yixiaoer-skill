package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var (
	categoriesType  string
	locationsQuery  string
	locationsType   string
	musicQuery      string
	goodsQuery      string
	collectionsType string
	challengesQuery string
	challengesType  string
	recordsPlatform string
	recordsLimit    string
	recordsStatus   string
)

func init() {
	categoriesCmd.Flags().StringVar(&categoriesType, "type", "video", "publish type")
	locationsCmd.Flags().StringVar(&locationsQuery, "query", "", "search keyword")
	locationsCmd.Flags().StringVar(&locationsQuery, "keyword", "", "search keyword")
	locationsCmd.Flags().StringVar(&locationsType, "type", "1", "location type")
	musicCmd.Flags().StringVar(&musicQuery, "query", "", "search keyword")
	musicCmd.Flags().StringVar(&musicQuery, "keyword", "", "search keyword")
	goodsCmd.Flags().StringVar(&goodsQuery, "query", "", "search keyword")
	goodsCmd.Flags().StringVar(&goodsQuery, "keyword", "", "search keyword")
	collectionsCmd.Flags().StringVar(&collectionsType, "type", "video", "publish type")
	challengesCmd.Flags().StringVar(&challengesQuery, "query", "", "search keyword")
	challengesCmd.Flags().StringVar(&challengesQuery, "keyword", "", "search keyword")
	challengesCmd.Flags().StringVar(&challengesType, "type", "video", "publish type")
	recordsCmd.Flags().StringVar(&recordsPlatform, "platform", "", "filter by platform")
	recordsCmd.Flags().StringVar(&recordsLimit, "limit", "10", "result limit")
	recordsCmd.Flags().StringVar(&recordsStatus, "status", "", "filter by status")
	recordsCmd.AddCommand(recordsListCmd)

	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(locationsCmd)
	rootCmd.AddCommand(musicCmd)
	rootCmd.AddCommand(goodsCmd)
	rootCmd.AddCommand(collectionsCmd)
	rootCmd.AddCommand(challengesCmd)
	rootCmd.AddCommand(recordsCmd)
	rootCmd.AddCommand(prepareCmd)
}

var categoriesCmd = &cobra.Command{
	Use:   "categories <account_id>",
	Short: "查询分类",
	Long:  "查询分类。\n\n当前支持平台：百家号、爱奇艺、哔哩哔哩、企鹅号、网易号、一点号、知乎、蜂网、AcFun。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuery(cmd, "categories", func(service queryflow.Service) (interface{}, error) {
			return service.Categories(args[0], categoriesType)
		})
	},
}

var locationsCmd = &cobra.Command{
	Use:   "locations <account_id>",
	Short: "查询位置",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuery(cmd, "locations", func(service queryflow.Service) (interface{}, error) {
			return service.Locations(args[0], locationsQuery, locationsType)
		})
	},
}

var musicCmd = &cobra.Command{
	Use:   "music <account_id>",
	Short: "查询音乐",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuery(cmd, "music", func(service queryflow.Service) (interface{}, error) {
			return service.Music(args[0], musicQuery)
		})
	},
}

var goodsCmd = &cobra.Command{
	Use:   "goods <account_id>",
	Short: "查询商品",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuery(cmd, "goods", func(service queryflow.Service) (interface{}, error) {
			return service.Goods(args[0], goodsQuery)
		})
	},
}

var collectionsCmd = &cobra.Command{
	Use:   "collections <account_id>",
	Short: "查询合集",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuery(cmd, "collections", func(service queryflow.Service) (interface{}, error) {
			return service.Collections(args[0], collectionsType)
		})
	},
}

var challengesCmd = &cobra.Command{
	Use:   "challenges <account_id>",
	Short: "查询话题/挑战",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuery(cmd, "challenges", func(service queryflow.Service) (interface{}, error) {
			return service.Challenges(args[0], challengesQuery, challengesType)
		})
	},
}

var recordsCmd = &cobra.Command{
	Use:   "records",
	Short: "查询发布记录",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRecordsList(cmd)
	},
}

var recordsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "列出发布记录",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRecordsList(cmd)
	},
}

func runRecordsList(cmd *cobra.Command) error {
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
	if action == "records" && recordsLimit == "" {
		return yxerrors.Usage("records limit must not be empty", nil).
			WithHint("请传入有效的 --limit 值，例如 10。")
	}
	result, err := query(queryflow.NewService())
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), action, result)
}
