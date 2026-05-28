package draft

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) Save(payload map[string]interface{}) (map[string]interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	body := cloneMap(payload)
	delete(body, "action")
	body["isDraft"] = true
	return client.New(cfg).SaveDraft(body)
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
