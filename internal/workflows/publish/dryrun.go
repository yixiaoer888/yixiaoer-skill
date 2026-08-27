package publish

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
)

type DryRunResult struct {
	Platform          string                               `json:"platform"`
	PublishType       string                               `json:"publishType"`
	PublishBody       map[string]interface{}               `json:"request"`
	PublishArgs       map[string]interface{}               `json:"publishArgs,omitempty"`
	PublishMode       string                               `json:"publishChannel"`
	PublishModeSource string                               `json:"publishChannelSource"`
	ClientID          string                               `json:"clientId,omitempty"`
	ClientIDSource    string                               `json:"clientIdSource"`
	RequestHash       string                               `json:"requestHash"`
	AccountIDs        []string                             `json:"accountIds,omitempty"`
	PlatformDraft     bool                                 `json:"platformDraft"`
	YixiaoerDraft     bool                                 `json:"yixiaoerDraft"`
	SchemaChecked     bool                                 `json:"schemaChecked"`
	RemoteChecks      bool                                 `json:"remoteChecks"`
	Normalizations    []publishmod.NormalizationEvent      `json:"normalizations,omitempty"`
	InferredFields    map[string]InferredField             `json:"inferredFields,omitempty"`
	ContentImages     []ArticleContentImageMaterialization `json:"contentImageMaterialization,omitempty"`
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
				"platform":                    result.Platform,
				"publishType":                 result.PublishType,
				"publishChannel":              result.PublishMode,
				"effectivePublishChannel":     result.PublishMode,
				"publishChannelSource":        result.PublishModeSource,
				"clientId":                    result.ClientID,
				"clientIdSource":              result.ClientIDSource,
				"requestHash":                 result.RequestHash,
				"accountIds":                  result.AccountIDs,
				"platformDraft":               result.PlatformDraft,
				"yixiaoerDraft":               result.YixiaoerDraft,
				"schemaChecked":               result.SchemaChecked,
				"remoteChecks":                result.RemoteChecks,
				"normalizations":              normalizationsForMeta(result.Normalizations),
				"inferredFields":              inferredFieldsForMeta(result.InferredFields),
				"contentImageMaterialization": contentImageMaterializationForMeta(result.ContentImages),
			},
		},
	}, nil
}

func normalizationsForMeta(events []publishmod.NormalizationEvent) []publishmod.NormalizationEvent {
	if events == nil {
		return []publishmod.NormalizationEvent{}
	}
	return events
}

func inferredFieldsForMeta(fields map[string]InferredField) map[string]InferredField {
	if fields == nil {
		return map[string]InferredField{}
	}
	return fields
}

func contentImageMaterializationForMeta(events []ArticleContentImageMaterialization) []ArticleContentImageMaterialization {
	if events == nil {
		return []ArticleContentImageMaterialization{}
	}
	return events
}

func (s Service) DryRun(input ExecuteInput) (DryRunResult, error) {
	prepared, err := s.Prepare(input, PrepareOptions{TraceNormalizations: true, RemoteChecks: RemoteChecksCloudWithKey})
	if err != nil {
		return DryRunResult{}, err
	}
	if err := AssertShoppingCartEntitlements(s.rt.Client, prepared.Payload, prepared.Platforms...); err != nil {
		return DryRunResult{}, err
	}

	return DryRunResult{
		Platform:          prepared.Platform,
		PublishType:       prepared.PublishType,
		PublishBody:       prepared.PublishBody,
		PublishArgs:       prepared.PublishArgs,
		PublishMode:       prepared.PublishMode,
		PublishModeSource: prepared.PublishModeSource,
		ClientID:          prepared.ClientID,
		ClientIDSource:    prepared.ClientIDSource,
		RequestHash:       requestHash(prepared.PublishBody),
		AccountIDs:        prepared.Preflight.AccountIDs,
		PlatformDraft:     isPlatformDraftPublish(prepared.PublishBody),
		YixiaoerDraft:     inferYixiaoerDraft(prepared.PublishBody),
		SchemaChecked:     true,
		RemoteChecks:      prepared.RemoteChecked || (shoppingCartEntitlementsSupported(prepared.Platforms) && len(ShoppingCartAccountIDs(prepared.Payload)) > 0),
		Normalizations:    prepared.Normalizations,
		InferredFields:    prepared.InferredFields,
		ContentImages:     previewArticleContentImageMaterialization(prepared.PublishBody, prepared.ContentBaseDir),
	}, nil
}

func requestHash(body map[string]interface{}) string {
	raw, err := json.Marshal(body)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
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
