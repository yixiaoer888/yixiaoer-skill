package api

import (
	"fmt"
	"strconv"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const accountsPageSize = 20

func (c *Client) Accounts(platform string) ([]map[string]interface{}, error) {
	return c.AccountsAll(platform, accountsPageSize)
}

func (c *Client) AccountsAll(platform string, size int) ([]map[string]interface{}, error) {
	if size <= 0 {
		size = accountsPageSize
	}

	var all []map[string]interface{}
	for page := 1; ; page++ {
		accounts, meta, err := c.AccountsPage(platform, page, size)
		if err != nil {
			return nil, err
		}
		all = append(all, accounts...)
		if !shouldFetchNextAccountsPage(meta, len(accounts), size) {
			return all, nil
		}
	}
}

func (c *Client) AccountsPage(platform string, page, size int) ([]map[string]interface{}, accountsPageMeta, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = accountsPageSize
	}

	baseEndpoint := "/v2/platform/accounts"
	params := map[string]string{
		"page": strconv.Itoa(page),
		"size": strconv.Itoa(size),
	}
	if platform != "" {
		params["platform"] = platformutil.ChineseName(platform)
	}

	endpoint := Query(baseEndpoint, params)
	var result map[string]interface{}
	if err := c.Get(endpoint, &result); err != nil {
		return nil, accountsPageMeta{}, err
	}

	data := DataOrSelf(result)
	accounts, err := normalizeAccounts(data)
	if err != nil {
		return nil, accountsPageMeta{}, err
	}
	return accounts, extractAccountsPageMeta(data), nil
}

func normalizeAccounts(data interface{}) ([]map[string]interface{}, error) {
	switch typed := data.(type) {
	case []interface{}:
		accounts := make([]map[string]interface{}, 0, len(typed))
		for _, item := range typed {
			if account, ok := item.(map[string]interface{}); ok {
				accounts = append(accounts, account)
			}
		}
		return accounts, nil
	case map[string]interface{}:
		for _, key := range []string{"data", "list", "records", "items", "rows"} {
			if nested, ok := typed[key]; ok {
				if accounts, err := normalizeAccounts(nested); err == nil {
					return accounts, nil
				}
			}
		}
		if hasAccountIdentity(typed) {
			return []map[string]interface{}{typed}, nil
		}
		return nil, yxerrors.Remote("unexpected accounts response", typed).
			WithCategory("remote_response")
	default:
		return nil, yxerrors.Remote("unexpected accounts response", typed).
			WithCategory("remote_response")
	}
}

func hasAccountIdentity(account map[string]interface{}) bool {
	for _, key := range []string{"platformAccountId", "id", "platformAccountName", "name", "nickname", "remark"} {
		if value, ok := account[key]; ok && value != nil && fmt.Sprint(value) != "" && fmt.Sprint(value) != "<nil>" {
			return true
		}
	}
	return false
}

type accountsPageMeta struct {
	total   int
	page    int
	size    int
	hasNext *bool
}

func extractAccountsPageMeta(data interface{}) accountsPageMeta {
	meta := accountsPageMeta{}
	typed, ok := data.(map[string]interface{})
	if !ok {
		return meta
	}

	meta.total = firstPositiveInt(typed, "totalSize")
	meta.page = firstPositiveInt(typed, "page", "pageNum", "current", "currentPage")
	meta.size = firstPositiveInt(typed, "size", "pageSize", "limit", "perPage")
	if pages := firstPositiveInt(typed, "totalPage"); pages > 0 {
		current := meta.page
		if current == 0 {
			current = 1
		}
		hasNext := current < pages
		meta.hasNext = &hasNext
	}
	if meta.hasNext == nil && (meta.total > 0 || meta.page > 0 || meta.size > 0) {
		current := meta.page
		if current == 0 {
			current = 1
		}
		size := meta.size
		if size == 0 {
			size = accountsPageSize
		}
		hasNext := current*size < meta.total
		meta.hasNext = &hasNext
	}
	return meta
}

func shouldFetchNextAccountsPage(meta accountsPageMeta, _ int, _ int) bool {
	if meta.hasNext != nil {
		return *meta.hasNext
	}
	return false
}

func firstPositiveInt(data map[string]interface{}, keys ...string) int {
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if parsed := positiveInt(value); parsed > 0 {
				return parsed
			}
		}
	}
	return 0
}

func positiveInt(value interface{}) int {
	switch typed := value.(type) {
	case float64:
		if typed > 0 {
			return int(typed)
		}
	case int:
		if typed > 0 {
			return typed
		}
	case int64:
		if typed > 0 {
			return int(typed)
		}
	case string:
		if typed == "" {
			return 0
		}
		if parsed, err := strconv.Atoi(typed); err == nil && parsed > 0 {
			return parsed
		}
	}
	return 0
}

func AccountID(account map[string]interface{}) string {
	for _, key := range []string{"platformAccountId", "id"} {
		if value, ok := account[key]; ok && value != nil {
			return fmt.Sprint(value)
		}
	}
	return ""
}

func AccountStatus(account map[string]interface{}) int {
	var value interface{}
	var ok bool
	for _, key := range []string{"status", "loginStatus"} {
		value, ok = account[key]
		if ok {
			break
		}
	}
	if !ok || value == nil {
		return 0
	}
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	default:
		return 0
	}
}
