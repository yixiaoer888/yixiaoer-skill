package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestReadPayloadMissingFileReturnsStructuredFileError(t *testing.T) {
	_, err := readPayload(filepath.Join(t.TempDir(), "yxer-missing-payload.json"))
	if err == nil {
		t.Fatal("expected missing payload error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yxerrors.Error, got %T: %v", err, err)
	}
	if typed.Category != "file_not_found" {
		t.Fatalf("expected file_not_found category, got %+v", typed)
	}
	if typed.Hint == "" {
		t.Fatalf("expected repair hint, got %+v", typed)
	}
}

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
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"抖音", "video", payloadPath})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	nextStep := data["nextStep"].(string)
	if !strings.Contains(nextStep, "--publish-channel local") || !strings.Contains(nextStep, "--client-id configured_client_1") {
		t.Fatalf("expected local publish nextStep, got %q", nextStep)
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
