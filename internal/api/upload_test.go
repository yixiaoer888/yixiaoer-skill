package api

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestUploadLocalImage(t *testing.T) {
	imageBytes := testPNG(t, 3, 2)
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "cover.png")
	if err := os.WriteFile(imagePath, imageBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	var uploaded []byte
	var uploadContentType string
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/storages/cloud-publish/upload-url":
			if got := r.Header.Get("Authorization"); got != "test-key" {
				t.Fatalf("unexpected authorization header: %s", got)
			}
			if got := r.URL.Query().Get("fileKey"); got != "cover.png" {
				t.Fatalf("unexpected fileKey: %s", got)
			}
			if got := r.URL.Query().Get("contentType"); got != "image/png" {
				t.Fatalf("unexpected contentType query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/cover.png",
					"key":        "uploaded/cover.png",
				},
			})
		case "/oss/cover.png":
			if r.Method != http.MethodPut {
				t.Fatalf("unexpected upload method: %s", r.Method)
			}
			uploadContentType = r.Header.Get("Content-Type")
			var err error
			uploaded, err = readAll(r)
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Upload(imagePath, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Key != "uploaded/cover.png" {
		t.Fatalf("unexpected key: %s", result.Key)
	}
	if result.ContentType != "image/png" || result.Bucket != "cloud-publish" || result.Format != "png" {
		t.Fatalf("unexpected result metadata: %+v", result)
	}
	if result.Width != 3 || result.Height != 2 {
		t.Fatalf("unexpected dimensions: %dx%d", result.Width, result.Height)
	}
	if uploadContentType != "image/png" {
		t.Fatalf("unexpected PUT content type: %s", uploadContentType)
	}
	if !bytes.Equal(uploaded, imageBytes) {
		t.Fatal("uploaded body did not match local file")
	}
}

func TestUploadURLImage(t *testing.T) {
	imageBytes := testPNG(t, 4, 5)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/source/photo.png":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write(imageBytes)
		case "/storages/material-library/upload-url":
			if got := r.URL.Query().Get("fileKey"); got != "photo.png" {
				t.Fatalf("unexpected fileKey: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/photo.png",
					"key":        "uploaded/photo.png",
				},
			})
		case "/oss/photo.png":
			if got := r.Header.Get("Content-Type"); got != "image/png" {
				t.Fatalf("unexpected PUT content type: %s", got)
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Upload(server.URL+"/source/photo.png", "material-library", false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Key != "uploaded/photo.png" || result.Bucket != "material-library" {
		t.Fatalf("unexpected upload result: %+v", result)
	}
	if result.Width != 4 || result.Height != 5 {
		t.Fatalf("unexpected dimensions: %dx%d", result.Width, result.Height)
	}
}

func TestDetectContentTypeVideoFallbacks(t *testing.T) {
	if got := DetectContentType("clip.mp4"); got != "video/mp4" {
		t.Fatalf("unexpected mp4 content type: %s", got)
	}
	if got := DetectContentType("asset.unknown"); got != "application/octet-stream" {
		t.Fatalf("unexpected unknown content type: %s", got)
	}
}

func testPNG(t *testing.T, width, height int) []byte {
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

func readAll(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.Body)
	return buffer.Bytes(), err
}
