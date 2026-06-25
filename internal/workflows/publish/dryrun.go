package publish

import (
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type DryRunResult struct {
	Platform      string                 `json:"platform"`
	PublishType   string                 `json:"publishType"`
	PublishBody   map[string]interface{} `json:"request"`
	PublishArgs   map[string]interface{} `json:"publishArgs,omitempty"`
	PublishMode   string                 `json:"publishChannel"`
	ClientID      string                 `json:"clientId,omitempty"`
	AccountIDs    []string               `json:"accountIds,omitempty"`
	PlatformDraft bool                   `json:"platformDraft"`
	YixiaoerDraft bool                   `json:"yixiaoerDraft"`
	SchemaChecked bool                   `json:"schemaChecked"`
}

func (s Service) DryRunEnvelope(input ExecuteInput) (EnvelopeResult, error) {
	result, err := s.DryRun(input)
	return s.wrapDryRunEnvelope(result, err)
}

func (s Service) wrapDryRunEnvelope(result DryRunResult, err error) (EnvelopeResult, error) {
	if err != nil {
		return EnvelopeResult{}, err
	}
	return EnvelopeResult{
		Action: "publish.dry-run",
		Data: map[string]interface{}{
			"dryRun":  true,
			"request": result.PublishBody,
			"meta": map[string]interface{}{
				"platform":       result.Platform,
				"publishType":    result.PublishType,
				"publishChannel": result.PublishMode,
				"clientId":       result.ClientID,
				"accountIds":     result.AccountIDs,
				"platformDraft":  result.PlatformDraft,
				"yixiaoerDraft":  result.YixiaoerDraft,
				"schemaChecked":  result.SchemaChecked,
			},
		},
	}, nil
}

func (s Service) DryRun(input ExecuteInput) (DryRunResult, error) {
	input.PublishType = publishmod.NormalizePublishType(input.PublishType)
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil {
		return DryRunResult{}, err
	}
	platforms := []string{platform}

	cfg := s.rt.Config
	resolvedPayload := cloneMap(input.Payload)
	channel, clientID, err := ResolvePublishMode(cfg, resolvedPayload, input.PositionalClientID, input.FlagChannel, input.FlagClientID)
	if err != nil {
		return DryRunResult{}, err
	}
	resolvedPayload["publishChannel"] = channel
	if clientID != "" {
		resolvedPayload["clientId"] = clientID
	} else {
		delete(resolvedPayload, "clientId")
	}
	if err := publishmod.RequireStandardPayload(resolvedPayload); err != nil {
		return DryRunResult{}, err
	}
	if err := publishmod.ResolveStandardPayloadResourceMetadata(resolvedPayload); err != nil {
		return DryRunResult{}, err
	}
	validator := schema.NewValidator(cfg.SchemaDir)
	topicPolicy := topicHTMLPolicyForPlatforms(validator, platforms, input.PublishType)
	publishArgs := publishmod.NormalizeStandardPayloadForSchemaValidation(input.PublishType, platforms, resolvedPayload)

	for _, platform := range platforms {
		result := validator.Validate(platform, input.PublishType, resolvedPayload)
		if !result.Valid {
			return DryRunResult{}, yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint("请根据对应平台 schema 修正 payload 字段后重试。").
				WithNextCommand("yxer schema fields <platform> <type>")
		}
	}
	preflight := publishmod.PreflightWithTopicHTMLPolicy(input.PublishType, platforms, payloadWithPublishMode(resolvedPayload, channel, clientID), topicPolicy)
	if len(preflight.Errors) > 0 {
		return DryRunResult{}, yxerrors.Usage("Publish preflight failed", preflight.Errors).
			WithHint("请先完成资源上传、账号校验，并确保发布参数中不包含外部 URL。")
	}

	body := BuildPublishBody(resolvedPayload, publishArgs, input.PublishType, platforms, channel, clientID)

	return DryRunResult{
		Platform:      platform,
		PublishType:   input.PublishType,
		PublishBody:   body,
		PublishArgs:   publishArgs,
		PublishMode:   channel,
		ClientID:      clientID,
		AccountIDs:    preflight.AccountIDs,
		PlatformDraft: isPlatformDraftPublish(body),
		YixiaoerDraft: inferYixiaoerDraft(body),
		SchemaChecked: true,
	}, nil
}

func isPlatformDraftPublish(body map[string]interface{}) bool {
	publishArgs, _ := body["publishArgs"].(map[string]interface{})
	if platformForm := weixinAccountArticlePlatformForm(publishArgs); platformForm != nil {
		switch value := platformForm["pubType"].(type) {
		case float64:
			return int(value) == 0
		case int:
			return value == 0
		}
	}
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	firstForm := firstObject(accountForms)
	firstCPF := objectField(firstForm, "contentPublishForm")
	switch value := firstCPF["pubType"].(type) {
	case float64:
		return int(value) == 0
	case int:
		return value == 0
	default:
		return false
	}
}
