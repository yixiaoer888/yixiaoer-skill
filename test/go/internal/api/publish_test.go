package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
			"statusCode": 0,
			"data":       map[string]interface{}{"taskSetId": "task_set_1"},
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
	if result["taskSetId"] != "task_set_1" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestPublishExtractsTaskSetIDFromStringData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"statusCode": 0,
			"data":       "task_set_string",
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Publish(map[string]interface{}{"publishType": "video"})
	if err != nil {
		t.Fatal(err)
	}
	if result["taskSetId"] != "task_set_string" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestPublishRejectsMissingTaskSetID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"statusCode": 0,
			"data":       nil,
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Publish(map[string]interface{}{"publishType": "video"}); err == nil {
		t.Fatal("expected error when response carries no taskSetId")
	}
}

func TestPublishRejectsNonZeroBusinessStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HTTP 200 but a non-zero business statusCode must be treated as failure.
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"statusCode": 4001,
			"message":    "账号代理不存在",
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	_, err := client.Publish(map[string]interface{}{"publishType": "video"})
	if err == nil {
		t.Fatal("expected business error for non-zero statusCode")
	}
	if !strings.Contains(err.Error(), "账号代理不存在") {
		t.Fatalf("expected business message, got %v", err)
	}
}
