package api

import (
	"fmt"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
)

type PrepareData struct {
	Platform        string                   `json:"platform"`
	Type            string                   `json:"type"`
	Accounts        []map[string]interface{} `json:"accounts"`
	Categories      interface{}              `json:"categories"`
	DefaultFormType string                   `json:"defaultFormType"`
	Workflow        string                   `json:"workflow"`
	DocsIndex       string                   `json:"docsIndex"`
	PlatformDoc     string                   `json:"platformDoc"`
	Schema          string                   `json:"schema"`
	RootSchema      string                   `json:"rootSchema"`
}

func (c *Client) Categories(accountID, publishType string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/categories", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
	}))
}

func (c *Client) Locations(accountID, keyword, locationType string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/location", accountID), map[string]string{
		"locationType": locationType,
		"keyWord":      keyword,
	}))
}

func (c *Client) Music(accountID, keyword string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/music", accountID), map[string]string{
		"keyWord": keyword,
	}))
}

func (c *Client) Goods(accountID, keyword string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/goods", accountID), map[string]string{
		"keyWord": keyword,
	}))
}

func (c *Client) Collections(accountID, publishType string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/collections", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
	}))
}

func (c *Client) Challenges(accountID, keyword, publishType string) (interface{}, error) {
	return c.queryData(Query(fmt.Sprintf("/platform-accounts/%s/challenges", accountID), map[string]string{
		"publishType": schemaTypeName(publishType),
		"keyWord":     keyword,
	}))
}

func (c *Client) Records(platform, limit, status string) (interface{}, error) {
	return c.queryData(Query("/v2/taskSets", map[string]string{
		"size":     limit,
		"platform": platformutil.ChineseName(platform),
		"status":   status,
	}))
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
		DocsIndex:       fmt.Sprintf("docs/publish/%s/index.md", publishType),
		PlatformDoc:     fmt.Sprintf("docs/publish/%s/%s.md", publishType, platform),
		Schema:          fmt.Sprintf("schemas/platforms/%s.%s.schema.json", platform, schemaTypeName(publishType)),
		RootSchema:      "schemas/publish.schema.json",
	}, nil
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
	switch publishType {
	case "image-text":
		return "imageText"
	default:
		return publishType
	}
}
