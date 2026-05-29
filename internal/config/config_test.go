package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveProjectDirPrefersAncestorOfWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, "schemas"))
	mustMkdirAll(t, filepath.Join(root, "workflows"))
	nested := filepath.Join(root, "cmd", "subdir")
	mustMkdirAll(t, nested)

	projectDir, err := resolveProjectDir(nested, "")
	if err != nil {
		t.Fatalf("resolveProjectDir returned error: %v", err)
	}
	if projectDir != root {
		t.Fatalf("projectDir = %q, want %q", projectDir, root)
	}
}

func TestResolveProjectDirFallsBackToExecutableDirectory(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, "schemas"))
	mustMkdirAll(t, filepath.Join(root, "workflows"))
	exeDir := filepath.Join(root, "bin")
	mustMkdirAll(t, exeDir)

	workingDir := t.TempDir()
	projectDir, err := resolveProjectDir(workingDir, exeDir)
	if err != nil {
		t.Fatalf("resolveProjectDir returned error: %v", err)
	}
	if projectDir != root {
		t.Fatalf("projectDir = %q, want %q", projectDir, root)
	}
}

func TestResolveProjectDirUsesWorkingDirectoryWhenNoProjectFound(t *testing.T) {
	workingDir := t.TempDir()

	projectDir, err := resolveProjectDir(workingDir, "")
	if err != nil {
		t.Fatalf("resolveProjectDir returned error: %v", err)
	}
	if projectDir != workingDir {
		t.Fatalf("projectDir = %q, want %q", projectDir, workingDir)
	}
}

func TestResolveProjectDirRejectsInvalidOverride(t *testing.T) {
	override := t.TempDir()
	t.Setenv("YIXIAOER_PROJECT_DIR", override)

	_, err := resolveProjectDir(t.TempDir(), "")
	if err == nil {
		t.Fatal("expected error for invalid override")
	}
}

func TestSaveLinkedAppStatePersistsConnection(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)

	gotPath, err := SaveLinkedAppState("yixiaoer", "acc_1", "测试账号", true)
	if err != nil {
		t.Fatalf("SaveLinkedAppState returned error: %v", err)
	}
	if gotPath != configPath {
		t.Fatalf("configPath = %q, want %q", gotPath, configPath)
	}

	fileCfg, err := loadFileConfig(configPath)
	if err != nil {
		t.Fatalf("loadFileConfig returned error: %v", err)
	}
	state := fileCfg.linkedAppState("yixiaoer")
	if !state.Connected {
		t.Fatal("expected connected state to persist")
	}
	if state.AccountID != "acc_1" {
		t.Fatalf("AccountID = %q, want %q", state.AccountID, "acc_1")
	}
	if state.AccountName != "测试账号" {
		t.Fatalf("AccountName = %q, want %q", state.AccountName, "测试账号")
	}
	if state.UpdatedAt == "" {
		t.Fatal("expected UpdatedAt to be populated")
	}
}

func TestSaveLinkedAppStateDisconnectClearsAccount(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)

	if _, err := SaveLinkedAppState("yixiaoer", "acc_1", "测试账号", true); err != nil {
		t.Fatalf("initial SaveLinkedAppState returned error: %v", err)
	}
	if _, err := SaveLinkedAppState("yixiaoer", "", "", false); err != nil {
		t.Fatalf("disconnect SaveLinkedAppState returned error: %v", err)
	}

	fileCfg, err := loadFileConfig(configPath)
	if err != nil {
		t.Fatalf("loadFileConfig returned error: %v", err)
	}
	state := fileCfg.linkedAppState("yixiaoer")
	if state.Connected {
		t.Fatal("expected disconnected state to persist")
	}
	if state.AccountID != "" || state.AccountName != "" {
		t.Fatalf("expected account info to be cleared, got %+v", state)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
}
