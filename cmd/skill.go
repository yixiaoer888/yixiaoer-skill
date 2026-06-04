package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/skillscheck"
)

func init() {
	rootCmd.AddCommand(skillCmd)
	skillCmd.AddCommand(skillShowCmd)
	skillCmd.AddCommand(skillCheckCmd)
	skillCmd.AddCommand(skillSyncCmd)
	skillSyncCmd.Flags().Bool("global", false, "install skill globally")
}

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "显示 AI agent 技能安装与同步信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSkillShow(cmd)
	},
}

var skillShowCmd = &cobra.Command{
	Use:   "show",
	Short: "输出当前项目技能包位置和安装命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSkillShow(cmd)
	},
}

var skillSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "安装或同步当前项目的 AI skill，并写入本地版本戳",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSkillSync(cmd)
	},
}

var skillCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "检查当前 skill 包内 Markdown 链接是否都能解析",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSkillCheck(cmd)
	},
}

func runSkillShow(cmd *cobra.Command) error {
	skillDir, err := skillscheck.DetectSkillDir()
	if err != nil {
		return err
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	status, err := skillscheck.Check(rootCmd.Version)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"name":         "yixiaoer",
		"version":      rootCmd.Version,
		"skillDir":     skillDir,
		"skillFile":    skillFile,
		"skillsStatus": status,
		"installMode":  "cli-first-skill-second",
		"summary": []string{
			"先安装 yxer CLI，再安装 yixiaoer skill，供 AI agent 读取技能后调用 yxer。",
			"技能文件负责约束工作流与命令选择，真正执行统一走 yxer CLI。",
		},
		"install": map[string]interface{}{
			"local": []string{
				"npx skills add \"" + skillDir + "\" -y",
			},
			"global": []string{
				"npx skills add \"" + skillDir + "\" -g -y",
			},
		},
		"sync": map[string]interface{}{
			"when": []string{
				"更新 yxer 版本后",
				"SKILL.md 或 references/workflows 发生变化后",
			},
			"commands": []string{
				"npx skills add \"" + skillDir + "\" -y",
				"npx skills add \"" + skillDir + "\" -g -y",
			},
		},
		"entrypoints": []string{
			"yxer doctor",
			"yxer config get",
			"yxer accounts",
			"yxer upload",
			"yxer validate",
			"yxer publish",
		},
	}

	return output.Success(cmd.OutOrStdout(), "skill.show", data)
}

func runSkillSync(cmd *cobra.Command) error {
	skillDir, err := skillscheck.DetectSkillDir()
	if err != nil {
		return err
	}
	globalInstall, _ := cmd.Flags().GetBool("global")

	return syncSkill(cmd, skillDir, globalInstall)
}

func runSkillCheck(cmd *cobra.Command) error {
	skillDir, err := skillscheck.DetectSkillDir()
	if err != nil {
		return err
	}
	report, checkErr := skillscheck.CheckSkillLinks(skillDir)
	if checkErr != nil {
		return fmt.Errorf("%w: %+v", checkErr, report.Issues)
	}
	return output.Success(cmd.OutOrStdout(), "skill.check", report)
}

func syncSkill(cmd *cobra.Command, skillDir string, globalInstall bool) error {
	npxPath, err := exec.LookPath("npx")
	if err != nil {
		return fmt.Errorf("npx not found in PATH; install Node.js and ensure npx is available")
	}

	args := []string{"-y", "skills", "add", skillDir, "-y"}
	if globalInstall {
		args = []string{"-y", "skills", "add", skillDir, "-g", "-y"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	execCmd := exec.CommandContext(ctx, npxPath, args...)
	execCmd.Stdout = cmd.OutOrStdout()
	execCmd.Stderr = cmd.ErrOrStderr()
	if err := execCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("skills sync timed out after 2m")
		}
		return err
	}

	if err := skillscheck.WriteStamp(rootCmd.Version); err != nil {
		return err
	}

	if cmd.Name() == "sync" && cmd.Parent() != nil && cmd.Parent().Name() == "skill" {
		data := map[string]interface{}{
			"name":      "yixiaoer",
			"version":   rootCmd.Version,
			"skillDir":  skillDir,
			"installed": true,
			"global":    globalInstall,
			"command":   append([]string{npxPath}, args...),
		}
		return output.Success(cmd.OutOrStdout(), "skill.sync", data)
	}

	return nil
}
