package api

import (
	"fmt"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func (c *Client) Accounts(platform string) ([]map[string]interface{}, error) {
	endpoint := "/v2/platform/accounts"
	if platform != "" {
		endpoint = Query(endpoint, map[string]string{"platform": platformutil.ChineseName(platform)})
	}
	var result map[string]interface{}
	if err := c.Get(endpoint, &result); err != nil {
		return nil, err
	}
	return normalizeAccounts(DataOrSelf(result))
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
