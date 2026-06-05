package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	queryflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/query"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var (
	categoriesType               string
	locationsQuery               string
	locationsKeyword             string
	locationsType                string
	locationsNextPage            string
	musicQuery                   string
	musicKeyword                 string
	musicCategoryID              string
	musicCategoryName            string
	musicNextPage                string
	activitiesType               string
	activitiesQuery              string
	activitiesKeyword            string
	activitiesCategoryID         string
	goodsQuery                   string
	goodsKeyword                 string
	goodsNextPage                string
	collectionsType              string
	miniAppsQuery                string
	miniAppsKeyword              string
	gamesQuery                   string
	gamesKeyword                 string
	hotEventsType                string
	challengesQuery              string
	challengesKeyword            string
	challengesType               string
	challengesNextPage           string
	recordsPlatform              string
	recordsLimit                 string
	recordsStatus                string
	accountOverviewPlatform      string
	accountOverviewName          string
	accountOverviewGroup         string
	accountOverviewLoginStatus   string
	accountOverviewMemberIDs     []string
	accountOverviewPage          int
	accountOverviewSize          int
	contentOverviewPlatform      string
	contentOverviewAccountID     string
	contentOverviewPublishUserID string
	contentOverviewType          string
	contentOverviewTitle         string
	contentOverviewPublishStart  string
	contentOverviewPublishEnd    string
	contentOverviewPage          int
	contentOverviewSize          int
	proxiesSize                  string
	updateAccountProxyID         string
	updateAccountKuaidailiArea   string
	updateAccountRemark          string
	updateAccountGroups          []string
	updateAccountDryRun          bool
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
	rootCmd.AddCommand(newMusicCategoriesCmd())
	rootCmd.AddCommand(newGoodsCmd())
	rootCmd.AddCommand(newCollectionsCmd())
	rootCmd.AddCommand(newMiniAppsCmd())
	rootCmd.AddCommand(newSyncAppsCmd())
	rootCmd.AddCommand(newGamesCmd())
	rootCmd.AddCommand(newHotEventsCmd())
	rootCmd.AddCommand(newGroupsCmd())
	rootCmd.AddCommand(newActivitiesCmd())
	rootCmd.AddCommand(newChallengesCmd())
	rootCmd.AddCommand(newRecordsCmd())
	rootCmd.AddCommand(newDetailsCmd())
	rootCmd.AddCommand(newAccountOverviewsCmd())
	rootCmd.AddCommand(newContentOverviewsCmd())
	rootCmd.AddCommand(newProxiesCmd())
	rootCmd.AddCommand(newProxyAreasCmd())
	rootCmd.AddCommand(newUpdateAccountCmd())
	queryCmd.AddCommand(newCategoriesCmd())
	queryCmd.AddCommand(newLocationsCmd())
	queryCmd.AddCommand(newMusicCmd())
	queryCmd.AddCommand(newMusicCategoriesCmd())
	queryCmd.AddCommand(newGoodsCmd())
	queryCmd.AddCommand(newCollectionsCmd())
	queryCmd.AddCommand(newMiniAppsCmd())
	queryCmd.AddCommand(newSyncAppsCmd())
	queryCmd.AddCommand(newGamesCmd())
	queryCmd.AddCommand(newHotEventsCmd())
	queryCmd.AddCommand(newGroupsCmd())
	queryCmd.AddCommand(newActivitiesCmd())
	queryCmd.AddCommand(newChallengesCmd())
	queryCmd.AddCommand(newRecordsCmd())
	queryCmd.AddCommand(newDetailsCmd())
	queryCmd.AddCommand(newAccountOverviewsCmd())
	queryCmd.AddCommand(newContentOverviewsCmd())
	queryCmd.AddCommand(newProxiesCmd())
	queryCmd.AddCommand(newProxyAreasCmd())
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
				return service.Locations(args[0], resolveQueryAlias(locationsQuery, locationsKeyword), locationsType, locationsNextPage)
			})
		},
	}
	cmd.Flags().StringVar(&locationsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&locationsKeyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&locationsType, "type", "1", "location type")
	cmd.Flags().StringVar(&locationsNextPage, "next-page", "", "pagination token from previous response")
	return cmd
}

func newMusicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "music <account_id>",
		Short: "查询音乐（兼容入口，推荐使用 yxer query music）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "music", func(service queryflow.Service) (interface{}, error) {
				return service.Music(args[0], resolveQueryAlias(musicQuery, musicKeyword), musicCategoryID, musicCategoryName, musicNextPage)
			})
		},
	}
	cmd.Flags().StringVar(&musicQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&musicKeyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&musicCategoryID, "category-id", "", "music category id")
	cmd.Flags().StringVar(&musicCategoryName, "category-name", "", "music category name")
	cmd.Flags().StringVar(&musicNextPage, "next-page", "", "pagination token from previous response")
	return cmd
}

func newMusicCategoriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "music-categories <account_id>",
		Short: "查询音乐分类（兼容入口，推荐使用 yxer query music-categories）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "music-categories", func(service queryflow.Service) (interface{}, error) {
				return service.MusicCategories(args[0])
			})
		},
	}
	return cmd
}

func newGoodsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goods <account_id>",
		Short: "查询商品（兼容入口，推荐使用 yxer query goods）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "goods", func(service queryflow.Service) (interface{}, error) {
				return service.Goods(args[0], resolveQueryAlias(goodsQuery, goodsKeyword), goodsNextPage)
			})
		},
	}
	cmd.Flags().StringVar(&goodsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&goodsKeyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&goodsNextPage, "next-page", "", "pagination token from previous response")
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

func newMiniAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miniapps <account_id>",
		Short: "查询小程序（兼容入口，推荐使用 yxer query miniapps）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "miniapps", func(service queryflow.Service) (interface{}, error) {
				return service.MiniApps(args[0], resolveQueryAlias(miniAppsQuery, miniAppsKeyword))
			})
		},
	}
	cmd.Flags().StringVar(&miniAppsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&miniAppsKeyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newSyncAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syncapps <account_id>",
		Short: "查询同步发布应用（兼容入口，推荐使用 yxer query syncapps）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "syncapps", func(service queryflow.Service) (interface{}, error) {
				return service.SyncApps(args[0])
			})
		},
	}
	return cmd
}

func newGamesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "games <account_id>",
		Short: "查询可挂载游戏（兼容入口，推荐使用 yxer query games）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "games", func(service queryflow.Service) (interface{}, error) {
				return service.Games(args[0], resolveQueryAlias(gamesQuery, gamesKeyword))
			})
		},
	}
	cmd.Flags().StringVar(&gamesQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&gamesKeyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newHotEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hot-events <account_id>",
		Short: "查询热点列表（兼容入口，推荐使用 yxer query hot-events）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "hot-events", func(service queryflow.Service) (interface{}, error) {
				return service.HotEvents(args[0], hotEventsType)
			})
		},
	}
	cmd.Flags().StringVar(&hotEventsType, "type", "video", "publish type")
	return cmd
}

func newGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups <account_id>",
		Short: "查询群聊列表（兼容入口，推荐使用 yxer query groups）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "groups", func(service queryflow.Service) (interface{}, error) {
				return service.Groups(args[0])
			})
		},
	}
	return cmd
}

func newActivitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activities <account_id>",
		Short: "查询征文/激励活动（兼容入口，推荐使用 yxer query activities）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "activities", func(service queryflow.Service) (interface{}, error) {
				return service.Activities(args[0], activitiesType, activitiesCategoryID, resolveQueryAlias(activitiesQuery, activitiesKeyword))
			})
		},
	}
	cmd.Flags().StringVar(&activitiesType, "type", "article", "publish type")
	cmd.Flags().StringVar(&activitiesCategoryID, "category-id", "", "category id")
	cmd.Flags().StringVar(&activitiesQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&activitiesKeyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newChallengesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "challenges <account_id>",
		Short: "查询话题/挑战（兼容入口，推荐使用 yxer query challenges）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "challenges", func(service queryflow.Service) (interface{}, error) {
				return service.Challenges(args[0], resolveQueryAlias(challengesQuery, challengesKeyword), challengesType, challengesNextPage)
			})
		},
	}
	cmd.Flags().StringVar(&challengesQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&challengesKeyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&challengesType, "type", "video", "publish type")
	cmd.Flags().StringVar(&challengesNextPage, "next-page", "", "pagination token from previous response")
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

func newDetailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "details <task_set_id>",
		Short: "查询发布任务详情（兼容入口，推荐使用 yxer query details）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "details", func(service queryflow.Service) (interface{}, error) {
				return service.Details(args[0])
			})
		},
	}
	return cmd
}

func newAccountOverviewsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-overviews",
		Short: "查询账号数据概览（兼容入口，推荐使用 yxer query account-overviews）",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(accountOverviewPlatform) == "" {
				return yxerrors.Usage("account overviews platform must not be empty", nil).
					WithHint("请传入 --platform，例如 --platform 抖音。")
			}
			opts := api.AccountOverviewOptions{
				Platform:    accountOverviewPlatform,
				Name:        accountOverviewName,
				Group:       accountOverviewGroup,
				LoginStatus: accountOverviewLoginStatus,
				MemberIDs:   accountOverviewMemberIDs,
				Page:        accountOverviewPage,
				Size:        accountOverviewSize,
			}
			return runQuery(cmd, "account-overviews", func(service queryflow.Service) (interface{}, error) {
				return service.AccountOverviews(opts)
			})
		},
	}
	cmd.Flags().StringVar(&accountOverviewPlatform, "platform", "", "platform name")
	cmd.Flags().StringVar(&accountOverviewName, "name", "", "account name keyword")
	cmd.Flags().StringVar(&accountOverviewGroup, "group", "", "group name")
	cmd.Flags().StringVar(&accountOverviewLoginStatus, "login-status", "", "login status")
	cmd.Flags().StringSliceVar(&accountOverviewMemberIDs, "member-id", nil, "member id; repeat or comma-separate for multiple")
	cmd.Flags().IntVar(&accountOverviewPage, "page", 1, "page number")
	cmd.Flags().IntVar(&accountOverviewSize, "size", 10, "page size")
	return cmd
}

func newContentOverviewsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content-overviews",
		Short: "查询作品数据概览（兼容入口，推荐使用 yxer query content-overviews）",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := api.ContentOverviewOptions{
				Platform:          contentOverviewPlatform,
				PlatformAccountID: contentOverviewAccountID,
				PublishUserID:     contentOverviewPublishUserID,
				Type:              contentOverviewType,
				Title:             contentOverviewTitle,
				PublishStartTime:  contentOverviewPublishStart,
				PublishEndTime:    contentOverviewPublishEnd,
				Page:              contentOverviewPage,
				Size:              contentOverviewSize,
			}
			return runQuery(cmd, "content-overviews", func(service queryflow.Service) (interface{}, error) {
				return service.ContentOverviews(opts)
			})
		},
	}
	cmd.Flags().StringVar(&contentOverviewPlatform, "platform", "", "platform name")
	cmd.Flags().StringVar(&contentOverviewAccountID, "account-id", "", "platform account id")
	cmd.Flags().StringVar(&contentOverviewPublishUserID, "publish-user-id", "", "publish user id")
	cmd.Flags().StringVar(&contentOverviewType, "type", "", "content type")
	cmd.Flags().StringVar(&contentOverviewTitle, "title", "", "title keyword")
	cmd.Flags().StringVar(&contentOverviewPublishStart, "publish-start-time", "", "publish start timestamp in milliseconds")
	cmd.Flags().StringVar(&contentOverviewPublishEnd, "publish-end-time", "", "publish end timestamp in milliseconds")
	cmd.Flags().IntVar(&contentOverviewPage, "page", 1, "page number")
	cmd.Flags().IntVar(&contentOverviewSize, "size", 10, "page size")
	return cmd
}

func newProxiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxies",
		Short: "查询代理列表（兼容入口，推荐使用 yxer query proxies）",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "proxies", func(service queryflow.Service) (interface{}, error) {
				return service.Proxies(proxiesSize)
			})
		},
	}
	cmd.Flags().StringVar(&proxiesSize, "size", "9999", "page size")
	return cmd
}

func newProxyAreasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy-areas",
		Short: "查询内置代理地区（兼容入口，推荐使用 yxer query proxy-areas）",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "proxy-areas", func(service queryflow.Service) (interface{}, error) {
				return service.ProxyAreas()
			})
		},
	}
	return cmd
}

func newUpdateAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-account <account_id>",
		Short: "更新账号代理或备注",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := updateAccountBody()
			if len(body) == 0 {
				return yxerrors.Usage("update account request must not be empty", nil).
					WithHint("请至少传入 --proxy-id、--kuaidaili-area、--remark 或 --group。")
			}
			if updateAccountDryRun {
				return output.Success(cmd.OutOrStdout(), "update-account.dry-run", map[string]interface{}{
					"dryRun":  true,
					"account": args[0],
					"request": body,
				})
			}
			rt, err := app.Load()
			if err != nil {
				return err
			}
			result, err := queryflow.NewService(rt).UpdateAccount(args[0], body)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "update-account", result)
		},
	}
	cmd.Flags().StringVar(&updateAccountProxyID, "proxy-id", "", "team proxy id")
	cmd.Flags().StringVar(&updateAccountKuaidailiArea, "kuaidaili-area", "", "built-in proxy area code")
	cmd.Flags().StringVar(&updateAccountRemark, "remark", "", "account remark")
	cmd.Flags().StringSliceVar(&updateAccountGroups, "group", nil, "group id; repeat or comma-separate for multiple")
	cmd.Flags().BoolVar(&updateAccountDryRun, "dry-run", false, "preview update request without performing the write")
	return cmd
}

func updateAccountBody() map[string]interface{} {
	body := map[string]interface{}{}
	if strings.TrimSpace(updateAccountProxyID) != "" {
		body["proxyId"] = updateAccountProxyID
	}
	if strings.TrimSpace(updateAccountKuaidailiArea) != "" {
		body["kuaidailiArea"] = updateAccountKuaidailiArea
	}
	if strings.TrimSpace(updateAccountRemark) != "" {
		body["remark"] = updateAccountRemark
	}
	if len(updateAccountGroups) > 0 {
		groups := make([]string, 0, len(updateAccountGroups))
		for _, group := range updateAccountGroups {
			if strings.TrimSpace(group) != "" {
				groups = append(groups, group)
			}
		}
		if len(groups) > 0 {
			body["groups"] = groups
		}
	}
	return body
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
