package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/skillscheck"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const (
	npmPackageName      = "@yixiaoermail/cli"
	npmInstallTimeout   = 10 * time.Minute
	verifyBinaryTimeout = 10 * time.Second
)

type updateInstallMethod string

const (
	updateInstallMethodNpm    updateInstallMethod = "npm"
	updateInstallMethodManual updateInstallMethod = "manual"
)

type updateDetectResult struct {
	Method       updateInstallMethod `json:"method"`
	ResolvedPath string              `json:"resolvedPath,omitempty"`
	NpmAvailable bool                `json:"npmAvailable"`
}

func (r updateDetectResult) CanAutoUpdate() bool {
	return r.Method == updateInstallMethodNpm && r.NpmAvailable
}

func (r updateDetectResult) ManualReason() string {
	if r.Method == updateInstallMethodNpm && !r.NpmAvailable {
		return "installed via npm, but npm is not available in PATH"
	}
	return "not installed via npm"
}

var (
	detectUpdateInstallMethod = defaultDetectUpdateInstallMethod
	runNpmGlobalInstall       = defaultRunNpmGlobalInstall
	verifyUpdatedBinary       = defaultVerifyUpdatedBinary
	prepareSelfReplace        = defaultPrepareSelfReplace
)

func init() {
	rootCmd.AddCommand(newUpdateCmd())
}

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "检查当前 CLI/skill 状态，并在 npm 安装场景下执行官方升级",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd)
		},
	}
	cmd.Flags().Bool("global", false, "install skill globally")
	cmd.Flags().Bool("check", false, "only check update status without syncing skill")
	return cmd
}

func runUpdate(cmd *cobra.Command) error {
	checkOnly, _ := cmd.Flags().GetBool("check")
	globalInstall, _ := cmd.Flags().GetBool("global")

	skillDir, err := skillscheck.DetectSkillDir()
	if err != nil {
		return err
	}
	skillVersion, err := skillscheck.SkillVersion(skillDir)
	if err != nil {
		return err
	}
	before, err := skillscheck.Check(skillVersion)
	if err != nil {
		return err
	}

	detect := detectUpdateInstallMethod()
	data := map[string]interface{}{
		"skillVersion": skillVersion,
		"cliVersion":   rootCmd.Version,
		"skillDir":     skillDir,
		"before":       before,
		"install": map[string]interface{}{
			"method":       string(detect.Method),
			"resolvedPath": detect.ResolvedPath,
			"npmAvailable": detect.NpmAvailable,
			"autoUpdate":   detect.CanAutoUpdate(),
		},
		"cliUpdate": map[string]interface{}{
			"supported": detect.CanAutoUpdate(),
			"command":   "yxer update",
			"fallbackCommands": []string{
				"npm install -g @yixiaoermail/cli@latest",
				"yxer --version",
				"yxer skill sync",
			},
			"changelog": "CHANGELOG.md",
		},
	}

	if checkOnly {
		data["action"] = "checked"
		if !detect.CanAutoUpdate() {
			data["cliUpdate"].(map[string]interface{})["message"] = detect.ManualReason()
		} else {
			data["cliUpdate"].(map[string]interface{})["message"] = "CLI was installed via npm and can be updated with `yxer update`."
		}
		return output.Success(cmd.OutOrStdout(), "update", data)
	}

	if !detect.CanAutoUpdate() {
		if err := syncSkill(cmd, skillDir, globalInstall); err != nil {
			return err
		}
		after, err := skillscheck.Check(skillVersion)
		if err != nil {
			return err
		}
		data["action"] = "manual_required"
		data["skillSync"] = map[string]interface{}{
			"ran":    true,
			"global": globalInstall,
		}
		data["after"] = after
		data["cliUpdate"].(map[string]interface{})["message"] = detect.ManualReason()
		return output.Success(cmd.OutOrStdout(), "update", data)
	}

	restore, err := prepareSelfReplace()
	if err != nil {
		return yxerrors.Internal("failed to prepare CLI self-update", err.Error()).
			WithCategory("update_error").
			WithHint("请关闭其他正在占用 yxer 的进程后重试；如果仍失败，可改用 `npm install -g @yixiaoermail/cli@latest`。").
			WithNextCommand("npm install -g @yixiaoermail/cli@latest")
	}

	if err := runNpmGlobalInstall("latest"); err != nil {
		restore()
		return yxerrors.Remote("failed to update CLI via npm", map[string]interface{}{
			"package": npmPackageName,
			"tag":     "latest",
			"error":   err.Error(),
		}).WithCategory("update_error").
			WithRetryable(true).
			WithHint("请检查 npm 网络与全局安装权限，然后重试 `yxer update`；也可先手动运行 `npm install -g @yixiaoermail/cli@latest`。").
			WithNextCommand("yxer update")
	}

	installedVersion, err := verifyUpdatedBinary()
	if err != nil {
		restore()
		return yxerrors.Remote("updated CLI binary verification failed", map[string]interface{}{
			"error": err.Error(),
		}).WithCategory("update_error").
			WithRetryable(true).
			WithHint("新版本已安装但验证失败；请先运行 `yxer --version` 检查当前全局命令，再按需重装 npm 包。").
			WithNextCommand("yxer --version")
	}

	updatedSkillDir, err := skillscheck.DetectSkillDir()
	if err != nil {
		return err
	}
	updatedSkillVersion, err := skillscheck.SkillVersion(updatedSkillDir)
	if err != nil {
		return err
	}

	if err := syncSkill(cmd, updatedSkillDir, globalInstall); err != nil {
		return err
	}

	after, err := skillscheck.Check(updatedSkillVersion)
	if err != nil {
		return err
	}

	data["action"] = "updated"
	data["skillVersion"] = updatedSkillVersion
	data["skillDir"] = updatedSkillDir
	data["skillSync"] = map[string]interface{}{
		"ran":    true,
		"global": globalInstall,
	}
	data["after"] = after
	data["cliUpdate"] = map[string]interface{}{
		"supported":        true,
		"command":          "yxer update",
		"installedVersion": installedVersion,
		"message":          "CLI update completed via npm and the packaged skill was re-synced.",
		"fallbackCommands": []string{
			"yxer --version",
			"yxer skill sync",
		},
		"changelog": "CHANGELOG.md",
	}

	return output.Success(cmd.OutOrStdout(), "update", data)
}

func defaultDetectUpdateInstallMethod() updateDetectResult {
	exe, err := os.Executable()
	if err != nil {
		return updateDetectResult{Method: updateInstallMethodManual}
	}
	resolved := exe
	if path, evalErr := filepath.EvalSymlinks(exe); evalErr == nil {
		resolved = path
	}

	result := updateDetectResult{
		Method:       updateInstallMethodManual,
		ResolvedPath: resolved,
	}
	if strings.Contains(strings.ToLower(resolved), "node_modules") {
		result.Method = updateInstallMethodNpm
	}
	if result.Method == updateInstallMethodNpm {
		_, err = exec.LookPath("npm")
		result.NpmAvailable = err == nil
	}
	return result
}

func defaultRunNpmGlobalInstall(tag string) error {
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), npmInstallTimeout)
	defer cancel()

	command := exec.CommandContext(ctx, npmPath, "install", "-g", npmPackageName+"@"+tag)
	command.Stdout = os.Stderr
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("npm install timed out after 10m")
		}
		return err
	}
	return nil
}

func defaultVerifyUpdatedBinary() (string, error) {
	exePath, err := exec.LookPath("yxer")
	if err != nil {
		exePath, err = os.Executable()
		if err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), verifyBinaryTimeout)
	defer cancel()

	var stdout bytes.Buffer
	command := exec.CommandContext(ctx, exePath, "--version")
	command.Stdout = &stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.New("binary verification timed out after 10s")
		}
		return "", err
	}

	var payload struct {
		OK      bool                   `json:"ok"`
		Action  string                 `json:"action"`
		Version string                 `json:"version"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		return "", err
	}
	if !payload.OK || payload.Action != "version" || strings.TrimSpace(payload.Version) == "" {
		return "", errors.New("unexpected version response shape")
	}
	return payload.Version, nil
}

func defaultPrepareSelfReplace() (func(), error) {
	if runtime.GOOS != "windows" {
		return func() {}, nil
	}

	exe, err := os.Executable()
	if err != nil {
		return func() {}, nil
	}
	resolved := exe
	if path, evalErr := filepath.EvalSymlinks(exe); evalErr == nil {
		resolved = path
	}
	oldPath := resolved + ".old"
	_ = os.Remove(oldPath)
	if err := os.Rename(resolved, oldPath); err != nil {
		return func() {}, err
	}

	restore := func() {
		if _, err := os.Stat(oldPath); err != nil {
			return
		}
		_ = os.Remove(resolved)
		_ = os.Rename(oldPath, resolved)
	}
	return restore, nil
}
