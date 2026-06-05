package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunRecordsListRequiresLimit(t *testing.T) {
	recordsPlatform = ""
	recordsLimit = ""
	recordsStatus = ""
	t.Cleanup(func() {
		recordsPlatform = ""
		recordsLimit = ""
		recordsStatus = ""
	})

	err := runRecordsList(testCobraCommand())
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
	cmd := &cobra.Command{Use: "locations"}
	cmd.Flags().StringVar(&locationsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&locationsKeyword, "keyword", "", "search keyword (alias for --query)")
	t.Cleanup(func() {
		locationsQuery = ""
		locationsKeyword = ""
	})

	if err := cmd.Flags().Parse([]string{"--keyword", "parks"}); err != nil {
		t.Fatal(err)
	}
	if locationsQuery != "" {
		t.Fatalf("expected primary query storage to remain empty, got %q", locationsQuery)
	}
	if locationsKeyword != "parks" {
		t.Fatalf("expected alias storage to capture keyword flag, got %q", locationsKeyword)
	}
}

func TestGoodsKeywordFlagUsesAliasStorage(t *testing.T) {
	cmd := &cobra.Command{Use: "goods"}
	cmd.Flags().StringVar(&goodsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&goodsKeyword, "keyword", "", "search keyword (alias for --query)")
	t.Cleanup(func() {
		goodsQuery = ""
		goodsKeyword = ""
	})

	if err := cmd.Flags().Parse([]string{"--keyword", "phone"}); err != nil {
		t.Fatal(err)
	}
	if goodsQuery != "" {
		t.Fatalf("expected primary query storage to remain empty, got %q", goodsQuery)
	}
	if goodsKeyword != "phone" {
		t.Fatalf("expected alias storage to capture keyword flag, got %q", goodsKeyword)
	}
}

func TestMiniAppsKeywordFlagUsesAliasStorage(t *testing.T) {
	cmd := &cobra.Command{Use: "miniapps"}
	cmd.Flags().StringVar(&miniAppsQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&miniAppsKeyword, "keyword", "", "search keyword (alias for --query)")
	t.Cleanup(func() {
		miniAppsQuery = ""
		miniAppsKeyword = ""
	})

	if err := cmd.Flags().Parse([]string{"--keyword", "抽奖"}); err != nil {
		t.Fatal(err)
	}
	if miniAppsQuery != "" {
		t.Fatalf("expected primary query storage to remain empty, got %q", miniAppsQuery)
	}
	if miniAppsKeyword != "抽奖" {
		t.Fatalf("expected alias storage to capture keyword flag, got %q", miniAppsKeyword)
	}
}

func TestGamesKeywordFlagUsesAliasStorage(t *testing.T) {
	cmd := &cobra.Command{Use: "games"}
	cmd.Flags().StringVar(&gamesQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&gamesKeyword, "keyword", "", "search keyword (alias for --query)")
	t.Cleanup(func() {
		gamesQuery = ""
		gamesKeyword = ""
	})

	if err := cmd.Flags().Parse([]string{"--keyword", "消消乐"}); err != nil {
		t.Fatal(err)
	}
	if gamesQuery != "" {
		t.Fatalf("expected primary query storage to remain empty, got %q", gamesQuery)
	}
	if gamesKeyword != "消消乐" {
		t.Fatalf("expected alias storage to capture keyword flag, got %q", gamesKeyword)
	}
}

func TestActivitiesKeywordFlagUsesAliasStorage(t *testing.T) {
	cmd := &cobra.Command{Use: "activities"}
	cmd.Flags().StringVar(&activitiesQuery, "query", "", "search keyword")
	cmd.Flags().StringVar(&activitiesKeyword, "keyword", "", "search keyword (alias for --query)")
	t.Cleanup(func() {
		activitiesQuery = ""
		activitiesKeyword = ""
	})

	if err := cmd.Flags().Parse([]string{"--keyword", "创作"}); err != nil {
		t.Fatal(err)
	}
	if activitiesQuery != "" {
		t.Fatalf("expected primary query storage to remain empty, got %q", activitiesQuery)
	}
	if activitiesKeyword != "创作" {
		t.Fatalf("expected alias storage to capture keyword flag, got %q", activitiesKeyword)
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

func TestLocationsCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newLocationsCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestQueryCommandExistsWithMiniAppsAndSyncAppsSubcommands(t *testing.T) {
	foundMiniApps := false
	foundSyncApps := false
	foundGames := false
	foundHotEvents := false
	foundGroups := false
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

func TestMiniAppsCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newMiniAppsCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestSyncAppsCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newSyncAppsCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestGamesCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newGamesCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestHotEventsCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newHotEventsCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestGroupsCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newGroupsCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestActivitiesCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newActivitiesCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestMusicCategoriesCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newMusicCategoriesCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestDetailsCommandMarkedAsCompatibilityEntry(t *testing.T) {
	cmd := newDetailsCmd()
	if !strings.Contains(cmd.Short, "兼容入口") {
		t.Fatalf("expected compatibility hint in short help, got %q", cmd.Short)
	}
}

func TestAccountOverviewsRequiresPlatform(t *testing.T) {
	accountOverviewPlatform = ""
	t.Cleanup(func() {
		accountOverviewPlatform = ""
	})

	err := newAccountOverviewsCmd().RunE(testCobraCommand(), nil)
	if err == nil {
		t.Fatal("expected missing platform error")
	}
	if !strings.Contains(err.Error(), "account overviews platform must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateAccountDryRunOutputsRequestBody(t *testing.T) {
	updateAccountKuaidailiArea = ""
	t.Cleanup(func() {
		updateAccountProxyID = ""
		updateAccountKuaidailiArea = ""
		updateAccountRemark = ""
		updateAccountGroups = nil
		updateAccountDryRun = false
	})

	var out bytes.Buffer
	cmd := newUpdateAccountCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"acc_1", "--proxy-id", "proxy_1", "--remark", "主账号", "--group", "group_1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
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
	updateAccountProxyID = ""
	updateAccountKuaidailiArea = ""
	updateAccountRemark = ""
	updateAccountGroups = nil
	updateAccountDryRun = true
	t.Cleanup(func() {
		updateAccountDryRun = false
	})

	err := newUpdateAccountCmd().RunE(testCobraCommand(), []string{"acc_1"})
	if err == nil {
		t.Fatal("expected empty request error")
	}
	if !strings.Contains(err.Error(), "update account request must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}
