package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/domain"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

var rootCmd = &cobra.Command{
	Use:     "yxer",
	Short:   "蚁小二多平台内容分发 CLI",
	Version: domain.SkillVersion,
}

func Execute() {
	code := ExecuteWithIO(os.Args[1:], os.Stdout, os.Stderr)
	if code != yxerrors.ExitOK {
		os.Exit(code)
	}
}

func ExecuteWithIO(args []string, stdout, stderr io.Writer) int {
	rootCmd.SetArgs(args)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	cmd, err := rootCmd.ExecuteC()
	if err != nil {
		return output.Error(stderr, structuredCommandError(cmd, err), "run command")
	}
	return yxerrors.ExitOK
}

func init() {
	rootCmd.SetVersionTemplate("yxer {{.Version}}\n")
	rootCmd.PersistentFlags().Bool("json", false, "output JSON")
	rootCmd.PersistentFlags().Bool("debug", false, "show debug logs")
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}

func structuredCommandError(cmd *cobra.Command, err error) error {
	var typed *yxerrors.Error
	if errors.As(err, &typed) {
		return err
	}
	if !strings.Contains(strings.ToLower(err.Error()), "unknown command") {
		if isCommandUsageError(err) {
			return commandUsageError(cmd, err)
		}
		return err
	}
	if cmd == nil {
		cmd = rootCmd
	}
	commandPath := cmd.CommandPath()
	if strings.TrimSpace(commandPath) == "" {
		commandPath = rootCmd.CommandPath()
	}
	if strings.TrimSpace(commandPath) == "" {
		commandPath = "yxer"
	}
	available := availableSubcommandNames(cmd)
	hint := fmt.Sprintf("请运行 %q 查看可用命令。", commandPath+" --help")
	if len(available) > 0 {
		hint = "可用子命令: " + strings.Join(available, ", ")
	}
	return yxerrors.Usage("unknown command", map[string]interface{}{
		"rawError":          err.Error(),
		"commandPath":       commandPath,
		"availableCommands": available,
	}).WithHint(hint).WithNextCommand(commandPath + " --help")
}

func commandUsageError(cmd *cobra.Command, err error) error {
	commandPath := "yxer"
	if cmd != nil && strings.TrimSpace(cmd.CommandPath()) != "" {
		commandPath = cmd.CommandPath()
	}
	return yxerrors.Usage(err.Error(), map[string]interface{}{
		"rawError":    err.Error(),
		"commandPath": commandPath,
	}).WithHint(fmt.Sprintf("请运行 %q 查看参数要求。", commandPath+" --help")).
		WithNextCommand(commandPath + " --help")
}

func isCommandUsageError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	usageFragments := []string{
		"accepts ",
		"requires ",
		"required flag",
		"unknown flag",
		"invalid argument",
		"flag needs an argument",
		"unknown shorthand flag",
	}
	for _, fragment := range usageFragments {
		if strings.Contains(message, fragment) {
			return true
		}
	}
	return false
}

func availableSubcommandNames(cmd *cobra.Command) []string {
	if cmd == nil {
		return nil
	}
	names := make([]string, 0, len(cmd.Commands()))
	for _, child := range cmd.Commands() {
		if child.Hidden || !child.IsAvailableCommand() {
			continue
		}
		name := child.Name()
		if name == "help" || name == "completion" {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
