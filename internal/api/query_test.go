package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestPlatformDocFileNameUsesShipinhaoAliasForImageText(t *testing.T) {
	if got := platformDocFileName("shipinhao", "imageText"); got != "shipinhao.md" {
		t.Fatalf("expected shipinhao imageText doc file, got %q", got)
	}
	if got := platformDocFileName("douyin", "imageText"); got != "douyin.md" {
		t.Fatalf("expected default platform doc file, got %q", got)
	}
}

func TestMiniAppsUsesExpectedEndpointAndKeywordQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/mini-apps" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("keyWord"); got != "抽奖" {
			t.Fatalf("unexpected keyword query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "miniapp_1", "name": "抽奖助手"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.MiniApps("acc_1", "抽奖")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one miniapp, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["name"] != "抽奖助手" {
		t.Fatalf("unexpected miniapp payload: %#v", first)
	}
}

func TestSyncAppsUsesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/sync-apps" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "sync_1", "name": "今日头条"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.SyncApps("acc_1")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one sync app, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["name"] != "今日头条" {
		t.Fatalf("unexpected sync app payload: %#v", first)
	}
}

func TestMusicCategoriesUsesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/music/category" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"yixiaoerId": "cat_1", "yixiaoerName": "流行"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.MusicCategories("acc_1")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one music category, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["yixiaoerName"] != "流行" {
		t.Fatalf("unexpected music category payload: %#v", first)
	}
}

func TestGamesUsesExpectedEndpointAndKeywordQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/games" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("keyWord"); got != "消消乐" {
			t.Fatalf("unexpected keyword query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "game_1", "name": "开心消消乐"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Games("acc_1", "消消乐")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one game, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["name"] != "开心消消乐" {
		t.Fatalf("unexpected game payload: %#v", first)
	}
}

func TestHotEventsUsesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/hot-events" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("publishType"); got != "video" {
			t.Fatalf("unexpected publishType query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "event_1", "name": "夏日热点"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.HotEvents("acc_1", "video")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one hot event, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["name"] != "夏日热点" {
		t.Fatalf("unexpected hot event payload: %#v", first)
	}
}

func TestDetailsUsesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/taskSets/task_set_1/tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "task_1", "status": "success"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Details("task_set_1")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one task, got %d", len(items))
	}
}

func TestAccountOverviewsUsesExpectedEndpointAndFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/overviews-v2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("platform"); got != "抖音" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		if got := r.URL.Query().Get("loginStatus"); got != "1" {
			t.Fatalf("unexpected loginStatus query: %s", got)
		}
		if got := r.URL.Query()["memberIds"]; len(got) != 2 || got[0] != "m1" || got[1] != "m2" {
			t.Fatalf("unexpected memberIds query: %#v", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"total": 1},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.AccountOverviews(AccountOverviewOptions{
		Platform:    "douyin",
		LoginStatus: "1",
		MemberIDs:   []string{"m1", "m2"},
		Page:        2,
		Size:        30,
	})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if data["total"] != float64(1) {
		t.Fatalf("unexpected overview data: %#v", data)
	}
}

func TestContentOverviewsUsesExpectedEndpointAndFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/contents/overviews" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("platform"); got != "小红书" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		if got := r.URL.Query().Get("platformAccountId"); got != "acc_1" {
			t.Fatalf("unexpected account query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"total": 1},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.ContentOverviews(ContentOverviewOptions{
		Platform:          "xhs",
		PlatformAccountID: "acc_1",
		Type:              "video",
		Page:              1,
		Size:              10,
	})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if data["total"] != float64(1) {
		t.Fatalf("unexpected content overview data: %#v", data)
	}
}

func TestProxiesAndProxyAreasUseExpectedEndpoints(t *testing.T) {
	paths := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		switch r.URL.Path {
		case "/proxys":
			if got := r.URL.Query().Get("size"); got != "50" {
				t.Fatalf("unexpected proxies size: %s", got)
			}
		case "/daili/areas":
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Proxies("50"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ProxyAreas(); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != "/proxys" || paths[1] != "/daili/areas" {
		t.Fatalf("unexpected paths: %#v", paths)
	}
}

func TestUpdateAccountUsesPatchEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/platform-accounts/acc_1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["proxyId"] != "proxy_1" {
			t.Fatalf("unexpected body: %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"ok": true}})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.UpdateAccount("acc_1", map[string]interface{}{"proxyId": "proxy_1"})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if data["ok"] != true {
		t.Fatalf("unexpected update result: %#v", data)
	}
}

func TestGroupsUsesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/group-chats" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"yixiaoerId": "group_1", "yixiaoerName": "品牌交流群"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Groups("acc_1")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one group, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["yixiaoerName"] != "品牌交流群" {
		t.Fatalf("unexpected group payload: %#v", first)
	}
}

func TestActivitiesUsesExpectedEndpointAndFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/activities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("publishType"); got != "video" {
			t.Fatalf("unexpected publishType query: %s", got)
		}
		if got := r.URL.Query().Get("categoryId"); got != "cat_1" {
			t.Fatalf("unexpected categoryId query: %s", got)
		}
		if got := r.URL.Query().Get("keyWord"); got != "创作" {
			t.Fatalf("unexpected keyword query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "activity_1", "name": "创作激励"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Activities("acc_1", "video", "cat_1", "创作")
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one activity, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["name"] != "创作激励" {
		t.Fatalf("unexpected activity payload: %#v", first)
	}
}
