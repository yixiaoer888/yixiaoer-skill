package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/domain"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
)

var rootCmd = &cobra.Command{
	Use:     "yxer",
	Short:   "蚁小二多平台内容分发 CLI",
	Version: domain.SkillVersion,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Error(os.Stderr, err, "run command")
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("yxer {{.Version}}\n")
	rootCmd.PersistentFlags().Bool("json", false, "output JSON")
	rootCmd.PersistentFlags().Bool("debug", false, "show debug logs")
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}

func wantJSON(cmd *cobra.Command) bool {
	value, _ := cmd.Flags().GetBool("json")
	if value {
		return true
	}
	value, _ = cmd.Root().PersistentFlags().GetBool("json")
	return value
}

func usageErr(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
