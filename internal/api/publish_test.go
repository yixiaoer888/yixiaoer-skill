package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestPublishPostsTaskSetBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/taskSets/v2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "test-key" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["publishType"] != "video" {
			t.Fatalf("unexpected publish body: %+v", body)
		}
		if body["publishChannel"] != "cloud" {
			t.Fatalf("expected publishChannel in API body: %+v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"taskSetId": "task_set_1"},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Publish(map[string]interface{}{
		"publishType":    "video",
		"publishChannel": "cloud",
	})
	if err != nil {
		t.Fatal(err)
	}
	data := result["data"].(map[string]interface{})
	if data["taskSetId"] != "task_set_1" {
		t.Fatalf("unexpected result: %+v", result)
	}
}
