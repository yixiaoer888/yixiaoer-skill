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
	SkillDir     string      `json:"skillDir"`
	FilesScanned int         `json:"filesScanned"`
	LinksChecked int         `json:"linksChecked"`
	InvalidLinks int         `json:"invalidLinks"`
	Issues       []LinkIssue `json:"issues,omitempty"`
}

type ValidationIssue struct {
	File    string `json:"file"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type FormatCheckReport struct {
	SkillDir         string            `json:"skillDir"`
	CheckedFile      string            `json:"checkedFile"`
	RequiredFields   []string          `json:"requiredFields"`
	RequiredSections []string          `json:"requiredSections"`
	InvalidFields    int               `json:"invalidFields"`
	MissingSections  int               `json:"missingSections"`
	Issues           []ValidationIssue `json:"issues,omitempty"`
}

type StructureCheckReport struct {
	SkillDir            string            `json:"skillDir"`
	RequiredPaths       []string          `json:"requiredPaths"`
	RequiredDomainDocs  []string          `json:"requiredDomainDocs"`
	CheckedDomainDocs   int               `json:"checkedDomainDocs"`
	MissingPaths        int               `json:"missingPaths"`
	InvalidDomainDocs   int               `json:"invalidDomainDocs"`
	MissingIntentRoutes int               `json:"missingIntentRoutes"`
	Issues              []ValidationIssue `json:"issues,omitempty"`
}

type PackageCheckReport struct {
	SkillDir      string               `json:"skillDir"`
	Valid         bool                 `json:"valid"`
	InvalidChecks int                  `json:"invalidChecks"`
	Format        FormatCheckReport    `json:"format"`
	Structure     StructureCheckReport `json:"structure"`
	Links         LinkCheckReport      `json:"links"`
}

var frontmatterPattern = regexp.MustCompile(`(?s)^---\n(.*?)\n---(?:\n|$)`)
var headingPattern = regexp.MustCompile(`(?m)^##\s+(.+?)\s*$`)

var requiredSkillFields = []string{
	"name",
	"version",
	"description",
	"metadata",
	"metadata.requires.bins",
	"metadata.cliHelp",
}

var requiredSkillSections = []string{
	"能力索引",
	"意图分流",
	"命令探索",
	"全局规则",
}

var requiredSkillPaths = []string{
	"SKILL.md",
	"QUICKSTART.md",
	"references/yixiaoer-shared.md",
	"references/domains/publish.md",
	"references/domains/accounts-and-env.md",
	"references/domains/draft-and-material.md",
	"references/domains/troubleshooting.md",
	"references/domains/install-and-sync.md",
}

var requiredDomainDocs = []string{
	"references/domains/publish.md",
	"references/domains/accounts-and-env.md",
	"references/domains/draft-and-material.md",
	"references/domains/troubleshooting.md",
	"references/domains/install-and-sync.md",
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

func CheckSkillFormat(skillDir string) (FormatCheckReport, error) {
	report := FormatCheckReport{
		SkillDir:         skillDir,
		CheckedFile:      "SKILL.md",
		RequiredFields:   append([]string(nil), requiredSkillFields...),
		RequiredSections: append([]string(nil), requiredSkillSections...),
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	raw, err := os.ReadFile(skillFile)
	if err != nil {
		report.Issues = append(report.Issues, ValidationIssue{
			File:    "SKILL.md",
			Rule:    "skill.file.exists",
			Message: err.Error(),
		})
		report.InvalidFields = len(report.Issues)
		return report, fmt.Errorf("skill format check failed")
	}

	content := normalizeMarkdown(raw)
	frontmatter, body, frontmatterErr := extractFrontmatter(content)
	if frontmatterErr != nil {
		report.Issues = append(report.Issues, ValidationIssue{
			File:    "SKILL.md",
			Rule:    "frontmatter",
			Message: frontmatterErr.Error(),
		})
		report.InvalidFields = len(report.Issues)
		return report, fmt.Errorf("skill format check failed")
	}

	for _, field := range requiredSkillFields {
		if !hasFrontmatterField(frontmatter, field) {
			report.Issues = append(report.Issues, ValidationIssue{
				File:    "SKILL.md",
				Rule:    field,
				Message: "missing required frontmatter field",
			})
		}
	}

	headings := collectHeadings(body)
	for _, section := range requiredSkillSections {
		if _, ok := headings[section]; !ok {
			report.Issues = append(report.Issues, ValidationIssue{
				File:    "SKILL.md",
				Rule:    "section:" + section,
				Message: "missing required section heading",
			})
			report.MissingSections++
		}
	}

	report.InvalidFields = len(report.Issues) - report.MissingSections
	if len(report.Issues) > 0 {
		return report, fmt.Errorf("skill format check failed")
	}
	return report, nil
}

func CheckSkillStructure(skillDir string) (StructureCheckReport, error) {
	report := StructureCheckReport{
		SkillDir:           skillDir,
		RequiredPaths:      append([]string(nil), requiredSkillPaths...),
		RequiredDomainDocs: append([]string(nil), requiredDomainDocs...),
	}

	for _, relPath := range requiredSkillPaths {
		fullPath := filepath.Join(skillDir, filepath.FromSlash(relPath))
		info, err := os.Stat(fullPath)
		if err != nil || info == nil || info.IsDir() {
			message := "missing required file"
			if err != nil {
				message = err.Error()
			}
			report.Issues = append(report.Issues, ValidationIssue{
				File:    relPath,
				Rule:    "path.exists",
				Message: message,
			})
			report.MissingPaths++
		}
	}

	for _, relPath := range requiredDomainDocs {
		fullPath := filepath.Join(skillDir, filepath.FromSlash(relPath))
		raw, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		report.CheckedDomainDocs++
		content := normalizeMarkdown(raw)
		if !containsAny(content, "## 读取顺序", "## 优先读取") {
			report.Issues = append(report.Issues, ValidationIssue{
				File:    relPath,
				Rule:    "domain.read-order",
				Message: "domain doc must define a read-order section",
			})
			report.InvalidDomainDocs++
		}
		if !strings.Contains(content, "## 常用命令") {
			report.Issues = append(report.Issues, ValidationIssue{
				File:    relPath,
				Rule:    "domain.commands",
				Message: "domain doc must define a common commands section",
			})
			report.InvalidDomainDocs++
		}
		if !containsAny(content, "## 决策规则", "## 规则", "## 排查规则", "## 强制门禁", "## 决策提示") {
			report.Issues = append(report.Issues, ValidationIssue{
				File:    relPath,
				Rule:    "domain.rules",
				Message: "domain doc must define a rules, gate, or decision section",
			})
			report.InvalidDomainDocs++
		}
		if !hasIntentMapping(content) {
			report.Issues = append(report.Issues, ValidationIssue{
				File:    relPath,
				Rule:    "domain.intent-routing",
				Message: "domain doc should include user intent routing cues such as 用户只说/用户明确说/适用范围",
			})
			report.MissingIntentRoutes++
		}
	}

	if len(report.Issues) > 0 {
		return report, fmt.Errorf("skill structure check failed")
	}
	return report, nil
}

func CheckSkillPackage(skillDir string) (PackageCheckReport, error) {
	report := PackageCheckReport{
		SkillDir: skillDir,
		Valid:    true,
	}

	formatReport, formatErr := CheckSkillFormat(skillDir)
	report.Format = formatReport
	if formatErr != nil {
		report.Valid = false
		report.InvalidChecks++
	}

	structureReport, structureErr := CheckSkillStructure(skillDir)
	report.Structure = structureReport
	if structureErr != nil {
		report.Valid = false
		report.InvalidChecks++
	}

	linkReport, linkErr := CheckSkillLinks(skillDir)
	report.Links = linkReport
	if linkErr != nil {
		report.Valid = false
		report.InvalidChecks++
	}

	if !report.Valid {
		return report, fmt.Errorf("skill package check failed")
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

func normalizeMarkdown(raw []byte) string {
	return strings.ReplaceAll(string(raw), "\r\n", "\n")
}

func extractFrontmatter(content string) (string, string, error) {
	matches := frontmatterPattern.FindStringSubmatch(content)
	if matches == nil {
		return "", "", errors.New("SKILL.md must start with valid YAML frontmatter")
	}
	return matches[1], strings.TrimSpace(content[len(matches[0]):]), nil
}

func hasFrontmatterField(frontmatter, field string) bool {
	switch field {
	case "metadata.requires.bins":
		return regexp.MustCompile(`(?m)^metadata:\s*$`).MatchString(frontmatter) &&
			regexp.MustCompile(`(?m)^[ \t]+requires:\s*$`).MatchString(frontmatter) &&
			regexp.MustCompile(`(?m)^[ \t]+bins:\s*\[[^\]]+\]`).MatchString(frontmatter)
	case "metadata.cliHelp":
		return regexp.MustCompile(`(?m)^metadata:\s*$`).MatchString(frontmatter) &&
			regexp.MustCompile(`(?m)^[ \t]+cliHelp:\s*.+$`).MatchString(frontmatter)
	default:
		return regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(field) + `:\s*.+?$|^` + regexp.QuoteMeta(field) + `:\s*$`).MatchString(frontmatter)
	}
}

func collectHeadings(body string) map[string]struct{} {
	headings := make(map[string]struct{})
	for _, match := range headingPattern.FindAllStringSubmatch(body, -1) {
		if len(match) < 2 {
			continue
		}
		headings[strings.TrimSpace(match[1])] = struct{}{}
	}
	return headings
}

func containsAny(content string, values ...string) bool {
	for _, value := range values {
		if strings.Contains(content, value) {
			return true
		}
	}
	return false
}

func hasIntentMapping(content string) bool {
	return containsAny(content, "用户只说", "用户明确说", "适用范围")
}
