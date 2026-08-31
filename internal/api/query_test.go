package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestClientAddsSetAPIKeyHintForUnauthorizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "账号登录失效",
			"code":    "UNAUTHORIZED",
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "expired-key", APIURL: server.URL})
	var out map[string]interface{}
	err := client.Get("/accounts", &out)
	if err == nil {
		t.Fatal("expected unauthorized error")
	}

	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yx error, got %T", err)
	}
	if typed.Hint == "" {
		t.Fatalf("expected unauthorized error to include api key hint: %#v", typed)
	}
	if typed.NextCommand != "yxer config set-api-key <apiKey>" {
		t.Fatalf("unexpected next command: %#v", typed.NextCommand)
	}
}

func TestPlatformDocFileNameUsesShipinhaoAliasForImageText(t *testing.T) {
	if got := platformDocFileName("shipinhao", "imageText"); got != "shipinhao.md" {
		t.Fatalf("expected shipinhao imageText doc file, got %q", got)
	}
	if got := platformDocFileName("douyin", "imageText"); got != "douyin.md" {
		t.Fatalf("expected default platform doc file, got %q", got)
	}
}

func TestPlatformDocFileNameUsesWeixinAccountImageTextDoc(t *testing.T) {
	for _, alias := range []string{"WeiXinGongZhongHao", "微信公众号", "weixin.account"} {
		if got := platformDocFileName(alias, "imageText"); got != "weixingongzhonghao.md" {
			t.Fatalf("platformDocFileName(%q, imageText) = %q", alias, got)
		}
	}
}

func TestClientWrapsInvalidJSONResponseAsRemoteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html>not json</html>"))
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	var out map[string]interface{}
	err := client.Get("/broken", &out)
	if err == nil {
		t.Fatal("expected invalid JSON response error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yx error, got %T", err)
	}
	if typed.Type != yxerrors.RemoteType || !typed.Retryable {
		t.Fatalf("unexpected structured error: %#v", typed)
	}
}

func TestGoodsDetailUsesExpectedEndpointAndPreservesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/platform-accounts/acc_1/goods-detail" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("url"); got != "https://haohuo.jinritemai.com/views/product/item2?id=1" {
			t.Fatalf("unexpected url query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{"sale_title": "product", "data": map[string]interface{}{"yixiaoerId": "goods_1", "raw": map[string]interface{}{"id": "1"}}},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.GoodsDetail("acc_1", "https://haohuo.jinritemai.com/views/product/item2?id=1")
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	items := data["list"].([]interface{})
	if len(items) != 1 || items[0].(map[string]interface{})["sale_title"] != "product" {
		t.Fatalf("unexpected goods detail data: %#v", data)
	}
}

func TestGoodsDetailReturnsStructuredRemoteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "invalid product link", "code": "INVALID_URL"})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	_, err := client.GoodsDetail("acc_1", "not-a-link")
	if err == nil {
		t.Fatal("expected goods detail remote error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yx error, got %T", err)
	}
	if typed.Type != yxerrors.RemoteType || typed.Message != "invalid product link" {
		t.Fatalf("unexpected structured error: %#v", typed)
	}
}

func TestEntitlementsUsesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/platform-accounts/acc_1/entitlements" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"shopping_cart": true, "group_shopping": false},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Entitlements("acc_1")
	if err != nil {
		t.Fatal(err)
	}
	entitlements := result.(map[string]interface{})
	if entitlements["shopping_cart"] != true || entitlements["group_shopping"] != false {
		t.Fatalf("unexpected entitlements: %#v", entitlements)
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

func TestMusicUsesDouyinRecommendationChartByDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/music" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("categoryId"); got != defaultDouyinMusicCategoryID {
			t.Fatalf("unexpected default category id: %q", got)
		}
		if got := r.URL.Query().Get("categoryName"); got != defaultDouyinMusicCategoryName {
			t.Fatalf("unexpected default category name: %q", got)
		}
		if _, ok := r.URL.Query()["keyWord"]; ok {
			t.Fatalf("default chart query must not include keyWord: %q", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Music("acc_1", "", "", "", ""); err != nil {
		t.Fatalf("Music() error = %v", err)
	}
}

func TestMusicUsesChartInsteadOfKeywordWhenBothAreProvided(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("categoryId"); got != "7088298745502646280" {
			t.Fatalf("unexpected category id: %q", got)
		}
		if got := r.URL.Query().Get("categoryName"); got != "热门榜" {
			t.Fatalf("unexpected category name: %q", got)
		}
		if _, ok := r.URL.Query()["keyWord"]; ok {
			t.Fatalf("chart query must take precedence over keyWord: %q", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Music("acc_1", "周杰伦", "7088298745502646280", "热门榜", ""); err != nil {
		t.Fatalf("Music() error = %v", err)
	}
}

func TestMusicUsesKeywordSearchWhenNoChartIsSelected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("keyWord"); got != "周杰伦" {
			t.Fatalf("unexpected keyWord: %q", got)
		}
		if _, ok := r.URL.Query()["categoryId"]; ok {
			t.Fatalf("keyword query must not include a chart: %q", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Music("acc_1", " 周杰伦 ", "", "", ""); err != nil {
		t.Fatalf("Music() error = %v", err)
	}
}

func TestDramaTasksAlwaysSendsKeyWordAndPreservesResponse(t *testing.T) {
	tests := []struct {
		name    string
		keyword string
	}{
		{name: "empty keyword", keyword: ""},
		{name: "keyword", keyword: "护妻"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/platform-accounts/acc_1/drama-tasks" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				values, ok := r.URL.Query()["keyWord"]
				if !ok || len(values) != 1 || values[0] != tt.keyword {
					t.Fatalf("unexpected keyWord query: %#v, raw query: %q", values, r.URL.RawQuery)
				}
				if tt.keyword == "" && r.URL.RawQuery != "keyWord=" {
					t.Fatalf("empty keyword must remain explicit in query: %q", r.URL.RawQuery)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[{"yixiaoerId":"event/1","yixiaoerImageUrl":"http://wxapp.tc.qq.com/cover","yixiaoerName":"风浪过后护妻安康"}]}`))
			}))
			defer server.Close()

			client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
			result, err := client.DramaTasks("acc_1", tt.keyword)
			if err != nil {
				t.Fatalf("DramaTasks() error = %v", err)
			}

			items, ok := result.([]interface{})
			if !ok || len(items) != 1 {
				t.Fatalf("unexpected result: %#v", result)
			}
			item, ok := items[0].(map[string]interface{})
			if !ok || item["yixiaoerName"] != "风浪过后护妻安康" {
				t.Fatalf("unexpected drama task: %#v", items[0])
			}
		})
	}
}

func TestDramaTasksReturnsStructuredRemoteErrorWithRecoveryCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "当前账号不支持剧集",
			"code":    "DRAMA_UNSUPPORTED",
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	_, err := client.DramaTasks("acc_1", "护妻")
	if err == nil {
		t.Fatal("expected drama task query error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured remote error, got %T", err)
	}
	if typed.Type != yxerrors.RemoteType || typed.Category != yxerrors.RemoteType || typed.Retryable {
		t.Fatalf("unexpected drama query error: %#v", typed)
	}
	details, ok := typed.Details.(map[string]interface{})
	if !ok || details["code"] != "DRAMA_UNSUPPORTED" {
		t.Fatalf("expected remote code in error details, got %#v", typed.Details)
	}
	if typed.Hint == "" || typed.NextCommand != "yxer query drama-tasks acc_1 --json" {
		t.Fatalf("expected drama recovery guidance, got %#v", typed)
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

func TestCategoriesBuildsTreeFromBilibiliOpenParentIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/categories" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("publishType"); got != "video" {
			t.Fatalf("unexpected publishType query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []map[string]interface{}{
					{
						"yixiaoerId":   "root_1",
						"yixiaoerName": "生活",
						"raw": map[string]interface{}{
							"id":       "root_1",
							"name":     "生活",
							"parentId": nil,
						},
					},
					{
						"yixiaoerId":   "child_1",
						"yixiaoerName": "日常",
						"raw": map[string]interface{}{
							"id":       "child_1",
							"name":     "日常",
							"parentId": "root_1",
						},
					},
					{
						"yixiaoerId":   "root_2",
						"yixiaoerName": "游戏",
						"raw": map[string]interface{}{
							"id":       "root_2",
							"name":     "游戏",
							"parentId": nil,
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Categories("acc_1", "video")
	if err != nil {
		t.Fatal(err)
	}
	payload, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected categories payload map, got %#v", result)
	}
	items, ok := payload["dataList"].([]interface{})
	if !ok {
		t.Fatalf("expected dataList, got %#v", payload["dataList"])
	}
	if len(items) != 2 {
		t.Fatalf("expected two root categories, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["yixiaoerId"] != "root_1" {
		t.Fatalf("unexpected first root category: %#v", first)
	}
	children, ok := first["child"].([]interface{})
	if !ok {
		t.Fatalf("expected child categories, got %#v", first["child"])
	}
	if len(children) != 1 {
		t.Fatalf("expected one child category, got %d", len(children))
	}
	child := children[0].(map[string]interface{})
	if child["yixiaoerId"] != "child_1" {
		t.Fatalf("unexpected child category: %#v", child)
	}
	raw := child["raw"].(map[string]interface{})
	if raw["parentId"] != "root_1" {
		t.Fatalf("expected raw payload to be preserved, got %#v", raw)
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

func TestAccountIncrementsUsesExpectedEndpointAndMilliseconds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/overview/incremental" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("startTime"); got != "1786636800000" {
			t.Fatalf("unexpected startTime: %s", got)
		}
		if got := r.URL.Query().Get("endTime"); got != "1787241599999" {
			t.Fatalf("unexpected endTime: %s", got)
		}
		if got := r.URL.Query().Get("group"); got != "group_1" {
			t.Fatalf("unexpected group: %s", got)
		}
		if got := r.URL.Query().Get("platform"); got != "小红书" {
			t.Fatalf("unexpected platform: %s", got)
		}
		if got := r.URL.Query().Get("name"); got != "皮蛋不圆" {
			t.Fatalf("unexpected account name: %s", got)
		}
		if got := r.Header.Get("Authorization"); got != "test-key" {
			t.Fatalf("unexpected authorization: %s", got)
		}
		if got := r.Header.Get("x-client"); got != "web" {
			t.Fatalf("unexpected x-client: %s", got)
		}
		if got := r.Header.Get("x-platform"); got != "windows" {
			t.Fatalf("unexpected x-platform: %s", got)
		}
		if got := r.Header.Get("x-version"); got != "2.7.3" {
			t.Fatalf("unexpected x-version: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"overviewData": map[string]interface{}{"newFans": 6},
				"trend":        []map[string]interface{}{{"date": "2026-08-14", "value": 2}},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	start, end, err := ShanghaiDateRange("2026-08-14", "2026-08-20")
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.AccountIncrements(AccountIncrementOptions{StartTime: start, EndTime: end, GroupID: "group_1", Platform: "xhs", Name: "皮蛋不圆"})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if data["overviewData"].(map[string]interface{})["newFans"] != float64(6) {
		t.Fatalf("unexpected incremental data: %#v", data)
	}
}

func TestDMMessageStatsUsesExpectedEndpointAndFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/social/dm-stats" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("startTime"); got != "1786636800000" {
			t.Fatalf("unexpected startTime: %s", got)
		}
		if got := r.URL.Query().Get("endTime"); got != "1787241599999" {
			t.Fatalf("unexpected endTime: %s", got)
		}
		if got := r.URL.Query().Get("platformAccountIds"); got != "a1,a2" {
			t.Fatalf("unexpected account IDs: %s", got)
		}
		if r.Header.Get("Authorization") != "test-key" || r.Header.Get("x-client") != "web" {
			t.Fatalf("expected API key and web headers")
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"summary": []map[string]interface{}{{"platformAccountId": "a1", "inCount": 2, "outCount": 1}},
				"trend":   []map[string]interface{}{{"date": "2026-08-15", "inCount": 2, "outCount": 1}},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.DMMessageStats(DMMessageStatsOptions{
		StartTime:          1786636800000,
		EndTime:            1787241599999,
		PlatformAccountIDs: []string{"a1", "", "a2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if len(data["trend"].([]interface{})) != 1 {
		t.Fatalf("unexpected DM trend: %#v", data)
	}
}

func TestManagedAccountsUsesExpectedEndpointAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/platform/accounts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("size") != "1000" {
			t.Fatalf("unexpected pagination: %s", r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "test-key" || r.Header.Get("x-version") != "2.7.3" {
			t.Fatalf("expected API key and client headers")
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"totalSize": 1,
				"data":      []map[string]interface{}{{"id": "a1", "principal": nil, "operators": []interface{}{}}},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.ManagedAccounts(ManagedAccountOptions{Page: 1, Size: 1000})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if data["totalSize"] != float64(1) {
		t.Fatalf("unexpected managed account data: %#v", data)
	}
}

func TestAccountIncrementDashboardPreservesIncrementalKeysAndAddsSections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/overview/incremental":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
				"accounts": []map[string]interface{}{{
					"platformAccountId": "a1", "status": 1, "platformAccountName": "账号一", "overviewUpdatedAt": 1787365834205,
				}}, "summary": map[string]interface{}{}, "trends": []interface{}{},
			}})
		case "/social/dm-stats":
			if got := r.URL.Query().Get("platformAccountIds"); got != "a1" {
				t.Fatalf("expected incremental account filter, got %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
				"summary": []map[string]interface{}{{"platformAccountId": "a1", "inCount": 4, "outCount": 2}},
				"trend":   []map[string]interface{}{{"date": "1970-01-01", "inCount": 4, "outCount": 2}},
			}})
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
				"page": 1, "size": 1000, "totalSize": 0, "totalPage": 1, "data": []interface{}{},
			}})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.AccountIncrementDashboard(AccountIncrementOptions{StartTime: 1, EndTime: 2})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	for _, key := range []string{"accounts", "summary", "trends", "dmMessageStats", "managedAccounts"} {
		if _, ok := data[key]; !ok {
			t.Fatalf("missing dashboard key %q: %#v", key, data)
		}
	}
	account := data["accounts"].([]interface{})[0].(map[string]interface{})
	if account["dmInCount"] != float64(4) || account["dmOutCount"] != float64(2) || account["statusLabel"] != "正常" || account["dataUpdatedAt"] != int64(1787365834205) || account["dataUpdatedTime"] == "" {
		t.Fatalf("unexpected enriched account: %#v", account)
	}
	dm := data["dmMessageStats"].(map[string]interface{})
	totals := dm["totals"].(map[string]interface{})
	if totals["inCount"] != float64(4) || totals["outCount"] != float64(2) {
		t.Fatalf("unexpected DM totals: %#v", totals)
	}
	accountTotals := dm["incrementalAccountTotals"].(map[string]interface{})
	if accountTotals["inCount"] != float64(4) || accountTotals["outCount"] != float64(2) {
		t.Fatalf("unexpected incremental account DM totals: %#v", accountTotals)
	}
	if len(dm["dailyTrend"].([]interface{})) != 1 {
		t.Fatalf("unexpected daily DM trend: %#v", dm["dailyTrend"])
	}
}

func TestShanghaiDateRangeRejectsInvalidAndReversedDates(t *testing.T) {
	if _, _, err := ShanghaiDateRange("2026/08/14", "2026-08-20"); err == nil {
		t.Fatal("expected invalid start date error")
	}
	if _, _, err := ShanghaiDateRange("2026-08-21", "2026-08-20"); err == nil {
		t.Fatal("expected reversed date error")
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

func TestMembersUsesExpectedEndpointAndFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/members" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("unexpected page query: %s", got)
		}
		if got := r.URL.Query().Get("size"); got != "20" {
			t.Fatalf("unexpected size query: %s", got)
		}
		if got := r.URL.Query()["statuses"]; len(got) != 2 || got[0] != "joined" || got[1] != "pending" {
			t.Fatalf("unexpected statuses query: %#v", got)
		}
		if got := r.URL.Query().Get("keyWords"); got != "张三" {
			t.Fatalf("unexpected keyWords query: %s", got)
		}
		if got := r.URL.Query().Get("role"); got != "member" {
			t.Fatalf("unexpected role query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "member_1", "name": "张三"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Members(MembersOptions{
		Page:     2,
		Size:     20,
		Statuses: []string{"joined", "pending"},
		KeyWords: "张三",
		Role:     "member",
	})
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one member, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["name"] != "张三" {
		t.Fatalf("unexpected member payload: %#v", first)
	}
}

func TestAccountGroupsUsesExpectedEndpointAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/groups" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("unexpected page query: %s", got)
		}
		if got := r.URL.Query().Get("size"); got != "20" {
			t.Fatalf("unexpected size query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "grp_1", "name": "核心账号组"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.AccountGroups(AccountGroupOptions{Page: 2, Size: 20})
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one account group, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["name"] != "核心账号组" {
		t.Fatalf("unexpected account group payload: %#v", first)
	}
}

func TestCreateAccountGroupUsesExpectedEndpointAndBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/groups" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["name"] != "核心账号组" {
			t.Fatalf("unexpected body: %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "grp_1",
				"name": "核心账号组",
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.CreateAccountGroup(map[string]interface{}{"name": "核心账号组"})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if data["id"] != "grp_1" || data["name"] != "核心账号组" {
		t.Fatalf("unexpected create account group result: %#v", data)
	}
}

func TestUpdateAccountGroupUsesExpectedEndpointAndBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/groups/grp_1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["name"] != "新版核心账号组" {
			t.Fatalf("unexpected body: %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "grp_1",
				"name": "新版核心账号组",
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.UpdateAccountGroup("grp_1", map[string]interface{}{"name": "新版核心账号组"})
	if err != nil {
		t.Fatal(err)
	}
	data := result.(map[string]interface{})
	if data["id"] != "grp_1" || data["name"] != "新版核心账号组" {
		t.Fatalf("unexpected update account group result: %#v", data)
	}
}

func TestDeleteAccountGroupUsesExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/groups/grp_1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.DeleteAccountGroup("grp_1")
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Fatalf("expected nil delete result for empty response, got %#v", result)
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

func TestPrepareUsesFirstOnlineAccountAcrossPagesForCategories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			switch r.URL.Query().Get("page") {
			case "1":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{
						"list": []map[string]interface{}{
							{"platformAccountId": "acc_1", "platformAccountName": "账号1", "status": 1},
							{"platformAccountId": "acc_2", "platformAccountName": "账号2", "status": 0},
						},
						"page":      1,
						"size":      50,
						"totalPage": 2,
						"totalSize": 3,
					},
				})
			case "2":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{
						"list": []map[string]interface{}{
							{"platformAccountId": "acc_3", "platformAccountName": "账号3", "status": 1},
						},
						"page":      2,
						"size":      50,
						"totalPage": 2,
						"totalSize": 3,
					},
				})
			default:
				t.Fatalf("unexpected page query: %s", r.URL.Query().Get("page"))
			}
		case "/platform-accounts/acc_1/categories":
			if got := r.URL.Query().Get("publishType"); got != "video" {
				t.Fatalf("unexpected publishType query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"id": "cat_1", "name": "分类1"},
				},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Prepare("小红书", "video")
	if err != nil {
		t.Fatal(err)
	}
	items, ok := result.Categories.([]interface{})
	if !ok {
		t.Fatalf("expected categories list, got %#v", result.Categories)
	}
	if len(items) != 1 {
		t.Fatalf("expected one category, got %d", len(items))
	}
	first := items[0].(map[string]interface{})
	if first["yixiaoerId"] != "cat_1" || first["yixiaoerName"] != "分类1" || first["raw"] == nil {
		t.Fatalf("unexpected category payload: %#v", first)
	}
	if len(result.Accounts) != 2 || result.Accounts[0]["platformAccountId"] != "acc_1" || result.Accounts[1]["platformAccountId"] != "acc_3" {
		t.Fatalf("expected online accounts in prepare response, got %#v", result.Accounts)
	}
}
