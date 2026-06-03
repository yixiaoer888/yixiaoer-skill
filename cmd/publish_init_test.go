package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestPublishInitCommandWritesTemplateFile(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)

	outputPath := filepath.Join(t.TempDir(), "douyin-video-payload.json")
	publishInitOutput = outputPath
	t.Cleanup(func() {
		publishInitOutput = ""
	})

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := publishInitCmd.RunE(cmd, []string{"抖音", "video"}); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["action"] != "publish" || payload["publishType"] != "video" {
		t.Fatalf("unexpected top-level template: %#v", payload)
	}
	args := payload["publishArgs"].(map[string]interface{})
	forms := args["accountForms"].([]interface{})
	form := forms[0].(map[string]interface{})
	if form["platformAccountId"] != "<platformAccountId>" {
		t.Fatalf("expected placeholder platformAccountId, got %#v", form["platformAccountId"])
	}
	cpf := form["contentPublishForm"].(map[string]interface{})
	if cpf["formType"] == nil || cpf["title"] == nil {
		t.Fatalf("expected required schema fields in template, got %#v", cpf)
	}
}

func TestPublishInitCommandPlacesArticleContentUnderPublishArgs(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)

	outputPath := filepath.Join(t.TempDir(), "zhihu-article-payload.json")
	publishInitOutput = outputPath
	t.Cleanup(func() {
		publishInitOutput = ""
	})

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := publishInitCmd.RunE(cmd, []string{"知乎", "article"}); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["desc"] == nil {
		t.Fatalf("expected article template to expose top-level desc, got %#v", payload)
	}
	args := payload["publishArgs"].(map[string]interface{})
	if args["content"] == nil {
		t.Fatalf("expected article content under publishArgs, got %#v", args)
	}
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	cpf := form["contentPublishForm"].(map[string]interface{})
	if _, exists := cpf["content"]; exists {
		t.Fatalf("did not expect article content inside contentPublishForm template, got %#v", cpf)
	}
}
