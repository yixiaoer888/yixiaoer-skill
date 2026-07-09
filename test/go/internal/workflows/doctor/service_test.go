package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestCheckIncludesPublishChannelReadiness(t *testing.T) {
	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectDir, "schemas"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("YIXIAOER_PROJECT_DIR", projectDir)
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)
	if _, err := config.SaveAPIKey("test-api-key"); err != nil {
		t.Fatal(err)
	}
	if _, err := config.SaveLocalClientID("client_123"); err != nil {
		t.Fatal(err)
	}

	rt, err := app.Load()
	if err != nil {
		t.Fatal(err)
	}

	checks, err := NewService(rt).Check()
	if err != nil {
		t.Fatal(err)
	}

	if checks["localPublishClientId"] != "client_123" {
		t.Fatalf("expected localPublishClientId in doctor checks, got %#v", checks["localPublishClientId"])
	}
	readiness, ok := checks["publishChannelReadiness"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected publishChannelReadiness object, got %#v", checks["publishChannelReadiness"])
	}
	local, ok := readiness["local"].(map[string]interface{})
	if !ok || local["configured"] != true {
		t.Fatalf("expected local readiness configured, got %#v", readiness["local"])
	}
}

func TestCheckSupportsReferencesWorkflowLayout(t *testing.T) {
	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectDir, "schemas"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "references", "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("YIXIAOER_PROJECT_DIR", projectDir)
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)
	if _, err := config.SaveAPIKey("test-api-key"); err != nil {
		t.Fatal(err)
	}

	rt, err := app.Load()
	if err != nil {
		t.Fatal(err)
	}

	checks, err := NewService(rt).Check()
	if err != nil {
		t.Fatal(err)
	}

	if checks["workflowsOK"] != true {
		t.Fatalf("expected workflowsOK true, got %#v", checks["workflowsOK"])
	}
	if checks["workflowDocsPath"] != filepath.Join(projectDir, "references", "workflows") {
		t.Fatalf("unexpected workflowDocsPath: %#v", checks["workflowDocsPath"])
	}
}
