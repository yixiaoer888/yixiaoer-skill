package config

import (
	"os"
	"path/filepath"
	"runtime"
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

func TestResolveProjectDirAcceptsReferencesWorkflowsLayout(t *testing.T) {
	root := t.TempDir()
	mustMkdirAll(t, filepath.Join(root, "schemas"))
	mustMkdirAll(t, filepath.Join(root, "references", "workflows"))
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

func TestResolveProjectDirRejectsInvalidOverride(t *testing.T) {
	override := t.TempDir()
	t.Setenv("YIXIAOER_PROJECT_DIR", override)

	_, err := resolveProjectDir(t.TempDir(), "")
	if err == nil {
		t.Fatal("expected error for invalid override")
	}
}

func TestWriteFileConfigUsesPrivatePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX file modes consistently")
	}
	configPath := filepath.Join(t.TempDir(), "config.json")

	if err := writeFileConfig(configPath, fileConfig{APIKey: "test-key"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("config mode = %o, want 600", got)
	}
}

func TestLoadFileConfigTightensExistingPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX file modes consistently")
	}
	configPath := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configPath, []byte(`{"apiKey":"test-key"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadFileConfig(configPath); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("config mode = %o, want 600", got)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
}
