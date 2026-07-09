package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestValidateCommandUsesConfiguredLocalClientID(t *testing.T) {
	withRepoRoot(t)
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)
	if _, err := config.SaveLocalClientID("configured_client_1"); err != nil {
		t.Fatal(err)
	}
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
		"publishChannel": "local",
		"publishArgs":    validPublishArgs(),
	})

	cmd := newValidateCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"抖音", "video", payloadPath})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateCommandUsesLocalFlags(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	cmd := newValidateCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"抖音", "video", payloadPath, "--publish-channel", "local", "--client-id", "flag_client_1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateCommandRejectsInnerBusinessFormPayload(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"formType":    "task",
		"title":       "只有内层表单",
		"description": "这类 payload 不该再通过 validate",
		"visibleType": float64(0),
	})

	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := newValidateCmd().RunE(cmd, []string{"小红书", "video", payloadPath})
	if err == nil {
		t.Fatal("expected standard payload error")
	}
	if !strings.Contains(err.Error(), "Standard publish payload is required") {
		t.Fatalf("expected standard payload error, got %v", err)
	}
}
