package doctor

import (
	"path/filepath"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestCheckIncludesPublishChannelReadiness(t *testing.T) {
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
