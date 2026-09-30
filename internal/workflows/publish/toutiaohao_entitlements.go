package publish

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const (
	toutiaohaoExclusiveSelectionDeclaration = 9
	toutiaohaoFirstSelectionDeclaration     = 10
	toutiaohaoExclusiveSelectionRight       = "right_exclusive_selection"
	toutiaohaoFirstSelectionRight           = "right_first_selection"
)

type toutiaohaoDeclarationRequirement struct {
	accountID   string
	declaration int
	right       string
}

// AssertToutiaohaoDeclarationEntitlements checks the account rights required
// by the two account-specific Toutiaohao declaration options. Ordinary
// declarations and article isFirst are independent of this check.
func AssertToutiaohaoDeclarationEntitlements(apiClient *api.Client, payload map[string]interface{}, platforms ...string) error {
	if !hasToutiaohaoPlatform(platforms) {
		return nil
	}
	requirements := toutiaohaoDeclarationRequirements(payload)
	if len(requirements) == 0 {
		return nil
	}
	if apiClient == nil {
		return yxerrors.Internal("Toutiaohao account-info client is unavailable", nil).
			WithCategory("toutiaohao_account_info")
	}

	accountInfoByID := make(map[string]interface{}, len(requirements))
	for _, requirement := range requirements {
		accountInfo, ok := accountInfoByID[requirement.accountID]
		if !ok {
			var err error
			accountInfo, err = apiClient.AccountInfo(requirement.accountID)
			if err != nil {
				return decorateToutiaohaoAccountInfoError(err, requirement.accountID)
			}
			accountInfoByID[requirement.accountID] = accountInfo
		}

		rightValue, found := toutiaohaoRightValue(accountInfo, requirement.right)
		if found && toutiaohaoRightEnabled(rightValue) {
			continue
		}
		return toutiaohaoMissingRightError(requirement, found, rightValue)
	}
	return nil
}

func hasToutiaohaoPlatform(platforms []string) bool {
	for _, platform := range platforms {
		if platformutil.CanonicalKey(platform) == "toutiaohao" {
			return true
		}
	}
	return false
}

func toutiaohaoDeclarationRequirements(payload map[string]interface{}) []toutiaohaoDeclarationRequirement {
	publishArgs := objectField(payload, "publishArgs")
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	seen := map[string]map[int]bool{}
	requirements := make([]toutiaohaoDeclarationRequirement, 0)
	for _, item := range accountForms {
		form, _ := item.(map[string]interface{})
		contentForm := objectField(form, "contentPublishForm")
		declaration, ok := toutiaohaoDeclarationNumber(contentForm["declaration"])
		if !ok {
			continue
		}
		right, ok := toutiaohaoDeclarationRight(declaration)
		if !ok {
			continue
		}
		accountID := stringField(form, "platformAccountId")
		if accountID == "" {
			accountID = stringField(form, "account_id")
		}
		if accountID == "" {
			continue
		}
		if seen[accountID] == nil {
			seen[accountID] = map[int]bool{}
		}
		if seen[accountID][declaration] {
			continue
		}
		seen[accountID][declaration] = true
		requirements = append(requirements, toutiaohaoDeclarationRequirement{
			accountID:   accountID,
			declaration: declaration,
			right:       right,
		})
	}
	return requirements
}

func toutiaohaoDeclarationNumber(value interface{}) (int, bool) {
	switch typed := value.(type) {
	case int:
		return typed, true
	case int8:
		return int(typed), true
	case int16:
		return int(typed), true
	case int32:
		return int(typed), true
	case int64:
		return int(typed), true
	case uint:
		return int(typed), true
	case uint8:
		return int(typed), true
	case uint16:
		return int(typed), true
	case uint32:
		return int(typed), true
	case uint64:
		if uint64(int(typed)) != typed {
			return 0, false
		}
		return int(typed), true
	case float64:
		if math.Trunc(typed) != typed {
			return 0, false
		}
		return int(typed), true
	case float32:
		value := float64(typed)
		if math.Trunc(value) != value {
			return 0, false
		}
		return int(typed), true
	case json.Number:
		parsed, err := strconv.Atoi(string(typed))
		return parsed, err == nil
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		return parsed, err == nil
	default:
		return 0, false
	}
}

func toutiaohaoDeclarationRight(declaration int) (string, bool) {
	switch declaration {
	case toutiaohaoExclusiveSelectionDeclaration:
		return toutiaohaoExclusiveSelectionRight, true
	case toutiaohaoFirstSelectionDeclaration:
		return toutiaohaoFirstSelectionRight, true
	default:
		return "", false
	}
}

func toutiaohaoRightValue(accountInfo interface{}, rightKey string) (interface{}, bool) {
	root, _ := accountInfo.(map[string]interface{})
	if root == nil {
		return nil, false
	}
	candidates := []interface{}{
		root,
		root["rights"],
		root["data"],
		root["accountInfo"],
	}
	if data, ok := root["data"].(map[string]interface{}); ok {
		candidates = append(candidates, data["rights"], data["accountInfo"])
	}
	for _, candidate := range candidates {
		object, ok := candidate.(map[string]interface{})
		if !ok || object == nil {
			continue
		}
		if value, exists := object[rightKey]; exists {
			return value, true
		}
		if rights, ok := object["rights"].(map[string]interface{}); ok {
			if value, exists := rights[rightKey]; exists {
				return value, true
			}
		}
	}
	return nil, false
}

func toutiaohaoRightEnabled(value interface{}) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		value := strings.ToLower(strings.TrimSpace(typed))
		return value == "1" || value == "true"
	case map[string]interface{}:
		if enabled, ok := typed["enabled"].(bool); ok && enabled {
			return true
		}
		if rightValue, ok := typed["value"]; ok && toutiaohaoRightValueEnabled(rightValue) {
			return true
		}
		for _, key := range []string{"status", "code"} {
			if flag, ok := typed[key]; ok && toutiaohaoRightStatusEnabled(flag) {
				return true
			}
		}
	default:
		return toutiaohaoRightNumericOne(value)
	}
	return false
}

func toutiaohaoRightValueEnabled(value interface{}) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	default:
		return toutiaohaoRightNumericOne(value)
	}
}

func toutiaohaoRightNumericOne(value interface{}) bool {
	switch typed := value.(type) {
	case float64:
		return typed == 1
	case float32:
		return typed == 1
	case int:
		return typed == 1
	case int8:
		return typed == 1
	case int16:
		return typed == 1
	case int32:
		return typed == 1
	case int64:
		return typed == 1
	case uint:
		return typed == 1
	case uint8:
		return typed == 1
	case uint16:
		return typed == 1
	case uint32:
		return typed == 1
	case uint64:
		return typed == 1
	default:
		return false
	}
}

func toutiaohaoRightStatusEnabled(value interface{}) bool {
	if toutiaohaoRightValueEnabled(value) {
		return true
	}
	text, ok := value.(string)
	return ok && strings.TrimSpace(text) == "1"
}

func decorateToutiaohaoAccountInfoError(err error, accountID string) error {
	var typed *yxerrors.Error
	if errors.As(err, &typed) {
		typed.WithCategory("toutiaohao_account_info").
			WithHint("查询头条号账号权益失败；请确认账号仍在线且 API 凭证有效后重试。")
		return typed
	}
	return yxerrors.Remote("Failed to query Toutiaohao account info", map[string]interface{}{
		"accountId": accountID,
		"cause":     err.Error(),
	}).WithCategory("toutiaohao_account_info").
		WithHint("查询头条号账号权益失败；请确认账号仍在线且 API 凭证有效后重试。")
}

func toutiaohaoMissingRightError(requirement toutiaohaoDeclarationRequirement, found bool, value interface{}) error {
	return yxerrors.New(yxerrors.ValidationType, "toutiaohao_declaration_entitlement_required", fmt.Sprintf("Toutiaohao declaration=%d requires account right %s", requirement.declaration, requirement.right), map[string]interface{}{
		"accountId":     requirement.accountID,
		"declaration":   requirement.declaration,
		"requiredRight": requirement.right,
		"rightFound":    found,
		"rightValue":    value,
	}).WithCategory("toutiaohao_declaration_entitlement").
		WithHint(fmt.Sprintf("账号 %s 没有 %s；请改用有对应权益的头条号账号，或将 declaration 改为 0/1/2/3/6/7/8。", requirement.accountID, requirement.right)).
		WithRetryable(false)
}
