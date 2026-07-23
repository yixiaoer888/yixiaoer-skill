package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPublishFormSessionCanSetNestedAndArrayValues(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	initial, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := setJSONPath(initial.Payload, "publishArgs.accountForms[0].contentPublishForm.title", "页面标题"); err != nil {
		t.Fatalf("direct path update: %v", err)
	}

	setTitle := newPublishFormSetCmd()
	valuePath := filepath.Join(t.TempDir(), "title.json")
	if err := os.WriteFile(valuePath, []byte(`"页面标题"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := setTitle.Flags().Set("value-file", valuePath); err != nil {
		t.Fatal(err)
	}
	setTitle.SetArgs([]string{sessionPath, "publishArgs.accountForms[0].contentPublishForm.title"})
	parsed, parseErr := readFormValue("", valuePath)
	if parseErr != nil {
		t.Fatal(parseErr)
	}
	checkSession, readErr := readPublishFormSession(sessionPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if pathErr := setJSONPath(cloneJSONMap(checkSession.Payload), "publishArgs.accountForms[0].contentPublishForm.title", parsed); pathErr != nil {
		t.Fatalf("same flow path update: %v", pathErr)
	}
	setTitle.SetOut(&bytes.Buffer{})
	if err := setTitle.Execute(); err != nil {
		t.Logf("set error: %#v", err)
		t.Fatal(err)
	}

	setTag := newPublishFormSetCmd()
	tagPath := filepath.Join(t.TempDir(), "tag.json")
	if err := os.WriteFile(tagPath, []byte(`"标签"`), 0o644); err != nil {
		t.Fatal(err)
	}
	setTag.SetArgs([]string{sessionPath, "publishArgs.accountForms[0].contentPublishForm.tags[0]", "--value-file", tagPath})
	setTag.SetOut(&bytes.Buffer{})
	if err := setTag.Execute(); err != nil {
		t.Fatal(err)
	}

	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	form := session.Payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if form["title"] != "页面标题" {
		t.Fatalf("unexpected title: %#v", form["title"])
	}
	if form["tags"].([]interface{})[0] != "标签" {
		t.Fatalf("unexpected tag: %#v", form["tags"])
	}
}

func TestPublishFormStartDryRunDoesNotWrite(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	var out bytes.Buffer
	cmd := newPublishFormStartCmd()
	cmd.SetArgs([]string{"抖音", "video", "--output", sessionPath, "--dry-run"})
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sessionPath); !os.IsNotExist(err) {
		t.Fatalf("expected dry-run not to write session, stat error=%v", err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "publish.form.start.dry-run" {
		t.Fatalf("unexpected action: %#v", response["action"])
	}
}

func TestPublishFormExportProducesStandardPayload(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	payloadPath := filepath.Join(t.TempDir(), "payload.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"知乎", "article", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	export := newPublishFormExportCmd()
	export.SetArgs([]string{sessionPath, "--output", payloadPath})
	export.SetOut(&bytes.Buffer{})
	if err := export.Execute(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(payloadPath)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["action"] != "publish" || payload["publishType"] != "article" {
		t.Fatalf("unexpected exported payload: %#v", payload)
	}
}

func TestReadFormValueAcceptsPlainTextAndCLIEnvelope(t *testing.T) {
	value, err := readFormValue("页面标题", "")
	if err != nil || value != "页面标题" {
		t.Fatalf("expected plain text value, got %#v, err=%v", value, err)
	}
	value, err = readFormValue(`{"ok":true,"data":{"yixiaoerId":"id_1"}}`, "")
	if err != nil {
		t.Fatal(err)
	}
	object := value.(map[string]interface{})
	if object["yixiaoerId"] != "id_1" {
		t.Fatalf("expected CLI data payload, got %#v", object)
	}
}
