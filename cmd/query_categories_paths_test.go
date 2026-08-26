package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCategoriesCommandPathsReturnsCompleteTreeAndLeafPaths(t *testing.T) {
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_categories/categories" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("publishType") != "video" {
			t.Fatalf("unexpected publishType: %s", r.URL.Query().Get("publishType"))
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []map[string]interface{}{
					{
						"yixiaoerId":   "43",
						"yixiaoerName": "生活",
						"raw":          map[string]interface{}{"id": 43},
						"child": []map[string]interface{}{
							{
								"yixiaoerId":   "58",
								"yixiaoerName": "生活百态",
								"raw":          map[string]interface{}{"id": 58, "channelId": 43},
								"child": []map[string]interface{}{
									{
										"yixiaoerId":   "99",
										"yixiaoerName": "访谈",
										"raw":          map[string]interface{}{"id": 99, "channelId": 58},
									},
								},
							},
						},
					},
				},
			},
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)

	cmd := newCategoriesCmd()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"acc_categories", "--type", "video", "--paths"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout must be JSON: %v; output=%s", err, stdout.String())
	}
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got %#v", response["data"])
	}
	if _, ok := data["categories"].([]interface{}); !ok {
		t.Fatalf("expected categories tree, got %#v", data["categories"])
	}
	paths, ok := data["paths"].([]interface{})
	if !ok || len(paths) != 1 {
		t.Fatalf("expected one category path, got %#v", data["paths"])
	}
	path := paths[0].(map[string]interface{})
	if path["path"] != "生活 > 生活百态 > 访谈" {
		t.Fatalf("unexpected category path: %#v", path["path"])
	}
	nodes := path["nodes"].([]interface{})
	if len(nodes) != 3 {
		t.Fatalf("expected three path nodes, got %#v", nodes)
	}
	if _, hasParentID := path["parentId"]; hasParentID {
		t.Fatal("did not expect two-level convenience fields on a three-level path")
	}
}

func TestCategoriesCommandDefaultOutputKeepsDataListShape(t *testing.T) {
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"dataList": []map[string]interface{}{
					{"yixiaoerId": "1", "yixiaoerName": "科技", "raw": map[string]interface{}{"id": 1}},
				},
			},
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)

	cmd := newCategoriesCmd()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"acc_categories"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout must be JSON: %v", err)
	}
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got %#v", response["data"])
	}
	if _, ok := data["dataList"].([]interface{}); !ok {
		t.Fatalf("expected default dataList output, got %#v", data)
	}
	if _, ok := data["paths"]; ok {
		t.Fatal("did not expect paths field without --paths")
	}
}
