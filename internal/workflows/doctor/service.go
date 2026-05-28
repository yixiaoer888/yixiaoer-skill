package doctor

import (
	"os"

	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (Service) Check() (map[string]interface{}, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	checks := map[string]interface{}{
		"apiUrl":        cfg.APIURL,
		"apiKeyPresent": cfg.APIKey != "",
		"schemaDir":     cfg.SchemaDir,
		"schemaDirOK":   PathExists(cfg.SchemaDir),
		"workflowsOK":   PathExists("workflows"),
	}
	if cfg.APIKey == "" {
		return nil, yxerrors.Auth("Missing YIXIAOER_API_KEY environment variable").
			WithHint("请先设置环境变量 YIXIAOER_API_KEY。").
			WithNextCommand("yxer doctor")
	}
	if !PathExists(cfg.SchemaDir) {
		return nil, yxerrors.Usage("schema directory not found", cfg.SchemaDir).
			WithHint("请确认当前目录是项目根目录，且 schemas 目录存在。").
			WithNextCommand("yxer schema list")
	}
	return checks, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
