package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const DefaultAPIURL = "https://www.yixiaoer.cn/api"

type Config struct {
	APIKey        string
	APIURL        string
	SchemaDir     string
	WorkDir       string
	ConfigPath    string
	LocalClientID string
}

type fileConfig struct {
	LocalClientID string `json:"localPublishClientId"`
}

func Load() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}
	apiURL := strings.TrimRight(os.Getenv("YIXIAOER_API_URL"), "/")
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}
	configPath, err := resolveConfigPath()
	if err != nil {
		return Config{}, err
	}
	fileCfg, err := loadFileConfig(configPath)
	if err != nil {
		return Config{}, err
	}
	return Config{
		APIKey:        os.Getenv("YIXIAOER_API_KEY"),
		APIURL:        apiURL,
		SchemaDir:     filepath.Join(wd, "schemas"),
		WorkDir:       wd,
		ConfigPath:    configPath,
		LocalClientID: strings.TrimSpace(fileCfg.LocalClientID),
	}, nil
}

func (c Config) RequireAPIKey() error {
	if c.APIKey == "" {
		return yxerrors.Auth("Missing YIXIAOER_API_KEY environment variable")
	}
	return nil
}

func SaveLocalClientID(clientID string) (string, error) {
	configPath, err := resolveConfigPath()
	if err != nil {
		return "", err
	}
	cfg, err := loadFileConfig(configPath)
	if err != nil {
		return "", err
	}
	cfg.LocalClientID = strings.TrimSpace(clientID)
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return "", err
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(configPath, raw, 0o644); err != nil {
		return "", err
	}
	return configPath, nil
}

func resolveConfigPath() (string, error) {
	if value := strings.TrimSpace(os.Getenv("YIXIAOER_CONFIG")); value != "" {
		return value, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".yxer", "config.json"), nil
}

func loadFileConfig(path string) (fileConfig, error) {
	var cfg fileConfig
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
