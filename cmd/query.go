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

func init() {
	rootCmd.AddCommand(newQueryCmd())
	rootCmd.AddCommand(newPrepareCmd())
}

func newQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "查询发布前置资源和发布记录",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newCategoriesCmd())
	cmd.AddCommand(newLocationsCmd())
	cmd.AddCommand(newMusicCmd())
	cmd.AddCommand(newMusicCategoriesCmd())
	cmd.AddCommand(newGoodsCmd())
	cmd.AddCommand(newGoodsDetailCmd())
	cmd.AddCommand(newEntitlementsCmd())
	cmd.AddCommand(newCollectionsCmd())
	cmd.AddCommand(newDramaTasksCmd())
	cmd.AddCommand(newMiniAppsCmd())
	cmd.AddCommand(newSyncAppsCmd())
	cmd.AddCommand(newGamesCmd())
	cmd.AddCommand(newHotEventsCmd())
	cmd.AddCommand(newGroupsCmd())
	cmd.AddCommand(newMembersCmd())
	cmd.AddCommand(newActivitiesCmd())
	cmd.AddCommand(newChallengesCmd())
	cmd.AddCommand(newRecordsCmd())
	cmd.AddCommand(newDetailsCmd())
	cmd.AddCommand(newAccountOverviewsCmd())
	cmd.AddCommand(newContentOverviewsCmd())
	cmd.AddCommand(newAccountIncrementsCmd())
	cmd.AddCommand(newProxiesCmd())
	cmd.AddCommand(newProxyAreasCmd())
	return cmd
}

func newCategoriesCmd() *cobra.Command {
	var publishType string
	var showPaths bool
	cmd := &cobra.Command{
		Use:   "categories <account_id>",
		Short: "查询分类",
		Long:  "查询分类。\n\n搜狐号视频分类支持使用 --paths 查看可直接发布的父子分类路径。\n\n当前支持平台：百家号、爱奇艺、哔哩哔哩、企鹅号、搜狐号、网易号、一点号、知乎、蜂网、AcFun。",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "categories", func(service queryflow.Service) (interface{}, error) {
				result, err := service.Categories(args[0], publishType)
				if err != nil || !showPaths {
					return result, err
				}
				return api.CategoryPathView(result)
			})
		},
	}
	cmd.Flags().StringVar(&publishType, "type", "video", "publish type")
	cmd.Flags().BoolVar(&showPaths, "paths", false, "show complete root-to-leaf category paths")
	return cmd
}

func newLocationsCmd() *cobra.Command {
	var query, keyword, locationType, nextPage string
	cmd := &cobra.Command{
		Use:   "locations <account_id>",
		Short: "查询位置",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "locations", func(service queryflow.Service) (interface{}, error) {
				return service.Locations(args[0], resolveQueryAlias(query, keyword), locationType, nextPage)
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&locationType, "type", "1", "location type")
	cmd.Flags().StringVar(&nextPage, "next-page", "", "pagination token from previous response")
	return cmd
}

func newMusicCmd() *cobra.Command {
	var query, keyword, categoryID, categoryName, nextPage string
	cmd := &cobra.Command{
		Use:   "music <account_id>",
		Short: "查询音乐",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "music", func(service queryflow.Service) (interface{}, error) {
				return service.Music(args[0], resolveQueryAlias(query, keyword), categoryID, categoryName, nextPage)
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&categoryID, "category-id", "", "music category id")
	cmd.Flags().StringVar(&categoryName, "category-name", "", "music category name")
	cmd.Flags().StringVar(&nextPage, "next-page", "", "pagination token from previous response")
	return cmd
}

func newMusicCategoriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "music-categories <account_id>",
		Short: "查询音乐分类",
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
	var query, keyword, nextPage string
	cmd := &cobra.Command{
		Use:   "goods <account_id>",
		Short: "查询商品",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "goods", func(service queryflow.Service) (interface{}, error) {
				return service.Goods(args[0], resolveQueryAlias(query, keyword), nextPage)
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&nextPage, "next-page", "", "pagination token from previous response")
	return cmd
}

func newGoodsDetailCmd() *cobra.Command {
	var productURL string
	cmd := &cobra.Command{
		Use:   "goods-detail <account_id>",
		Short: "Parse a product link into a shopping-cart product",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(productURL) == "" {
				return yxerrors.Usage("goods detail url must not be empty", nil).
					WithHint("Provide a product link with --url.")
			}
			return runQuery(cmd, "goods-detail", func(service queryflow.Service) (interface{}, error) {
				return service.GoodsDetail(args[0], strings.TrimSpace(productURL))
			})
		},
	}
	cmd.Flags().StringVar(&productURL, "url", "", "product link to parse")
	return cmd
}

func newEntitlementsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entitlements <account_id>",
		Short: "Query shopping-cart and group-shopping entitlements",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "entitlements", func(service queryflow.Service) (interface{}, error) {
				return service.Entitlements(args[0])
			})
		},
	}
	return cmd
}

func newCollectionsCmd() *cobra.Command {
	var publishType string
	cmd := &cobra.Command{
		Use:   "collections <account_id>",
		Short: "查询合集",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "collections", func(service queryflow.Service) (interface{}, error) {
				return service.Collections(args[0], publishType)
			})
		},
	}
	cmd.Flags().StringVar(&publishType, "type", "video", "publish type")
	return cmd
}

func newDramaTasksCmd() *cobra.Command {
	var query, keyword string
	cmd := &cobra.Command{
		Use:   "drama-tasks <account_id>",
		Short: "查询视频号剧集",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "drama-tasks", func(service queryflow.Service) (interface{}, error) {
				return service.DramaTasks(args[0], resolveQueryAlias(query, keyword))
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newMiniAppsCmd() *cobra.Command {
	var query, keyword string
	cmd := &cobra.Command{
		Use:   "miniapps <account_id>",
		Short: "查询小程序",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "miniapps", func(service queryflow.Service) (interface{}, error) {
				return service.MiniApps(args[0], resolveQueryAlias(query, keyword))
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newSyncAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syncapps <account_id>",
		Short: "查询同步发布应用",
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
	var query, keyword string
	cmd := &cobra.Command{
		Use:   "games <account_id>",
		Short: "查询可挂载游戏",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "games", func(service queryflow.Service) (interface{}, error) {
				return service.Games(args[0], resolveQueryAlias(query, keyword))
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newHotEventsCmd() *cobra.Command {
	var publishType string
	cmd := &cobra.Command{
		Use:   "hot-events <account_id>",
		Short: "查询热点列表",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "hot-events", func(service queryflow.Service) (interface{}, error) {
				return service.HotEvents(args[0], publishType)
			})
		},
	}
	cmd.Flags().StringVar(&publishType, "type", "video", "publish type")
	return cmd
}

func newGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups <account_id>",
		Short: "查询群聊列表",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "groups", func(service queryflow.Service) (interface{}, error) {
				return service.Groups(args[0])
			})
		},
	}
	return cmd
}

func newMembersCmd() *cobra.Command {
	opts := membersOptions{}
	cmd := &cobra.Command{
		Use:   "members",
		Short: "查询成员列表",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := api.MembersOptions{
				Page:     opts.Page,
				Size:     opts.Size,
				Statuses: opts.Statuses,
				KeyWords: resolveQueryAlias(opts.Query, opts.Keyword),
				Role:     opts.Role,
			}
			return runQuery(cmd, "members", func(service queryflow.Service) (interface{}, error) {
				return service.Members(input)
			})
		},
	}
	cmd.Flags().IntVar(&opts.Page, "page", 1, "page number")
	cmd.Flags().IntVar(&opts.Size, "size", 10, "page size")
	cmd.Flags().StringSliceVar(&opts.Statuses, "status", nil, "member status; repeat or comma-separate for multiple")
	cmd.Flags().StringVar(&opts.Query, "query", "", "member name or phone keyword")
	cmd.Flags().StringVar(&opts.Keyword, "keyword", "", "member name or phone keyword (alias for --query)")
	cmd.Flags().StringVar(&opts.Role, "role", "", "member role: master, admin, member")
	return cmd
}

type membersOptions struct {
	Page     int
	Size     int
	Statuses []string
	Query    string
	Keyword  string
	Role     string
}

func newActivitiesCmd() *cobra.Command {
	var publishType, query, keyword, categoryID string
	cmd := &cobra.Command{
		Use:   "activities <account_id>",
		Short: "查询征文/激励活动",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "activities", func(service queryflow.Service) (interface{}, error) {
				return service.Activities(args[0], publishType, categoryID, resolveQueryAlias(query, keyword))
			})
		},
	}
	cmd.Flags().StringVar(&publishType, "type", "article", "publish type")
	cmd.Flags().StringVar(&categoryID, "category-id", "", "category id")
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}

func newChallengesCmd() *cobra.Command {
	var query, keyword, publishType, nextPage string
	cmd := &cobra.Command{
		Use:   "challenges <account_id>",
		Short: "查询话题/挑战",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "challenges", func(service queryflow.Service) (interface{}, error) {
				return service.Challenges(args[0], resolveQueryAlias(query, keyword), publishType, nextPage)
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	cmd.Flags().StringVar(&publishType, "type", "video", "publish type")
	cmd.Flags().StringVar(&nextPage, "next-page", "", "pagination token from previous response")
	return cmd
}

func newRecordsCmd() *cobra.Command {
	opts := recordsOptions{}
	cmd := &cobra.Command{
		Use:   "records",
		Short: "查询发布记录",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecordsListWithOptions(cmd, opts)
		},
	}
	cmd.Flags().StringVar(&opts.Platform, "platform", "", "filter by platform")
	cmd.Flags().StringVar(&opts.Limit, "limit", "", "result limit (required)")
	cmd.Flags().StringVar(&opts.Status, "status", "", "filter by status")
	cmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "列出发布记录",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecordsListWithOptions(cmd, opts)
		},
	})
	return cmd
}

type recordsOptions struct {
	Platform string
	Limit    string
	Status   string
}

func runRecordsListWithOptions(cmd *cobra.Command, opts recordsOptions) error {
	if strings.TrimSpace(opts.Limit) == "" {
		return yxerrors.Usage("records limit must not be empty", nil).
			WithHint("请传入有效的 --limit 值，例如 10。")
	}
	return runQuery(cmd, "records.list", func(service queryflow.Service) (interface{}, error) {
		return service.Records(opts.Platform, opts.Limit, opts.Status)
	})
}

func newDetailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "details <task_set_id>",
		Short: "查询发布任务详情",
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
	opts := accountOverviewOptions{}
	cmd := &cobra.Command{
		Use:   "account-overviews",
		Short: "查询账号数据概览",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(opts.Platform) == "" {
				return yxerrors.Usage("account overviews platform must not be empty", nil).
					WithHint("请传入 --platform，例如 --platform 抖音。")
			}
			input := api.AccountOverviewOptions{
				Platform:    opts.Platform,
				Name:        opts.Name,
				Group:       opts.Group,
				LoginStatus: opts.LoginStatus,
				MemberIDs:   opts.MemberIDs,
				Page:        opts.Page,
				Size:        opts.Size,
			}
			return runQuery(cmd, "account-overviews", func(service queryflow.Service) (interface{}, error) {
				return service.AccountOverviews(input)
			})
		},
	}
	cmd.Flags().StringVar(&opts.Platform, "platform", "", "platform name")
	cmd.Flags().StringVar(&opts.Name, "name", "", "account name keyword")
	cmd.Flags().StringVar(&opts.Group, "group", "", "group name")
	cmd.Flags().StringVar(&opts.LoginStatus, "login-status", "", "login status")
	cmd.Flags().StringSliceVar(&opts.MemberIDs, "member-id", nil, "member id; repeat or comma-separate for multiple")
	cmd.Flags().IntVar(&opts.Page, "page", 1, "page number")
	cmd.Flags().IntVar(&opts.Size, "size", 10, "page size")
	return cmd
}

type accountOverviewOptions struct {
	Platform    string
	Name        string
	Group       string
	LoginStatus string
	MemberIDs   []string
	Page        int
	Size        int
}

func newContentOverviewsCmd() *cobra.Command {
	opts := contentOverviewOptions{}
	cmd := &cobra.Command{
		Use:   "content-overviews",
		Short: "查询作品数据概览",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := api.ContentOverviewOptions{
				Platform:          opts.Platform,
				PlatformAccountID: opts.AccountID,
				PublishUserID:     opts.PublishUserID,
				Type:              opts.Type,
				Title:             opts.Title,
				PublishStartTime:  opts.PublishStart,
				PublishEndTime:    opts.PublishEnd,
				Page:              opts.Page,
				Size:              opts.Size,
			}
			return runQuery(cmd, "content-overviews", func(service queryflow.Service) (interface{}, error) {
				return service.ContentOverviews(input)
			})
		},
	}
	cmd.Flags().StringVar(&opts.Platform, "platform", "", "platform name")
	cmd.Flags().StringVar(&opts.AccountID, "account-id", "", "platform account id")
	cmd.Flags().StringVar(&opts.PublishUserID, "publish-user-id", "", "publish user id")
	cmd.Flags().StringVar(&opts.Type, "type", "", "content type")
	cmd.Flags().StringVar(&opts.Title, "title", "", "title keyword")
	cmd.Flags().StringVar(&opts.PublishStart, "publish-start-time", "", "publish start timestamp in milliseconds")
	cmd.Flags().StringVar(&opts.PublishEnd, "publish-end-time", "", "publish end timestamp in milliseconds")
	cmd.Flags().IntVar(&opts.Page, "page", 1, "page number")
	cmd.Flags().IntVar(&opts.Size, "size", 10, "page size")
	return cmd
}

type contentOverviewOptions struct {
	Platform      string
	AccountID     string
	PublishUserID string
	Type          string
	Title         string
	PublishStart  string
	PublishEnd    string
	Page          int
	Size          int
}

func newProxiesCmd() *cobra.Command {
	var size string
	cmd := &cobra.Command{
		Use:   "proxies",
		Short: "查询代理列表",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "proxies", func(service queryflow.Service) (interface{}, error) {
				return service.Proxies(size)
			})
		},
	}
	cmd.Flags().StringVar(&size, "size", "9999", "page size")
	return cmd
}

func newProxyAreasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy-areas",
		Short: "查询内置代理地区",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "proxy-areas", func(service queryflow.Service) (interface{}, error) {
				return service.ProxyAreas()
			})
		},
	}
	return cmd
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
