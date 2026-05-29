package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const DefaultAPIURL = "https://www.yixiaoer.cn/api"

type Config struct {
	APIKey        string
	APIURL        string
	ProjectDir    string
	SchemaDir     string
	WorkDir       string
	ConfigPath    string
	LocalClientID string
	LinkedApp     LinkedAppState
}

type fileConfig struct {
	APIKey        string                 `json:"apiKey"`
	LocalClientID string                 `json:"localPublishClientId"`
	LinkedApps    map[string]linkedAppKV `json:"linkedApps,omitempty"`
}

type LinkedAppState struct {
	AppID       string `json:"appId"`
	Connected   bool   `json:"connected"`
	AccountID   string `json:"accountId,omitempty"`
	AccountName string `json:"accountName,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

type linkedAppKV struct {
	Connected   bool   `json:"connected"`
	AccountID   string `json:"accountId,omitempty"`
	AccountName string `json:"accountName,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

func Load() (Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}
	exeDir := ""
	if exePath, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exePath)
	}
	projectDir, err := resolveProjectDir(cwd, exeDir)
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
		APIKey:        strings.TrimSpace(fileCfg.APIKey),
		APIURL:        apiURL,
		ProjectDir:    projectDir,
		SchemaDir:     filepath.Join(projectDir, "schemas"),
		WorkDir:       cwd,
		ConfigPath:    configPath,
		LocalClientID: strings.TrimSpace(fileCfg.LocalClientID),
		LinkedApp:     fileCfg.linkedAppState("yixiaoer"),
	}, nil
}

func (c Config) RequireAPIKey() error {
	if c.APIKey == "" {
		return yxerrors.Auth("Missing apiKey configuration").
			WithHint("请先执行 yxer config set-api-key <apiKey> 完成 CLI 初始化。").
			WithNextCommand("yxer config set-api-key <apiKey>")
	}
	return nil
}

func SaveAPIKey(apiKey string) (string, error) {
	configPath, err := resolveConfigPath()
	if err != nil {
		return "", err
	}
	cfg, err := loadFileConfig(configPath)
	if err != nil {
		return "", err
	}
	cfg.APIKey = strings.TrimSpace(apiKey)
	if err := writeFileConfig(configPath, cfg); err != nil {
		return "", err
	}
	return configPath, nil
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
	if err := writeFileConfig(configPath, cfg); err != nil {
		return "", err
	}
	return configPath, nil
}

func SaveLinkedAppState(appID, accountID, accountName string, connected bool) (string, error) {
	configPath, err := resolveConfigPath()
	if err != nil {
		return "", err
	}
	cfg, err := loadFileConfig(configPath)
	if err != nil {
		return "", err
	}
	if cfg.LinkedApps == nil {
		cfg.LinkedApps = map[string]linkedAppKV{}
	}
	cfg.LinkedApps[strings.TrimSpace(appID)] = linkedAppKV{
		Connected:   connected,
		AccountID:   strings.TrimSpace(accountID),
		AccountName: strings.TrimSpace(accountName),
		UpdatedAt:   strings.TrimSpace(nowRFC3339()),
	}
	if !connected {
		cfg.LinkedApps[strings.TrimSpace(appID)] = linkedAppKV{
			Connected: false,
			UpdatedAt: strings.TrimSpace(nowRFC3339()),
		}
	}
	if err := writeFileConfig(configPath, cfg); err != nil {
		return "", err
	}
	return configPath, nil
}

func writeFileConfig(configPath string, cfg fileConfig) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(configPath, raw, 0o644); err != nil {
		return err
	}
	return nil
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

func (cfg fileConfig) linkedAppState(appID string) LinkedAppState {
	state := LinkedAppState{AppID: appID}
	if cfg.LinkedApps == nil {
		return state
	}
	raw, ok := cfg.LinkedApps[appID]
	if !ok {
		return state
	}
	state.Connected = raw.Connected
	state.AccountID = strings.TrimSpace(raw.AccountID)
	state.AccountName = strings.TrimSpace(raw.AccountName)
	state.UpdatedAt = strings.TrimSpace(raw.UpdatedAt)
	return state
}

var nowRFC3339 = func() string {
	return time.Now().Format(time.RFC3339)
}

func resolveProjectDir(cwd, exeDir string) (string, error) {
	if value := strings.TrimSpace(os.Getenv("YIXIAOER_PROJECT_DIR")); value != "" {
		abs, err := filepath.Abs(value)
		if err != nil {
			return "", err
		}
		if !isProjectDir(abs) {
			return "", yxerrors.Usage("project directory not found", abs).
				WithHint("请确认 YIXIAOER_PROJECT_DIR 指向项目根目录，且 schemas 和 workflows 目录存在。")
		}
		return abs, nil
	}
	if dir, ok := findProjectDirFrom(cwd); ok {
		return dir, nil
	}
	if exeDir != "" {
		if dir, ok := findProjectDirFrom(exeDir); ok {
			return dir, nil
		}
	}
	abs, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func findProjectDirFrom(start string) (string, bool) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for {
		if isProjectDir(current) {
			return current, true
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func isProjectDir(path string) bool {
	return isDir(filepath.Join(path, "schemas")) && isDir(filepath.Join(path, "workflows"))
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
