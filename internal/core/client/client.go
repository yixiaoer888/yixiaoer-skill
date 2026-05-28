package client

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
)

type Client = api.Client
type UploadResult = api.UploadResult
type PrepareData = api.PrepareData

func New(cfg config.Config) *Client {
	return api.NewClient(cfg)
}

func Query(endpoint string, params map[string]string) string {
	return api.Query(endpoint, params)
}

func DataOrSelf(value map[string]interface{}) interface{} {
	return api.DataOrSelf(value)
}

func AccountID(account map[string]interface{}) string {
	return api.AccountID(account)
}

func AccountStatus(account map[string]interface{}) int {
	return api.AccountStatus(account)
}

func FilterOnlineAccounts(accounts []map[string]interface{}) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, account := range accounts {
		if AccountStatus(account) == 1 {
			filtered = append(filtered, account)
		}
	}
	return filtered
}
