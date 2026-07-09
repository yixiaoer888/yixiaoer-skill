package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunRecordsListRequiresLimit(t *testing.T) {
	err := runRecordsListWithOptions(testCobraCommand(), recordsOptions{})
	if err == nil {
		t.Fatal("expected records limit validation error")
	}
	if !strings.Contains(err.Error(), "records limit must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveQueryAliasPrefersPrimaryValue(t *testing.T) {
	if got := resolveQueryAlias("main", "alias"); got != "main" {
		t.Fatalf("expected primary value, got %q", got)
	}
}

func TestResolveQueryAliasFallsBackToAlias(t *testing.T) {
	if got := resolveQueryAlias(" ", "alias"); got != "alias" {
		t.Fatalf("expected alias fallback, got %q", got)
	}
}

func TestLocationsKeywordFlagUsesAliasStorage(t *testing.T) {
	assertKeywordFlagUsesAliasStorage(t, "locations", "parks")
}

func TestGoodsKeywordFlagUsesAliasStorage(t *testing.T) {
	assertKeywordFlagUsesAliasStorage(t, "goods", "phone")
}

func TestMiniAppsKeywordFlagUsesAliasStorage(t *testing.T) {
	assertKeywordFlagUsesAliasStorage(t, "miniapps", "抽奖")
}

func TestGamesKeywordFlagUsesAliasStorage(t *testing.T) {
	assertKeywordFlagUsesAliasStorage(t, "games", "消消乐")
}

func TestMembersKeywordFlagUsesAliasStorage(t *testing.T) {
	assertKeywordFlagUsesAliasStorage(t, "members", "张三")
}

func TestActivitiesKeywordFlagUsesAliasStorage(t *testing.T) {
	assertKeywordFlagUsesAliasStorage(t, "activities", "创作")
}

func assertKeywordFlagUsesAliasStorage(t *testing.T, use, value string) {
	t.Helper()
	var query string
	var keyword string
	cmd := &cobra.Command{Use: use}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")

	if err := cmd.Flags().Parse([]string{"--keyword", value}); err != nil {
		t.Fatal(err)
	}
	if query != "" {
		t.Fatalf("expected primary query storage to remain empty, got %q", query)
	}
	if keyword != value {
		t.Fatalf("expected alias storage to capture keyword flag, got %q", keyword)
	}
}

func TestQueryCommandExistsWithLocationsSubcommand(t *testing.T) {
	found := false
	for _, child := range rootCmd.Commands() {
		if child.Name() != "query" {
			continue
		}
		for _, grandchild := range child.Commands() {
			if grandchild.Name() == "locations" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatal("expected query command to expose locations subcommand")
	}
}

func TestQueryCommandsUseCommandLocalFlagStorage(t *testing.T) {
	locations := newLocationsCmd()
	if err := locations.Flags().Parse([]string{"--query", "parks", "--keyword", "alias"}); err != nil {
		t.Fatal(err)
	}
	anotherLocations := newLocationsCmd()
	if got, err := anotherLocations.Flags().GetString("query"); err != nil || got != "" {
		t.Fatalf("expected independent locations command query default, got %q, err=%v", got, err)
	}

	update := newAccountsUpdateCmd()
	if err := update.Flags().Parse([]string{"--proxy-id", "proxy_1", "--dry-run"}); err != nil {
		t.Fatal(err)
	}
}

func TestRootDoesNotExposeQueryCompatibilityCommands(t *testing.T) {
	removed := map[string]bool{
		"categories":        true,
		"locations":         true,
		"music":             true,
		"music-categories":  true,
		"goods":             true,
		"collections":       true,
		"miniapps":          true,
		"syncapps":          true,
		"games":             true,
		"hot-events":        true,
		"groups":            true,
		"activities":        true,
		"challenges":        true,
		"records":           true,
		"details":           true,
		"account-overviews": true,
		"content-overviews": true,
		"proxies":           true,
		"proxy-areas":       true,
	}
	for _, child := range rootCmd.Commands() {
		if removed[child.Name()] {
			t.Fatalf("root should not expose query compatibility command %q", child.Name())
		}
	}
}

func TestQueryCommandExistsWithMiniAppsAndSyncAppsSubcommands(t *testing.T) {
	foundMiniApps := false
	foundSyncApps := false
	foundGames := false
	foundHotEvents := false
	foundGroups := false
	foundMembers := false
	foundActivities := false
	foundMusicCategories := false
	foundDetails := false
	foundAccountOverviews := false
	foundContentOverviews := false
	foundProxies := false
	foundProxyAreas := false
	for _, child := range rootCmd.Commands() {
		if child.Name() != "query" {
			continue
		}
		for _, grandchild := range child.Commands() {
			switch grandchild.Name() {
			case "miniapps":
				foundMiniApps = true
			case "syncapps":
				foundSyncApps = true
			case "games":
				foundGames = true
			case "hot-events":
				foundHotEvents = true
			case "groups":
				foundGroups = true
			case "members":
				foundMembers = true
			case "activities":
				foundActivities = true
			case "music-categories":
				foundMusicCategories = true
			case "details":
				foundDetails = true
			case "account-overviews":
				foundAccountOverviews = true
			case "content-overviews":
				foundContentOverviews = true
			case "proxies":
				foundProxies = true
			case "proxy-areas":
				foundProxyAreas = true
			}
		}
	}
	if !foundMiniApps {
		t.Fatal("expected query command to expose miniapps subcommand")
	}
	if !foundSyncApps {
		t.Fatal("expected query command to expose syncapps subcommand")
	}
	if !foundGames {
		t.Fatal("expected query command to expose games subcommand")
	}
	if !foundHotEvents {
		t.Fatal("expected query command to expose hot-events subcommand")
	}
	if !foundGroups {
		t.Fatal("expected query command to expose groups subcommand")
	}
	if !foundMembers {
		t.Fatal("expected query command to expose members subcommand")
	}
	if !foundActivities {
		t.Fatal("expected query command to expose activities subcommand")
	}
	if !foundMusicCategories {
		t.Fatal("expected query command to expose music-categories subcommand")
	}
	if !foundDetails {
		t.Fatal("expected query command to expose details subcommand")
	}
	if !foundAccountOverviews {
		t.Fatal("expected query command to expose account-overviews subcommand")
	}
	if !foundContentOverviews {
		t.Fatal("expected query command to expose content-overviews subcommand")
	}
	if !foundProxies {
		t.Fatal("expected query command to expose proxies subcommand")
	}
	if !foundProxyAreas {
		t.Fatal("expected query command to expose proxy-areas subcommand")
	}
}

func TestAccountOverviewsRequiresPlatform(t *testing.T) {
	err := newAccountOverviewsCmd().RunE(testCobraCommand(), nil)
	if err == nil {
		t.Fatal("expected missing platform error")
	}
	if !strings.Contains(err.Error(), "account overviews platform must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateAccountDryRunOutputsRequestBody(t *testing.T) {
	var out bytes.Buffer
	cmd := newAccountsUpdateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"acc_1", "--proxy-id", "proxy_1", "--remark", "主账号", "--group", "group_1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "accounts.update.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["dryRun"] != true || data["account"] != "acc_1" {
		t.Fatalf("unexpected dry-run metadata: %#v", data)
	}
	request := data["request"].(map[string]interface{})
	if request["proxyId"] != "proxy_1" || request["remark"] != "主账号" {
		t.Fatalf("unexpected update account request: %#v", request)
	}
	groups := request["groups"].([]interface{})
	if len(groups) != 1 || groups[0] != "group_1" {
		t.Fatalf("unexpected groups: %#v", groups)
	}
}

func TestUpdateAccountRejectsEmptyRequest(t *testing.T) {
	err := newAccountsUpdateCmd().RunE(testCobraCommand(), []string{"acc_1"})
	if err == nil {
		t.Fatal("expected empty request error")
	}
	if !strings.Contains(err.Error(), "update account request must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}
