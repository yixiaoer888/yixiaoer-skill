package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestMaterialCommandExposesListAndMoveByNameSubcommands(t *testing.T) {
	cmd := newMaterialCmd()
	found := map[string]bool{}
	for _, child := range cmd.Commands() {
		found[child.Name()] = true
	}
	for _, name := range []string{"list", "move-by-name"} {
		if !found[name] {
			t.Fatalf("expected material to expose %s subcommand", name)
		}
	}
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
	if request["groupId"] != "group_1" {
		t.Fatalf("unexpected move request: %#v", request)
	}
	materialIDs, ok := request["materialIds"].([]interface{})
	if !ok || len(materialIDs) != 1 || materialIDs[0] != "material_1" {
		t.Fatalf("unexpected move material IDs: %#v", request["materialIds"])
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

func TestMaterialMoveByNameDryRunResolvesExactMaterialID(t *testing.T) {
	withRepoRoot(t)
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/material" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("fileName") != "demo.png" || r.URL.Query().Get("type") != "image" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"records": []interface{}{map[string]interface{}{
					"id":       "material_1",
					"fileName": "demo.png",
					"type":     "image",
				}},
			},
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)

	var out bytes.Buffer
	cmd := newMaterialMoveByNameCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"demo.png", "--group-id", "group_1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "material.move-by-name.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["materialId"] != "material_1" || data["fileName"] != "demo.png" || data["groupId"] != "group_1" {
		t.Fatalf("unexpected response data: %#v", data)
	}
	request := data["request"].(map[string]interface{})
	materialIDs := request["materialIds"].([]interface{})
	if len(materialIDs) != 1 || materialIDs[0] != "material_1" {
		t.Fatalf("unexpected request: %#v", request)
	}
}

func TestMaterialMoveByNameRejectsDuplicateMatches(t *testing.T) {
	withRepoRoot(t)
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"id": "material_1", "fileName": "demo.png"},
					map[string]interface{}{"id": "material_2", "fileName": "demo.png"},
				},
			},
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)

	cmd := newMaterialMoveByNameCmd()
	cmd.SetArgs([]string{"demo.png", "--group-id", "group_1", "--dry-run"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "multiple materials matched") {
		t.Fatalf("expected duplicate match error, got %v", err)
	}
}
