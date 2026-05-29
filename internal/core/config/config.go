package config

import base "github.com/yixiaoer/yixiaoer-skill/internal/config"

const DefaultAPIURL = base.DefaultAPIURL

type Config = base.Config

func Load() (Config, error) {
	return base.Load()
}

func SaveAPIKey(apiKey string) (string, error) {
	return base.SaveAPIKey(apiKey)
}

func SaveLocalClientID(clientID string) (string, error) {
	return base.SaveLocalClientID(clientID)
}

func SaveLinkedAppState(appID, accountID, accountName string, connected bool) (string, error) {
	return base.SaveLinkedAppState(appID, accountID, accountName, connected)
}
