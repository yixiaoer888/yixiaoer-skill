package upload

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) Upload(pathOrURL, bucket string) (client.UploadResult, error) {
	cfg, err := config.Load()
	if err != nil {
		return client.UploadResult{}, err
	}
	return client.New(cfg).Upload(pathOrURL, bucket)
}
