package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigInitSavesAPIKeyOnly(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)

	configInitAPIKey = "test-api-key"
	configInitLocalClientID = ""
	configInitBindApp = false
	configInitAccountID = ""
	configInitAccountName = ""

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := configInitCmd.RunE(cmd, nil); err != nil {
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
	if data["boundApp"] != false {
		t.Fatalf("expected boundApp=false, got %#v", data["boundApp"])
	}
}

func TestConfigInitCanBindLinkedApp(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)

	configInitAPIKey = "test-api-key"
	configInitLocalClientID = "local_client_1"
	configInitBindApp = true
	configInitAccountID = "acc_1"
	configInitAccountName = "主账号"

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := configInitCmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["boundApp"] != true {
		t.Fatalf("expected boundApp=true, got %#v", data["boundApp"])
	}
	if data["localPublishClientId"] != "local_client_1" {
		t.Fatalf("expected local client id to persist, got %#v", data["localPublishClientId"])
	}
	linkedApp := data["linkedApp"].(map[string]interface{})
	if linkedApp["connected"] != true {
		t.Fatalf("expected linkedApp.connected=true, got %#v", linkedApp["connected"])
	}
	if linkedApp["accountId"] != "acc_1" {
		t.Fatalf("expected linkedApp.accountId=acc_1, got %#v", linkedApp["accountId"])
	}
}
