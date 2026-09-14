package publish

import (
	"encoding/json"
	"testing"
)

func TestBuildPublishBodyMapsTaobaoGuangheScheduledTime(t *testing.T) {
	cpf := map[string]interface{}{"formType": "task", "scheduledTime": float64(1760000000000)}
	publishArgs := map[string]interface{}{"accountForms": []interface{}{map[string]interface{}{"platformAccountId": "acc_1", "contentPublishForm": cpf}}}
	payload := map[string]interface{}{"action": "publish", "publishArgs": publishArgs}
	body := BuildPublishBody(payload, publishArgs, "video", []string{"淘宝光合"}, "cloud", "")
	got := body["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if got["prePubTime"] != float64(1760000000000) {
		t.Fatalf("expected prePubTime milliseconds, got %#v", got)
	}
	if _, exists := got["scheduledTime"]; exists {
		t.Fatalf("scheduledTime must be removed from gateway body: %#v", got)
	}
}

func TestBuildPublishBodyPreservesTaobaoGuangheGoodsSnapshot(t *testing.T) {
	goods := []interface{}{map[string]interface{}{"yixiaoerId": "goods-1", "yixiaoerName": "商品", "price": float64(99), "raw": map[string]interface{}{"nested": map[string]interface{}{"token": "preserve"}}}}
	cpf := map[string]interface{}{"shopping_cart": goods}
	publishArgs := map[string]interface{}{"accountForms": []interface{}{map[string]interface{}{"contentPublishForm": cpf}}}
	want, _ := json.Marshal(goods)
	body := BuildPublishBody(map[string]interface{}{"publishArgs": publishArgs}, publishArgs, "video", []string{"淘宝光合"}, "cloud", "")
	gotGoods := body["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})["shopping_cart"]
	got, _ := json.Marshal(gotGoods)
	if string(got) != string(want) {
		t.Fatalf("goods snapshot changed: got %s want %s", got, want)
	}
}

func TestBuildPublishBodyDoesNotMapOtherPlatformScheduledTime(t *testing.T) {
	cpf := map[string]interface{}{"scheduledTime": float64(1760000000000)}
	publishArgs := map[string]interface{}{"accountForms": []interface{}{map[string]interface{}{"contentPublishForm": cpf}}}
	body := BuildPublishBody(map[string]interface{}{"publishArgs": publishArgs}, publishArgs, "video", []string{"抖音"}, "cloud", "")
	got := body["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if got["scheduledTime"] != float64(1760000000000) || got["prePubTime"] != nil {
		t.Fatalf("unexpected non-Taobao schedule mapping: %#v", got)
	}
}
