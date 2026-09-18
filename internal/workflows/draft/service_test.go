package draft

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestSaveSendsNormalizedWeixinArticleCover(t *testing.T) {
	var gotBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/taskSets/drafts" {
			t.Fatalf("unexpected draft request: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"id": "draft_weixin_article_1"},
		})
	}))
	defer server.Close()

	payload := map[string]interface{}{
		"action":      "publish",
		"publishType": "article",
		"platforms":   []interface{}{"微信公众号"},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_weixin_draft_1",
				},
			},
			"platformForms": map[string]interface{}{
				"微信公众号": map[string]interface{}{
					"articles": []interface{}{
						map[string]interface{}{
							"title":   "公众号草稿标题",
							"content": "<p>公众号草稿正文</p>",
							"type":    float64(1),
							"cover": map[string]interface{}{
								"key": "wx-draft-cover-key",
								"raw": map[string]interface{}{"source": "upload"},
							},
						},
					},
					"notifySubscribers": float64(0),
					"pubType":           float64(0),
				},
			},
		},
	}

	_, err := NewService(app.New(config.Config{APIKey: "test-key", APIURL: server.URL})).Save(payload)
	if err != nil {
		t.Fatal(err)
	}
	if gotBody["coverKey"] != "wx-draft-cover-key" {
		t.Fatalf("expected draft request coverKey, got %#v", gotBody["coverKey"])
	}
	cover := gotBody["cover"].(map[string]interface{})
	if cover["key"] != "wx-draft-cover-key" {
		t.Fatalf("expected draft request cover, got %#v", cover)
	}
	if gotBody["isDraft"] != true {
		t.Fatalf("expected isDraft=true, got %#v", gotBody["isDraft"])
	}
	if _, exists := gotBody["action"]; exists {
		t.Fatalf("did not expect action in draft request: %#v", gotBody)
	}
}
