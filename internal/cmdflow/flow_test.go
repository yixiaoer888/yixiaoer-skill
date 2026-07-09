package cmdflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunStopsWhenValidateFails(t *testing.T) {
	wantErr := errors.New("validation failed")
	dryRunCalled := false
	executeCalled := false

	err := Run(testCommand(), false, Flow{
		Validate: func() error { return wantErr },
		DryRun: func() (Result, error) {
			dryRunCalled = true
			return Result{}, nil
		},
		Execute: func() (Result, error) {
			executeCalled = true
			return Result{}, nil
		},
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected validate error, got %v", err)
	}
	if dryRunCalled {
		t.Fatal("dry-run handler should not be called after validate failure")
	}
	if executeCalled {
		t.Fatal("execute handler should not be called after validate failure")
	}
}

func TestRunDryRunSkipsExecute(t *testing.T) {
	cmd, stdout := testCommandWithOutput()
	executeCalled := false

	err := Run(cmd, true, Flow{
		Validate: func() error { return nil },
		DryRun: func() (Result, error) {
			return Result{Action: "example.dry-run", Data: map[string]interface{}{"dryRun": true}}, nil
		},
		Execute: func() (Result, error) {
			executeCalled = true
			return Result{Action: "example", Data: nil}, nil
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executeCalled {
		t.Fatal("execute handler should not be called during dry-run")
	}
	envelope := decodeEnvelope(t, stdout.Bytes())
	if envelope["action"] != "example.dry-run" {
		t.Fatalf("unexpected action: %v", envelope["action"])
	}
	data, _ := envelope["data"].(map[string]interface{})
	if data["dryRun"] != true {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestRunExecuteWritesSuccessEnvelope(t *testing.T) {
	cmd, stdout := testCommandWithOutput()

	err := Run(cmd, false, Flow{
		Validate: func() error { return nil },
		Execute: func() (Result, error) {
			return Result{Action: "example", Data: map[string]interface{}{"ok": "yes"}}, nil
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	envelope := decodeEnvelope(t, stdout.Bytes())
	if envelope["action"] != "example" {
		t.Fatalf("unexpected action: %v", envelope["action"])
	}
	data, _ := envelope["data"].(map[string]interface{})
	if data["ok"] != "yes" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func testCommand() *cobra.Command {
	cmd, _ := testCommandWithOutput()
	return cmd
}

func testCommandWithOutput() (*cobra.Command, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(stdout)
	return cmd, stdout
}

func decodeEnvelope(t *testing.T, raw []byte) map[string]interface{} {
	t.Helper()
	var envelope map[string]interface{}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("decode envelope: %v\n%s", err, string(raw))
	}
	return envelope
}
