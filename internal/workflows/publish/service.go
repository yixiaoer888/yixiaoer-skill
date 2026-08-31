package publish

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

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
	PublishType                 string
	PlatformInput               string
	PayloadPath                 string
	Payload                     map[string]interface{}
	ContentFile                 string
	PositionalClientID          string
	FlagChannel                 string
	FlagClientID                string
	AutoFallbackLocal           bool
	ContinueOnContentImageError bool
}

type Service struct {
	rt *app.Runtime
}

type RemoteCheckMode string

const (
	RemoteChecksNone         RemoteCheckMode = "none"
	RemoteChecksRequired     RemoteCheckMode = "required"
	RemoteChecksCloudWithKey RemoteCheckMode = "cloud_with_api_key"
)

type PrepareOptions struct {
	TraceNormalizations bool
	RemoteChecks        RemoteCheckMode
}

type PreparedPublish struct {
	Platform          string
	Platforms         []string
	PublishType       string
	Payload           map[string]interface{}
	PublishArgs       map[string]interface{}
	PublishMode       string
	PublishModeSource string
	ClientID          string
	ClientIDSource    string
	Preflight         publishmod.PreflightResult
	Normalizations    []publishmod.NormalizationEvent
	PublishBody       map[string]interface{}
	ContentBaseDir    string
	InferredFields    map[string]InferredField
	RemoteChecked     bool
}

type InferredField struct {
	Value      interface{} `json:"value"`
	SourcePath string      `json:"sourcePath"`
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func (s Service) DeletePublishedTask(taskID string) (map[string]interface{}, error) {
	return s.rt.Client.DeletePublishedTask(taskID)
}

func (s Service) Prepare(input ExecuteInput, opts PrepareOptions) (PreparedPublish, error) {
	input.PublishType = publishmod.NormalizePublishType(input.PublishType)
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil {
		return PreparedPublish{}, err
	}
	platforms := []string{platform}
	cfg := s.rt.Config
	resolvedPayload := cloneMap(input.Payload)
	mode, err := ResolvePublishModeDetailed(cfg, resolvedPayload, input.PositionalClientID, input.FlagChannel, input.FlagClientID)
	if err != nil {
		return PreparedPublish{}, err
	}
	channel, clientID := mode.Channel, mode.ClientID
	resolvedPayload["publishChannel"] = channel
	if clientID != "" {
		resolvedPayload["clientId"] = clientID
	} else {
		delete(resolvedPayload, "clientId")
	}
	if err := publishmod.RequireStandardPayload(resolvedPayload); err != nil {
		return PreparedPublish{}, err
	}
	contentBaseDir := ""
	if contentFile := strings.TrimSpace(input.ContentFile); contentFile != "" {
		if input.PublishType != "article" {
			return PreparedPublish{}, yxerrors.Usage("content file is only supported for article publish", map[string]interface{}{
				"publishType": input.PublishType,
				"path":        contentFile,
			}).WithCategory("article_content_source").
				WithHint("--content-file 仅用于 article 发布；视频和图文请继续使用 payload.json 中的资源字段。")
		}
		source, err := loadArticleContentSource(contentFile)
		if err != nil {
			return PreparedPublish{}, err
		}
		publishArgs := objectField(resolvedPayload, "publishArgs")
		publishArgs["content"] = source.Content
		contentBaseDir = source.BaseDir
	}
	if err := publishmod.ResolveStandardPayloadResourceMetadata(resolvedPayload); err != nil {
		return PreparedPublish{}, err
	}
	if platformutil.CanonicalKey(platform) == "souhuhao" && input.PublishType == "video" {
		resolvedPayload, err = s.rt.Client.NormalizeSohuhaoVideoPayloadForPlatform(resolvedPayload, platform)
		if err != nil {
			return PreparedPublish{}, err
		}
	}
	validator := schema.NewValidator(cfg.SchemaDir)
	topicPolicy := topicHTMLPolicyForPlatforms(validator, platforms, input.PublishType)
	var normalizations []publishmod.NormalizationEvent
	publishArgs := publishmod.NormalizeStandardPayloadForSchemaValidationWithTrace(input.PublishType, platforms, resolvedPayload, &normalizations)

	for _, platform := range platforms {
		result, err := validator.ValidateStrict(platform, input.PublishType, resolvedPayload)
		if err != nil {
			return PreparedPublish{}, schemaUnavailableError(platform, input.PublishType, cfg.SchemaDir, err)
		}
		if !result.Valid {
			return PreparedPublish{}, yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint(schemaValidationHintForPlatform(platform, result.Errors)).
				WithNextCommand(fmt.Sprintf("yxer validate %s %s <payload.json>", platform, input.PublishType))
		}
	}
	preflight := publishmod.PreflightWithTopicHTMLPolicyAndTrace(input.PublishType, platforms, payloadWithPublishMode(resolvedPayload, channel, clientID), topicPolicy, &normalizations)
	if len(preflight.Errors) > 0 {
		return PreparedPublish{}, yxerrors.Usage("Publish preflight failed", preflight.Errors).
			WithHint("请先完成资源上传、账号校验，并确保发布参数中不包含外部 URL。")
	}

	remoteChecked := false
	if shouldPrepareRemoteCheck(opts.RemoteChecks, channel, cfg, preflight.AccountIDs) {
		accountsByID, err := ResolveTargetAccounts(s.rt.Client, platforms, preflight.AccountIDs)
		if err != nil {
			return PreparedPublish{}, err
		}
		if err := AssertCloudChannelReady(channel, platforms, accountsByID); err != nil {
			return PreparedPublish{}, err
		}
		remoteChecked = true
	}

	body, inferredFields := BuildPublishBodyWithInferred(resolvedPayload, publishArgs, input.PublishType, platforms, channel, clientID)
	if err := validateInstagramMediaKeys(platform, input.PublishType, body); err != nil {
		return PreparedPublish{}, err
	}

	if !opts.TraceNormalizations {
		normalizations = nil
	}
	return PreparedPublish{
		Platform:          platform,
		Platforms:         platforms,
		PublishType:       input.PublishType,
		Payload:           resolvedPayload,
		PublishArgs:       publishArgs,
		PublishMode:       channel,
		PublishModeSource: mode.ChannelSource,
		ClientID:          clientID,
		ClientIDSource:    mode.ClientIDSource,
		Preflight:         preflight,
		Normalizations:    normalizations,
		PublishBody:       body,
		ContentBaseDir:    contentBaseDir,
		InferredFields:    inferredFields,
		RemoteChecked:     remoteChecked,
	}, nil
}

func shouldPrepareRemoteCheck(mode RemoteCheckMode, channel string, cfg config.Config, accountIDs []string) bool {
	if len(accountIDs) == 0 {
		return false
	}
	switch mode {
	case RemoteChecksRequired:
		return true
	case RemoteChecksCloudWithKey:
		return channel == "cloud" && cfg.APIKey != ""
	default:
		return false
	}
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
	apiClient := s.rt.Client
	coverCompressionEvents := []CoverCompressionEvent{}
	if payload, events, err := materializeShipinhaoCoverCompression(apiClient, input); err != nil {
		return nil, err
	} else if len(events) > 0 {
		input.Payload = payload
		coverCompressionEvents = events
	}
	prepared, err := s.Prepare(input, PrepareOptions{RemoteChecks: RemoteChecksRequired})
	if err != nil {
		return nil, err
	}
	if err := AssertShoppingCartEntitlements(apiClient, prepared.Payload, prepared.Platforms...); err != nil {
		return nil, err
	}
	cfg := s.rt.Config
	if events, err := materializeArticleContentImages(apiClient, prepared.PublishBody, prepared.ContentBaseDir, input.ContinueOnContentImageError); err != nil {
		return nil, contentImageMaterializationPromptError(input, prepared, events, err)
	}
	result, err := apiClient.Publish(prepared.PublishBody)
	if err == nil {
		attachCoverCompressionEvents(result, coverCompressionEvents)
		return result, nil
	}
	if mapped := mapInstagramMediaFetchError(prepared.Platform, prepared.PublishType, err); mapped != nil {
		return nil, mapped
	}
	if !shouldOfferLocalPublishRetry(err, prepared.PublishMode) {
		return nil, err
	}
	if !input.AutoFallbackLocal {
		return nil, buildLocalFallbackError(prepared.Platform, prepared.PublishType, prepared.ClientID, err)
	}
	localChannel, localClientID, resolveErr := ResolvePublishMode(cfg, prepared.Payload, "", "local", "")
	if resolveErr != nil {
		return nil, resolveErr
	}
	prepared.Payload["publishChannel"] = localChannel
	prepared.Payload["clientId"] = localClientID
	body := BuildPublishBody(prepared.Payload, prepared.PublishArgs, prepared.PublishType, prepared.Platforms, localChannel, localClientID)
	if events, err := materializeArticleContentImages(apiClient, body, prepared.ContentBaseDir, input.ContinueOnContentImageError); err != nil {
		return nil, contentImageMaterializationPromptError(input, prepared, events, err)
	}
	result, err = apiClient.Publish(body)
	if err == nil {
		attachCoverCompressionEvents(result, coverCompressionEvents)
	}
	return result, err
}

func attachCoverCompressionEvents(result map[string]interface{}, events []CoverCompressionEvent) {
	if len(events) == 0 || result == nil {
		return
	}
	result["mediaProcessing"] = map[string]interface{}{
		"coverCompression": events,
	}
}

func contentImageMaterializationPromptError(input ExecuteInput, prepared PreparedPublish, events []ArticleContentImageMaterialization, cause error) error {
	nextCommand := continueContentImagePublishCommand(input, prepared)
	return yxerrors.Remote("article content image materialization failed; confirmation is required before publishing", map[string]interface{}{
		"events": events,
		"cause":  cause.Error(),
	}).WithCategory("article_content_image_materialization_confirmation").
		WithHint("文章正文中的部分图片无法转存为稳定地址。CLI 已中止发布，避免把第三方可能不可访问的图片地址发出去；如确认可以保留原图地址继续发布，请重新执行 nextCommand。").
		WithNextCommand(nextCommand).
		WithRetryable(true)
}

func continueContentImagePublishCommand(input ExecuteInput, prepared PreparedPublish) string {
	publishType := firstNonEmptyString(input.PublishType, prepared.PublishType, "<type>")
	platform := firstNonEmptyString(input.PlatformInput, prepared.Platform, "<platform>")
	payloadPath := firstNonEmptyString(input.PayloadPath, "<payload.json>")
	parts := []string{"yxer", "publish", publishType, platform, payloadPath}
	if contentFile := strings.TrimSpace(input.ContentFile); contentFile != "" {
		parts = append(parts, "--content-file", quotePublishCommandArg(contentFile))
	}
	if input.FlagChannel != "" {
		parts = append(parts, "--publish-channel", input.FlagChannel)
	}
	if input.FlagClientID != "" {
		parts = append(parts, "--client-id", input.FlagClientID)
	}
	if input.AutoFallbackLocal {
		parts = append(parts, "--auto-fallback-local")
	}
	parts = append(parts, "--continue-on-content-image-error")
	return strings.Join(parts, " ")
}

func quotePublishCommandArg(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	if strings.ContainsAny(value, " \t\"'") {
		return strconv.Quote(value)
	}
	return value
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
		dst[key] = clonePublishValue(value)
	}
	return dst
}

func clonePublishValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			out[key] = clonePublishValue(item)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(typed))
		for index, item := range typed {
			out[index] = clonePublishValue(item)
		}
		return out
	case []map[string]interface{}:
		out := make([]map[string]interface{}, len(typed))
		for index, item := range typed {
			out[index], _ = clonePublishValue(item).(map[string]interface{})
		}
		return out
	default:
		return value
	}
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
	stripArticleContentFromForms(body, platforms)
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
			{value: stringField(contentPublishForm, "desc"), sourcePath: "publishArgs.accountForms[0].contentPublishForm.desc"},
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

func stripArticleContentFromForms(body map[string]interface{}, platforms []string) {
	if publishmod.NormalizePublishType(stringField(body, "publishType")) != "article" {
		return
	}
	if articleKeepsContentInForm(platforms) {
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

func articleKeepsContentInForm(platforms []string) bool {
	for _, platform := range platforms {
		if platformutil.CanonicalKey(platform) == "jianshu" {
			return true
		}
	}
	return false
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
	// resolvedPayload is already a deep clone of the caller's payload. Keep the
	// nested publish structures shared here so preflight normalizations (for
	// example topic HTML) are reflected in the final request body, while mode
	// fields remain isolated at the outer envelope.
	withMode := make(map[string]interface{}, len(payload)+2)
	for key, value := range payload {
		withMode[key] = value
	}
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

type PublishModeResolution struct {
	Channel        string
	ClientID       string
	ChannelSource  string
	ClientIDSource string
}

func ResolvePublishMode(cfg config.Config, payload map[string]interface{}, positionalClientID, flagChannel, flagClientID string) (string, string, error) {
	resolution, err := ResolvePublishModeDetailed(cfg, payload, positionalClientID, flagChannel, flagClientID)
	if err != nil {
		return "", "", err
	}
	return resolution.Channel, resolution.ClientID, nil
}

func ResolvePublishModeDetailed(cfg config.Config, payload map[string]interface{}, positionalClientID, flagChannel, flagClientID string) (PublishModeResolution, error) {
	channel := "cloud"
	channelSource := "default"
	clientID := ""
	clientIDSource := "none"
	payloadChannel := ""
	if value, ok := payload["publishChannel"]; ok {
		payloadChannel = strings.TrimSpace(fmt.Sprint(value))
		channel = payloadChannel
		channelSource = "payload"
	}
	if value, ok := payload["clientId"]; ok {
		clientID = strings.TrimSpace(fmt.Sprint(value))
		if clientID != "" {
			clientIDSource = "payload"
		}
	}
	if strings.TrimSpace(positionalClientID) != "" {
		if strings.TrimSpace(flagChannel) != "local" && payloadChannel != "local" {
			return PublishModeResolution{}, yxerrors.Usage("positional clientId requires local publish channel", []string{
				`The fourth positional clientId is deprecated and no longer switches publishChannel implicitly.`,
				`Use flags: yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>`,
			}).
				WithHint("请显式传入 --publish-channel local --client-id <clientId>，避免误把云发布切成本机发布。").
				WithNextCommand("yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>")
		}
		channel = "local"
		channelSource = "positional"
		clientID = strings.TrimSpace(positionalClientID)
		clientIDSource = "positional"
	}
	if strings.TrimSpace(flagChannel) != "" {
		channel = strings.TrimSpace(flagChannel)
		channelSource = "flag"
	}
	if strings.TrimSpace(flagClientID) != "" {
		clientID = strings.TrimSpace(flagClientID)
		clientIDSource = "flag"
	}
	switch channel {
	case "", "cloud":
		channel = "cloud"
		clientID = ""
		clientIDSource = "none"
	case "local":
		if clientID == "" {
			clientID = strings.TrimSpace(cfg.LocalClientID)
			if clientID != "" {
				clientIDSource = "config"
			}
		}
		if clientID == "" {
			return PublishModeResolution{}, yxerrors.Usage(`clientId is required when publishChannel is "local"`, []string{
				`Run: yxer config set-local-client-id <clientId>`,
				`Or pass flags: yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>`,
			}).
				WithHint("本机发布必须指定 clientId，可通过配置或命令参数提供。").
				WithNextCommand("yxer config set-local-client-id <clientId>")
		}
	default:
		return PublishModeResolution{}, yxerrors.Usage(`publishChannel must be "cloud" or "local"`, []string{
			fmt.Sprintf("got %q", channel),
		}).
			WithHint(`publishChannel 仅支持 "cloud" 或 "local"。`)
	}
	return PublishModeResolution{
		Channel:        channel,
		ClientID:       clientID,
		ChannelSource:  channelSource,
		ClientIDSource: clientIDSource,
	}, nil
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
		WithHint("当前账号云发布失败，可改用本机发布；CLI 不会默认自动重试，如需授权自动回退，请显式传入 --auto-fallback-local。").
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
			WithNextCommand("yxer accounts list <platform> --status 1 --json")
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

func validateInstagramMediaKeys(platform, publishType string, body map[string]interface{}) error {
	if platformutil.ChineseName(platform) != "Instagram" || publishmod.NormalizePublishType(publishType) != "video" {
		return nil
	}
	accountForms, _ := objectField(body, "publishArgs")["accountForms"].([]interface{})
	for index, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		video := objectField(form, "video")
		key := stringField(video, "key")
		if key == "" || isASCIIString(key) {
			continue
		}
		return yxerrors.Usage("Instagram publish requires an ASCII-only uploaded video key", map[string]interface{}{
			"path":  fmt.Sprintf("publishArgs.accountForms[%d].video.key", index),
			"value": key,
		}).WithCategory("instagram_media_key").
			WithHint("Instagram/Meta 拉取视频时会对带中文字符的媒体 URL 返回 HTTP 400。请重新上传视频，确保原文件名使用英文、数字、连字符或下划线。").
			WithNextCommand("yxer upload --file <ascii_video_name>.mp4")
	}
	return nil
}

func mapInstagramMediaFetchError(platform, publishType string, err error) error {
	if platformutil.ChineseName(platform) != "Instagram" || publishmod.NormalizePublishType(publishType) != "video" || err == nil {
		return nil
	}
	message := strings.ToLower(err.Error())
	if !strings.Contains(message, "media could not be fetched") && !strings.Contains(message, "video download failed") {
		return nil
	}
	var typed *yxerrors.Error
	if errors.As(err, &typed) {
		return typed.WithCategory("instagram_media_fetch").
			WithHint("Instagram/Meta 在创建发布容器时无法拉取视频。若资源 URL 或 key 中包含中文字符，请将原视频改成英文文件名后重新执行 yxer upload，再更新 payload 中的 video.key。").
			WithNextCommand("yxer upload --file <ascii_video_name>.mp4")
	}
	return yxerrors.Remote("Instagram could not fetch the uploaded video media", err.Error()).
		WithCategory("instagram_media_fetch").
		WithHint("Instagram/Meta 在创建发布容器时无法拉取视频。若资源 URL 或 key 中包含中文字符，请将原视频改成英文文件名后重新执行 yxer upload，再更新 payload 中的 video.key。").
		WithNextCommand("yxer upload --file <ascii_video_name>.mp4")
}

func isASCIIString(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return true
	}
	if decoded, err := url.PathUnescape(value); err == nil {
		value = decoded
	}
	for len(value) > 0 {
		r, size := utf8.DecodeRuneInString(value)
		if r > 127 {
			return false
		}
		value = value[size:]
	}
	return true
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
