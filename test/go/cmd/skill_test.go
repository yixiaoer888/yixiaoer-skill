package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestSkillMarkdownLinksExist(t *testing.T) {
	withRepoRoot(t)
	for _, skillPath := range []string{
		filepath.Join("skills", "yixiaoer", "SKILL.md"),
		filepath.Join("skills", "yixiaoer", "references", "yixiaoer-shared.md"),
	} {
		raw, err := os.ReadFile(skillPath)
		if err != nil {
			t.Fatal(err)
		}
		if len(raw) == 0 {
			t.Fatalf("expected %s to be non-empty", skillPath)
		}
	}
}

func TestSkillCheckCommandSuccess(t *testing.T) {
	skillDir := createSkillCheckFixture(t)
	t.Setenv("YIXIAOER_SKILL_DIR", skillDir)

	cmd := testCobraCommand()
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	if err := runSkillCheck(cmd); err != nil {
		t.Fatal(err)
	}

	var response struct {
		OK   bool `json:"ok"`
		Data struct {
			Valid         bool `json:"valid"`
			InvalidChecks int  `json:"invalidChecks"`
			Links         struct {
				InvalidLinks int `json:"invalidLinks"`
			} `json:"links"`
		} `json:"data"`
	}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !response.OK {
		t.Fatal("response.OK = false, want true")
	}
	if !response.Data.Valid {
		t.Fatal("response.Data.Valid = false, want true")
	}
	if response.Data.InvalidChecks != 0 {
		t.Fatalf("InvalidChecks = %d, want 0", response.Data.InvalidChecks)
	}
	if response.Data.Links.InvalidLinks != 0 {
		t.Fatalf("InvalidLinks = %d, want 0", response.Data.Links.InvalidLinks)
	}
}

func TestSkillCheckCommandReturnsStructuredError(t *testing.T) {
	skillDir := filepath.Join(t.TempDir(), "skills", "yixiaoer")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("YIXIAOER_SKILL_DIR", skillDir)

	err := runSkillCheck(testCobraCommand())
	if err == nil {
		t.Fatal("expected skill check error")
	}

	typed, ok := err.(*yxerrors.Error)
	if !ok {
		t.Fatalf("expected *yxerrors.Error, got %T", err)
	}
	if typed.Code != yxerrors.UsageErr {
		t.Fatalf("unexpected error code: %s", typed.Code)
	}
	if typed.Category != "skill_validation" {
		t.Fatalf("unexpected error category: %s", typed.Category)
	}
	if typed.NextCommand != "yxer skill check" {
		t.Fatalf("unexpected next command: %s", typed.NextCommand)
	}
	if typed.Details == nil {
		t.Fatal("expected structured report details on error")
	}
}

func TestResolveSkillsRunnerPrefersSkillsBinary(t *testing.T) {
	binDir := t.TempDir()
	createTestExecutable(t, binDir, "skills")
	createTestExecutable(t, binDir, "npx")
	t.Setenv("PATH", binDir)

	runner, err := resolveSkillsRunner(false, `C:\skill\dir`)
	if err != nil {
		t.Fatalf("resolveSkillsRunner returned error: %v", err)
	}
	if filepath.Base(runner.Path) != executableName("skills") {
		t.Fatalf("runner.Path = %q, want %q", runner.Path, executableName("skills"))
	}
	if len(runner.Args) < 2 || runner.Args[0] != "add" {
		t.Fatalf("runner.Args = %#v, want skills add args", runner.Args)
	}
}

func TestResolveSkillsRunnerFallsBackToNpx(t *testing.T) {
	binDir := t.TempDir()
	createTestExecutable(t, binDir, "npx")
	t.Setenv("PATH", binDir)

	runner, err := resolveSkillsRunner(true, `C:\skill\dir`)
	if err != nil {
		t.Fatalf("resolveSkillsRunner returned error: %v", err)
	}
	if filepath.Base(runner.Path) != executableName("npx") {
		t.Fatalf("runner.Path = %q, want %q", runner.Path, executableName("npx"))
	}
	if len(runner.Args) < 4 || runner.Args[0] != "-y" || runner.Args[1] != "skills" {
		t.Fatalf("runner.Args = %#v, want npx skills args", runner.Args)
	}
}

func createSkillCheckFixture(t *testing.T) string {
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
	if err := os.WriteFile(filepath.Join(skillDir, "references", "workflows", "publish-workflow.md"), []byte("# Workflow\n"), 0o644); err != nil {
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

func createTestExecutable(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, executableName(name))
	content := []byte("")
	if runtime.GOOS == "windows" {
		content = []byte("@echo off\r\n")
	}
	if err := os.WriteFile(path, content, 0o755); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", path, err)
	}
}

func executableName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".cmd"
	}
	return name
}
