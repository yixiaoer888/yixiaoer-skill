package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDramaTasksCommandUsesJSONOutputContract(t *testing.T) {
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_drama/drama-tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("keyWord") != "护妻" {
			t.Fatalf("unexpected keyWord: %q", r.URL.Query().Get("keyWord"))
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{{
				"yixiaoerId":       "event/1",
				"yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover",
				"yixiaoerName":     "风浪过后护妻安康",
			}},
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)
	withRepoRoot(t)

	var stdout bytes.Buffer
	cmd := newDramaTasksCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"acc_drama", "--query", "护妻"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout must contain JSON only: %v; output=%s", err, stdout.String())
	}
	if response["ok"] != true || response["action"] != "drama-tasks" {
		t.Fatalf("unexpected output envelope: %#v", response)
	}
	items := response["data"].([]interface{})
	if len(items) != 1 || items[0].(map[string]interface{})["yixiaoerId"] != "event/1" {
		t.Fatalf("unexpected drama output data: %#v", response["data"])
	}
}

func TestDramaTasksCommandWritesStructuredRemoteErrorToStderr(t *testing.T) {
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "当前账号不支持剧集",
			"code":    "DRAMA_UNSUPPORTED",
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)
	withRepoRoot(t)

	var stdout, stderr bytes.Buffer
	code := ExecuteWithIO([]string{"query", "drama-tasks", "acc_drama", "--json"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected non-zero exit code for remote drama error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected stdout to stay empty on error, got %q", stdout.String())
	}
	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatalf("stderr must contain structured JSON error: %v; output=%s", err, stderr.String())
	}
	errorObject := response["error"].(map[string]interface{})
	if errorObject["category"] != "remote_error" || errorObject["nextCommand"] != "yxer query drama-tasks acc_drama --json" {
		t.Fatalf("unexpected drama remote error envelope: %#v", errorObject)
	}
	details := errorObject["details"].(map[string]interface{})
	if details["code"] != "DRAMA_UNSUPPORTED" {
		t.Fatalf("expected remote code in details, got %#v", details)
	}
}
