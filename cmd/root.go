package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	rootCmd.SetVersionTemplate(`{"ok":true,"action":"version","version":"{{.Version}}","data":{"version":"{{.Version}}"}}` + "\n")
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = output.Success(cmd.OutOrStdout(), "help", commandHelpData(cmd))
	})
	rootCmd.PersistentFlags().Bool("json", false, "output JSON")
	rootCmd.PersistentFlags().Bool("debug", false, "show debug logs")
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}

func commandHelpData(cmd *cobra.Command) map[string]interface{} {
	if cmd == nil {
		cmd = rootCmd
	}
	data := map[string]interface{}{
		"commandPath": cmd.CommandPath(),
		"use":         cmd.UseLine(),
		"short":       cmd.Short,
		"subcommands": availableSubcommandNames(cmd),
		"flags":       visibleFlagNames(cmd),
	}
	if cmd.Long != "" {
		data["long"] = cmd.Long
	}
	return data
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
	unknown := unknownCommandName(err.Error())
	suggestions := commandSuggestions(unknown, available)
	helpCommand := commandPath + " --help"
	nextCommand := helpCommand
	hint := fmt.Sprintf("请运行 %q 查看可用命令。", helpCommand)
	if len(suggestions) > 0 {
		suggestedCommands := qualifySubcommands(commandPath, suggestions)
		nextCommand = suggestedCommands[0]
		hint = fmt.Sprintf("未知命令 %q。你可能想运行: %s；也可运行 %q 查看可用命令。", unknown, strings.Join(suggestedCommands, ", "), helpCommand)
	} else if len(available) > 0 {
		hint = fmt.Sprintf("未知命令 %q。可用子命令: %s；也可运行 %q 查看帮助。", unknown, strings.Join(available, ", "), helpCommand)
	}
	return yxerrors.Usage("unknown command", map[string]interface{}{
		"rawError":          err.Error(),
		"commandPath":       commandPath,
		"unknownCommand":    unknown,
		"availableCommands": available,
		"suggestions":       suggestions,
		"helpCommand":       helpCommand,
	}).WithHint(hint).WithNextCommand(nextCommand)
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

func unknownCommandName(message string) string {
	start := strings.Index(message, `"`)
	if start < 0 {
		return ""
	}
	rest := message[start+1:]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return ""
	}
	return rest[:end]
}

func commandSuggestions(input string, available []string) []string {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" || len(available) == 0 {
		return nil
	}
	type scored struct {
		name  string
		score int
	}
	var matches []scored
	for _, name := range available {
		lower := strings.ToLower(name)
		score := commandSuggestionScore(input, lower)
		if score < 0 {
			continue
		}
		matches = append(matches, scored{name: name, score: score})
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].score == matches[j].score {
			return matches[i].name < matches[j].name
		}
		return matches[i].score < matches[j].score
	})
	limit := 3
	if len(matches) < limit {
		limit = len(matches)
	}
	result := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, matches[i].name)
	}
	return result
}

func commandSuggestionScore(input, candidate string) int {
	if input == candidate {
		return 0
	}
	if strings.TrimSuffix(candidate, "s") == input || strings.TrimSuffix(input, "s") == candidate {
		return 1
	}
	if strings.HasPrefix(candidate, input) || strings.HasPrefix(input, candidate) {
		return 2
	}
	distance := levenshteinDistance(input, candidate)
	threshold := len(input) / 3
	if threshold < 2 {
		threshold = 2
	}
	if distance <= threshold {
		return 10 + distance
	}
	return -1
}

func qualifySubcommands(commandPath string, names []string) []string {
	commands := make([]string, 0, len(names))
	for _, name := range names {
		commands = append(commands, strings.TrimSpace(commandPath+" "+name))
	}
	return commands
}

func levenshteinDistance(a, b string) int {
	ar := []rune(a)
	br := []rune(b)
	if len(ar) == 0 {
		return len(br)
	}
	if len(br) == 0 {
		return len(ar)
	}
	prev := make([]int, len(br)+1)
	curr := make([]int, len(br)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ar); i++ {
		curr[0] = i
		for j := 1; j <= len(br); j++ {
			cost := 0
			if ar[i-1] != br[j-1] {
				cost = 1
			}
			curr[j] = minInt(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(br)]
}

func minInt(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, value := range values[1:] {
		if value < min {
			min = value
		}
	}
	return min
}

func visibleFlagNames(cmd *cobra.Command) []string {
	if cmd == nil {
		return nil
	}
	seen := map[string]bool{}
	var names []string
	addFlag := func(name string) {
		if strings.TrimSpace(name) == "" || seen[name] {
			return
		}
		seen[name] = true
		names = append(names, name)
	}
	cmd.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			addFlag(flag.Name)
		}
	})
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			addFlag(flag.Name)
		}
	})
	sort.Strings(names)
	return names
}
