package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestTaobaoGuangheQueryCommandOutputContract(t *testing.T) {
	withRepoRoot(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"dataList": []interface{}{}, "nextPage": nil}})
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)
	var out bytes.Buffer
	cmd := newTaobaoGuangheGoodsCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"acc_1", "--type", "video"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["ok"] != true || response["action"] != "taobao-guanghe-goods" || response["version"] == "" || response["data"] == nil {
		t.Fatalf("unexpected output contract: %#v", response)
	}
}

func TestTaobaoGuangheQueryCommandReturnsStructuredTypeError(t *testing.T) {
	err := newTaobaoGuangheGoodsTabsCmd().RunE(testCobraCommand(), []string{"acc_1"})
	if err == nil {
		t.Fatal("expected missing content type error")
	}
}

func TestTaobaoGuangheQueryFailureUsesStructuredStderr(t *testing.T) {
	withRepoRoot(t)
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "商品会话暂不可用", "code": "GOODS_SESSION_UNAVAILABLE"})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)
	var stdout, stderr bytes.Buffer
	code := ExecuteWithIO([]string{"query", "taobao-guanghe-goods", "acc_1", "--type", "video", "--json"}, &stdout, &stderr)
	if code == 0 || stdout.Len() != 0 {
		t.Fatalf("expected failure with empty stdout, code=%d stdout=%q", code, stdout.String())
	}
	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatalf("stderr must contain JSON: %v; output=%s", err, stderr.String())
	}
	errorObject := response["error"].(map[string]interface{})
	if errorObject["code"] != "GOODS_SESSION_UNAVAILABLE" || errorObject["category"] != "taobao_guanghe_goods_query" {
		t.Fatalf("unexpected error envelope: %#v", errorObject)
	}
}

func TestTaobaoGuanghePublishFormSourceContract(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	start := newPublishFormStartCmd()
	start.SetOut(&bytes.Buffer{})
	start.SetArgs([]string{"淘宝光合", "video", "--output", sessionPath})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	forms, _ := publishFormAccountForms(session.Payload)
	forms[0]["platformAccountId"] = "acc_1"
	if err := writePublishFormSession(sessionPath, session); err != nil {
		t.Fatal(err)
	}

	value := `{"data":{"dataList":[{"yixiaoerId":"goods_1","yixiaoerName":"商品一","price":99,"raw":{"id":"goods_1","secret":"preserved"}}]}}`
	choose := newPublishFormChooseCmd()
	choose.SetOut(&bytes.Buffer{})
	choose.SetArgs([]string{sessionPath, "shopping_cart", "--value", value, "--id", "goods_1", "--source-command", "yxer query taobao-guanghe-goods acc_1 --type video --json"})
	if err := choose.Execute(); err != nil {
		t.Fatal(err)
	}
	session, _ = readPublishFormSession(sessionPath)
	forms, _ = publishFormAccountForms(session.Payload)
	items := forms[0]["contentPublishForm"].(map[string]interface{})["shopping_cart"].([]interface{})
	if len(items) != 1 || items[0].(map[string]interface{})["price"] != float64(99) {
		t.Fatalf("unexpected selected goods: %#v", items)
	}
	if _, err := validatePublishFormProvenance(session); err != nil {
		t.Fatal(err)
	}
	items[0].(map[string]interface{})["price"] = float64(100)
	if _, err := validatePublishFormProvenance(session); err == nil {
		t.Fatal("expected object drift to fail provenance validation")
	}
}

func TestTaobaoGuanghePublishFormRejectsWrongSources(t *testing.T) {
	session := publishFormSession{Platform: "淘宝光合", Type: "imageText"}
	for _, tc := range []struct{ name, command, target, code string }{
		{"wrong command", "yxer query goods acc_1 --type imageText --json", "acc_1", "taobao_guanghe_goods_invalid"},
		{"wrong type", "yxer query taobao-guanghe-goods acc_1 --type video --json", "acc_1", "taobao_guanghe_goods_type_mismatch"},
		{"missing target", "yxer query taobao-guanghe-goods acc_1 --type imageText --json", "<account_id>", "taobao_guanghe_goods_account_mismatch"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var typed *yxerrors.Error
			err := validateTaobaoGuangheFormSource(session, "shopping_cart", tc.command, tc.target)
			if !errors.As(err, &typed) || typed.Code != tc.code || typed.Category != "taobao_guanghe_goods_source" {
				t.Fatalf("unexpected error: %#v", err)
			}
		})
	}
	if err := validatePublishFormSourceAccount("yxer query taobao-guanghe-goods acc_other --type imageText --json", "acc_1"); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("expected cross-account error, got %v", err)
	}
}

func TestTaobaoGuanghePublishFormSelectsUpToSixGoods(t *testing.T) {
	session := publishFormSession{Platform: "淘宝光合", Type: "video"}
	items := make([]interface{}, 7)
	ids := make([]string, 7)
	for i := range items {
		ids[i] = string(rune('a' + i))
		items[i] = map[string]interface{}{"yixiaoerId": ids[i], "yixiaoerName": "商品", "raw": map[string]interface{}{"id": ids[i]}}
	}
	selected, candidates, err := selectPublishFormCandidatesForSession(session, map[string]interface{}{"dataList": items}, "shopping_cart", -1, "", ids[:6])
	if err != nil || len(candidates) != 7 || len(selected.([]interface{})) != 6 {
		t.Fatalf("unexpected multi-selection selected=%#v candidates=%d err=%v", selected, len(candidates), err)
	}
	var typed *yxerrors.Error
	_, _, err = selectPublishFormCandidatesForSession(session, map[string]interface{}{"dataList": items}, "shopping_cart", -1, "", ids)
	if !errors.As(err, &typed) || typed.Code != "taobao_guanghe_goods_limit" {
		t.Fatalf("expected goods limit error, got %#v", err)
	}
}
