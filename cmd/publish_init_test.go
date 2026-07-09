package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPublishInitCommandWritesTemplateFile(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)

	outputPath := filepath.Join(t.TempDir(), "douyin-video-payload.json")

	var out bytes.Buffer
	cmd := newPublishInitCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"抖音", "video", "--output", outputPath})

	if err := cmd.Execute(); err != nil {
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
	video := form["video"].(map[string]interface{})
	if video["duration"] == nil {
		t.Fatalf("expected video resource placeholder with duration, got %#v", video)
	}
	if form["cover"] == nil || form["coverKey"] == nil {
		t.Fatalf("expected account-level cover and coverKey placeholders, got %#v", form)
	}
	cpf := form["contentPublishForm"].(map[string]interface{})
	if cpf["formType"] == nil || cpf["title"] == nil {
		t.Fatalf("expected required schema fields in template, got %#v", cpf)
	}
	if _, exists := cpf["video"]; exists {
		t.Fatalf("did not expect video resource under contentPublishForm template, got %#v", cpf)
	}
}

func TestPublishInitCommandPlacesArticleContentUnderPublishArgs(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)

	outputPath := filepath.Join(t.TempDir(), "zhihu-article-payload.json")

	var out bytes.Buffer
	cmd := newPublishInitCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"知乎", "article", "--output", outputPath})

	if err := cmd.Execute(); err != nil {
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
