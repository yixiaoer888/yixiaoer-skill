package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestSaveDraft(t *testing.T) {
	var gotMethod string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/taskSets/drafts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":      "draft_1",
				"isDraft": true,
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.SaveDraft(map[string]interface{}{
		"publishType": "video",
		"isDraft":     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPut {
		t.Fatalf("unexpected method: %s", gotMethod)
	}
	if gotBody["publishType"] != "video" || gotBody["isDraft"] != true {
		t.Fatalf("unexpected request body: %+v", gotBody)
	}
	data := DataOrSelf(result).(map[string]interface{})
	if data["id"] != "draft_1" {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestMaterial(t *testing.T) {
	var gotMethod string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/material" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":       "material_1",
				"filePath": gotBody["filePath"],
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Material(map[string]interface{}{
		"filePath": "uploaded/demo.png",
		"fileName": "demo.png",
		"width":    1280,
		"height":   720,
		"type":     "image",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("unexpected method: %s", gotMethod)
	}
	if gotBody["filePath"] != "uploaded/demo.png" || gotBody["type"] != "image" {
		t.Fatalf("unexpected request body: %+v", gotBody)
	}
	data := DataOrSelf(result).(map[string]interface{})
	if data["id"] != "material_1" {
		t.Fatalf("unexpected response: %+v", result)
	}
}
