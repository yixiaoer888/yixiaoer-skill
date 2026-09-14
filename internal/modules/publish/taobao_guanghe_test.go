package publish

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func taobaoPayload(publishType string, goods []interface{}) map[string]interface{} {
	cpf := map[string]interface{}{"formType": "task", "shopping_cart": goods}
	form := map[string]interface{}{"platformAccountId": "acc_1", "contentPublishForm": cpf}
	if publishType == "video" {
		form["video"] = map[string]interface{}{"key": "video.mp4", "size": 1, "width": 1080, "height": 1920}
		form["cover"] = map[string]interface{}{"key": "cover.jpg", "size": 1, "width": 1080, "height": 1920}
		form["coverKey"] = "cover.jpg"
	} else {
		cpf["title"] = "标题"
		form["images"] = []interface{}{map[string]interface{}{"key": "image.jpg", "size": 1, "width": 720, "height": 720}}
	}
	return map[string]interface{}{"action": "publish", "publishType": publishType, "platforms": []interface{}{"淘宝光合"}, "publishArgs": map[string]interface{}{"accountForms": []interface{}{form}}}
}

func taobaoGoods(count int) []interface{} {
	items := make([]interface{}, count)
	for i := range items {
		items[i] = map[string]interface{}{"yixiaoerId": string(rune('a' + i)), "yixiaoerName": "商品", "raw": map[string]interface{}{"nested": map[string]interface{}{"url": "https://example.com/goods"}}}
	}
	return items
}

func TestValidateTaobaoGuangheGoodsLimitsAndShape(t *testing.T) {
	for _, count := range []int{0, 1, 6} {
		if err := ValidateTaobaoGuanghePublish("淘宝光合", "video", "cloud", taobaoPayload("video", taobaoGoods(count))); err != nil {
			t.Fatalf("count %d should pass: %v", count, err)
		}
	}
	for _, tc := range []struct {
		name  string
		goods []interface{}
		code  string
	}{
		{"seven", taobaoGoods(7), "taobao_guanghe_goods_limit"},
		{"missing raw", []interface{}{map[string]interface{}{"yixiaoerId": "1", "yixiaoerName": "商品"}}, "taobao_guanghe_goods_invalid"},
		{"duplicate", []interface{}{taobaoGoods(1)[0], taobaoGoods(1)[0]}, "taobao_guanghe_goods_invalid"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var typed *yxerrors.Error
			err := ValidateTaobaoGuanghePublish("淘宝光合", "video", "cloud", taobaoPayload("video", tc.goods))
			if !errors.As(err, &typed) || typed.Code != tc.code || typed.Retryable {
				t.Fatalf("unexpected error: %#v", err)
			}
		})
	}
}

func TestTaobaoGuangheShoppingCartIsNotNormalized(t *testing.T) {
	payload := taobaoPayload("video", taobaoGoods(1))
	cpf := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	want, _ := json.Marshal(cpf["shopping_cart"])
	NormalizePlatformSpecificFields("video", []string{"淘宝光合"}, payload["publishArgs"].(map[string]interface{}))
	got, _ := json.Marshal(cpf["shopping_cart"])
	if string(got) != string(want) {
		t.Fatalf("shopping cart changed: got %s want %s", got, want)
	}
}

func TestTaobaoGuangheImageTextAndChannelRules(t *testing.T) {
	payload := taobaoPayload("imageText", nil)
	form := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	cpf := form["contentPublishForm"].(map[string]interface{})
	cpf["title"] = ""
	if err := ValidateTaobaoGuanghePublish("淘宝光合", "imageText", "cloud", payload); err == nil {
		t.Fatal("expected empty content error")
	}
	cpf["desc"] = "正文"
	form["images"].([]interface{})[0].(map[string]interface{})["width"] = 719
	if err := ValidateTaobaoGuanghePublish("淘宝光合", "imageText", "cloud", payload); err == nil {
		t.Fatal("expected image dimensions error")
	}
	if err := ValidateTaobaoGuanghePublish("淘宝光合", "video", "local", taobaoPayload("video", nil)); err != nil {
		t.Fatalf("local publish should be supported: %v", err)
	}
	if err := ValidateTaobaoGuanghePublish("淘宝光合", "article", "cloud", taobaoPayload("video", nil)); err == nil {
		t.Fatal("expected article error")
	}
}
