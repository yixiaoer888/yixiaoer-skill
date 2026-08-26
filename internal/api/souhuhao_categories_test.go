package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestCategoryPathViewBuildsRecursivePaths(t *testing.T) {
	tree := []interface{}{
		souhuhaoCanonicalCategory("43", "生活", map[string]interface{}{"id": 43}, []interface{}{
			souhuhaoCanonicalCategory("58", "生活百态", map[string]interface{}{"id": 58, "channelId": 43}, []interface{}{
				souhuhaoCanonicalCategory("99", "访谈", map[string]interface{}{"id": 99, "channelId": 58}, nil),
			}),
		}),
	}

	view, err := CategoryPathView(tree)
	if err != nil {
		t.Fatal(err)
	}
	data, ok := view.(map[string]interface{})
	if !ok {
		t.Fatalf("expected category path view map, got %#v", view)
	}
	paths, ok := data["paths"].([]interface{})
	if !ok || len(paths) != 1 {
		t.Fatalf("expected one leaf path, got %#v", data["paths"])
	}
	path, ok := paths[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected path object, got %#v", paths[0])
	}
	if path["path"] != "生活 > 生活百态 > 访谈" {
		t.Fatalf("unexpected path: %#v", path["path"])
	}
	nodes, ok := path["nodes"].([]interface{})
	if !ok || len(nodes) != 3 {
		t.Fatalf("expected three path nodes, got %#v", path["nodes"])
	}
	if path["category"] == nil {
		t.Fatalf("expected complete category path, got %#v", path)
	}
	if _, ok := data["categories"].([]interface{}); !ok {
		t.Fatalf("expected categories to be a tree, got %#v", data["categories"])
	}
}

func TestNormalizeSohuhaoVideoPayloadBuildsNestedRawAndParentPath(t *testing.T) {
	var categoriesCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categoriesCalls++
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/platform-accounts/acc_sh_1/categories" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("publishType"); got != "video" {
			t.Fatalf("unexpected publishType query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []interface{}{
					souhuhaoCanonicalCategory("43", "生活", map[string]interface{}{"id": 43, "cmsChannelId": 206}, []interface{}{
						souhuhaoCanonicalCategory("58", "生活百态", map[string]interface{}{"id": 58, "channelId": 43, "name": "生活百态"}, nil),
					}),
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	payload := map[string]interface{}{
		"publishType": "video",
		"platforms":   []interface{}{"搜狐号"},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_sh_1",
					"contentPublishForm": map[string]interface{}{
						"category": []interface{}{
							map[string]interface{}{
								"id":   "58",
								"text": "生活百态",
								"raw":  map[string]interface{}{"id": 58, "channelId": 43},
							},
						},
					},
				},
			},
		},
	}
	original, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.NormalizeSohuhaoVideoPayload(payload)
	if err != nil {
		t.Fatal(err)
	}
	if categoriesCalls != 1 {
		t.Fatalf("expected one categories query, got %d", categoriesCalls)
	}
	forms := got["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})
	cpf := forms[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	categories := cpf["category"].([]interface{})
	if len(categories) != 2 {
		t.Fatalf("expected parent and child category path, got %#v", categories)
	}
	for i, expected := range []string{"43", "58"} {
		item := categories[i].(map[string]interface{})
		if item["id"] != expected {
			t.Fatalf("expected category %d id %q, got %#v", i, expected, item["id"])
		}
		canonical := item["raw"].(map[string]interface{})
		if canonical["yixiaoerId"] != expected {
			t.Fatalf("expected category %d canonical id %q, got %#v", i, expected, canonical["yixiaoerId"])
		}
		platformRaw := canonical["raw"].(map[string]interface{})
		if platformRaw["id"] == nil {
			t.Fatalf("expected category %d platform raw id, got %#v", i, platformRaw)
		}
	}
	updated, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, updated) {
		t.Fatalf("expected caller payload to remain unchanged, before=%s after=%s", original, updated)
	}
}

func TestNormalizeSohuhaoVideoPayloadRejectsInvalidCategoryDataWithStructuredError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []interface{}{
					souhuhaoCanonicalCategory("43", "生活", map[string]interface{}{"name": "生活"}, nil),
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	_, err := client.NormalizeSohuhaoVideoPayload(map[string]interface{}{
		"publishType": "video",
		"platforms":   []interface{}{"搜狐号"},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{map[string]interface{}{
				"platformAccountId": "acc_invalid",
				"contentPublishForm": map[string]interface{}{
					"category": []interface{}{map[string]interface{}{"id": "43", "text": "生活"}},
				},
			}},
		},
	})
	if err == nil {
		t.Fatal("expected invalid category data error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yx error, got %T", err)
	}
	if typed.Code != "sohuhao_category_data_invalid" || typed.Category != "sohuhao_category" || typed.Message != "搜狐号视频分类数据无效" {
		t.Fatalf("unexpected category data error: %#v", typed)
	}
	if typed.Retryable || typed.Hint == "" || typed.NextCommand != "yxer query categories acc_invalid --type video --json" {
		t.Fatalf("unexpected error guidance: %#v", typed)
	}
	details, ok := typed.Details.(map[string]interface{})
	if !ok || details["accountId"] != "acc_invalid" || details["cause"] == nil {
		t.Fatalf("expected account and cause details, got %#v", typed.Details)
	}
	available, ok := details["availablePaths"].([]interface{})
	if !ok || len(available) != 0 {
		t.Fatalf("expected empty available paths, got %#v", details["availablePaths"])
	}
}

func TestNormalizeSohuhaoVideoPayloadRejectsAmbiguousName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []interface{}{
					souhuhaoCanonicalCategory("1", "生活", map[string]interface{}{"id": 1}, nil),
					souhuhaoCanonicalCategory("2", "生活", map[string]interface{}{"id": 2}, nil),
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	_, err := client.NormalizeSohuhaoVideoPayload(map[string]interface{}{
		"publishType": "video",
		"platforms":   []interface{}{"souhuhao"},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{map[string]interface{}{
				"platformAccountId": "acc_ambiguous",
				"contentPublishForm": map[string]interface{}{
					"category": []interface{}{map[string]interface{}{"text": "生活"}},
				},
			}},
		},
	})
	if err == nil {
		t.Fatal("expected ambiguous category error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yx error, got %T", err)
	}
	if typed.Code != "sohuhao_category_ambiguous" || typed.Message != "搜狐号视频分类名称不唯一" || typed.Retryable {
		t.Fatalf("unexpected ambiguous category error: %#v", typed)
	}
	details, ok := typed.Details.(map[string]interface{})
	if !ok || details["accountId"] != "acc_ambiguous" || details["requested"] == nil {
		t.Fatalf("expected ambiguity details, got %#v", typed.Details)
	}
	available, ok := details["availablePaths"].([]interface{})
	if !ok || len(available) != 2 {
		t.Fatalf("expected both available paths, got %#v", details["availablePaths"])
	}
}

func TestNormalizeSohuhaoVideoPayloadRejectsConflictingIdentity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []interface{}{
					souhuhaoCanonicalCategory("58", "生活百态", map[string]interface{}{"id": 58}, nil),
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	_, err := client.NormalizeSohuhaoVideoPayload(map[string]interface{}{
		"publishType": "video",
		"platforms":   []interface{}{"搜狐号"},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{map[string]interface{}{
				"platformAccountId": "acc_conflict",
				"contentPublishForm": map[string]interface{}{
					"category": []interface{}{map[string]interface{}{
						"id":   "58",
						"text": "生活百态",
						"raw":  map[string]interface{}{"id": 59, "name": "生活百态"},
					}},
				},
			}},
		},
	})
	if err == nil {
		t.Fatal("expected conflicting identity error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yx error, got %T", err)
	}
	if typed.Code != "sohuhao_category_invalid" || typed.Message != "搜狐号视频分类参数无效" || typed.Retryable {
		t.Fatalf("unexpected conflicting identity error: %#v", typed)
	}
}

func TestNormalizeSohuhaoVideoPayloadSkipsOtherPlatforms(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("category query should not run for another platform")
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	payload := map[string]interface{}{
		"publishType": "video",
		"platforms":   []interface{}{"百家号"},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{map[string]interface{}{
				"platformAccountId": "acc_other",
				"contentPublishForm": map[string]interface{}{
					"category": []interface{}{map[string]interface{}{"id": "1", "text": "科技"}},
				},
			}},
		},
	}
	original, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.NormalizeSohuhaoVideoPayload(payload)
	if err != nil {
		t.Fatal(err)
	}
	updated, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, updated) {
		t.Fatalf("expected non-Sohu payload to remain unchanged, before=%s after=%s", original, updated)
	}
}

func souhuhaoCanonicalCategory(id, name string, raw map[string]interface{}, children []interface{}) map[string]interface{} {
	category := map[string]interface{}{
		"yixiaoerId":   id,
		"yixiaoerName": name,
		"raw":          raw,
	}
	if children != nil {
		category["child"] = children
	}
	return category
}
