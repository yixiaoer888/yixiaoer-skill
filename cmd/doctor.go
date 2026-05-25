package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "检查本地配置和目录",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		checks := map[string]interface{}{
			"apiUrl":        cfg.APIURL,
			"apiKeyPresent": cfg.APIKey != "",
			"schemaDir":     cfg.SchemaDir,
			"schemaDirOK":   pathExists(cfg.SchemaDir),
			"workflowsOK":   pathExists("workflows"),
		}
		if cfg.APIKey == "" {
			return yxerrors.Auth("Missing YIXIAOER_API_KEY environment variable")
		}
		if !pathExists(cfg.SchemaDir) {
			return yxerrors.Usage("schema directory not found", cfg.SchemaDir)
		}
		return output.Success(cmd.OutOrStdout(), "doctor", checks)
	},
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
