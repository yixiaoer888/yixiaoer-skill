package draft

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
)

type Service struct {
	rt *app.Runtime
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func (s Service) Save(payload map[string]interface{}) (map[string]interface{}, error) {
	body := cloneMap(payload)
	delete(body, "action")
	body["isDraft"] = true
	return s.rt.Client.SaveDraft(body)
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
