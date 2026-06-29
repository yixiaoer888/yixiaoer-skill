package skillscheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckMissingStamp(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(dir, "config.json"))

	status, err := Check("3.0.0")
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if status.State != "missing" {
		t.Fatalf("status.State = %q, want missing", status.State)
	}
	if status.Sync {
		t.Fatal("status.Sync = true, want false")
	}
}

func TestSkillVersion(t *testing.T) {
	skillDir := createValidSkillFixture(t)

	version, err := SkillVersion(skillDir)
	if err != nil {
		t.Fatalf("SkillVersion returned error: %v", err)
	}
	if version != "3.1.1" {
		t.Fatalf("SkillVersion = %q, want 3.1.1", version)
	}
}

func TestSkillVersionMissing(t *testing.T) {
	skillDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# missing frontmatter\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := SkillVersion(skillDir); err == nil {
		t.Fatal("SkillVersion error = nil, want non-nil")
	}
}

func TestWriteAndCheckStamp(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(dir, "config.json"))

	if err := WriteStamp("3.1.1"); err != nil {
		t.Fatalf("WriteStamp returned error: %v", err)
	}

	status, err := Check("3.1.1")
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if status.State != "in_sync" {
		t.Fatalf("status.State = %q, want in_sync", status.State)
	}
	if !status.Sync {
		t.Fatal("status.Sync = false, want true")
	}
}

func TestNoticeStale(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(dir, "config.json"))

	path, err := StampPath()
	if err != nil {
		t.Fatalf("StampPath returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte("2.9.0\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	notice, err := Notice("3.1.1", `C:\repo\skills\yixiaoer`)
	if err != nil {
		t.Fatalf("Notice returned error: %v", err)
	}
	if notice == nil {
		t.Fatal("notice = nil, want non-nil")
	}
	if got := notice["state"]; got != "stale" {
		t.Fatalf("notice[state] = %v, want stale", got)
	}
	if got := notice["command"]; got != "yxer skill sync" {
		t.Fatalf("notice[command] = %v, want yxer skill sync", got)
	}
	fallbacks, ok := notice["fallbackCommands"].([]string)
	if !ok {
		t.Fatalf("notice[fallbackCommands] type = %T, want []string", notice["fallbackCommands"])
	}
	if len(fallbacks) != 2 {
		t.Fatalf("len(fallbackCommands) = %d, want 2", len(fallbacks))
	}
}

func TestFindSkillDirFromPackagedLayout(t *testing.T) {
	root := t.TempDir()
	skillDir := filepath.Join(root, "skills", "yixiaoer")
	distDir := filepath.Join(root, "dist")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, ok := findSkillDirFrom(distDir)
	if !ok {
		t.Fatal("findSkillDirFrom(distDir) = not found, want found")
	}
	if got != skillDir {
		t.Fatalf("findSkillDirFrom(distDir) = %q, want %q", got, skillDir)
	}
}

func TestCheckSkillLinksSuccess(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "yixiaoer")
	refsDir := filepath.Join(skillDir, "references")
	if err := os.MkdirAll(refsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("[ref](./references/doc.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refsDir, "doc.md"), []byte("[home](../SKILL.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := CheckSkillLinks(skillDir)
	if err != nil {
		t.Fatalf("CheckSkillLinks returned error: %v", err)
	}
	if report.FilesScanned != 2 {
		t.Fatalf("FilesScanned = %d, want 2", report.FilesScanned)
	}
	if report.LinksChecked != 2 {
		t.Fatalf("LinksChecked = %d, want 2", report.LinksChecked)
	}
	if report.InvalidLinks != 0 {
		t.Fatalf("InvalidLinks = %d, want 0", report.InvalidLinks)
	}
}

func TestCheckSkillLinksMissingTarget(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "yixiaoer")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("[missing](./references/missing.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := CheckSkillLinks(skillDir)
	if err == nil {
		t.Fatal("CheckSkillLinks error = nil, want non-nil")
	}
	if report.InvalidLinks != 1 {
		t.Fatalf("InvalidLinks = %d, want 1", report.InvalidLinks)
	}
	if len(report.Issues) != 1 {
		t.Fatalf("len(Issues) = %d, want 1", len(report.Issues))
	}
	if report.Issues[0].File != "SKILL.md" {
		t.Fatalf("issue file = %q, want SKILL.md", report.Issues[0].File)
	}
}

func TestCheckSkillFormatSuccess(t *testing.T) {
	skillDir := createValidSkillFixture(t)

	report, err := CheckSkillFormat(skillDir)
	if err != nil {
		t.Fatalf("CheckSkillFormat returned error: %v", err)
	}
	if report.InvalidFields != 0 {
		t.Fatalf("InvalidFields = %d, want 0", report.InvalidFields)
	}
	if report.MissingSections != 0 {
		t.Fatalf("MissingSections = %d, want 0", report.MissingSections)
	}
}

func TestCheckSkillFormatMissingCliHelp(t *testing.T) {
	skillDir := createValidSkillFixture(t)
	content := `---
name: yixiaoer
version: 3.1.1
description: "skill"
metadata:
  requires:
    bins: ["yxer"]
---

# 蚁小二 Skill

## 能力索引

## 意图分流

## 命令探索

## 全局规则
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := CheckSkillFormat(skillDir)
	if err == nil {
		t.Fatal("CheckSkillFormat error = nil, want non-nil")
	}
	if report.InvalidFields == 0 {
		t.Fatal("InvalidFields = 0, want > 0")
	}
}

func TestCheckSkillStructureSuccess(t *testing.T) {
	skillDir := createValidSkillFixture(t)

	report, err := CheckSkillStructure(skillDir)
	if err != nil {
		t.Fatalf("CheckSkillStructure returned error: %v", err)
	}
	if report.MissingPaths != 0 {
		t.Fatalf("MissingPaths = %d, want 0", report.MissingPaths)
	}
	if report.InvalidDomainDocs != 0 {
		t.Fatalf("InvalidDomainDocs = %d, want 0", report.InvalidDomainDocs)
	}
	if report.MissingIntentRoutes != 0 {
		t.Fatalf("MissingIntentRoutes = %d, want 0", report.MissingIntentRoutes)
	}
}

func TestCheckSkillStructureMissingQuickstart(t *testing.T) {
	skillDir := createValidSkillFixture(t)
	if err := os.Remove(filepath.Join(skillDir, "QUICKSTART.md")); err != nil {
		t.Fatal(err)
	}

	report, err := CheckSkillStructure(skillDir)
	if err == nil {
		t.Fatal("CheckSkillStructure error = nil, want non-nil")
	}
	if report.MissingPaths == 0 {
		t.Fatal("MissingPaths = 0, want > 0")
	}
}

func TestCheckSkillPackageSuccess(t *testing.T) {
	skillDir := createValidSkillFixture(t)

	report, err := CheckSkillPackage(skillDir)
	if err != nil {
		t.Fatalf("CheckSkillPackage returned error: %v", err)
	}
	if !report.Valid {
		t.Fatal("report.Valid = false, want true")
	}
	if report.InvalidChecks != 0 {
		t.Fatalf("InvalidChecks = %d, want 0", report.InvalidChecks)
	}
}

func createValidSkillFixture(t *testing.T) string {
	t.Helper()

	skillDir := filepath.Join(t.TempDir(), "skills", "yixiaoer")
	dirs := []string{
		skillDir,
		filepath.Join(skillDir, "references"),
		filepath.Join(skillDir, "references", "domains"),
		filepath.Join(skillDir, "references", "workflows"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	skillContent := `---
name: yixiaoer
version: 3.1.1
description: "通过 yxer CLI 操作蚁小二"
metadata:
  requires:
    bins: ["yxer"]
  cliHelp: "yxer --help"
---

# 蚁小二 Skill

## 能力索引

- 发布：[./references/domains/publish.md](./references/domains/publish.md)

## 意图分流

| 用户意图 | 先读 |
| --- | --- |
| 帮我发 | [publish](./references/domains/publish.md) |

## 命令探索

` + "```bash\nyxer --help\n```\n" + `

## 全局规则

- 先读 [shared](./references/yixiaoer-shared.md)
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "QUICKSTART.md"), []byte("# Quickstart\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "references", "yixiaoer-shared.md"), []byte("# Shared\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	workflowDoc := "# Workflow\n"
	if err := os.WriteFile(filepath.Join(skillDir, "references", "workflows", "publish-workflow.md"), []byte(workflowDoc), 0o644); err != nil {
		t.Fatal(err)
	}

	domainContent := `# Domain

适用范围：用户只说“帮我发”时，先进入本域。

## 读取顺序

1. [workflow](../workflows/publish-workflow.md)

## 常用命令

` + "```bash\nyxer doctor\n```\n" + `

## 规则

- 用户明确说本机发布时，改走本机流程
`
	for _, name := range []string{
		"publish.md",
		"accounts-and-env.md",
		"draft-and-material.md",
		"troubleshooting.md",
		"install-and-sync.md",
	} {
		if err := os.WriteFile(filepath.Join(skillDir, "references", "domains", name), []byte(domainContent), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	return skillDir
}
