package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestErrorIncludesMachineFields(t *testing.T) {
	var out bytes.Buffer

	exitCode := Error(&out, yxerrors.Usage("bad payload", []string{"missing title"}).
		WithCategory("validation").
		WithHint("fill title").
		WithNextCommand("yxer validate xhs imageText payload.json"), "run command")
	if exitCode != yxerrors.ExitValidation {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.UsageErr {
		t.Fatalf("unexpected code: %#v", errorObj)
	}
	if errorObj["category"] != "validation" {
		t.Fatalf("unexpected category: %#v", errorObj)
	}
	if errorObj["nextCommand"] != "yxer validate xhs imageText payload.json" {
		t.Fatalf("unexpected nextCommand: %#v", errorObj)
	}
}

func TestErrorNormalizesBareErrorAsInternal(t *testing.T) {
	var out bytes.Buffer

	exitCode := Error(&out, errors.New("boom"), "run command")
	if exitCode != yxerrors.ExitInternal {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.InternalErr {
		t.Fatalf("unexpected code: %#v", errorObj)
	}
	if errorObj["type"] != yxerrors.InternalType {
		t.Fatalf("unexpected type: %#v", errorObj)
	}
	if errorObj["details"] != "boom" {
		t.Fatalf("unexpected details: %#v", errorObj)
	}
}

func TestSuccessDoesNotEscapeAngleBrackets(t *testing.T) {
	var out bytes.Buffer

	err := Success(&out, "schema.get", map[string]interface{}{
		"queryCommands": map[string]string{
			"location": "yxer query locations <account_id> [--query 关键词]",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	got := out.String()
	if strings.Contains(got, "\\u003c") || strings.Contains(got, "\\u003e") {
		t.Fatalf("expected angle brackets to remain literal, got %q", got)
	}
	if !strings.Contains(got, "<account_id>") {
		t.Fatalf("expected literal placeholder in output, got %q", got)
	}
}
