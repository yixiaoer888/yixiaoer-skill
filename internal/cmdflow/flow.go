package cmdflow

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type Result struct {
	Action string
	Data   interface{}
}

type Flow struct {
	Validate func() error
	DryRun   func() (Result, error)
	Execute  func() (Result, error)
}

func Run(cmd *cobra.Command, dryRun bool, flow Flow) error {
	if flow.Validate != nil {
		if err := flow.Validate(); err != nil {
			return err
		}
	}

	if dryRun {
		if flow.DryRun == nil {
			return yxerrors.Usage("dry-run is not supported for this command", nil)
		}
		result, err := flow.DryRun()
		if err != nil {
			return err
		}
		return output.Success(cmd.OutOrStdout(), result.Action, result.Data)
	}

	if flow.Execute == nil {
		return yxerrors.Internal("command execute handler is not configured", nil)
	}
	result, err := flow.Execute()
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), result.Action, result.Data)
}
