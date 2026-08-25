package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPublishDeleteDryRunOutputsRequest(t *testing.T) {
	var out bytes.Buffer
	cmd := newPublishDeleteCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"task_1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "publish.delete.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["dryRun"] != true || data["taskId"] != "task_1" {
		t.Fatalf("unexpected dry-run metadata: %#v", data)
	}
	request := data["request"].(map[string]interface{})
	if request["method"] != "DELETE" || request["path"] != "/tasks/task_1/publish" {
		t.Fatalf("unexpected delete request: %#v", request)
	}
}

func TestPublishDeleteCallsExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/tasks/task_1/publish" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"statusCode": 0})
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	var out bytes.Buffer
	cmd := newPublishDeleteCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"task_1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["ok"] != true || response["action"] != "publish.delete" {
		t.Fatalf("unexpected delete envelope: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["taskId"] != "task_1" {
		t.Fatalf("unexpected delete data: %#v", data)
	}
	remote := data["response"].(map[string]interface{})
	if remote["statusCode"] != float64(0) {
		t.Fatalf("unexpected delete response: %#v", remote)
	}
}

func TestPublishDeleteRejectsEmptyTaskID(t *testing.T) {
	err := newPublishDeleteCmd().RunE(testCobraCommand(), []string{" "})
	if err == nil {
		t.Fatal("expected empty task id error")
	}
	if !strings.Contains(err.Error(), "published task id must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}
