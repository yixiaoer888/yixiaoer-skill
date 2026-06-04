package cmd

import "testing"

func TestFilterAccountsByNameAndStatus(t *testing.T) {
	accounts := []map[string]interface{}{
		{"platformAccountId": "acc_1", "name": "主账号", "status": float64(1)},
		{"platformAccountId": "acc_2", "nickname": "备用账号", "status": float64(0)},
		{"platformAccountId": "acc_3", "remark": "主账号二", "status": float64(1)},
	}

	filtered := filterAccounts(accounts, "主账号", 1)
	if len(filtered) != 2 {
		t.Fatalf("expected two matching accounts, got %d", len(filtered))
	}
	if accountName(filtered[0]) != "主账号" {
		t.Fatalf("unexpected first account name: %s", accountName(filtered[0]))
	}
}

func TestFilterAccountsUsesPlatformAccountNameAndLoginStatus(t *testing.T) {
	accounts := []map[string]interface{}{
		{"platformAccountId": "acc_1", "platformAccountName": "抖音主账号", "loginStatus": float64(1)},
		{"platformAccountId": "acc_2", "platformAccountName": "快手主账号", "loginStatus": float64(1)},
	}

	filtered := filterAccounts(accounts, "抖音", 1)
	if len(filtered) != 1 {
		t.Fatalf("expected one matching account, got %d", len(filtered))
	}
	if accountName(filtered[0]) != "抖音主账号" {
		t.Fatalf("unexpected account name: %s", accountName(filtered[0]))
	}
}

func TestFilterAccountsSortsOnlineAccountsFirst(t *testing.T) {
	accounts := []map[string]interface{}{
		{"platformAccountId": "acc_9", "platformAccountName": "吹牛不算牛y", "loginStatus": float64(2)},
		{"platformAccountId": "acc_10", "platformAccountName": "Max8862", "loginStatus": float64(1)},
		{"platformAccountId": "acc_2", "platformAccountName": "Alpha", "loginStatus": float64(1)},
	}

	filtered := filterAccounts(accounts, "", -1)
	if len(filtered) != 3 {
		t.Fatalf("expected three accounts, got %d", len(filtered))
	}
	if accountName(filtered[0]) != "Alpha" {
		t.Fatalf("unexpected first account: %s", accountName(filtered[0]))
	}
	if accountName(filtered[1]) != "Max8862" {
		t.Fatalf("unexpected second account: %s", accountName(filtered[1]))
	}
	if accountName(filtered[2]) != "吹牛不算牛y" {
		t.Fatalf("unexpected third account: %s", accountName(filtered[2]))
	}
}
