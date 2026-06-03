package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestQueryEndpoints(t *testing.T) {
	tests := []struct {
		name      string
		call      func(*Client) (interface{}, error)
		path      string
		queryWant map[string]string
	}{
		{
			name: "categories",
			call: func(c *Client) (interface{}, error) {
				return c.Categories("acc_1", "article")
			},
			path:      "/platform-accounts/acc_1/categories",
			queryWant: map[string]string{"publishType": "article"},
		},
		{
			name: "locations",
			call: func(c *Client) (interface{}, error) {
				return c.Locations("acc_1", "上海", "2")
			},
			path:      "/platform-accounts/acc_1/location",
			queryWant: map[string]string{"locationType": "2", "keyWord": "上海"},
		},
		{
			name: "music",
			call: func(c *Client) (interface{}, error) {
				return c.Music("acc_1", "周杰伦")
			},
			path:      "/platform-accounts/acc_1/music",
			queryWant: map[string]string{"keyWord": "周杰伦"},
		},
		{
			name: "goods",
			call: func(c *Client) (interface{}, error) {
				return c.Goods("acc_1", "口红")
			},
			path:      "/platform-accounts/acc_1/goods",
			queryWant: map[string]string{"keyWord": "口红"},
		},
		{
			name: "collections",
			call: func(c *Client) (interface{}, error) {
				return c.Collections("acc_1", "video")
			},
			path:      "/platform-accounts/acc_1/collections",
			queryWant: map[string]string{"publishType": "video"},
		},
		{
			name: "challenges",
			call: func(c *Client) (interface{}, error) {
				return c.Challenges("acc_1", "旅行", "video")
			},
			path:      "/platform-accounts/acc_1/challenges",
			queryWant: map[string]string{"publishType": "video", "keyWord": "旅行"},
		},
		{
			name: "records",
			call: func(c *Client) (interface{}, error) {
				return c.Records("xhs", "20", "failed")
			},
			path:      "/v2/taskSets",
			queryWant: map[string]string{"platform": "小红书", "size": "20", "status": "failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				for key, want := range tt.queryWant {
					if got := r.URL.Query().Get(key); got != want {
						t.Fatalf("unexpected %s query: %s", key, got)
					}
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{"id": "item_1"},
					},
				})
			}))
			defer server.Close()

			client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
			data, err := tt.call(client)
			if err != nil {
				t.Fatal(err)
			}
			items, ok := data.([]interface{})
			if !ok || len(items) != 1 {
				t.Fatalf("unexpected data: %#v", data)
			}
		})
	}
}

func TestPrepareReturnsOnlineAccountsAndOptionalCategories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			if got := r.URL.Query().Get("platform"); got != "小红书" {
				t.Fatalf("unexpected platform query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_online", "status": 1},
					{"platformAccountId": "acc_offline", "status": 0},
				},
			})
		case "/platform-accounts/acc_online/categories":
			if got := r.URL.Query().Get("publishType"); got != "video" {
				t.Fatalf("unexpected publishType query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"id": "cat_1"}},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	data, err := client.Prepare("xhs", "video")
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Accounts) != 1 {
		t.Fatalf("expected one online account, got %d", len(data.Accounts))
	}
	if data.Schema != "schemas/platforms/xhs.video.schema.json" {
		t.Fatalf("unexpected schema path: %s", data.Schema)
	}
	if data.RootSchema != "schemas/publish.schema.json" {
		t.Fatalf("unexpected root schema path: %s", data.RootSchema)
	}
	if data.Categories == nil {
		t.Fatal("expected categories")
	}
}

func TestQueryDataAcceptsTopLevelArrayResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "cat_1"},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	data, err := client.Categories("acc_1", "video")
	if err != nil {
		t.Fatal(err)
	}
	items, ok := data.([]interface{})
	if !ok || len(items) != 1 {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestQueryDataPreservesCompleteDynamicObjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/music" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"yixiaoerId":   "music_1",
					"yixiaoerName": "完整音乐对象",
					"playUrl":      "https://example.invalid/preview.mp3",
					"url":          "https://example.invalid/music",
					"duration":     240,
					"artist":       "artist_1",
					"raw": map[string]interface{}{
						"id":      "music_1",
						"nested":  map[string]interface{}{"value": "kept"},
						"extUrl":  "https://example.invalid/raw",
						"numbers": []interface{}{1, 2, 3},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	data, err := client.Music("acc_1", "完整")
	if err != nil {
		t.Fatal(err)
	}
	items, ok := data.([]interface{})
	if !ok || len(items) != 1 {
		t.Fatalf("unexpected data: %#v", data)
	}
	item, ok := items[0].(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected item type: %#v", items[0])
	}
	if item["playUrl"] != "https://example.invalid/preview.mp3" || item["url"] != "https://example.invalid/music" {
		t.Fatalf("expected query result URLs to be preserved, got %#v", item)
	}
	raw, ok := item["raw"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected raw object, got %#v", item["raw"])
	}
	nested, ok := raw["nested"].(map[string]interface{})
	if !ok || nested["value"] != "kept" {
		t.Fatalf("expected nested raw fields to be preserved, got %#v", raw)
	}
}

func TestImageTextQueriesNormalizePublishType(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) (interface{}, error)
		path string
	}{
		{
			name: "categories",
			call: func(c *Client) (interface{}, error) {
				return c.Categories("acc_1", "imageText")
			},
			path: "/platform-accounts/acc_1/categories",
		},
		{
			name: "collections",
			call: func(c *Client) (interface{}, error) {
				return c.Collections("acc_1", "imageText")
			},
			path: "/platform-accounts/acc_1/collections",
		},
		{
			name: "challenges",
			call: func(c *Client) (interface{}, error) {
				return c.Challenges("acc_1", "旅行", "imageText")
			},
			path: "/platform-accounts/acc_1/challenges",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				if got := r.URL.Query().Get("publishType"); got != "imageText" {
					t.Fatalf("expected normalized publishType imageText, got %s", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []map[string]interface{}{
						{"id": "item_1"},
					},
				})
			}))
			defer server.Close()

			client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
			if _, err := tt.call(client); err != nil {
				t.Fatal(err)
			}
		})
	}
}
