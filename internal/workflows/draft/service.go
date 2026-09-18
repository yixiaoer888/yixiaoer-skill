package draft

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
)

type Service struct {
	rt *app.Runtime
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func (s Service) Save(payload map[string]interface{}) (map[string]interface{}, error) {
	body := publishflow.BuildDraftBody(payload)
	return s.rt.Client.SaveDraft(body)
}
