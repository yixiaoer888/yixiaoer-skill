package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestExecuteWithIOUnknownCommandWritesStructuredErrorToStderr(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := ExecuteWithIO([]string{"definitely-not-a-command"}, &stdout, &stderr)
	if code != yxerrors.ExitValidation {
		t.Fatalf("expected validation exit code, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected stdout to stay empty on error, got %q", stdout.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatalf("stderr should contain JSON error envelope, got %q: %v", stderr.String(), err)
	}
	if response["ok"] != false {
		t.Fatalf("expected ok=false response, got %#v", response)
	}
	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.UsageErr {
		t.Fatalf("expected usage error, got %#v", errorObj)
	}
	if errorObj["message"] != "unknown command" {
		t.Fatalf("expected unknown command message, got %#v", errorObj)
	}
	details := errorObj["details"].(map[string]interface{})
	if _, ok := details["availableCommands"].([]interface{}); !ok {
		t.Fatalf("expected availableCommands in details, got %#v", details)
	}
}

func TestExecuteWithIOArgumentErrorIsValidationError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := ExecuteWithIO([]string{"publish"}, &stdout, &stderr)
	if code != yxerrors.ExitValidation {
		t.Fatalf("expected validation exit code, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected stdout to stay empty on error, got %q", stdout.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatalf("stderr should contain JSON error envelope, got %q: %v", stderr.String(), err)
	}
	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.UsageErr {
		t.Fatalf("expected usage error, got %#v", errorObj)
	}
	if errorObj["type"] != yxerrors.ValidationType {
		t.Fatalf("expected validation type, got %#v", errorObj)
	}
	if errorObj["nextCommand"] != "yxer publish --help" {
		t.Fatalf("expected publish help next command, got %#v", errorObj)
	}
}
