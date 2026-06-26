package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/skillscheck"
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("global", false, "install skill globally")
	updateCmd.Flags().Bool("check", false, "only check update status without syncing skill")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "检查当前 CLI/skill 状态，并同步 AI skill",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdate(cmd)
	},
}

func runUpdate(cmd *cobra.Command) error {
	skillDir, err := skillscheck.DetectSkillDir()
	if err != nil {
		return err
	}

	checkOnly, _ := cmd.Flags().GetBool("check")
	globalInstall, _ := cmd.Flags().GetBool("global")

	skillVersion, err := skillscheck.SkillVersion(skillDir)
	if err != nil {
		return err
	}

	before, err := skillscheck.Check(skillVersion)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"skillVersion": skillVersion,
		"cliVersion":   rootCmd.Version,
		"skillDir":     skillDir,
		"before":       before,
		"cliUpdate": map[string]interface{}{
			"supported": false,
			"message":   "当前命令会检查 CLI/skill 状态并同步 skill；如果 CLI 通过 npm 安装，推荐使用 npm 全局更新到最新版本。",
			"commands": []string{
				"npm install -g @yixiaoermail/cli@latest",
				"yxer --version",
				"yxer skill sync",
			},
			"changelog": "CHANGELOG.md",
		},
	}

	if checkOnly {
		data["action"] = "checked"
		return output.Success(cmd.OutOrStdout(), "update", data)
	}

	if err := syncSkill(cmd, skillDir, globalInstall); err != nil {
		return err
	}

	after, err := skillscheck.Check(skillVersion)
	if err != nil {
		return err
	}

	data["action"] = "updated"
	data["skillSync"] = map[string]interface{}{
		"ran":    true,
		"global": globalInstall,
	}
	data["after"] = after

	return output.Success(cmd.OutOrStdout(), "update", data)
}
