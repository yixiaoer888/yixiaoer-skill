package upload

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
)

type Service struct {
	rt *app.Runtime
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func (s Service) Upload(pathOrURL, bucket string, autoMeta bool) (api.UploadResult, error) {
	return s.rt.Client.Upload(pathOrURL, bucket, autoMeta)
}

func (s Service) UploadWithOptions(pathOrURL, bucket string, autoMeta bool, opts api.UploadOptions) (api.UploadResult, error) {
	return s.rt.Client.UploadWithOptions(pathOrURL, bucket, autoMeta, opts)
}
