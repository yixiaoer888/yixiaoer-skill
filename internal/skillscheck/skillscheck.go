package skillscheck

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const stampFile = "skills.stamp"

type Status struct {
	Current string `json:"current"`
	Target  string `json:"target"`
	Sync    bool   `json:"inSync"`
	State   string `json:"state"`
}

func StampPath() (string, error) {
	baseDir, err := resolveBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, stampFile), nil
}

func ReadStamp() (string, error) {
	path, err := StampPath()
	if err != nil {
		return "", err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

func WriteStamp(version string) error {
	path, err := StampPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strings.TrimSpace(version)+"\n"), 0o644)
}

func Check(currentVersion string) (Status, error) {
	stamp, err := ReadStamp()
	if err != nil {
		return Status{}, err
	}
	status := Status{
		Current: stamp,
		Target:  strings.TrimSpace(currentVersion),
		Sync:    stamp != "" && normalizeVersion(stamp) == normalizeVersion(currentVersion),
		State:   "stale",
	}
	switch {
	case stamp == "":
		status.State = "missing"
		status.Sync = false
	case status.Sync:
		status.State = "in_sync"
	default:
		status.State = "stale"
	}
	return status, nil
}

func Notice(currentVersion, skillDir string) (map[string]interface{}, error) {
	status, err := Check(currentVersion)
	if err != nil {
		return nil, err
	}
	if status.Sync {
		return nil, nil
	}
	command := `npx skills add "` + skillDir + `" -y`
	return map[string]interface{}{
		"type":    "skills_out_of_sync",
		"current": status.Current,
		"target":  status.Target,
		"state":   status.State,
		"message": "本地 AI skill 与当前 yxer 版本未同步，请重新安装技能包。",
		"command": command,
	}, nil
}

func normalizeVersion(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "v")
}

func DetectSkillDir() (string, error) {
	if value := strings.TrimSpace(os.Getenv("YIXIAOER_SKILL_DIR")); value != "" {
		return value, nil
	}
	if wd, err := os.Getwd(); err == nil {
		if dir, ok := findSkillDirFrom(wd); ok {
			return dir, nil
		}
	}
	if exePath, err := os.Executable(); err == nil {
		if dir, ok := findSkillDirFrom(filepath.Dir(exePath)); ok {
			return dir, nil
		}
	}
	return "", errors.New(`failed to locate "skills/yixiaoer"; set YIXIAOER_SKILL_DIR explicitly`)
}

func findSkillDirFrom(start string) (string, bool) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for {
		candidate := filepath.Join(current, "skills", "yixiaoer")
		if isSkillDir(candidate) {
			return candidate, true
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func isSkillDir(path string) bool {
	info, err := os.Stat(filepath.Join(path, "SKILL.md"))
	return err == nil && !info.IsDir()
}

func resolveBaseDir() (string, error) {
	if value := strings.TrimSpace(os.Getenv("YIXIAOER_CONFIG")); value != "" {
		return filepath.Dir(value), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".yxer"), nil
}
