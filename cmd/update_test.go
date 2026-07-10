package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/skillscheck"
)

func TestRunUpdateCheckReportsAutoUpdateCapability(t *testing.T) {
	reset := stubUpdateDependencies(t)
	defer reset()

	detectUpdateInstallMethod = func() updateDetectResult {
		return updateDetectResult{
			Method:       updateInstallMethodNpm,
			ResolvedPath: `C:\npm\node_modules\@yixiaoermail\cli\bin-native\yxer.exe`,
			NpmAvailable: true,
		}
	}

	skillDir := writeTempSkill(t, "3.2.4")
	t.Setenv("YIXIAOER_SKILL_DIR", skillDir)
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(t.TempDir(), "config.json"))

	var stdout bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.Flags().Bool("global", false, "")
	cmd.Flags().Bool("check", true, "")

	if err := runUpdate(cmd); err != nil {
		t.Fatalf("runUpdate returned error: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout should contain JSON, got %q: %v", stdout.String(), err)
	}
	data := response["data"].(map[string]interface{})
	install := data["install"].(map[string]interface{})
	if install["method"] != string(updateInstallMethodNpm) {
		t.Fatalf("install.method = %v, want %q", install["method"], updateInstallMethodNpm)
	}
	if install["autoUpdate"] != true {
		t.Fatalf("install.autoUpdate = %v, want true", install["autoUpdate"])
	}
}

func TestRunUpdateAutoUpdateSyncsSkillAndReportsInstalledVersion(t *testing.T) {
	reset := stubUpdateDependencies(t)
	defer reset()

	detectUpdateInstallMethod = func() updateDetectResult {
		return updateDetectResult{
			Method:       updateInstallMethodNpm,
			ResolvedPath: `C:\npm\node_modules\@yixiaoermail\cli\bin-native\yxer.exe`,
			NpmAvailable: true,
		}
	}

	var npmRan bool
	runNpmGlobalInstall = func(tag string) error {
		npmRan = true
		if tag != "latest" {
			t.Fatalf("tag = %q, want latest", tag)
		}
		return nil
	}
	verifyUpdatedBinary = func() (string, error) {
		return "3.2.5", nil
	}
	prepareSelfReplace = func() (func(), error) {
		return func() {}, nil
	}

	skillDir := writeTempSkill(t, "3.2.4")
	t.Setenv("YIXIAOER_SKILL_DIR", skillDir)
	configPath := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)

	skillsLog := filepath.Join(t.TempDir(), "skills-log.txt")
	skillsScript := writeSkillsStub(t, skillsLog)
	originalPath := os.Getenv("PATH")
	t.Setenv("PATH", filepath.Dir(skillsScript)+string(os.PathListSeparator)+originalPath)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.Flags().Bool("global", false, "")
	cmd.Flags().Bool("check", false, "")

	if err := runUpdate(cmd); err != nil {
		t.Fatalf("runUpdate returned error: %v", err)
	}
	if !npmRan {
		t.Fatal("expected npm update runner to execute")
	}

	stampPath, err := skillscheck.StampPath()
	if err != nil {
		t.Fatalf("StampPath returned error: %v", err)
	}
	rawStamp, err := os.ReadFile(stampPath)
	if err != nil {
		t.Fatalf("expected stamp file to be written: %v", err)
	}
	if got := strings.TrimSpace(string(rawStamp)); got != "3.2.4" {
		t.Fatalf("stamp = %q, want 3.2.4", got)
	}

	logContent, err := os.ReadFile(skillsLog)
	if err != nil {
		t.Fatalf("expected skills stub log: %v", err)
	}
	if !strings.Contains(string(logContent), "add "+skillDir+" -y") {
		t.Fatalf("unexpected skills invocation log: %q", string(logContent))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout should contain JSON, got %q: %v", stdout.String(), err)
	}
	data := response["data"].(map[string]interface{})
	cliUpdate := data["cliUpdate"].(map[string]interface{})
	if cliUpdate["installedVersion"] != "3.2.5" {
		t.Fatalf("installedVersion = %v, want 3.2.5", cliUpdate["installedVersion"])
	}
}

func TestRunUpdateManualInstallFallsBackToSkillSync(t *testing.T) {
	reset := stubUpdateDependencies(t)
	defer reset()

	detectUpdateInstallMethod = func() updateDetectResult {
		return updateDetectResult{
			Method:       updateInstallMethodManual,
			ResolvedPath: `/usr/local/bin/yxer`,
		}
	}

	skillDir := writeTempSkill(t, "3.2.4")
	t.Setenv("YIXIAOER_SKILL_DIR", skillDir)
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(t.TempDir(), "config.json"))

	skillsLog := filepath.Join(t.TempDir(), "skills-log.txt")
	skillsScript := writeSkillsStub(t, skillsLog)
	originalPath := os.Getenv("PATH")
	t.Setenv("PATH", filepath.Dir(skillsScript)+string(os.PathListSeparator)+originalPath)

	var stdout bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.Flags().Bool("global", false, "")
	cmd.Flags().Bool("check", false, "")

	if err := runUpdate(cmd); err != nil {
		t.Fatalf("runUpdate returned error: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout should contain JSON, got %q: %v", stdout.String(), err)
	}
	data := response["data"].(map[string]interface{})
	if data["action"] != "manual_required" {
		t.Fatalf("action = %v, want manual_required", data["action"])
	}
}

func stubUpdateDependencies(t *testing.T) func() {
	t.Helper()
	origDetect := detectUpdateInstallMethod
	origRunNpm := runNpmGlobalInstall
	origVerify := verifyUpdatedBinary
	origPrepare := prepareSelfReplace
	return func() {
		detectUpdateInstallMethod = origDetect
		runNpmGlobalInstall = origRunNpm
		verifyUpdatedBinary = origVerify
		prepareSelfReplace = origPrepare
	}
}

func writeTempSkill(t *testing.T, version string) string {
	t.Helper()
	root := t.TempDir()
	skillDir := filepath.Join(root, "skills", "yixiaoer")
	if err := os.MkdirAll(filepath.Join(skillDir, "references", "domains"), 0o755); err != nil {
		t.Fatalf("MkdirAll skillDir: %v", err)
	}
	content := "---\nname: yixiaoer\nversion: " + version + "\ndescription: test\nmetadata:\n  requires:\n    bins: [yxer]\n  cliHelp: yxer --help\n---\n\n## 能力索引\n\n## 意图分流\n\n## 命令探索\n\n## 全局规则\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile SKILL.md: %v", err)
	}
	return skillDir
}

func writeSkillsStub(t *testing.T, logPath string) string {
	t.Helper()
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "skills")
	script := "#!/bin/sh\nprintf '%s ' \"$@\" > \"" + filepath.ToSlash(logPath) + "\"\n"
	if runtime.GOOS == "windows" {
		scriptPath += ".cmd"
		script = "@echo off\r\nsetlocal\r\n>\"" + logPath + "\" echo %*\r\nexit /b 0\r\n"
	}
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("WriteFile skills stub: %v", err)
	}
	return scriptPath
}
