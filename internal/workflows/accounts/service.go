package accounts

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) List(platform, name string, status int) ([]map[string]interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	accounts, err := client.New(cfg).Accounts(platform)
	if err != nil {
		return nil, err
	}
	return FilterAccounts(accounts, name, status), nil
}

func FilterAccounts(accounts []map[string]interface{}, name string, status int) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, account := range accounts {
		if name != "" && !strings.Contains(AccountName(account), name) {
			continue
		}
		if status >= 0 && client.AccountStatus(account) != status {
			continue
		}
		filtered = append(filtered, account)
	}
	sortAccounts(filtered)
	return filtered
}

func sortAccounts(accounts []map[string]interface{}) {
	sort.SliceStable(accounts, func(i, j int) bool {
		left := accounts[i]
		right := accounts[j]

		leftStatus := client.AccountStatus(left)
		rightStatus := client.AccountStatus(right)
		if leftStatus != rightStatus {
			return leftStatus == 1
		}

		leftName := strings.ToLower(AccountName(left))
		rightName := strings.ToLower(AccountName(right))
		if leftName != rightName {
			return leftName < rightName
		}

		leftID := strings.ToLower(client.AccountID(left))
		rightID := strings.ToLower(client.AccountID(right))
		return leftID < rightID
	})
}

func AccountName(account map[string]interface{}) string {
	for _, key := range []string{"platformAccountName", "name", "nickname", "remark", "platformAuthorId"} {
		if value := fmt.Sprint(account[key]); value != "" && value != "<nil>" {
			return value
		}
	}
	return "未命名"
}
