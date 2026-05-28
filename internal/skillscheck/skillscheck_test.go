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
