package publish

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/core/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type ExecuteInput struct {
	PublishType        string
	PlatformInput      string
	Payload            map[string]interface{}
	PositionalClientID string
	FlagChannel        string
	FlagClientID       string
}

type Service struct{}

var (
	PromptInput  io.Reader = os.Stdin
	PromptOutput io.Writer = os.Stdout
)

func NewService() Service {
	return Service{}
}

func (Service) Execute(input ExecuteInput) (map[string]interface{}, error) {
	input.PublishType = NormalizePublishType(input.PublishType)
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil {
		return nil, err
	}
	platforms := []string{platform}
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	resolvedPayload := cloneMap(input.Payload)
	channel, clientID, err := ResolvePublishMode(cfg, resolvedPayload, input.PositionalClientID, input.FlagChannel, input.FlagClientID)
	if err != nil {
		return nil, err
	}
	resolvedPayload["publishChannel"] = channel
	if clientID != "" {
		resolvedPayload["clientId"] = clientID
	} else {
		delete(resolvedPayload, "clientId")
	}
	NormalizeStandardPublishArgs(ExtractPublishArgs(resolvedPayload))
	publishArgs := ExtractPublishArgs(resolvedPayload)

	validator := schema.NewValidator(cfg.SchemaDir)
	for _, platform := range platforms {
		result := validator.Validate(platform, input.PublishType, resolvedPayload)
		if !result.Valid {
			return nil, yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint("请根据对应平台 schema 修正 payload 字段后重试。").
				WithNextCommand(fmt.Sprintf("yxer validate %s %s <payload.json>", platform, input.PublishType))
		}
	}
	preflight := Preflight(input.PublishType, platforms, resolvedPayload)
	if len(preflight.Errors) > 0 {
		return nil, yxerrors.Usage("Publish preflight failed", preflight.Errors).
			WithHint("请先完成资源上传、账号校验，并确保发布参数中不包含外部 URL。")
	}

	apiClient := client.New(cfg)
	if err := AssertAccountsOnline(apiClient, platforms, preflight.AccountIDs); err != nil {
		return nil, err
	}

	body := BuildPublishBody(resolvedPayload, publishArgs, input.PublishType, platforms)
	body["publishChannel"] = channel
	if clientID != "" {
		body["clientId"] = clientID
	} else {
		delete(body, "clientId")
	}
	result, err := apiClient.Publish(body)
	if err == nil {
		return result, nil
	}
	if !shouldOfferLocalPublishRetry(err, channel) {
		return nil, err
	}
	confirmed, confirmErr := confirmLocalPublishRetry(platform)
	if confirmErr != nil {
		return nil, confirmErr
	}
	if !confirmed {
		return nil, err
	}
	localChannel, localClientID, resolveErr := ResolvePublishMode(cfg, resolvedPayload, "", "local", "")
	if resolveErr != nil {
		return nil, resolveErr
	}
	resolvedPayload["publishChannel"] = localChannel
	resolvedPayload["clientId"] = localClientID
	body = BuildPublishBody(resolvedPayload, publishArgs, input.PublishType, platforms)
	body["publishChannel"] = localChannel
	body["clientId"] = localClientID
	return apiClient.Publish(body)
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return map[string]interface{}{}
	}
	dst := make(map[string]interface{}, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func BuildPublishBody(payload, publishArgs map[string]interface{}, publishType string, platforms []string) map[string]interface{} {
	body := map[string]interface{}{
		"publishType":    publishType,
		"platforms":      platforms,
		"publishArgs":    publishArgs,
		"publishChannel": "cloud",
	}
	if _, ok := payload["publishArgs"].(map[string]interface{}); ok {
		for key, value := range payload {
			if key == "action" {
				continue
			}
			if key == "publishArgs" {
				body[key] = publishArgs
				continue
			}
			body[key] = value
		}
		body["publishType"] = publishType
		body["platforms"] = platforms
		if _, ok := body["publishChannel"]; !ok {
			body["publishChannel"] = "cloud"
		}
		normalizePublishEnvelope(body, publishArgs, publishType)
		return body
	}
	copyOptionalPublishFields(body, payload)
	normalizePublishEnvelope(body, publishArgs, publishType)
	return body
}

func normalizePublishEnvelope(body, publishArgs map[string]interface{}, publishType string) {
	if body == nil {
		return
	}
	if publishArgs == nil {
		publishArgs = map[string]interface{}{}
	}
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	firstForm := firstObject(accountForms)
	firstCPF := objectField(firstForm, "contentPublishForm")

	if _, ok := body["cover"]; !ok {
		if cover := firstNonNil(
			publishArgs["cover"],
			firstForm["cover"],
			firstCPF["cover"],
		); cover != nil {
			body["cover"] = cover
		}
	}
	if stringField(body, "coverKey") == "" {
		if coverKey := firstNonEmptyString(
			stringField(publishArgs, "coverKey"),
			stringField(firstForm, "coverKey"),
			stringField(firstCPF, "coverKey"),
			stringField(objectField(body, "cover"), "key"),
		); coverKey != "" {
			body["coverKey"] = coverKey
		}
	}
	if stringField(body, "desc") == "" {
		if desc := inferOuterDesc(publishType, publishArgs, firstCPF); desc != "" {
			body["desc"] = desc
		}
	}
	if _, ok := body["isDraft"]; !ok {
		body["isDraft"] = false
	}
	if _, ok := body["isAppContent"]; !ok {
		body["isAppContent"] = false
	}
}

func inferOuterDesc(publishType string, publishArgs, contentPublishForm map[string]interface{}) string {
	switch NormalizePublishType(publishType) {
	case "article":
		return firstNonEmptyString(
			stringField(contentPublishForm, "title"),
			stringField(contentPublishForm, "description"),
			stringField(publishArgs, "content"),
			stringField(contentPublishForm, "content"),
		)
	case "video", "imageText":
		return firstNonEmptyString(
			stringField(contentPublishForm, "description"),
			stringField(contentPublishForm, "title"),
			stringField(publishArgs, "content"),
			stringField(contentPublishForm, "content"),
		)
	default:
		return firstNonEmptyString(
			stringField(contentPublishForm, "description"),
			stringField(contentPublishForm, "title"),
			stringField(publishArgs, "content"),
			stringField(contentPublishForm, "content"),
		)
	}
}

func firstObject(items []interface{}) map[string]interface{} {
	for _, item := range items {
		if obj, ok := item.(map[string]interface{}); ok {
			return obj
		}
	}
	return nil
}

func firstNonNil(values ...interface{}) interface{} {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && value != "<nil>" {
			return value
		}
	}
	return ""
}

func ResolvePublishMode(cfg config.Config, payload map[string]interface{}, positionalClientID, flagChannel, flagClientID string) (string, string, error) {
	channel := "cloud"
	clientID := ""
	if value, ok := payload["publishChannel"]; ok {
		channel = strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := payload["clientId"]; ok {
		clientID = strings.TrimSpace(fmt.Sprint(value))
	}
	if strings.TrimSpace(positionalClientID) != "" {
		channel = "local"
		clientID = strings.TrimSpace(positionalClientID)
	}
	if strings.TrimSpace(flagChannel) != "" {
		channel = strings.TrimSpace(flagChannel)
	}
	if strings.TrimSpace(flagClientID) != "" {
		clientID = strings.TrimSpace(flagClientID)
	}
	switch channel {
	case "", "cloud":
		channel = "cloud"
		clientID = ""
	case "local":
		if clientID == "" {
			clientID = strings.TrimSpace(cfg.LocalClientID)
		}
		if clientID == "" {
			return "", "", yxerrors.Usage(`clientId is required when publishChannel is "local"`, []string{
				`Run: yxer config set-local-client-id <clientId>`,
				`Or pass a fourth positional argument: yxer publish video <platform> payload.json <clientId>`,
				`Or pass flags: yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>`,
			}).
				WithHint("本机发布必须指定 clientId，可通过配置或命令参数提供。").
				WithNextCommand("yxer config set-local-client-id <clientId>")
		}
	default:
		return "", "", yxerrors.Usage(`publishChannel must be "cloud" or "local"`, []string{
			fmt.Sprintf("got %q", channel),
		}).
			WithHint(`publishChannel 仅支持 "cloud" 或 "local"。`)
	}
	return channel, clientID, nil
}

func shouldOfferLocalPublishRetry(err error, channel string) bool {
	if strings.TrimSpace(channel) != "cloud" {
		return false
	}
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "账号代理不存在") || strings.Contains(strings.ToLower(message), "proxy")
}

func confirmLocalPublishRetry(platform string) (bool, error) {
	if PromptOutput != nil {
		if _, err := fmt.Fprintf(PromptOutput, "%s 账号未设置代理，是否改为走本机发布？[y/N]: ", platform); err != nil {
			return false, err
		}
	}
	reader := bufio.NewReader(PromptInput)
	answer, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

func copyOptionalPublishFields(dst, src map[string]interface{}) {
	for _, field := range []string{"cover", "coverKey", "taskSetId", "desc", "isDraft", "isAppContent"} {
		if value, ok := src[field]; ok {
			dst[field] = value
		}
	}
	if channel, ok := src["publishChannel"]; ok {
		dst["publishChannel"] = channel
	}
	if clientID, ok := src["clientId"]; ok {
		dst["clientId"] = clientID
	}
}

func SplitPlatforms(value string) []string {
	var platforms []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			platforms = append(platforms, platformutil.ChineseName(item))
		}
	}
	return platforms
}

func SinglePlatform(value string) (string, error) {
	platforms := SplitPlatforms(value)
	if len(platforms) != 1 {
		return "", yxerrors.Usage("publish supports exactly one platform per command", []string{
			`Use Agent orchestration for multi-platform publishing: run "yxer accounts", "yxer schema get", "yxer validate", and "yxer publish" once per platform.`,
			`Example: yxer publish imageText xhs xhs-payload.json; then yxer publish imageText zhihu zhihu-payload.json`,
		}).
			WithHint("单次 publish 命令只支持一个平台，请拆分成多次调用。")
	}
	return platforms[0], nil
}

func AssertAccountsOnline(apiClient *client.Client, platforms []string, accountIDs []string) error {
	wanted := map[string]bool{}
	for _, id := range accountIDs {
		wanted[id] = true
	}
	found := map[string]map[string]interface{}{}
	for _, platform := range platforms {
		accounts, err := apiClient.Accounts(platform)
		if err != nil {
			return err
		}
		for _, account := range accounts {
			id := client.AccountID(account)
			if wanted[id] {
				found[id] = account
			}
		}
	}
	var errors []string
	for id := range wanted {
		account, ok := found[id]
		if !ok {
			errors = append(errors, "account "+id+": not found in target platform account list")
			continue
		}
		if status := client.AccountStatus(account); status != 1 {
			errors = append(errors, fmt.Sprintf("account %s: status=%d; publish requires status=1", id, status))
		}
	}
	if len(errors) > 0 {
		return yxerrors.Usage("Account preflight failed", errors).
			WithHint("请先运行账号查询，确认目标账号存在且状态为在线。").
			WithNextCommand("yxer accounts <platform>")
	}
	return nil
}

func resolveAccountID(apiClient *client.Client, platform, selector string) (string, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return "", yxerrors.Usage("account is required in flags mode", nil).
			WithHint("请传入 --account，值可以是账号 ID、名称或昵称。").
			WithNextCommand(fmt.Sprintf("yxer accounts %s", platform))
	}
	accounts, err := apiClient.Accounts(platform)
	if err != nil {
		return "", err
	}
	var exact []map[string]interface{}
	var fuzzy []map[string]interface{}
	for _, account := range accounts {
		id := strings.TrimSpace(client.AccountID(account))
		name := accountDisplayName(account)
		if id == selector || name == selector {
			exact = append(exact, account)
			continue
		}
		if strings.Contains(name, selector) {
			fuzzy = append(fuzzy, account)
		}
	}
	candidates := exact
	if len(candidates) == 0 {
		candidates = fuzzy
	}
	if len(candidates) == 1 {
		if client.AccountStatus(candidates[0]) != 1 {
			return "", yxerrors.Usage("selected account is offline", selector).
				WithHint("目标账号当前不在线，请先检查账号状态。").
				WithNextCommand(fmt.Sprintf("yxer accounts %s", platform))
		}
		return client.AccountID(candidates[0]), nil
	}
	if len(candidates) == 0 {
		return "", yxerrors.Usage("account not found", selector).
			WithHint("未找到匹配账号，请先查询平台账号列表。").
			WithNextCommand(fmt.Sprintf("yxer accounts %s", platform))
	}
	return "", yxerrors.Usage("account selector is ambiguous", selector).
		WithHint("匹配到多个账号，请改用更精确的名称或直接传账号 ID。").
		WithNextCommand(fmt.Sprintf("yxer accounts %s --name %s", platform, selector))
}

func accountDisplayName(account map[string]interface{}) string {
	for _, key := range []string{"platformAccountName", "name", "nickname", "remark", "platformAuthorId"} {
		if value := strings.TrimSpace(fmt.Sprint(account[key])); value != "" && value != "<nil>" {
			return value
		}
	}
	return ""
}
