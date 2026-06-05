package query

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

func (s Service) Categories(accountID, publishType string) (interface{}, error) {
	return s.rt.Client.Categories(accountID, publishType)
}

func (s Service) Locations(accountID, keyword, locationType string) (interface{}, error) {
	return s.rt.Client.Locations(accountID, keyword, locationType)
}

func (s Service) Music(accountID, keyword string) (interface{}, error) {
	return s.rt.Client.Music(accountID, keyword)
}

func (s Service) Goods(accountID, keyword string) (interface{}, error) {
	return s.rt.Client.Goods(accountID, keyword)
}

func (s Service) Collections(accountID, publishType string) (interface{}, error) {
	return s.rt.Client.Collections(accountID, publishType)
}

func (s Service) Challenges(accountID, keyword, publishType string) (interface{}, error) {
	return s.rt.Client.Challenges(accountID, keyword, publishType)
}

func (s Service) Records(platform, limit, status string) (interface{}, error) {
	return s.rt.Client.Records(platform, limit, status)
}

func (s Service) Prepare(platform, publishType string) (api.PrepareData, error) {
	return s.rt.Client.Prepare(platform, publishType)
}
