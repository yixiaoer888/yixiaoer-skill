package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type PrepareData struct {
	Platform        string                   `json:"platform"`
	Type            string                   `json:"type"`
	Accounts        []map[string]interface{} `json:"accounts,omitempty"`
	Categories      interface{}              `json:"categories"`
	DefaultFormType string                   `json:"defaultFormType"`
	Workflow        string                   `json:"workflow"`
	DocsIndex       string                   `json:"docsIndex"`
	PlatformDoc     string                   `json:"platformDoc"`
	Schema          string                   `json:"schema"`
	RootSchema      string                   `json:"rootSchema"`
	Form            interface{}              `json:"form,omitempty"`
}

func (c *Client) Categories(accountID, publishType string) (interface{}, error) {
	result, err := c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/categories", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
	}))
	if err != nil {
		return nil, err
	}
	return normalizeCategoryTree(result), nil
}

func (c *Client) Locations(accountID, keyword, locationType, nextPage string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/location", accountID), map[string]string{
		"locationType": locationType,
		"keyWord":      keyword,
		"nextPage":     nextPage,
	}))
}

func (c *Client) Music(accountID, keyword, categoryID, categoryName, nextPage string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/music", accountID), map[string]string{
		"keyWord":      keyword,
		"categoryId":   categoryID,
		"categoryName": categoryName,
		"nextPage":     nextPage,
	}))
}

func (c *Client) MusicCategories(accountID string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/music/category", accountID), nil))
}

func (c *Client) Goods(accountID, keyword, nextPage string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/goods", accountID), map[string]string{
		"keyWord":  keyword,
		"nextPage": nextPage,
	}))
}

func (c *Client) GoodsDetail(accountID, productURL string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/goods-detail", accountID), map[string]string{
		"url": productURL,
	}))
}

func (c *Client) Entitlements(accountID string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/entitlements", accountID), nil))
}

func (c *Client) Collections(accountID, publishType string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/collections", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
	}))
}

func (c *Client) MiniApps(accountID, keyword string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/mini-apps", accountID), map[string]string{
		"keyWord": keyword,
	}))
}

func (c *Client) SyncApps(accountID string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/sync-apps", accountID), nil))
}

func (c *Client) Games(accountID, keyword string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/games", accountID), map[string]string{
		"keyWord": keyword,
	}))
}

func (c *Client) HotEvents(accountID, publishType string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/hot-events", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
	}))
}

func (c *Client) Groups(accountID string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/group-chats", accountID), nil))
}

type MembersOptions struct {
	Page     int
	Size     int
	Statuses []string
	KeyWords string
	Role     string
}

func (c *Client) Members(opts MembersOptions) (interface{}, error) {
	values := url.Values{}
	setIfPositive(values, "page", opts.Page)
	setIfPositive(values, "size", opts.Size)
	setIfNotEmpty(values, "keyWords", opts.KeyWords)
	setIfNotEmpty(values, "role", opts.Role)
	for _, status := range opts.Statuses {
		if status != "" {
			values.Add("statuses", status)
		}
	}
	return c.queryData(QueryValues("/members", values))
}

type AccountGroupOptions struct {
	Page int
	Size int
}

func (c *Client) AccountGroups(opts AccountGroupOptions) (interface{}, error) {
	values := url.Values{}
	setIfPositive(values, "page", opts.Page)
	setIfPositive(values, "size", opts.Size)
	return c.queryData(QueryValues("/groups", values))
}

func (c *Client) CreateAccountGroup(body map[string]interface{}) (interface{}, error) {
	var result interface{}
	if err := c.Post("/groups", body, &result); err != nil {
		return nil, err
	}
	if typed, ok := result.(map[string]interface{}); ok {
		return DataOrSelf(typed), nil
	}
	return result, nil
}

func (c *Client) UpdateAccountGroup(groupID string, body map[string]interface{}) (interface{}, error) {
	var result interface{}
	if err := c.Patch(fmt.Sprintf("/groups/%s", groupID), body, &result); err != nil {
		return nil, err
	}
	if typed, ok := result.(map[string]interface{}); ok {
		return DataOrSelf(typed), nil
	}
	return result, nil
}

func (c *Client) DeleteAccountGroup(groupID string) (interface{}, error) {
	var result interface{}
	if err := c.Delete(fmt.Sprintf("/groups/%s", groupID), &result); err != nil {
		return nil, err
	}
	if typed, ok := result.(map[string]interface{}); ok {
		return DataOrSelf(typed), nil
	}
	return result, nil
}

func (c *Client) Activities(accountID, publishType, categoryID, keyword string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/activities", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
		"categoryId":  categoryID,
		"keyWord":     keyword,
	}))
}

func (c *Client) Challenges(accountID, keyword, publishType, nextPage string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/challenges", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
		"keyWord":     keyword,
		"nextPage":    nextPage,
	}))
}

func (c *Client) Records(platform, limit, status string) (interface{}, error) {
	return c.queryData(Query("/v2/taskSets", map[string]string{
		"size":     limit,
		"platform": platformutil.ChineseName(platform),
		"status":   status,
	}))
}

func (c *Client) Details(taskSetID string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/v2/taskSets/%s/tasks", taskSetID), nil))
}

type AccountOverviewOptions struct {
	Platform    string
	Name        string
	Group       string
	LoginStatus string
	MemberIDs   []string
	Page        int
	Size        int
}

func (c *Client) AccountOverviews(opts AccountOverviewOptions) (interface{}, error) {
	values := url.Values{}
	values.Set("platform", platformutil.ChineseName(opts.Platform))
	setIfNotEmpty(values, "name", opts.Name)
	setIfNotEmpty(values, "group", opts.Group)
	setIfNotEmpty(values, "loginStatus", opts.LoginStatus)
	setIfPositive(values, "page", opts.Page)
	setIfPositive(values, "size", opts.Size)
	for _, id := range opts.MemberIDs {
		if id != "" {
			values.Add("memberIds", id)
		}
	}
	return c.queryData(QueryValues("/platform-accounts/overviews-v2", values))
}

type ContentOverviewOptions struct {
	Platform          string
	PlatformAccountID string
	PublishUserID     string
	Type              string
	Title             string
	PublishStartTime  string
	PublishEndTime    string
	Page              int
	Size              int
}

// AccountIncrementOptions describes the date window accepted by the
// account-data incremental endpoint. The endpoint expects Unix milliseconds.
type AccountIncrementOptions struct {
	StartTime int64
	EndTime   int64
	GroupID   string
	Platform  string
	Name      string
}

func (c *Client) AccountIncrements(opts AccountIncrementOptions) (interface{}, error) {
	values := url.Values{}
	values.Set("startTime", strconv.FormatInt(opts.StartTime, 10))
	values.Set("endTime", strconv.FormatInt(opts.EndTime, 10))
	setIfNotEmpty(values, "group", opts.GroupID)
	setIfNotEmpty(values, "platform", platformutil.ChineseName(opts.Platform))
	setIfNotEmpty(values, "name", opts.Name)
	return c.queryDataWithHeaders(QueryValues("/overview/incremental", values), overviewHeaders())
}

// DMMessageStatsOptions describes the private-message statistics endpoint.
// Platform account IDs are optional; omitting them requests the team trend.
type DMMessageStatsOptions struct {
	StartTime          int64
	EndTime            int64
	PlatformAccountIDs []string
	Platform           string
}

func (c *Client) DMMessageStats(opts DMMessageStatsOptions) (interface{}, error) {
	values := url.Values{}
	values.Set("startTime", strconv.FormatInt(opts.StartTime, 10))
	values.Set("endTime", strconv.FormatInt(opts.EndTime, 10))
	ids := make([]string, 0, len(opts.PlatformAccountIDs))
	for _, id := range opts.PlatformAccountIDs {
		if strings.TrimSpace(id) != "" {
			ids = append(ids, strings.TrimSpace(id))
		}
	}
	setIfNotEmpty(values, "platformAccountIds", strings.Join(ids, ","))
	setIfNotEmpty(values, "platform", platformutil.ChineseName(opts.Platform))
	return c.queryDataWithHeaders(QueryValues("/social/dm-stats", values), overviewHeaders())
}

// ManagedAccountOptions controls the all-platform account management detail
// query used by the overview dashboard.
type ManagedAccountOptions struct {
	Platform string
	Page     int
	Size     int
}

func (c *Client) ManagedAccounts(opts ManagedAccountOptions) (interface{}, error) {
	values := url.Values{}
	setIfNotEmpty(values, "platform", platformutil.ChineseName(opts.Platform))
	setIfPositive(values, "page", opts.Page)
	setIfPositive(values, "size", opts.Size)
	return c.queryDataWithHeaders(QueryValues("/v2/platform/accounts", values), overviewHeaders())
}

// AccountIncrementDashboard combines the three read-only sections displayed
// by the web overview while preserving the existing incremental data keys.
func (c *Client) AccountIncrementDashboard(opts AccountIncrementOptions) (interface{}, error) {
	increments, err := c.AccountIncrements(opts)
	if err != nil {
		return nil, err
	}
	accountIDs := incrementalAccountIDs(increments)
	var dmStats interface{} = map[string]interface{}{
		"summary": []interface{}{},
		"trend":   []interface{}{},
	}
	if len(accountIDs) > 0 {
		dmStats, err = c.DMMessageStats(DMMessageStatsOptions{
			StartTime:          opts.StartTime,
			EndTime:            opts.EndTime,
			PlatformAccountIDs: accountIDs,
		})
		if err != nil {
			return nil, err
		}
	}
	managed, err := c.ManagedAccounts(ManagedAccountOptions{Page: 1, Size: 1000})
	if err != nil {
		return nil, err
	}
	data, ok := increments.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"incremental":     increments,
			"dmMessageStats":  dmStats,
			"managedAccounts": managed,
		}, nil
	}
	enrichIncrementalAccounts(data, dmStats)
	normalizeDMMessageStats(dmStats, data, opts.StartTime, opts.EndTime)
	data["dmMessageStats"] = dmStats
	data["managedAccounts"] = managed
	return data, nil
}

func incrementalAccountIDs(value interface{}) []string {
	data, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	accounts, ok := data["accounts"].([]interface{})
	if !ok {
		return nil
	}
	ids := make([]string, 0, len(accounts))
	for _, item := range accounts {
		account, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := account["platformAccountId"].(string); ok && id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func enrichIncrementalAccounts(data map[string]interface{}, dmStats interface{}) {
	statsByAccount := map[string]map[string]interface{}{}
	if stats, ok := dmStats.(map[string]interface{}); ok {
		if summary, ok := stats["summary"].([]interface{}); ok {
			for _, item := range summary {
				row, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				if id, ok := row["platformAccountId"].(string); ok && id != "" {
					statsByAccount[id] = row
				}
			}
		}
	}
	accounts, ok := data["accounts"].([]interface{})
	if !ok {
		return
	}
	for _, item := range accounts {
		account, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := account["platformAccountId"].(string)
		stats := statsByAccount[id]
		account["dmInCount"] = numberOrZero(stats, "inCount")
		account["dmOutCount"] = numberOrZero(stats, "outCount")
		if status, ok := account["status"].(float64); ok {
			if status == 1 {
				account["statusLabel"] = "正常"
			} else {
				account["statusLabel"] = "失效"
			}
		}
		if updatedAt := accountUpdateTimestamp(account); updatedAt > 0 {
			account["dataUpdatedAt"] = updatedAt
			account["dataUpdatedTime"] = time.UnixMilli(updatedAt).
				In(time.FixedZone("Asia/Shanghai", 8*60*60)).
				Format("2006-01-02 15:04:05")
		}
	}
}

func accountUpdateTimestamp(account map[string]interface{}) int64 {
	for _, key := range []string{"overviewUpdatedAt", "updatedAt"} {
		value := numberOrZero(account, key)
		if value > 0 {
			return int64(value)
		}
	}
	return 0
}

func normalizeDMMessageStats(value interface{}, incrementalData map[string]interface{}, startTime, endTime int64) {
	stats, ok := value.(map[string]interface{})
	if !ok {
		return
	}
	var received, sent float64
	if summary, ok := stats["summary"].([]interface{}); ok {
		for _, item := range summary {
			row, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			received += numberOrZero(row, "inCount")
			sent += numberOrZero(row, "outCount")
		}
	}
	stats["totals"] = map[string]interface{}{"inCount": received, "outCount": sent}
	var incrementalReceived, incrementalSent float64
	if accounts, ok := incrementalData["accounts"].([]interface{}); ok {
		for _, item := range accounts {
			row, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			incrementalReceived += numberOrZero(row, "dmInCount")
			incrementalSent += numberOrZero(row, "dmOutCount")
		}
	}
	stats["incrementalAccountTotals"] = map[string]interface{}{
		"inCount":  incrementalReceived,
		"outCount": incrementalSent,
	}

	trendByDate := map[string]map[string]interface{}{}
	if trend, ok := stats["trend"].([]interface{}); ok {
		for _, item := range trend {
			row, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if date, ok := row["date"].(string); ok {
				trendByDate[date] = row
			}
		}
	}
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	start := time.UnixMilli(startTime).In(location)
	end := time.UnixMilli(endTime).In(location)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, location)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, location)
	daily := make([]interface{}, 0)
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		date := day.Format("2006-01-02")
		row := trendByDate[date]
		daily = append(daily, map[string]interface{}{
			"date":     date,
			"inCount":  numberOrZero(row, "inCount"),
			"outCount": numberOrZero(row, "outCount"),
		})
	}
	stats["dailyTrend"] = daily
}

func numberOrZero(row map[string]interface{}, key string) float64 {
	if row == nil {
		return 0
	}
	switch value := row[key].(type) {
	case float64:
		return value
	case float32:
		return float64(value)
	case int:
		return float64(value)
	case int64:
		return float64(value)
	default:
		return 0
	}
}

func overviewHeaders() map[string]string {
	return map[string]string{
		"x-client":   "web",
		"x-platform": "windows",
		"x-version":  "2.7.3",
	}
}

// ShanghaiDateRange converts inclusive YYYY-MM-DD dates into the inclusive
// millisecond range used by the web endpoint.
func ShanghaiDateRange(startDate, endDate string) (int64, int64, error) {
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	start, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(startDate), location)
	if err != nil {
		return 0, 0, yxerrors.Usage(fmt.Sprintf("invalid start date %q: use YYYY-MM-DD", startDate), nil).
			WithHint("日期格式必须为 YYYY-MM-DD。")
	}
	end, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(endDate), location)
	if err != nil {
		return 0, 0, yxerrors.Usage(fmt.Sprintf("invalid end date %q: use YYYY-MM-DD", endDate), nil).
			WithHint("日期格式必须为 YYYY-MM-DD。")
	}
	if start.After(end) {
		return 0, 0, yxerrors.Usage(fmt.Sprintf("start date %q must not be after end date %q", startDate, endDate), nil).
			WithHint("开始日期不能晚于结束日期。")
	}
	return start.UnixMilli(), end.Add(24*time.Hour - time.Millisecond).UnixMilli(), nil
}

func (c *Client) ContentOverviews(opts ContentOverviewOptions) (interface{}, error) {
	values := url.Values{}
	setIfNotEmpty(values, "platform", platformutil.ChineseName(opts.Platform))
	setIfNotEmpty(values, "platformAccountId", opts.PlatformAccountID)
	setIfNotEmpty(values, "publishUserId", opts.PublishUserID)
	setIfNotEmpty(values, "type", opts.Type)
	setIfNotEmpty(values, "title", opts.Title)
	setIfNotEmpty(values, "publishStartTime", opts.PublishStartTime)
	setIfNotEmpty(values, "publishEndTime", opts.PublishEndTime)
	setIfPositive(values, "page", opts.Page)
	setIfPositive(values, "size", opts.Size)
	return c.queryData(QueryValues("/contents/overviews", values))
}

func (c *Client) Proxies(size string) (interface{}, error) {
	if size == "" {
		size = "9999"
	}
	return c.queryData(Query("/proxys", map[string]string{"size": size}))
}

func (c *Client) ProxyAreas() (interface{}, error) {
	return c.queryData(Query("/daili/areas", nil))
}

func (c *Client) UpdateAccount(accountID string, body map[string]interface{}) (interface{}, error) {
	var result interface{}
	if err := c.Patch(fmt.Sprintf("/platform-accounts/%s", accountID), body, &result); err != nil {
		return nil, err
	}
	if typed, ok := result.(map[string]interface{}); ok {
		return DataOrSelf(typed), nil
	}
	return result, nil
}

func (c *Client) Prepare(platform, publishType string) (PrepareData, error) {
	accounts, err := c.Accounts(platform)
	if err != nil {
		return PrepareData{}, err
	}
	onlineAccounts := filterOnlineAccounts(accounts)

	var categories interface{}
	if len(onlineAccounts) > 0 && (publishType == "video" || publishType == "article") {
		if result, err := c.Categories(AccountID(onlineAccounts[0]), publishType); err == nil {
			categories = result
		}
	}

	schemaPlatform := platform
	if platformutil.CanonicalKey(platform) == "weixin.account" {
		schemaPlatform = "weixin.account"
	}

	return PrepareData{
		Platform:        platform,
		Type:            publishType,
		Accounts:        onlineAccounts,
		Categories:      categories,
		DefaultFormType: "task",
		Workflow:        fmt.Sprintf("workflows/publish-%s.md", publishType),
		DocsIndex:       fmt.Sprintf("skills/yixiaoer/references/platforms/%s/index.md", publishType),
		PlatformDoc:     fmt.Sprintf("skills/yixiaoer/references/platforms/%s/%s", publishType, platformDocFileName(platform, publishType)),
		Schema:          fmt.Sprintf("schemas/platforms/%s.%s.schema.json", schemaPlatform, schemaTypeName(publishType)),
		RootSchema:      "schemas/publish.schema.json",
	}, nil
}

func setIfNotEmpty(values url.Values, key, value string) {
	if value != "" {
		values.Set(key, value)
	}
}

func setIfPositive(values url.Values, key string, value int) {
	if value > 0 {
		values.Set(key, strconv.Itoa(value))
	}
}

func platformDocFileName(platform, publishType string) string {
	if publishType == "imageText" {
		switch platformutil.CanonicalKey(platform) {
		case "shipinhao":
			return "shipinhao.md"
		case "weixin.account":
			return "weixingongzhonghao.md"
		}
	}
	return platform + ".md"
}

func (c *Client) queryData(endpoint string) (interface{}, error) {
	return c.queryDataWithHeaders(endpoint, nil)
}

func (c *Client) queryDataWithHeaders(endpoint string, headers map[string]string) (interface{}, error) {
	var result interface{}
	if err := c.GetWithHeaders(endpoint, headers, &result); err != nil {
		return nil, err
	}
	if typed, ok := result.(map[string]interface{}); ok {
		return DataOrSelf(typed), nil
	}
	return result, nil
}

func normalizeCategoryTree(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		if rawList, ok := typed["dataList"].([]interface{}); ok {
			rawList = normalizeCategoryItems(rawList)
			tree := buildCategoryTreeFromParentIDs(rawList)
			if tree != nil {
				out := cloneInterfaceMap(typed)
				out["dataList"] = tree
				return out
			}
			out := cloneInterfaceMap(typed)
			out["dataList"] = rawList
			return out
		}
	case []interface{}:
		typed = normalizeCategoryItems(typed)
		if tree := buildCategoryTreeFromParentIDs(typed); tree != nil {
			return tree
		}
		return typed
	}
	return value
}

// normalizeCategoryItems exposes the stable query shape. Category endpoints
// have historically returned either raw {id,name} objects or the frontend
// {id,text,raw} wrapper; callers should only have to consume the query object
// {yixiaoerId,yixiaoerName,raw}.
func normalizeCategoryItems(items []interface{}) []interface{} {
	out := make([]interface{}, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			out = append(out, item)
			continue
		}
		out = append(out, normalizeCategoryItem(obj))
	}
	return out
}

func normalizeCategoryItem(obj map[string]interface{}) map[string]interface{} {
	if obj == nil {
		return obj
	}
	// Unwrap the publish/frontend representation when its raw value is already
	// the canonical query object.
	if raw, ok := obj["raw"].(map[string]interface{}); ok && isCategoryQueryObject(raw) {
		canonical := cloneInterfaceMap(raw)
		copyCategoryChildren(obj, canonical)
		return canonical
	}
	if isCategoryQueryObject(obj) {
		canonical := cloneInterfaceMap(obj)
		copyCategoryChildren(obj, canonical)
		return canonical
	}

	id := stringField(obj, "id")
	if id == "" {
		id = stringField(obj, "value")
	}
	name := stringField(obj, "name")
	if name == "" {
		name = stringField(obj, "text")
	}
	if name == "" {
		name = stringField(obj, "label")
	}
	if id == "" || name == "" {
		return obj
	}
	canonical := map[string]interface{}{
		"yixiaoerId":   id,
		"yixiaoerName": name,
		"raw":          cloneInterfaceMap(obj),
	}
	copyCategoryChildren(obj, canonical)
	return canonical
}

func isCategoryQueryObject(obj map[string]interface{}) bool {
	if obj == nil {
		return false
	}
	_, hasID := obj["yixiaoerId"]
	_, hasName := obj["yixiaoerName"]
	_, hasRaw := obj["raw"]
	return hasID && hasName && hasRaw
}

func copyCategoryChildren(from, to map[string]interface{}) {
	for _, key := range []string{"child", "children"} {
		if children, ok := from[key].([]interface{}); ok {
			to[key] = normalizeCategoryItems(children)
		}
	}
}

func buildCategoryTreeFromParentIDs(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}
	childrenByParent := map[string][]map[string]interface{}{}
	roots := []map[string]interface{}{}
	hasParentID := false
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return nil
		}
		id := stringField(obj, "yixiaoerId")
		if id == "" {
			return nil
		}
		parentID, hasParent := categoryParentID(obj)
		if hasParent {
			hasParentID = true
		}
		cloned := cloneInterfaceMap(obj)
		delete(cloned, "child")
		delete(cloned, "children")
		if parentID == "" {
			roots = append(roots, cloned)
			continue
		}
		childrenByParent[parentID] = append(childrenByParent[parentID], cloned)
	}
	if !hasParentID {
		return nil
	}
	var attachChildren func(map[string]interface{})
	attachChildren = func(item map[string]interface{}) {
		id := stringField(item, "yixiaoerId")
		children := childrenByParent[id]
		if len(children) == 0 {
			return
		}
		out := make([]interface{}, 0, len(children))
		for _, child := range children {
			attachChildren(child)
			out = append(out, child)
		}
		item["child"] = out
	}
	out := make([]interface{}, 0, len(roots))
	for _, root := range roots {
		attachChildren(root)
		out = append(out, root)
	}
	return out
}

func categoryParentID(item map[string]interface{}) (string, bool) {
	raw, ok := item["raw"].(map[string]interface{})
	if !ok {
		return "", false
	}
	parent, ok := raw["parentId"]
	if !ok || parent == nil {
		return "", ok
	}
	return strings.TrimSpace(fmt.Sprint(parent)), true
}

func stringField(item map[string]interface{}, key string) string {
	value, ok := item[key]
	if !ok || value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func cloneInterfaceMap(input map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}

func filterOnlineAccounts(accounts []map[string]interface{}) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, account := range accounts {
		if AccountStatus(account) == 1 {
			filtered = append(filtered, account)
		}
	}
	return filtered
}

func schemaTypeName(publishType string) string {
	return publishType
}
