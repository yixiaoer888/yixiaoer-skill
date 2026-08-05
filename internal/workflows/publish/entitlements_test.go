package publish

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestShoppingCartAccountIDsSkipsFormsWithoutShoppingCart(t *testing.T) {
	payload := entitlementTestPayload("acc_cart", true)
	forms := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})
	forms = append(forms, map[string]interface{}{
		"platformAccountId":  "acc_plain",
		"contentPublishForm": map[string]interface{}{},
	})
	payload["publishArgs"].(map[string]interface{})["accountForms"] = forms

	ids := ShoppingCartAccountIDs(payload)
	if len(ids) != 1 || ids[0] != "acc_cart" {
		t.Fatalf("unexpected shopping-cart account ids: %#v", ids)
	}
}

func TestAssertShoppingCartEntitlementsRejectsDeniedAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_1/entitlements" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":{"shopping_cart":false,"group_shopping":true}}`))
	}))
	defer server.Close()

	client := api.NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	err := AssertShoppingCartEntitlements(client, entitlementTestPayload("acc_1", true))
	if err == nil {
		t.Fatal("expected denied shopping-cart entitlement")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yx error, got %T", err)
	}
	if typed.Category != "shopping_cart_entitlement" || typed.NextCommand != "yxer query entitlements acc_1" {
		t.Fatalf("unexpected entitlement error: %#v", typed)
	}
}

func TestAssertShoppingCartEntitlementsSkipsPayloadWithoutShoppingCart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("entitlements endpoint should not be called without shopping_cart")
	}))
	defer server.Close()

	client := api.NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if err := AssertShoppingCartEntitlements(client, entitlementTestPayload("acc_1", false)); err != nil {
		t.Fatal(err)
	}
}

func entitlementTestPayload(accountID string, includeShoppingCart bool) map[string]interface{} {
	contentForm := map[string]interface{}{}
	if includeShoppingCart {
		contentForm["shopping_cart"] = []interface{}{map[string]interface{}{"data": map[string]interface{}{"yixiaoerId": "goods_1"}}}
	}
	return map[string]interface{}{
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId":  accountID,
					"contentPublishForm": contentForm,
				},
			},
		},
	}
}
