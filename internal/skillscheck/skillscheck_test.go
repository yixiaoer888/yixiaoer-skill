package skillscheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckMissingStamp(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(dir, "config.json"))

	status, err := Check("3.0.0")
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if status.State != "missing" {
		t.Fatalf("status.State = %q, want missing", status.State)
	}
	if status.Sync {
		t.Fatal("status.Sync = true, want false")
	}
}

func TestWriteAndCheckStamp(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(dir, "config.json"))

	if err := WriteStamp("3.0.0"); err != nil {
		t.Fatalf("WriteStamp returned error: %v", err)
	}

	status, err := Check("3.0.0")
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if status.State != "in_sync" {
		t.Fatalf("status.State = %q, want in_sync", status.State)
	}
	if !status.Sync {
		t.Fatal("status.Sync = false, want true")
	}
}

func TestNoticeStale(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(dir, "config.json"))

	path, err := StampPath()
	if err != nil {
		t.Fatalf("StampPath returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte("2.9.0\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	notice, err := Notice("3.0.0", `C:\repo\skills\yixiaoer`)
	if err != nil {
		t.Fatalf("Notice returned error: %v", err)
	}
	if notice == nil {
		t.Fatal("notice = nil, want non-nil")
	}
	if got := notice["state"]; got != "stale" {
		t.Fatalf("notice[state] = %v, want stale", got)
	}
}

func TestCheckSkillLinksSuccess(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "yixiaoer")
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

	report, err := CheckSkillLinks(skillDir)
	if err != nil {
		t.Fatalf("CheckSkillLinks returned error: %v", err)
	}
	if report.FilesScanned != 2 {
		t.Fatalf("FilesScanned = %d, want 2", report.FilesScanned)
	}
	if report.LinksChecked != 2 {
		t.Fatalf("LinksChecked = %d, want 2", report.LinksChecked)
	}
	if report.InvalidLinks != 0 {
		t.Fatalf("InvalidLinks = %d, want 0", report.InvalidLinks)
	}
}

func TestCheckSkillLinksMissingTarget(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "yixiaoer")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("[missing](./references/missing.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := CheckSkillLinks(skillDir)
	if err == nil {
		t.Fatal("CheckSkillLinks error = nil, want non-nil")
	}
	if report.InvalidLinks != 1 {
		t.Fatalf("InvalidLinks = %d, want 1", report.InvalidLinks)
	}
	if len(report.Issues) != 1 {
		t.Fatalf("len(Issues) = %d, want 1", len(report.Issues))
	}
	if report.Issues[0].File != "SKILL.md" {
		t.Fatalf("issue file = %q, want SKILL.md", report.Issues[0].File)
	}
}
