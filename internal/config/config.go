package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const DefaultAPIURL = "https://www.yixiaoer.cn/api"

type Config struct {
	APIKey    string
	APIURL    string
	SchemaDir string
	WorkDir   string
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
	return Config{
		APIKey:    os.Getenv("YIXIAOER_API_KEY"),
		APIURL:    apiURL,
		SchemaDir: filepath.Join(wd, "schemas"),
		WorkDir:   wd,
	}, nil
}

func (c Config) RequireAPIKey() error {
	if c.APIKey == "" {
		return yxerrors.Auth("Missing YIXIAOER_API_KEY environment variable")
	}
	return nil
}
