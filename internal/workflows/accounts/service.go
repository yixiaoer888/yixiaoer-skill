package accounts

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
)

type Service struct {
	rt *app.Runtime
}

type ListOptions struct {
	Page int
	Size int
	All  bool
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func (s Service) List(platform, name string, status int) ([]map[string]interface{}, error) {
	return s.ListWithOptions(platform, name, status, ListOptions{})
}

func (s Service) ListWithOptions(platform, name string, status int, opts ListOptions) ([]map[string]interface{}, error) {
	apiClient := s.rt.Client
	var err error
	var accounts []map[string]interface{}
	if opts.All {
		accounts, err = apiClient.AccountsAll(platform, opts.Size)
	} else {
		accounts, _, err = apiClient.AccountsPage(platform, opts.Page, opts.Size)
	}
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
		if status >= 0 && api.AccountStatus(account) != status {
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

		leftStatus := api.AccountStatus(left)
		rightStatus := api.AccountStatus(right)
		if leftStatus != rightStatus {
			return leftStatus == 1
		}

		leftName := strings.ToLower(AccountName(left))
		rightName := strings.ToLower(AccountName(right))
		if leftName != rightName {
			return leftName < rightName
		}

		leftID := strings.ToLower(api.AccountID(left))
		rightID := strings.ToLower(api.AccountID(right))
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
