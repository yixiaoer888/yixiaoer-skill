package cmd

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
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

func TestUploadFlagDefaultEnablesAutoMeta(t *testing.T) {
	if uploadCmd.Flag("auto-meta").DefValue != "true" {
		t.Fatalf("expected upload --auto-meta default to be true, got %q", uploadCmd.Flag("auto-meta").DefValue)
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

func TestMaterialAddDryRunExtractsImageMetadata(t *testing.T) {
	withRepoRoot(t)
	imagePath := filepath.Join(t.TempDir(), "material.png")
	if err := os.WriteFile(imagePath, testPNGBytesWithSize(t, 7, 9), 0o644); err != nil {
		t.Fatal(err)
	}
	materialFilePath = imagePath
	materialThumbPath = ""
	materialType = ""
	materialDryRun = true
	t.Cleanup(func() {
		materialFilePath = ""
		materialThumbPath = ""
		materialType = ""
		materialDryRun = false
	})

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := materialAddCmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	request := data["request"].(map[string]interface{})
	upload := data["upload"].(map[string]interface{})
	if request["width"] != float64(7) || request["height"] != float64(9) {
		t.Fatalf("expected material add dry-run to extract dimensions, got request=%#v", request)
	}
	if upload["width"] != float64(7) || upload["height"] != float64(9) {
		t.Fatalf("expected upload preview dimensions, got upload=%#v", upload)
	}
	if request["type"] != "image" {
		t.Fatalf("expected inferred image type, got request=%#v", request)
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

func testPNGBytesWithSize(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{R: 20, G: 120, B: 200, A: 255})
		}
	}
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
