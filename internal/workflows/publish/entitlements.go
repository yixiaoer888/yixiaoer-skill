package publish

import (
	"fmt"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

// ShoppingCartAccountIDs returns only accounts whose form contains a non-empty
// shopping_cart. Group-shopping entitlement is intentionally out of scope.
func ShoppingCartAccountIDs(payload map[string]interface{}) []string {
	publishArgs := objectField(payload, "publishArgs")
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	seen := map[string]bool{}
	var accountIDs []string
	for _, item := range accountForms {
		form, _ := item.(map[string]interface{})
		contentForm := objectField(form, "contentPublishForm")
		if !hasShoppingCartValue(contentForm["shopping_cart"]) {
			continue
		}
		accountID := stringField(form, "platformAccountId")
		if accountID == "" {
			accountID = stringField(form, "account_id")
		}
		if accountID != "" && !seen[accountID] {
			seen[accountID] = true
			accountIDs = append(accountIDs, accountID)
		}
	}
	return accountIDs
}

func hasShoppingCartValue(value interface{}) bool {
	switch typed := value.(type) {
	case []interface{}:
		return len(typed) > 0
	case map[string]interface{}:
		return len(typed) > 0
	default:
		return false
	}
}

// shoppingCartEntitlementsSupported reports whether the generic entitlement
// endpoint can be used for the target platform. Duoduo Video accepts the
// manual goods_id cart object, but does not expose this endpoint.
func shoppingCartEntitlementsSupported(platforms []string) bool {
	if len(platforms) == 0 {
		return true
	}
	for _, platform := range platforms {
		if platformutil.CanonicalKey(platform) == "duoduoshipin" {
			return false
		}
	}
	return true
}

func AssertShoppingCartEntitlements(apiClient *api.Client, payload map[string]interface{}, platforms ...string) error {
	if !shoppingCartEntitlementsSupported(platforms) {
		return nil
	}
	accountIDs := ShoppingCartAccountIDs(payload)
	for _, accountID := range accountIDs {
		result, err := apiClient.Entitlements(accountID)
		if err != nil {
			return err
		}
		entitlements, _ := result.(map[string]interface{})
		shoppingCart, _ := entitlements["shopping_cart"].(bool)
		if shoppingCart {
			continue
		}
		return yxerrors.Usage("Shopping cart entitlement check failed", map[string]interface{}{
			"accountId":     accountID,
			"shopping_cart": false,
		}).
			WithCategory("shopping_cart_entitlement").
			WithHint("This account does not have shopping-cart publishing entitlement. Remove shopping_cart or use an entitled account.").
			WithNextCommand(fmt.Sprintf("yxer query entitlements %s", accountID))
	}
	return nil
}
