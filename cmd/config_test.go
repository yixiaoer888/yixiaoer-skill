package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestConfigInitSavesAPIKeyOnly(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)

	var out bytes.Buffer
	cmd := newConfigInitCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--api-key", "test-api-key"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["apiKeyPresent"] != true {
		t.Fatalf("expected apiKeyPresent=true, got %#v", data["apiKeyPresent"])
	}
}

func TestConfigInitKeepsFlagsCommandLocal(t *testing.T) {
	first := newConfigInitCmd()
	if err := first.Flags().Parse([]string{"--api-key", "first-key"}); err != nil {
		t.Fatal(err)
	}

	second := newConfigInitCmd()
	if second.Flag("api-key").Value.String() != "" {
		t.Fatalf("expected fresh config init command to have empty api-key flag, got %q", second.Flag("api-key").Value.String())
	}
}
