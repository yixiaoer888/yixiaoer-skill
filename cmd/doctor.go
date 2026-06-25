package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/skillscheck"
	doctorflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/doctor"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "检查本地配置和目录",
	RunE: func(cmd *cobra.Command, args []string) error {
		rt, err := app.Load()
		if err != nil {
			return err
		}
		checks, err := doctorflow.NewService(rt).Check()
		if err != nil {
			return err
		}
		skillDir, err := skillscheck.DetectSkillDir()
		if err != nil {
			return output.SuccessWithNotice(cmd.OutOrStdout(), "doctor", checks, map[string]interface{}{
				"skills": map[string]interface{}{
					"type":    "skills_path_unresolved",
					"target":  "unknown",
					"state":   "unknown",
					"message": `未能自动定位 "skills/yixiaoer" 目录；如需检查或同步 skill，请设置 YIXIAOER_SKILL_DIR。`,
				},
			})
		}
		skillVersion, err := skillscheck.SkillVersion(skillDir)
		if err != nil {
			return err
		}
		notice, err := skillscheck.Notice(skillVersion, skillDir)
		if err != nil {
			return err
		}
		return output.SuccessWithNotice(cmd.OutOrStdout(), "doctor", checks, map[string]interface{}{
			"skills": notice,
		})
	},
}

func pathExists(path string) bool {
	return doctorflow.PathExists(path)
}
