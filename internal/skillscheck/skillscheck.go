package skillscheck

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const stampFile = "skills.stamp"

var markdownLinkPattern = regexp.MustCompile(`!?\[[^\]]*\]\(([^)]+)\)`)

type Status struct {
	Current string `json:"current"`
	Target  string `json:"target"`
	Sync    bool   `json:"inSync"`
	State   string `json:"state"`
}

type LinkIssue struct {
	File   string `json:"file"`
	Link   string `json:"link"`
	Target string `json:"target"`
	Error  string `json:"error"`
}

type LinkCheckReport struct {
	SkillDir      string      `json:"skillDir"`
	FilesScanned  int         `json:"filesScanned"`
	LinksChecked  int         `json:"linksChecked"`
	InvalidLinks  int         `json:"invalidLinks"`
	Issues        []LinkIssue `json:"issues,omitempty"`
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

func CheckSkillLinks(skillDir string) (LinkCheckReport, error) {
	report := LinkCheckReport{SkillDir: skillDir}
	root, err := filepath.Abs(skillDir)
	if err != nil {
		return report, err
	}
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(path)) != ".md" {
			return nil
		}
		report.FilesScanned++
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		matches := markdownLinkPattern.FindAllStringSubmatch(string(raw), -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			link := strings.TrimSpace(match[1])
			if !shouldCheckLocalLink(link) {
				continue
			}
			report.LinksChecked++
			targetPath := resolveLocalLink(path, link)
			info, statErr := os.Stat(targetPath)
			if statErr != nil || info == nil {
				report.Issues = append(report.Issues, LinkIssue{
					File:   relativeTo(root, path),
					Link:   link,
					Target: relativeTo(root, targetPath),
					Error:  statErr.Error(),
				})
				continue
			}
		}
		return nil
	})
	if err != nil {
		return report, err
	}
	sort.Slice(report.Issues, func(i, j int) bool {
		if report.Issues[i].File == report.Issues[j].File {
			return report.Issues[i].Link < report.Issues[j].Link
		}
		return report.Issues[i].File < report.Issues[j].File
	})
	report.InvalidLinks = len(report.Issues)
	if report.InvalidLinks > 0 {
		return report, fmt.Errorf("skill link check found %d invalid links", report.InvalidLinks)
	}
	return report, nil
}

func shouldCheckLocalLink(link string) bool {
	if link == "" {
		return false
	}
	if strings.HasPrefix(link, "#") {
		return false
	}
	lower := strings.ToLower(link)
	for _, prefix := range []string{"http://", "https://", "mailto:", "tel:"} {
		if strings.HasPrefix(lower, prefix) {
			return false
		}
	}
	return true
}

func resolveLocalLink(baseFile, link string) string {
	link = strings.TrimSpace(link)
	link = strings.Trim(strings.TrimPrefix(link, "<"), ">")
	if idx := strings.Index(link, "#"); idx >= 0 {
		link = link[:idx]
	}
	link = filepath.FromSlash(link)
	if filepath.IsAbs(link) {
		return filepath.Clean(link)
	}
	return filepath.Clean(filepath.Join(filepath.Dir(baseFile), link))
}

func relativeTo(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}
