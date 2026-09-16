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

func TestTaobaoGuangheGoodsTabsMapsContentType(t *testing.T) {
	for _, tc := range []struct{ publishType, pageType string }{{"video", "video"}, {"imageText", "photo"}} {
		t.Run(tc.publishType, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/platform-accounts/acc_1/taobao-guanghe/goods-tabs" || r.URL.Query().Get("pageType") != tc.pageType {
					t.Fatalf("unexpected request: %s?%s", r.URL.Path, r.URL.RawQuery)
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{map[string]interface{}{"source": "mine"}}})
			}))
			defer server.Close()
			result, err := NewClient(config.Config{APIKey: "test", APIURL: server.URL}).TaobaoGuangheGoodsTabs("acc_1", tc.publishType)
			if err != nil || len(result.([]interface{})) != 1 {
				t.Fatalf("unexpected result=%#v err=%v", result, err)
			}
		})
	}
}

func TestTaobaoGuangheGoodsPreservesCursorAndQueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.URL.Path != "/platform-accounts/acc_1/taobao-guanghe/goods" || q.Get("pageType") != "photo" || q.Get("keyword") != "中文 商品" || q.Get("nextPage") != "opaque+/==" || q.Get("filterValue") != "f1" || q.Get("secondLevelFilterValue") != "f2" || q.Get("source") != "selected-source" {
			t.Fatalf("unexpected request: %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"dataList": []interface{}{}, "nextPage": nil}})
	}))
	defer server.Close()
	_, err := NewClient(config.Config{APIKey: "test", APIURL: server.URL}).TaobaoGuangheGoods("acc_1", TaobaoGuangheGoodsOptions{PublishType: "imageText", Keyword: "中文 商品", NextPage: "opaque+/==", FilterValue: "f1", SecondFilterValue: "f2", Source: "selected-source"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTaobaoGuangheGoodsRejectsUnsupportedTypeBeforeRequest(t *testing.T) {
	client := NewClient(config.Config{})
	_, err := client.TaobaoGuangheGoodsTabs("acc_1", "article")
	var typed *yxerrors.Error
	if !errors.As(err, &typed) || typed.Code != "taobao_guanghe_invalid_content_type" || typed.Category != "taobao_guanghe_goods" || typed.Retryable {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestTaobaoGuangheGoodsRequiresExplicitSourceBeforeRequest(t *testing.T) {
	client := NewClient(config.Config{})
	_, err := client.TaobaoGuangheGoods("acc_1", TaobaoGuangheGoodsOptions{PublishType: "video"})
	var typed *yxerrors.Error
	if !errors.As(err, &typed) || typed.Code != "taobao_guanghe_goods_source_required" || typed.Category != "taobao_guanghe_goods_source" || typed.Retryable {
		t.Fatalf("unexpected error: %#v", err)
	}
	if typed.NextCommand != "yxer query taobao-guanghe-goods-tabs acc_1 --type video --json" {
		t.Fatalf("unexpected next command: %q", typed.NextCommand)
	}
}
