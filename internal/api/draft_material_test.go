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

func TestMoveMaterial(t *testing.T) {
	var gotMethod string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/material/batch/set-group" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"movedTotal": 1,
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.MoveMaterial("material_1", map[string]interface{}{"groupId": "group_1"})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("unexpected method: %s", gotMethod)
	}
	if len(gotBody) != 2 || gotBody["groupId"] != "group_1" {
		t.Fatalf("unexpected request body: %+v", gotBody)
	}
	materialIDs, ok := gotBody["materialIds"].([]interface{})
	if !ok || len(materialIDs) != 1 || materialIDs[0] != "material_1" {
		t.Fatalf("unexpected material IDs: %+v", gotBody["materialIds"])
	}
	data := DataOrSelf(result).(map[string]interface{})
	if data["movedTotal"] != float64(1) {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestMaterialGroups(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/material/groups" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("size") != "25" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		gotMethod = r.Method
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{map[string]interface{}{"id": "group_1"}}})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.MaterialGroups(MaterialGroupOptions{Page: 2, Size: 25})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("unexpected method: %s", gotMethod)
	}
	items, ok := result.([]interface{})
	if !ok || len(items) != 1 {
		t.Fatalf("unexpected response: %#v", result)
	}
}

func TestMaterials(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/material" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("size") != "25" || r.URL.Query().Get("fileName") != "demo.png" || r.URL.Query().Get("type") != "image" || r.URL.Query().Get("groupId") != "group_1" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		gotMethod = r.Method
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"records": []interface{}{map[string]interface{}{"id": "material_1", "fileName": "demo.png"}},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Materials(MaterialListOptions{
		Page:     2,
		Size:     25,
		FileName: "demo.png",
		Type:     "image",
		GroupID:  "group_1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("unexpected method: %s", gotMethod)
	}
	page := result.(map[string]interface{})
	records := page["records"].([]interface{})
	if len(records) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
