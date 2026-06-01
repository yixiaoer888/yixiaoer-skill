package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
)

func TestUploadDryRunPreviewUsesExplicitFileFlag(t *testing.T) {
	uploadFile = "C:\\tmp\\cover.png"
	uploadURL = ""
	uploadBucket = "cloud-publish"
	uploadDryRun = true
	uploadAutoMeta = true
	t.Cleanup(func() {
		uploadFile = ""
		uploadURL = ""
		uploadBucket = ""
		uploadDryRun = false
		uploadAutoMeta = false
	})

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := uploadCmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	request := data["request"].(map[string]interface{})
	if request["source"] != "C:\\tmp\\cover.png" || request["sourceType"] != "file" {
		t.Fatalf("unexpected dry-run upload preview: %#v", request)
	}
	if request["autoMeta"] != true {
		t.Fatalf("expected autoMeta flag in dry-run upload preview, got %#v", request)
	}
}

func TestMaterialCreateDryRunBuildsPreviewBody(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"filePath":  "material-library/demo.png",
		"fileName":  "demo.png",
		"width":     100,
		"height":    200,
		"type":      "image",
		"thumbPath": "material-library/demo-thumb.png",
	})
	materialCreateDryRun = true
	t.Cleanup(func() {
		materialCreateDryRun = false
	})

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := materialCreateCmd.RunE(cmd, []string{payloadPath}); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	request := data["request"].(map[string]interface{})
	if request["fileName"] != "demo.png" || request["type"] != "image" {
		t.Fatalf("unexpected material create dry-run request: %#v", request)
	}
}

func TestDraftSaveDryRunAddsDraftFlag(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action": "publish",
		"title":  "草稿标题",
	})
	draftDryRun = true
	t.Cleanup(func() {
		draftDryRun = false
	})

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := draftSaveCmd.RunE(cmd, []string{payloadPath}); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	request := data["request"].(map[string]interface{})
	if request["isDraft"] != true {
		t.Fatalf("expected dry-run draft payload to include isDraft=true, got %#v", request)
	}
	if _, ok := request["action"]; ok {
		t.Fatalf("expected action to be removed in dry-run draft payload, got %#v", request)
	}
}
