package publish

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
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
	SchemaChecked bool                   `json:"schemaChecked"`
}

func (Service) DryRun(input ExecuteInput) (DryRunResult, error) {
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil {
		return DryRunResult{}, err
	}
	platforms := []string{platform}
	publishArgs := ExtractPublishArgs(input.Payload)

	cfg, err := config.Load()
	if err != nil {
		return DryRunResult{}, err
	}
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
	NormalizeStandardPublishArgs(ExtractPublishArgs(resolvedPayload))

	validator := schema.NewValidator(cfg.SchemaDir)
	for _, platform := range platforms {
		result := validator.Validate(platform, input.PublishType, resolvedPayload)
		if !result.Valid {
			return DryRunResult{}, yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint("请根据对应平台 schema 修正 payload 字段后重试。").
				WithNextCommand("yxer schema get <platform> <type>")
		}
	}
	preflight := Preflight(input.PublishType, platforms, resolvedPayload)
	if len(preflight.Errors) > 0 {
		return DryRunResult{}, yxerrors.Usage("Publish preflight failed", preflight.Errors).
			WithHint("请先完成资源上传、账号校验，并确保发布参数中不包含外部 URL。")
	}

	body := BuildPublishBody(resolvedPayload, publishArgs, input.PublishType, platforms)
	body["publishChannel"] = channel
	if clientID != "" {
		body["clientId"] = clientID
	} else {
		delete(body, "clientId")
	}

	return DryRunResult{
		Platform:      platform,
		PublishType:   input.PublishType,
		PublishBody:   body,
		PublishArgs:   publishArgs,
		PublishMode:   channel,
		ClientID:      clientID,
		AccountIDs:    preflight.AccountIDs,
		SchemaChecked: true,
	}, nil
}
