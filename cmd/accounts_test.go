package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

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

func TestAccountsListSubcommandHasListFlags(t *testing.T) {
	accountsCmd := newAccountsCmd()
	listCmd := findChildCommand(t, accountsCmd, "list")
	if listCmd.Flags().Lookup("name") == nil {
		t.Fatal("expected accounts list to expose --name flag")
	}
	if listCmd.Flags().Lookup("status") == nil {
		t.Fatal("expected accounts list to expose --status flag")
	}
	if listCmd.Flags().Lookup("page") == nil {
		t.Fatal("expected accounts list to expose --page flag")
	}
	if listCmd.Flags().Lookup("size") == nil {
		t.Fatal("expected accounts list to expose --size flag")
	}
	if listCmd.Flags().Lookup("all") == nil {
		t.Fatal("expected accounts list to expose --all flag")
	}
}

func TestRunAccountsListRejectsInvalidPage(t *testing.T) {
	err := runAccountsListWithOptions(&cobra.Command{}, nil, accountsListOptions{
		Status: -1,
		Page:   0,
		Size:   20,
	})
	if err == nil {
		t.Fatal("expected error for invalid page")
	}
	if !strings.Contains(err.Error(), "accounts page must be greater than 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAccountsListRejectsInvalidSize(t *testing.T) {
	err := runAccountsListWithOptions(&cobra.Command{}, nil, accountsListOptions{
		Status: -1,
		Page:   1,
		Size:   0,
	})
	if err == nil {
		t.Fatal("expected error for invalid size")
	}
	if !strings.Contains(err.Error(), "accounts size must be greater than 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAccountsCommandExposesUpdateSubcommand(t *testing.T) {
	cmd := newAccountsCmd()
	updateCmd := findChildCommand(t, cmd, "update")
	if updateCmd.Use != "update <account_id>" {
		t.Fatalf("unexpected update command use: %q", updateCmd.Use)
	}
}

func TestAccountsUpdateRejectsListFilterFlags(t *testing.T) {
	cmd := newAccountsCmd()
	cmd.SetArgs([]string{"update", "acc_1", "--proxy-id", "proxy_1", "--name", "ignored", "--dry-run"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected accounts update to reject list filter flags")
	}
	if !strings.Contains(err.Error(), "unknown flag: --name") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func findChildCommand(t *testing.T, parent *cobra.Command, name string) *cobra.Command {
	t.Helper()
	for _, child := range parent.Commands() {
		if child.Name() == name {
			return child
		}
	}
	t.Fatalf("expected %s child command", name)
	return nil
}
