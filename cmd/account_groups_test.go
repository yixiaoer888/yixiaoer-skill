package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestAccountGroupCommandExistsWithResourceSubcommands(t *testing.T) {
	var groupCmdFound bool
	foundList := false
	foundCreate := false
	foundUpdate := false
	foundDelete := false

	for _, child := range rootCmd.Commands() {
		if child.Name() != "account-group" {
			continue
		}
		groupCmdFound = true
		for _, grandchild := range child.Commands() {
			switch grandchild.Name() {
			case "list":
				foundList = true
			case "create":
				foundCreate = true
			case "update":
				foundUpdate = true
			case "delete":
				foundDelete = true
			}
		}
	}

	if !groupCmdFound {
		t.Fatal("expected root command to expose account-group resource command")
	}
	if !foundList || !foundCreate || !foundUpdate || !foundDelete {
		t.Fatalf("expected account-group to expose list/create/update/delete, got list=%v create=%v update=%v delete=%v", foundList, foundCreate, foundUpdate, foundDelete)
	}
}

func TestAccountGroupListCommandUsesStructuredAction(t *testing.T) {
	cmd := newAccountGroupListCmd()
	if cmd.Name() != "list" {
		t.Fatalf("unexpected command name: %s", cmd.Name())
	}
	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "ls" {
		t.Fatalf("unexpected aliases: %#v", cmd.Aliases)
	}
	if cmd.Flags().Lookup("page") == nil {
		t.Fatal("expected account-group list to expose --page flag")
	}
	if cmd.Flags().Lookup("size") == nil {
		t.Fatal("expected account-group list to expose --size flag")
	}
}

func TestBuildAccountGroupBodyRejectsEmptyName(t *testing.T) {
	_, err := buildAccountGroupBody("   ", accountGroupUpdateOptions{})
	if err == nil {
		t.Fatal("expected empty account group name error")
	}
	if !strings.Contains(err.Error(), "account group name must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateAccountGroupDryRunOutputsRequestBody(t *testing.T) {
	var out bytes.Buffer
	cmd := newAccountGroupCreateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"核心账号组", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "account-group.create.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["dryRun"] != true {
		t.Fatalf("unexpected dry-run metadata: %#v", data)
	}
	request := data["request"].(map[string]interface{})
	if request["name"] != "核心账号组" {
		t.Fatalf("unexpected request body: %#v", request)
	}
}

func TestCreateAccountGroupDryRunSupportsVisibleScopeAndUsers(t *testing.T) {
	var out bytes.Buffer
	cmd := newAccountGroupCreateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"核心账号组", "--visible-scope", "specific", "--visible-user", "user_1", "--visible-user", "user_2", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	request := response["data"].(map[string]interface{})["request"].(map[string]interface{})
	if request["visibleScope"] != "specific" {
		t.Fatalf("unexpected visibleScope: %#v", request)
	}
	users := request["visibleUsers"].([]interface{})
	if len(users) != 2 || users[0] != "user_1" || users[1] != "user_2" {
		t.Fatalf("unexpected visibleUsers: %#v", users)
	}
}

func TestUpdateAccountGroupDryRunOutputsRequestBody(t *testing.T) {
	var out bytes.Buffer
	cmd := newAccountGroupUpdateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"grp_1", "新版核心账号组", "--visible-scope", "specific", "--visible-user", "user_1", "--visible-user", "user_2", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "account-group.update.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["dryRun"] != true || data["groupId"] != "grp_1" {
		t.Fatalf("unexpected dry-run metadata: %#v", data)
	}
	request := data["request"].(map[string]interface{})
	if request["name"] != "新版核心账号组" {
		t.Fatalf("unexpected request body: %#v", request)
	}
	if request["visibleScope"] != "specific" {
		t.Fatalf("unexpected visibleScope: %#v", request)
	}
	users := request["visibleUsers"].([]interface{})
	if len(users) != 2 || users[0] != "user_1" || users[1] != "user_2" {
		t.Fatalf("unexpected visibleUsers: %#v", users)
	}
}

func TestUpdateAccountGroupRejectsEmptyID(t *testing.T) {
	err := newAccountGroupUpdateCmd().RunE(testCobraCommand(), []string{" ", "新版核心账号组"})
	if err == nil {
		t.Fatal("expected empty account group id error")
	}
	if !strings.Contains(err.Error(), "account group id must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildAccountGroupBodyRejectsInvalidVisibleScope(t *testing.T) {
	_, err := buildAccountGroupBody("新版核心账号组", accountGroupUpdateOptions{
		VisibleScope: "private",
	})
	if err == nil {
		t.Fatal("expected invalid visible scope error")
	}
	if !strings.Contains(err.Error(), "account group visible scope must be all or specific") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildAccountGroupBodyRejectsMissingVisibleUsersForSpecificScope(t *testing.T) {
	_, err := buildAccountGroupBody("新版核心账号组", accountGroupUpdateOptions{
		VisibleScope: "specific",
	})
	if err == nil {
		t.Fatal("expected missing visibleUsers error")
	}
	if !strings.Contains(err.Error(), "account group visibleUsers must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildAccountGroupBodySupportsAllScopeWithoutUsers(t *testing.T) {
	body, err := buildAccountGroupBody("新版核心账号组", accountGroupUpdateOptions{
		VisibleScope: "all",
	})
	if err != nil {
		t.Fatal(err)
	}
	if body["visibleScope"] != "all" {
		t.Fatalf("unexpected body: %#v", body)
	}
	if _, ok := body["visibleUsers"]; ok {
		t.Fatalf("did not expect visibleUsers for all scope: %#v", body)
	}
}

func TestDeleteAccountGroupDryRunOutputsRequestBody(t *testing.T) {
	var out bytes.Buffer
	cmd := newAccountGroupDeleteCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"grp_1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "account-group.delete.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["dryRun"] != true || data["groupId"] != "grp_1" {
		t.Fatalf("unexpected dry-run metadata: %#v", data)
	}
}

func TestDeleteAccountGroupRejectsEmptyID(t *testing.T) {
	err := newAccountGroupDeleteCmd().RunE(testCobraCommand(), []string{" "})
	if err == nil {
		t.Fatal("expected empty account group id error")
	}
	if !strings.Contains(err.Error(), "account group id must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAccountGroupCreateCommandUsesResourceShape(t *testing.T) {
	cmd := newAccountGroupCreateCmd()
	if cmd.Use != "create <name>" {
		t.Fatalf("unexpected use: %q", cmd.Use)
	}
}

func TestAccountGroupUpdateCommandUsesResourceShape(t *testing.T) {
	cmd := newAccountGroupUpdateCmd()
	if cmd.Use != "update <group_id> <name>" {
		t.Fatalf("unexpected use: %q", cmd.Use)
	}
}

func TestAccountGroupDeleteCommandUsesResourceShape(t *testing.T) {
	cmd := newAccountGroupDeleteCmd()
	if cmd.Use != "delete <group_id>" {
		t.Fatalf("unexpected use: %q", cmd.Use)
	}
}
