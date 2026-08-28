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
