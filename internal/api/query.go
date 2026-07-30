package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
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

	return PrepareData{
		Platform:        platform,
		Type:            publishType,
		Accounts:        onlineAccounts,
		Categories:      categories,
		DefaultFormType: "task",
		Workflow:        fmt.Sprintf("workflows/publish-%s.md", publishType),
		DocsIndex:       fmt.Sprintf("skills/yixiaoer/references/platforms/%s/index.md", publishType),
		PlatformDoc:     fmt.Sprintf("skills/yixiaoer/references/platforms/%s/%s", publishType, platformDocFileName(platform, publishType)),
		Schema:          fmt.Sprintf("schemas/platforms/%s.%s.schema.json", platform, schemaTypeName(publishType)),
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
	if publishType == "imageText" && platform == "shipinhao" {
		return "shipinhao.md"
	}
	return platform + ".md"
}

func (c *Client) queryData(endpoint string) (interface{}, error) {
	var result interface{}
	if err := c.Get(endpoint, &result); err != nil {
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
			tree := buildCategoryTreeFromParentIDs(rawList)
			if tree != nil {
				out := cloneInterfaceMap(typed)
				out["dataList"] = tree
				return out
			}
		}
	case []interface{}:
		if tree := buildCategoryTreeFromParentIDs(typed); tree != nil {
			return tree
		}
	}
	return value
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
