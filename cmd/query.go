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
	var publishType string
	cmd := &cobra.Command{
		Use:   "categories <account_id>",
		Short: "查询分类",
		Long:  "查询分类。\n\n当前支持平台：百家号、爱奇艺、哔哩哔哩、企鹅号、网易号、一点号、知乎、蜂网、AcFun。",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "categories", func(service queryflow.Service) (interface{}, error) {
				return service.Categories(args[0], publishType)
			})
		},
	}
	cmd.Flags().StringVar(&publishType, "type", "video", "publish type")
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

func runRecordsList(cmd *cobra.Command) error {
	return runRecordsListWithOptions(cmd, recordsOptions{
		Platform: recordsPlatform,
		Limit:    recordsLimit,
		Status:   recordsStatus,
	})
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

func newUpdateAccountCmd() *cobra.Command {
	opts := updateAccountOptions{}
	cmd := &cobra.Command{
		Use:   "update-account <account_id>",
		Short: "更新账号代理或备注",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := opts.body()
			if len(body) == 0 {
				return yxerrors.Usage("update account request must not be empty", nil).
					WithHint("请至少传入 --proxy-id、--kuaidaili-area、--remark 或 --group。")
			}
			if opts.DryRun {
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
	cmd.Flags().StringVar(&opts.ProxyID, "proxy-id", "", "team proxy id")
	cmd.Flags().StringVar(&opts.KuaidailiArea, "kuaidaili-area", "", "built-in proxy area code")
	cmd.Flags().StringVar(&opts.Remark, "remark", "", "account remark")
	cmd.Flags().StringSliceVar(&opts.Groups, "group", nil, "group id; repeat or comma-separate for multiple")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview update request without performing the write")
	return cmd
}

func updateAccountBody() map[string]interface{} {
	return updateAccountOptions{
		ProxyID:       updateAccountProxyID,
		KuaidailiArea: updateAccountKuaidailiArea,
		Remark:        updateAccountRemark,
		Groups:        updateAccountGroups,
	}.body()
}

type updateAccountOptions struct {
	ProxyID       string
	KuaidailiArea string
	Remark        string
	Groups        []string
	DryRun        bool
}

func (opts updateAccountOptions) body() map[string]interface{} {
	body := map[string]interface{}{}
	if strings.TrimSpace(opts.ProxyID) != "" {
		body["proxyId"] = opts.ProxyID
	}
	if strings.TrimSpace(opts.KuaidailiArea) != "" {
		body["kuaidailiArea"] = opts.KuaidailiArea
	}
	if strings.TrimSpace(opts.Remark) != "" {
		body["remark"] = opts.Remark
	}
	if len(opts.Groups) > 0 {
		groups := make([]string, 0, len(opts.Groups))
		for _, group := range opts.Groups {
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
