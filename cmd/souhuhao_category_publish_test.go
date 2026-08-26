package cmd

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestSouhuhaoDryRunResolvesCategoryPathBeforeBuildingRequest(t *testing.T) {
	withRepoRoot(t)
	configureAPIKey(t, "test-key")
	var categoryCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_sh_dry_run/categories" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		categoryCalls++
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []map[string]interface{}{
					{
						"yixiaoerId":   "43",
						"yixiaoerName": "生活",
						"raw":          map[string]interface{}{"id": 43, "cmsChannelId": 206},
						"child": []map[string]interface{}{
							{
								"yixiaoerId":   "58",
								"yixiaoerName": "生活百态",
								"raw":          map[string]interface{}{"id": 58, "channelId": 43},
							},
						},
					},
				},
			},
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)

	service := publishflow.NewService(testRuntime(t))
	result, err := service.DryRun(publishflow.ExecuteInput{
		PublishType:   "video",
		PlatformInput: "搜狐号",
		Payload: map[string]interface{}{
			"action":         "publish",
			"publishType":    "video",
			"platforms":      []interface{}{"搜狐号"},
			"publishChannel": "local",
			"clientId":       "client_1",
			"publishArgs": map[string]interface{}{
				"video": map[string]interface{}{
					"key":      "video-key",
					"size":     float64(1024),
					"width":    float64(720),
					"height":   float64(1280),
					"duration": float64(31),
				},
				"accountForms": []interface{}{map[string]interface{}{
					"platformAccountId": "acc_sh_dry_run",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(720),
						"height": float64(720),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "搜狐号分类测试标题",
						"description": "这是搜狐号分类测试描述。",
						"tags":        []interface{}{"生活"},
						"declaration": float64(0),
						"pubType":     float64(1),
						"category": []interface{}{map[string]interface{}{
							"id":   "58",
							"text": "生活百态",
							"raw":  map[string]interface{}{"id": 58, "channelId": 43},
						}},
					},
				}},
			},
		},
	})
	if err != nil {
		t.Fatalf("dry-run failed: %#v", err)
	}
	if categoryCalls != 1 {
		t.Fatalf("expected one category query, got %d", categoryCalls)
	}
	args := result.PublishBody["publishArgs"].(map[string]interface{})
	forms := args["accountForms"].([]interface{})
	form := forms[0].(map[string]interface{})
	cpf := form["contentPublishForm"].(map[string]interface{})
	categories := cpf["category"].([]interface{})
	if len(categories) != 2 {
		t.Fatalf("expected complete parent-child category path, got %#v", categories)
	}
	for index, expected := range []string{"43", "58"} {
		item := categories[index].(map[string]interface{})
		canonical := item["raw"].(map[string]interface{})
		platformRaw := canonical["raw"].(map[string]interface{})
		if item["id"] != expected || canonical["yixiaoerId"] != expected || platformRaw["id"] == nil {
			t.Fatalf("expected nested category at index %d, got %#v", index, item)
		}
	}
}

func TestSouhuhaoDryRunStopsBeforePublishWhenCategoryQueryFails(t *testing.T) {
	withRepoRoot(t)
	configureAPIKey(t, "test-key")
	publishCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/platform-accounts/acc_sh_query_error/categories":
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "category gateway unavailable"})
		case "/taskSets/v2":
			publishCalls++
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)

	service := publishflow.NewService(testRuntime(t))
	_, err := service.DryRun(publishflow.ExecuteInput{
		PublishType:   "video",
		PlatformInput: "搜狐号",
		Payload: map[string]interface{}{
			"action":         "publish",
			"publishType":    "video",
			"platforms":      []interface{}{"搜狐号"},
			"publishChannel": "local",
			"clientId":       "client_1",
			"publishArgs": map[string]interface{}{
				"accountForms": []interface{}{map[string]interface{}{
					"platformAccountId": "acc_sh_query_error",
					"contentPublishForm": map[string]interface{}{
						"category": []interface{}{map[string]interface{}{"id": "58", "text": "生活百态"}},
					},
				}},
			},
		},
	})
	if err == nil {
		t.Fatal("expected category query error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured category error, got %T", err)
	}
	if typed.Code != "sohuhao_category_query_failed" || typed.Category != "sohuhao_category" || !typed.Retryable {
		t.Fatalf("unexpected category query error: %#v", typed)
	}
	if typed.Message != "搜狐号视频分类查询失败" || typed.Hint == "" || typed.NextCommand != "yxer query categories acc_sh_query_error --type video --json" {
		t.Fatalf("unexpected category query guidance: %#v", typed)
	}
	details, ok := typed.Details.(map[string]interface{})
	if !ok || details["accountId"] != "acc_sh_query_error" || details["cause"] == nil {
		t.Fatalf("expected query failure details, got %#v", typed.Details)
	}
	requested, ok := details["requested"].([]interface{})
	if !ok || len(requested) != 1 {
		t.Fatalf("expected requested category details, got %#v", details["requested"])
	}
	available, ok := details["availablePaths"].([]interface{})
	if !ok || len(available) != 0 {
		t.Fatalf("expected empty available paths on query failure, got %#v", details["availablePaths"])
	}
	if publishCalls != 0 {
		t.Fatalf("expected no publish request, got %d", publishCalls)
	}
}
