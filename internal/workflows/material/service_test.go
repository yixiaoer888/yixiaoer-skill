package material

import (
	"strings"
	"testing"
)

func TestBuildMoveBodyTrimsGroupID(t *testing.T) {
	body := BuildMoveBody(" material_1 ", MoveInput{GroupID: " group_1 "})
	if len(body) != 2 || body["groupId"] != "group_1" {
		t.Fatalf("unexpected move body: %#v", body)
	}
	materialIDs, ok := body["materialIds"].([]string)
	if !ok || len(materialIDs) != 1 || materialIDs[0] != "material_1" {
		t.Fatalf("unexpected material IDs: %#v", body["materialIds"])
	}
}

func TestValidateMoveInputRejectsEmptyGroupID(t *testing.T) {
	err := ValidateMoveInput(MoveInput{GroupID: "  "})
	if err == nil {
		t.Fatal("expected missing material group id error")
	}
	if !strings.Contains(err.Error(), "material group id must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFindExactFileNameMatchesSupportsPaginatedResults(t *testing.T) {
	result := map[string]interface{}{
		"records": []interface{}{
			map[string]interface{}{"id": "material_1", "fileName": "cover.png", "type": "image"},
			map[string]interface{}{"id": "material_2", "fileName": "other.png", "type": "image"},
			map[string]interface{}{"id": "material_3", "fileName": "cover.png", "filePath": "material/cover.png"},
		},
	}
	matches := FindExactFileNameMatches(result, "cover.png")
	if len(matches) != 2 {
		t.Fatalf("expected two exact matches, got %#v", matches)
	}
	if matches[0].ID != "material_1" || matches[1].ID != "material_3" {
		t.Fatalf("unexpected matches: %#v", matches)
	}
}

func TestFindExactFileNameMatchesRejectsPartialNames(t *testing.T) {
	result := []interface{}{
		map[string]interface{}{"id": "material_1", "fileName": "cover.png"},
		map[string]interface{}{"id": "material_2", "fileName": "cover-2.png"},
	}
	if matches := FindExactFileNameMatches(result, "cover"); len(matches) != 0 {
		t.Fatalf("expected no partial-name matches, got %#v", matches)
	}
}
