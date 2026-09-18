package draft

import publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"

func PreviewSave(payload map[string]interface{}) map[string]interface{} {
	return publishflow.BuildDraftBody(payload)
}
