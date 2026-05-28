package query

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) Categories(accountID, publishType string) (interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return client.New(cfg).Categories(accountID, publishType)
}

func (Service) Locations(accountID, keyword, locationType string) (interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return client.New(cfg).Locations(accountID, keyword, locationType)
}

func (Service) Music(accountID, keyword string) (interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return client.New(cfg).Music(accountID, keyword)
}

func (Service) Goods(accountID, keyword string) (interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return client.New(cfg).Goods(accountID, keyword)
}

func (Service) Collections(accountID, publishType string) (interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return client.New(cfg).Collections(accountID, publishType)
}

func (Service) Challenges(accountID, keyword, publishType string) (interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return client.New(cfg).Challenges(accountID, keyword, publishType)
}

func (Service) Records(platform, limit, status string) (interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return client.New(cfg).Records(platform, limit, status)
}

func (Service) Prepare(platform, publishType string) (client.PrepareData, error) {
	cfg, err := config.Load()
	if err != nil {
		return client.PrepareData{}, err
	}
	return client.New(cfg).Prepare(platform, publishType)
}
