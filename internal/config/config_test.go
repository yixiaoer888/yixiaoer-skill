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

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
}
