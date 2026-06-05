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

func (s Service) Locations(accountID, keyword, locationType, nextPage string) (interface{}, error) {
	return s.rt.Client.Locations(accountID, keyword, locationType, nextPage)
}

func (s Service) Music(accountID, keyword, categoryID, categoryName, nextPage string) (interface{}, error) {
	return s.rt.Client.Music(accountID, keyword, categoryID, categoryName, nextPage)
}

func (s Service) MusicCategories(accountID string) (interface{}, error) {
	return s.rt.Client.MusicCategories(accountID)
}

func (s Service) Goods(accountID, keyword, nextPage string) (interface{}, error) {
	return s.rt.Client.Goods(accountID, keyword, nextPage)
}

func (s Service) Collections(accountID, publishType string) (interface{}, error) {
	return s.rt.Client.Collections(accountID, publishType)
}

func (s Service) MiniApps(accountID, keyword string) (interface{}, error) {
	return s.rt.Client.MiniApps(accountID, keyword)
}

func (s Service) SyncApps(accountID string) (interface{}, error) {
	return s.rt.Client.SyncApps(accountID)
}

func (s Service) Games(accountID, keyword string) (interface{}, error) {
	return s.rt.Client.Games(accountID, keyword)
}

func (s Service) HotEvents(accountID, publishType string) (interface{}, error) {
	return s.rt.Client.HotEvents(accountID, publishType)
}

func (s Service) Groups(accountID string) (interface{}, error) {
	return s.rt.Client.Groups(accountID)
}

func (s Service) Activities(accountID, publishType, categoryID, keyword string) (interface{}, error) {
	return s.rt.Client.Activities(accountID, publishType, categoryID, keyword)
}

func (s Service) Challenges(accountID, keyword, publishType, nextPage string) (interface{}, error) {
	return s.rt.Client.Challenges(accountID, keyword, publishType, nextPage)
}

func (s Service) Records(platform, limit, status string) (interface{}, error) {
	return s.rt.Client.Records(platform, limit, status)
}

func (s Service) Details(taskSetID string) (interface{}, error) {
	return s.rt.Client.Details(taskSetID)
}

func (s Service) AccountOverviews(opts api.AccountOverviewOptions) (interface{}, error) {
	return s.rt.Client.AccountOverviews(opts)
}

func (s Service) ContentOverviews(opts api.ContentOverviewOptions) (interface{}, error) {
	return s.rt.Client.ContentOverviews(opts)
}

func (s Service) Proxies(size string) (interface{}, error) {
	return s.rt.Client.Proxies(size)
}

func (s Service) ProxyAreas() (interface{}, error) {
	return s.rt.Client.ProxyAreas()
}

func (s Service) UpdateAccount(accountID string, body map[string]interface{}) (interface{}, error) {
	return s.rt.Client.UpdateAccount(accountID, body)
}

func (s Service) Prepare(platform, publishType string) (api.PrepareData, error) {
	return s.rt.Client.Prepare(platform, publishType)
}
