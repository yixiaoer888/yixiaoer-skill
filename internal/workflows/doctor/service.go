package doctor

import (
	"os"
	"path/filepath"

	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type Service struct {
	rt *app.Runtime
}

func NewService(rt *app.Runtime) Service {
	return Service{rt: rt}
}

func (s Service) Check() (map[string]interface{}, error) {
	cfg := s.rt.Config
	checks := map[string]interface{}{
		"projectDir":           cfg.ProjectDir,
		"workDir":              cfg.WorkDir,
		"apiUrl":               cfg.APIURL,
		"apiKeyPresent":        cfg.APIKey != "",
		"schemaDir":            cfg.SchemaDir,
		"schemaDirOK":          PathExists(cfg.SchemaDir),
		"workflowsOK":          hasWorkflowDocs(cfg.ProjectDir),
		"workflowDocsPath":     workflowDocsPath(cfg.ProjectDir),
		"localPublishClientId": cfg.LocalClientID,
		"publishChannelReadiness": map[string]interface{}{
			"cloud": map[string]interface{}{
				"configured": cfg.APIKey != "",
				"note":       "云发布是否需要平台代理，取决于目标平台账号配置；正式发布前建议结合 accounts/prepare 结果确认。",
			},
			"local": map[string]interface{}{
				"configured": cfg.LocalClientID != "",
				"clientId":   cfg.LocalClientID,
				"note":       "本机发布需要可用 clientId，且蚁小二客户端在线。",
			},
		},
	}
	if cfg.APIKey == "" {
		return nil, yxerrors.Auth("Missing apiKey configuration").
			WithHint("请先执行 yxer config set-api-key <apiKey> 完成 CLI 初始化。").
			WithNextCommand("yxer config set-api-key <apiKey>")
	}
	if !PathExists(cfg.SchemaDir) {
		return nil, yxerrors.Usage("schema directory not found", cfg.SchemaDir).
			WithHint("请确认当前目录位于项目内，或 yxer 可执行文件位于项目目录中；也可显式设置 YIXIAOER_PROJECT_DIR。").
			WithNextCommand("yxer schema list")
	}
	return checks, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasWorkflowDocs(projectDir string) bool {
	return PathExists(filepath.Join(projectDir, "workflows")) || PathExists(filepath.Join(projectDir, "references", "workflows"))
}

func workflowDocsPath(projectDir string) string {
	workflowsDir := filepath.Join(projectDir, "workflows")
	if PathExists(workflowsDir) {
		return workflowsDir
	}
	return filepath.Join(projectDir, "references", "workflows")
}
