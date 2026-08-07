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

func TestPublishWrapsCategoryOnlyAtAPIBoundary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		args := body["publishArgs"].(map[string]interface{})
		forms := args["accountForms"].([]interface{})
		form := forms[0].(map[string]interface{})
		cpf := form["contentPublishForm"].(map[string]interface{})
		category := cpf["category"].([]interface{})[0].(map[string]interface{})
		if category["id"] != "1000013" || category["text"] != "科技互联网" {
			t.Fatalf("expected publish category wrapper, got %#v", category)
		}
		raw := category["raw"].(map[string]interface{})
		if raw["yixiaoerId"] != "1000013" || raw["yixiaoerName"] != "科技互联网" {
			t.Fatalf("expected canonical query object under raw, got %#v", raw)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"statusCode": 0, "data": "task_1"})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	body := map[string]interface{}{
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{map[string]interface{}{
				"contentPublishForm": map[string]interface{}{
					"category": []interface{}{map[string]interface{}{
						"yixiaoerId": "1000013", "yixiaoerName": "科技互联网",
						"raw": map[string]interface{}{"id": "1000013", "name": "科技互联网"},
					}},
				},
			}},
		},
	}
	if _, err := client.Publish(body); err != nil {
		t.Fatal(err)
	}
	category := body["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})["category"].([]interface{})[0].(map[string]interface{})
	if _, exists := category["id"]; exists {
		t.Fatal("expected caller payload to remain canonical; API conversion should use the wire boundary")
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
