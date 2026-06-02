package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSkillMarkdownLinksExist(t *testing.T) {
	withRepoRoot(t)
	for _, skillPath := range []string{
		filepath.Join("skills", "yixiaoer", "SKILL.md"),
		filepath.Join("skills", "yixiaoer", "references", "yixiaoer-shared.md"),
	} {
		raw, err := os.ReadFile(skillPath)
		if err != nil {
			t.Fatal(err)
		}
		if len(raw) == 0 {
			t.Fatalf("expected %s to be non-empty", skillPath)
		}
	}
}

func TestSkillCheckCommandSuccess(t *testing.T) {
	skillDir := filepath.Join(t.TempDir(), "skills", "yixiaoer")
	refsDir := filepath.Join(skillDir, "references")
	if err := os.MkdirAll(refsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("[ref](./references/doc.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refsDir, "doc.md"), []byte("[home](../SKILL.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("YIXIAOER_SKILL_DIR", skillDir)

	cmd := testCobraCommand()
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	if err := runSkillCheck(cmd); err != nil {
		t.Fatal(err)
	}

	var response struct {
		OK   bool `json:"ok"`
		Data struct {
			FilesScanned int `json:"filesScanned"`
			InvalidLinks int `json:"invalidLinks"`
		} `json:"data"`
	}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !response.OK {
		t.Fatal("response.OK = false, want true")
	}
	if response.Data.FilesScanned == 0 {
		t.Fatal("FilesScanned = 0, want > 0")
	}
	if response.Data.InvalidLinks != 0 {
		t.Fatalf("InvalidLinks = %d, want 0", response.Data.InvalidLinks)
	}
}
