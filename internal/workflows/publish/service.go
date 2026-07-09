package publish

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	accountsflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/accounts"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type ExecuteInput struct {
	PublishType        string
	PlatformInput      string
	Payload            map[string]interface{}
	PositionalClientID string
	FlagChannel        string
	FlagClientID       string
	AutoFallbackLocal  bool
}

type Service struct {
	rt *app.Runtime
}

type InferredField struct {
	Value      interface{} `json:"value"`
	SourcePath string      `json:"sourcePath"`
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func topicHTMLPolicyForPlatforms(validator schema.Validator, platforms []string, publishType string) publishmod.TopicHTMLPolicy {
	policy := publishmod.TopicHTMLPolicy{}
	for _, platform := range platforms {
		doc, err := validator.Schema(platform, publishType)
		if err != nil {
			continue
		}
		for key, fields := range publishmod.TopicHTMLPolicyFromSchema(platform, doc.Properties) {
			policy[key] = fields
		}
	}
	if len(policy) == 0 {
		return nil
	}
	return policy
}

func (s Service) Execute(input ExecuteInput) (map[string]interface{}, error) {
	input.PublishType = publishmod.NormalizePublishType(input.PublishType)
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil {
		return nil, err
	}
	platforms := []string{platform}
	cfg := s.rt.Config
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
	if err := publishmod.RequireStandardPayload(resolvedPayload); err != nil {
		return nil, err
	}
	if err := publishmod.ResolveStandardPayloadResourceMetadata(resolvedPayload); err != nil {
		return nil, err
	}
	validator := schema.NewValidator(cfg.SchemaDir)
	topicPolicy := topicHTMLPolicyForPlatforms(validator, platforms, input.PublishType)
	publishArgs := publishmod.NormalizeStandardPayloadForSchemaValidation(input.PublishType, platforms, resolvedPayload)

	for _, platform := range platforms {
		result, err := validator.ValidateStrict(platform, input.PublishType, resolvedPayload)
		if err != nil {
			return nil, schemaUnavailableError(platform, input.PublishType, cfg.SchemaDir, err)
		}
		if !result.Valid {
			return nil, yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint(schemaValidationHint(result.Errors)).
				WithNextCommand(fmt.Sprintf("yxer validate %s %s <payload.json>", platform, input.PublishType))
		}
	}
	preflight := publishmod.PreflightWithTopicHTMLPolicy(input.PublishType, platforms, resolvedPayload, topicPolicy)
	if len(preflight.Errors) > 0 {
		return nil, yxerrors.Usage("Publish preflight failed", preflight.Errors).
			WithHint("请先完成资源上传、账号校验，并确保发布参数中不包含外部 URL。")
	}

	apiClient := s.rt.Client
	accountsByID, err := ResolveTargetAccounts(apiClient, platforms, preflight.AccountIDs)
	if err != nil {
		return nil, err
	}
	if err := AssertCloudChannelReady(channel, platforms, accountsByID); err != nil {
		return nil, err
	}

	body := BuildPublishBody(resolvedPayload, publishArgs, input.PublishType, platforms, channel, clientID)
	result, err := apiClient.Publish(body)
	if err == nil {
		return result, nil
	}
	if !shouldOfferLocalPublishRetry(err, channel) {
		return nil, err
	}
	if !input.AutoFallbackLocal {
		return nil, buildLocalFallbackError(platform, input.PublishType, clientID, err)
	}
	localChannel, localClientID, resolveErr := ResolvePublishMode(cfg, resolvedPayload, "", "local", "")
	if resolveErr != nil {
		return nil, resolveErr
	}
	resolvedPayload["publishChannel"] = localChannel
	resolvedPayload["clientId"] = localClientID
	body = BuildPublishBody(resolvedPayload, publishArgs, input.PublishType, platforms, localChannel, localClientID)
	return apiClient.Publish(body)
}

func (s Service) ExecuteEnvelope(input ExecuteInput) (EnvelopeResult, error) {
	result, err := s.Execute(input)
	return s.wrapExecuteEnvelope(result, err)
}

func (s Service) wrapExecuteEnvelope(result map[string]interface{}, err error) (EnvelopeResult, error) {
	if err != nil {
		return EnvelopeResult{}, err
	}
	return EnvelopeResult{
		Action: "publish",
		Data:   result,
	}, nil
}

func schemaUnavailableError(platform, publishType, schemaDir string, err error) error {
	return yxerrors.Usage("Schema file is required for publish validation", map[string]interface{}{
		"platform":    platform,
		"publishType": publishType,
		"schemaDir":   schemaDir,
		"cause":       err.Error(),
	}).WithHint("请确认 YIXIAOER_PROJECT_DIR 指向包含 schemas/platforms 的项目根目录，且目标平台和类型的 schema 文件存在。").
		WithNextCommand("yxer schema list")
}

func SchemaUnavailableForCommand(platform, publishType, schemaDir string, err error) error {
	return schemaUnavailableError(platform, publishType, schemaDir, err)
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

func BuildPublishBody(payload, publishArgs map[string]interface{}, publishType string, platforms []string, channel, clientID string) map[string]interface{} {
	body, _ := BuildPublishBodyWithInferred(payload, publishArgs, publishType, platforms, channel, clientID)
	return body
}

func BuildPublishBodyWithInferred(payload, publishArgs map[string]interface{}, publishType string, platforms []string, channel, clientID string) (map[string]interface{}, map[string]InferredField) {
	body := map[string]interface{}{
		"publishType":    publishType,
		"platforms":      platforms,
		"publishArgs":    publishArgs,
		"publishChannel": channel,
	}
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
	applyPublishMode(body, channel, clientID)
	stripArticleContentFromForms(body)
	inferred := normalizePublishEnvelope(body, publishArgs, publishType)
	return body, inferred
}

func applyPublishMode(body map[string]interface{}, channel, clientID string) {
	if channel == "" {
		channel = "cloud"
	}
	body["publishChannel"] = channel
	if channel == "local" && clientID != "" {
		body["clientId"] = clientID
		return
	}
	delete(body, "clientId")
}

func normalizePublishEnvelope(body, publishArgs map[string]interface{}, publishType string) map[string]InferredField {
	inferred := map[string]InferredField{}
	if body == nil {
		return inferred
	}
	if publishArgs == nil {
		publishArgs = map[string]interface{}{}
	}
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	firstForm := firstObject(accountForms)
	firstCPF := objectField(firstForm, "contentPublishForm")
	weixinPlatformForm := weixinAccountArticlePlatformForm(publishArgs)
	firstWeixinArticle := firstWeixinAccountArticle(weixinPlatformForm)

	if _, ok := body["cover"]; !ok {
		candidates := []fieldCandidate{
			{value: publishArgs["cover"], sourcePath: "publishArgs.cover"},
			{value: firstForm["cover"], sourcePath: "publishArgs.accountForms[0].cover"},
			{value: firstCPF["cover"], sourcePath: "publishArgs.accountForms[0].contentPublishForm.cover"},
			{value: firstWeixinArticle["cover"], sourcePath: "publishArgs.platformForms[微信公众号].articles[0].cover"},
		}
		if cover, ok := firstNonNilCandidate(candidates); ok {
			body["cover"] = cover.value
			inferred["cover"] = InferredField{Value: cover.value, SourcePath: cover.sourcePath}
		}
	}
	if stringField(body, "coverKey") == "" {
		candidates := []fieldCandidate{
			{value: stringField(publishArgs, "coverKey"), sourcePath: "publishArgs.coverKey"},
			{value: stringField(firstForm, "coverKey"), sourcePath: "publishArgs.accountForms[0].coverKey"},
			{value: stringField(firstCPF, "coverKey"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.coverKey"},
			{value: stringField(firstWeixinArticle, "coverKey"), sourcePath: "publishArgs.platformForms[微信公众号].articles[0].coverKey"},
			{value: stringField(objectField(firstWeixinArticle, "cover"), "key"), sourcePath: "publishArgs.platformForms[微信公众号].articles[0].cover.key"},
			{value: stringField(objectField(body, "cover"), "key"), sourcePath: "cover.key"},
		}
		if coverKey, ok := firstNonEmptyStringCandidate(candidates); ok {
			body["coverKey"] = coverKey.value
			inferred["coverKey"] = InferredField{Value: coverKey.value, SourcePath: coverKey.sourcePath}
		}
	}
	if stringField(body, "desc") == "" {
		if desc, ok := inferOuterDesc(publishType, publishArgs, firstCPF, firstWeixinArticle); ok {
			body["desc"] = desc.value
			inferred["desc"] = InferredField{Value: desc.value, SourcePath: desc.sourcePath}
		}
	}
	if _, ok := body["isDraft"]; !ok {
		body["isDraft"] = inferYixiaoerDraft(body)
		inferred["isDraft"] = InferredField{Value: body["isDraft"], SourcePath: "default"}
	}
	if _, ok := body["isAppContent"]; !ok {
		body["isAppContent"] = false
		inferred["isAppContent"] = InferredField{Value: false, SourcePath: "default"}
	}
	return inferred
}

func inferYixiaoerDraft(body map[string]interface{}) bool {
	if body == nil {
		return false
	}
	if value, ok := body["isDraft"]; ok {
		if draft, ok := value.(bool); ok {
			return draft
		}
	}
	return false
}

func inferOuterDesc(publishType string, publishArgs, contentPublishForm, weixinArticle map[string]interface{}) (fieldCandidate, bool) {
	switch publishmod.NormalizePublishType(publishType) {
	case "article":
		return firstNonEmptyStringCandidate([]fieldCandidate{
			{value: stringField(weixinArticle, "digest"), sourcePath: "publishArgs.platformForms[微信公众号].articles[0].digest"},
			{value: stringField(weixinArticle, "title"), sourcePath: "publishArgs.platformForms[微信公众号].articles[0].title"},
			{value: stringField(weixinArticle, "content"), sourcePath: "publishArgs.platformForms[微信公众号].articles[0].content"},
			{value: stringField(contentPublishForm, "desc"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.desc"},
			{value: stringField(contentPublishForm, "title"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.title"},
			{value: stringField(publishArgs, "content"), sourcePath: "publishArgs.content"},
			{value: stringField(contentPublishForm, "content"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.content"},
		})
	case "video", "imageText":
		return firstNonEmptyStringCandidate([]fieldCandidate{
			{value: stringField(contentPublishForm, "description"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.description"},
			{value: stringField(contentPublishForm, "title"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.title"},
			{value: stringField(publishArgs, "content"), sourcePath: "publishArgs.content"},
			{value: stringField(contentPublishForm, "content"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.content"},
		})
	default:
		return firstNonEmptyStringCandidate([]fieldCandidate{
			{value: stringField(contentPublishForm, "description"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.description"},
			{value: stringField(contentPublishForm, "title"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.title"},
			{value: stringField(publishArgs, "content"), sourcePath: "publishArgs.content"},
			{value: stringField(contentPublishForm, "content"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.content"},
		})
	}
}

func stripArticleContentFromForms(body map[string]interface{}) {
	if publishmod.NormalizePublishType(stringField(body, "publishType")) != "article" {
		return
	}
	publishArgs, _ := body["publishArgs"].(map[string]interface{})
	if publishArgs == nil {
		return
	}
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	for _, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}
		delete(cpf, "content")
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

func weixinAccountArticlePlatformForm(publishArgs map[string]interface{}) map[string]interface{} {
	if publishArgs == nil {
		return nil
	}
	platformForms, _ := publishArgs["platformForms"].(map[string]interface{})
	if platformForms == nil {
		return nil
	}
	for _, key := range []string{"微信公众号", "weixin.account"} {
		form, _ := platformForms[key].(map[string]interface{})
		if form != nil {
			return form
		}
	}
	return nil
}

func firstWeixinAccountArticle(platformForm map[string]interface{}) map[string]interface{} {
	if platformForm == nil {
		return nil
	}
	articles, _ := platformForm["articles"].([]interface{})
	for _, item := range articles {
		if article, ok := item.(map[string]interface{}); ok {
			return article
		}
	}
	return nil
}

func objectField(obj map[string]interface{}, key string) map[string]interface{} {
	if obj == nil {
		return nil
	}
	value, _ := obj[key].(map[string]interface{})
	return value
}

func stringField(obj map[string]interface{}, key string) string {
	if obj == nil {
		return ""
	}
	value := obj[key]
	if value == nil {
		return ""
	}
	text := strings.TrimSpace(fmt.Sprint(value))
	if text == "<nil>" {
		return ""
	}
	return text
}

func firstNonNil(values ...interface{}) interface{} {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

type fieldCandidate struct {
	value      interface{}
	sourcePath string
}

func firstNonNilCandidate(values []fieldCandidate) (fieldCandidate, bool) {
	for _, value := range values {
		if value.value != nil {
			return value, true
		}
	}
	return fieldCandidate{}, false
}

func firstNonEmptyStringCandidate(values []fieldCandidate) (fieldCandidate, bool) {
	for _, value := range values {
		text := strings.TrimSpace(fmt.Sprint(value.value))
		if text != "" && text != "<nil>" {
			value.value = text
			return value, true
		}
	}
	return fieldCandidate{}, false
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

func payloadWithPublishMode(payload map[string]interface{}, channel, clientID string) map[string]interface{} {
	withMode := cloneMap(payload)
	if channel == "" {
		channel = "cloud"
	}
	withMode["publishChannel"] = channel
	if clientID != "" {
		withMode["clientId"] = clientID
	} else {
		delete(withMode, "clientId")
	}
	return withMode
}

func ResolvePublishMode(cfg config.Config, payload map[string]interface{}, positionalClientID, flagChannel, flagClientID string) (string, string, error) {
	channel := "cloud"
	clientID := ""
	payloadChannel := ""
	if value, ok := payload["publishChannel"]; ok {
		payloadChannel = strings.TrimSpace(fmt.Sprint(value))
		channel = payloadChannel
	}
	if value, ok := payload["clientId"]; ok {
		clientID = strings.TrimSpace(fmt.Sprint(value))
	}
	if strings.TrimSpace(positionalClientID) != "" {
		if strings.TrimSpace(flagChannel) != "local" && payloadChannel != "local" {
			return "", "", yxerrors.Usage("positional clientId requires local publish channel", []string{
				`The fourth positional clientId is deprecated and no longer switches publishChannel implicitly.`,
				`Use flags: yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>`,
			}).
				WithHint("请显式传入 --publish-channel local --client-id <clientId>，避免误把云发布切成本机发布。").
				WithNextCommand("yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>")
		}
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
	var typed *yxerrors.Error
	if errors.As(err, &typed) {
		if strings.EqualFold(remoteErrorCode(typed), "PROXY_NOT_CONFIGURED") {
			return true
		}
	}
	message := err.Error()
	return strings.Contains(message, "账号代理不存在") || strings.Contains(strings.ToLower(message), "proxy")
}

func remoteErrorCode(err *yxerrors.Error) string {
	if err == nil {
		return ""
	}
	details, _ := err.Details.(map[string]interface{})
	if details == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(details["code"]))
}

func buildLocalFallbackError(platform, publishType, clientID string, cause error) error {
	nextCommand := fmt.Sprintf("yxer publish %s %s <payload.json> --publish-channel local --client-id <clientId>", publishType, platform)
	if strings.TrimSpace(clientID) != "" {
		nextCommand = fmt.Sprintf("yxer publish %s %s <payload.json> --publish-channel local --client-id %s", publishType, platform, clientID)
	}
	return yxerrors.Remote("cloud publish failed; local publish fallback is available", cause.Error()).
		WithCategory("publish_channel_fallback").
		WithHint("当前账号云发布失败，可改用本机发布；如需自动回退，请显式传入 --auto-fallback-local。").
		WithNextCommand(nextCommand)
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
			`Use Agent orchestration for multi-platform publishing: run "yxer accounts", "yxer schema fields", "yxer validate", and "yxer publish" once per platform; add "yxer schema get" only when you need the full skeleton.`,
			`Example: yxer publish imageText xhs xhs-payload.json; then yxer publish imageText zhihu zhihu-payload.json`,
		}).
			WithHint("单次 publish 命令只支持一个平台，请拆分成多次调用。")
	}
	return platforms[0], nil
}

func AssertAccountsOnline(apiClient *api.Client, platforms []string, accountIDs []string) error {
	_, err := ResolveTargetAccounts(apiClient, platforms, accountIDs)
	return err
}

func ResolveTargetAccounts(apiClient *api.Client, platforms []string, accountIDs []string) (map[string]map[string]interface{}, error) {
	wanted := map[string]bool{}
	for _, id := range accountIDs {
		wanted[id] = true
	}
	found := map[string]map[string]interface{}{}
	for _, platform := range platforms {
		accounts, err := apiClient.Accounts(platform)
		if err != nil {
			return nil, err
		}
		for _, account := range accounts {
			id := api.AccountID(account)
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
		if status := api.AccountStatus(account); status != 1 {
			errors = append(errors, fmt.Sprintf("account %s: status=%d; publish requires status=1", id, status))
		}
	}
	if len(errors) > 0 {
		return nil, yxerrors.Usage("Account preflight failed", errors).
			WithHint("请先运行账号查询，确认目标账号存在且状态为在线。").
			WithNextCommand("yxer accounts <platform>")
	}
	return found, nil
}

func AssertCloudChannelReady(channel string, platforms []string, accountsByID map[string]map[string]interface{}) error {
	if strings.TrimSpace(channel) != "cloud" {
		return nil
	}
	if !requiresCloudProxy(platforms) {
		return nil
	}
	var missing []string
	for _, account := range accountsByID {
		if !accountHasCloudProxy(account) {
			missing = append(missing, accountsflow.AccountName(account))
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return yxerrors.Usage("Cloud publish preflight failed", []string{
		"以下账号未设置代理，请先设置代理：" + strings.Join(missing, "、"),
	}).
		WithHint("当前账号云发布缺少代理配置，建议先配置代理，或改用本机发布。").
		WithNextCommand("yxer publish <type> <platform> <payload.json> --publish-channel local --client-id <clientId>")
}

func requiresCloudProxy(platforms []string) bool {
	for _, platform := range platforms {
		switch platformutil.ChineseName(platform) {
		case "视频号":
			return true
		}
	}
	return false
}

func accountHasCloudProxy(account map[string]interface{}) bool {
	if account == nil {
		return false
	}
	if proxyInfo, ok := account["proxyInfo"].(map[string]interface{}); ok && len(proxyInfo) > 0 {
		return true
	}
	for _, key := range []string{"proxyId", "kuaidailiArea"} {
		if value := strings.TrimSpace(fmt.Sprint(account[key])); value != "" && value != "<nil>" {
			return true
		}
	}
	return false
}
