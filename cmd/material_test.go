package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestMaterialCommandExposesMoveSubcommand(t *testing.T) {
	cmd := newMaterialCmd()
	for _, child := range cmd.Commands() {
		if child.Name() == "move" {
			return
		}
	}
	t.Fatal("expected material to expose move subcommand")
}

func TestMaterialCommandExposesGroupsSubcommand(t *testing.T) {
	cmd := newMaterialCmd()
	for _, child := range cmd.Commands() {
		if child.Name() == "groups" {
			if len(child.Aliases) != 1 || child.Aliases[0] != "group-list" {
				t.Fatalf("unexpected material groups aliases: %#v", child.Aliases)
			}
			return
		}
	}
	t.Fatal("expected material to expose groups subcommand")
}

func TestMaterialMoveDryRunOutputsStableRequest(t *testing.T) {
	var out bytes.Buffer
	cmd := newMaterialMoveCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"material_1", "--group-id", "group_1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "material.move.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["dryRun"] != true || data["materialId"] != "material_1" || data["groupId"] != "group_1" {
		t.Fatalf("unexpected dry-run data: %#v", data)
	}
	request := data["request"].(map[string]interface{})
	if len(request) != 1 || request["groupId"] != "group_1" {
		t.Fatalf("unexpected move request: %#v", request)
	}
}

func TestMaterialMoveRequiresGroupID(t *testing.T) {
	err := newMaterialMoveCmd().RunE(testCobraCommand(), []string{"material_1"})
	if err == nil {
		t.Fatal("expected missing material group id error")
	}
	if !strings.Contains(err.Error(), "material group id must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMaterialMoveRejectsEmptyMaterialID(t *testing.T) {
	err := newMaterialMoveCmd().RunE(testCobraCommand(), []string{" "})
	if err == nil {
		t.Fatal("expected empty material id error")
	}
	if !strings.Contains(err.Error(), "material id must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}
