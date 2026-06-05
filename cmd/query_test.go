package cmd

import (
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
