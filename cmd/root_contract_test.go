package cmd

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestExecuteWithIOUnknownCommandWritesStructuredErrorToStderr(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := ExecuteWithIO([]string{"definitely-not-a-command"}, &stdout, &stderr)
	if code != yxerrors.ExitValidation {
		t.Fatalf("expected validation exit code, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected stdout to stay empty on error, got %q", stdout.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatalf("stderr should contain JSON error envelope, got %q: %v", stderr.String(), err)
	}
	if response["ok"] != false {
		t.Fatalf("expected ok=false response, got %#v", response)
	}
	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.UsageErr {
		t.Fatalf("expected usage error, got %#v", errorObj)
	}
	if errorObj["message"] != "unknown command" {
		t.Fatalf("expected unknown command message, got %#v", errorObj)
	}
	details := errorObj["details"].(map[string]interface{})
	if _, ok := details["availableCommands"].([]interface{}); !ok {
		t.Fatalf("expected availableCommands in details, got %#v", details)
	}
}

func TestExecuteWithIOArgumentErrorIsValidationError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := ExecuteWithIO([]string{"publish"}, &stdout, &stderr)
	if code != yxerrors.ExitValidation {
		t.Fatalf("expected validation exit code, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected stdout to stay empty on error, got %q", stdout.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatalf("stderr should contain JSON error envelope, got %q: %v", stderr.String(), err)
	}
	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.UsageErr {
		t.Fatalf("expected usage error, got %#v", errorObj)
	}
	if errorObj["type"] != yxerrors.ValidationType {
		t.Fatalf("expected validation type, got %#v", errorObj)
	}
	if errorObj["nextCommand"] != "yxer publish --help" {
		t.Fatalf("expected publish help next command, got %#v", errorObj)
	}
}

func TestExecuteWithIOUnknownFlagWritesStructuredErrorToStderr(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := ExecuteWithIO([]string{"accounts", "list", "--definitely-unknown"}, &stdout, &stderr)
	if code != yxerrors.ExitValidation {
		t.Fatalf("expected validation exit code, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected stdout to stay empty on error, got %q", stdout.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatalf("stderr should contain JSON error envelope, got %q: %v", stderr.String(), err)
	}
	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.UsageErr {
		t.Fatalf("expected usage error, got %#v", errorObj)
	}
	if errorObj["nextCommand"] != "yxer accounts list --help" {
		t.Fatalf("expected accounts list help next command, got %#v", errorObj)
	}
}

func TestExecuteWithIOHelpWritesStructuredJSONToStdout(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := ExecuteWithIO([]string{"config", "--help"}, &stdout, &stderr)
	if code != yxerrors.ExitOK {
		t.Fatalf("expected ok exit code, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected stderr to stay empty, got %q", stderr.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout should contain JSON help envelope, got %q: %v", stdout.String(), err)
	}
	if response["action"] != "help" {
		t.Fatalf("expected help action, got %#v", response["action"])
	}
	data := response["data"].(map[string]interface{})
	if data["commandPath"] != "yxer config" {
		t.Fatalf("unexpected help command path: %#v", data["commandPath"])
	}
}

func TestExecuteWithIOVersionWritesStructuredJSONToStdout(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := ExecuteWithIO([]string{"--version"}, &stdout, &stderr)
	if code != yxerrors.ExitOK {
		t.Fatalf("expected ok exit code, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected stderr to stay empty, got %q", stderr.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout should contain JSON version envelope, got %q: %v", stdout.String(), err)
	}
	if response["action"] != "version" || response["ok"] != true {
		t.Fatalf("unexpected version response: %#v", response)
	}
}

func TestRootCommandTreeExposesStableTopLevelGroups(t *testing.T) {
	expected := []string{
		"account-group",
		"accounts",
		"config",
		"doctor",
		"draft",
		"material",
		"prepare",
		"publish",
		"query",
		"schema",
		"skill",
		"update",
		"upload",
		"validate",
	}
	for _, name := range expected {
		if !rootHasCommand(name) {
			t.Fatalf("expected root command to expose %q", name)
		}
	}
	if rootHasCommand("update-account") {
		t.Fatal("expected legacy update-account command to be removed")
	}
}

func TestVisibleCommandTreeSnapshot(t *testing.T) {
	expected := map[string][]string{
		"yxer":          {"account-group", "accounts", "config", "doctor", "draft", "material", "prepare", "publish", "query", "schema", "skill", "update", "upload", "validate"},
		"account-group": {"create", "delete", "list", "update"},
		"accounts":      {"list", "update"},
		"config":        {"get", "init", "set-api-key", "set-local-client-id"},
		"draft":         {"save"},
		"material":      {"add", "create"},
		"publish":       {"init"},
		"query": {
			"account-overviews",
			"activities",
			"categories",
			"challenges",
			"collections",
			"content-overviews",
			"details",
			"games",
			"goods",
			"groups",
			"hot-events",
			"locations",
			"members",
			"miniapps",
			"music",
			"music-categories",
			"proxies",
			"proxy-areas",
			"records",
			"syncapps",
		},
		"records": {"list"},
		"schema":  {"catalog", "fields", "get", "list"},
		"skill":   {"check", "show", "sync"},
	}

	actual := map[string][]string{
		"yxer":          visibleSubcommandNames(rootCmd),
		"account-group": visibleSubcommandNames(mustRootChild(t, "account-group")),
		"accounts":      visibleSubcommandNames(mustRootChild(t, "accounts")),
		"config":        visibleSubcommandNames(mustRootChild(t, "config")),
		"draft":         visibleSubcommandNames(mustRootChild(t, "draft")),
		"material":      visibleSubcommandNames(mustRootChild(t, "material")),
		"publish":       visibleSubcommandNames(mustRootChild(t, "publish")),
		"query":         visibleSubcommandNames(mustRootChild(t, "query")),
		"records":       visibleSubcommandNames(mustChild(t, mustRootChild(t, "query"), "records")),
		"schema":        visibleSubcommandNames(mustRootChild(t, "schema")),
		"skill":         visibleSubcommandNames(mustRootChild(t, "skill")),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("visible command tree changed\nexpected: %#v\nactual:   %#v", expected, actual)
	}
}

func rootHasCommand(name string) bool {
	for _, child := range rootCmd.Commands() {
		if child.Name() == name {
			return true
		}
	}
	return false
}

func mustRootChild(t *testing.T, name string) *cobra.Command {
	t.Helper()
	return mustChild(t, rootCmd, name)
}

func mustChild(t *testing.T, parent *cobra.Command, name string) *cobra.Command {
	t.Helper()
	for _, child := range parent.Commands() {
		if child.Name() == name {
			return child
		}
	}
	t.Fatalf("expected %q to contain child %q", parent.CommandPath(), name)
	return nil
}

func visibleSubcommandNames(cmd *cobra.Command) []string {
	names := make([]string, 0, len(cmd.Commands()))
	for _, child := range cmd.Commands() {
		if child.Hidden || !child.IsAvailableCommand() || child.Name() == "help" || child.Name() == "completion" {
			continue
		}
		names = append(names, child.Name())
	}
	sort.Strings(names)
	return names
}
