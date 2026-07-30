package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestUploadDryRunPreviewUsesExplicitFileFlag(t *testing.T) {
	var out bytes.Buffer
	cmd := newUploadCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--file", "C:\\tmp\\cover.png", "--dry-run"})

	if err := cmd.Execute(); err != nil {
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

func TestUploadDryRunShowsShipinhaoCoverCompressionLimit(t *testing.T) {
	var out bytes.Buffer
	cmd := newUploadCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--file", "C:\\tmp\\cover.png", "--platform", "视频号", "--usage", "cover", "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	processing := data["mediaProcessing"].(map[string]interface{})
	if processing["maxImageBytes"] != float64(512*1024) {
		t.Fatalf("expected shipinhao cover compression limit, got %#v", processing)
	}
}

func TestUploadFlagDefaultEnablesAutoMeta(t *testing.T) {
	cmd := newUploadCmd()
	if cmd.Flag("auto-meta").DefValue != "true" {
		t.Fatalf("expected upload --auto-meta default to be true, got %q", cmd.Flag("auto-meta").DefValue)
	}
}

func TestUploadRejectsResourceTypeSubcommand(t *testing.T) {
	_, err := resolveUploadSource([]string{"video", "demo.mp4"}, uploadOptions{})
	if err == nil {
		t.Fatal("expected upload source error")
	}
	if !strings.Contains(err.Error(), "upload accepts exactly one file path or URL") {
		t.Fatalf("expected hint to mention incorrect upload form, got %v", err)
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured yxerror, got %T", err)
	}
	if strings.Contains(typed.Hint, "upload video") {
		t.Fatalf("expected generic upload hint, got %q", typed.Hint)
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
	var out bytes.Buffer
	cmd := newMaterialCreateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{payloadPath, "--dry-run"})

	if err := cmd.Execute(); err != nil {
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
	var out bytes.Buffer
	cmd := newMaterialAddCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--file", imagePath, "--dry-run"})

	if err := cmd.Execute(); err != nil {
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
	var out bytes.Buffer
	cmd := newDraftSaveCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{payloadPath, "--dry-run"})

	if err := cmd.Execute(); err != nil {
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
