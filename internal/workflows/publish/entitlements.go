package publish

import (
	"fmt"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
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
		shoppingCart, _ := contentForm["shopping_cart"].([]interface{})
		if len(shoppingCart) == 0 {
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

func AssertShoppingCartEntitlements(apiClient *api.Client, payload map[string]interface{}) error {
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
