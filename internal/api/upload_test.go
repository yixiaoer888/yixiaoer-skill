package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
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

func TestUploadCompressesImageWhenMaxImageBytesIsSet(t *testing.T) {
	imageBytes := noisyPNG(t, 1200, 900)
	if len(imageBytes) <= 512*1024 {
		t.Fatalf("test image must exceed 512KB, got %d", len(imageBytes))
	}
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "cover.png")
	if err := os.WriteFile(imagePath, imageBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	var uploaded []byte
	var uploadContentType string
	var requestedFileKey string
	var requestedSize string
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/storages/cloud-publish/upload-url":
			requestedFileKey = r.URL.Query().Get("fileKey")
			requestedSize = r.URL.Query().Get("size")
			if got := r.URL.Query().Get("contentType"); got != "image/jpeg" {
				t.Fatalf("unexpected compressed contentType query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/cover.jpg",
					"key":        "uploaded/cover.jpg",
				},
			})
		case "/oss/cover.jpg":
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
	result, err := client.UploadWithOptions(imagePath, "", false, UploadOptions{MaxImageBytes: 512 * 1024})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Compressed {
		t.Fatalf("expected compressed upload result, got %+v", result)
	}
	if result.OriginalSize != int64(len(imageBytes)) {
		t.Fatalf("unexpected original size: %+v", result)
	}
	if result.Size > 512*1024 || len(uploaded) > 512*1024 {
		t.Fatalf("expected compressed body <= 512KB, result=%d uploaded=%d", result.Size, len(uploaded))
	}
	if result.ContentType != "image/jpeg" || result.Format != "jpg" || uploadContentType != "image/jpeg" {
		t.Fatalf("unexpected compressed metadata/content type: result=%+v put=%s", result, uploadContentType)
	}
	if requestedFileKey != "cover.jpg" {
		t.Fatalf("expected compressed fileKey to use jpg extension, got %q", requestedFileKey)
	}
	if requestedSize != strconv.FormatInt(result.Size, 10) {
		t.Fatalf("expected upload-url size query to match compressed result size %d, got %q", result.Size, requestedSize)
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

func TestUploadURLImageFallsBackToStorageProxy(t *testing.T) {
	imageBytes := testPNG(t, 4, 5)
	var proxyCalls int
	var uploadCalls int
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/source/blocked.png":
			http.Error(w, "forbidden", http.StatusForbidden)
		case "/storages/proxy-url":
			proxyCalls++
			if got := r.URL.Query().Get("url"); got != server.URL+"/source/blocked.png" {
				t.Fatalf("unexpected proxy source url: %s", got)
			}
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write(imageBytes)
		case "/storages/material-library/upload-url":
			if got := r.URL.Query().Get("fileKey"); got != "blocked.png" {
				t.Fatalf("unexpected fileKey: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/blocked.png",
					"key":        "uploaded/blocked.png",
				},
			})
		case "/oss/blocked.png":
			uploadCalls++
			if got := r.Header.Get("Content-Type"); got != "image/png" {
				t.Fatalf("unexpected PUT content type: %s", got)
			}
			uploaded, err := readAll(r)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(uploaded, imageBytes) {
				t.Fatal("uploaded body did not match proxy response")
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	result, err := client.Upload(server.URL+"/source/blocked.png", "material-library", false)
	if err != nil {
		t.Fatal(err)
	}
	if proxyCalls != 1 || uploadCalls != 1 {
		t.Fatalf("expected one proxy call and one upload call, got proxy=%d upload=%d", proxyCalls, uploadCalls)
	}
	if result.Key != "uploaded/blocked.png" || result.Bucket != "material-library" {
		t.Fatalf("unexpected upload result: %+v", result)
	}
	if result.Width != 4 || result.Height != 5 {
		t.Fatalf("unexpected dimensions: %dx%d", result.Width, result.Height)
	}
}

func TestUploadLocalFileUsesASCIISafeObjectName(t *testing.T) {
	imageBytes := testPNG(t, 3, 2)
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "飞书20250424-172618.png")
	if err := os.WriteFile(imagePath, imageBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	var requestedFileKey string
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/storages/cloud-publish/upload-url":
			requestedFileKey = r.URL.Query().Get("fileKey")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/video.png",
					"key":        "uploaded/video.png",
				},
			})
		case "/oss/video.png":
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Upload(imagePath, "", false); err != nil {
		t.Fatal(err)
	}
	if requestedFileKey == "" {
		t.Fatal("expected upload-url request to include fileKey")
	}
	if requestedFileKey == "飞书20250424-172618.png" {
		t.Fatalf("expected fileKey to be normalized, got %q", requestedFileKey)
	}
	for _, r := range requestedFileKey {
		if r > 127 {
			t.Fatalf("expected ASCII-only fileKey, got %q", requestedFileKey)
		}
	}
	if filepath.Ext(requestedFileKey) != ".png" {
		t.Fatalf("expected normalized fileKey to preserve extension, got %q", requestedFileKey)
	}
}

func TestInspectUploadRejectsOversizedRemoteFile(t *testing.T) {
	var sourceRead bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sourceRead = true
		w.Header().Set("Content-Length", "5368709121")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	_, _, err := InspectUpload(server.URL+"/huge.bin", false)
	if err == nil {
		t.Fatal("expected oversized remote file error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured error, got %T: %v", err, err)
	}
	if typed.Code != yxerrors.UsageErr {
		t.Fatalf("expected usage error, got %+v", typed)
	}
	if !sourceRead {
		t.Fatal("expected source handler to be called")
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

func TestUploadHTTPClientDoesNotUseFixedTotalTimeout(t *testing.T) {
	if uploadHTTPClient.Timeout != 0 {
		t.Fatalf("expected upload HTTP client to avoid a fixed total timeout, got %s", uploadHTTPClient.Timeout)
	}
}

func TestInspectUploadLocalImage(t *testing.T) {
	imageBytes := testPNG(t, 6, 4)
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "preview.png")
	if err := os.WriteFile(imagePath, imageBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	result, fileName, err := InspectUpload(imagePath, true)
	if err != nil {
		t.Fatal(err)
	}
	if fileName != "preview.png" {
		t.Fatalf("unexpected fileName: %s", fileName)
	}
	if result.ContentType != "image/png" || result.Format != "png" {
		t.Fatalf("unexpected metadata: %+v", result)
	}
	if result.Width != 6 || result.Height != 4 {
		t.Fatalf("unexpected dimensions: %dx%d", result.Width, result.Height)
	}
	if result.Size != int64(len(imageBytes)) {
		t.Fatalf("unexpected size: %d", result.Size)
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

func noisyPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	seed := uint32(1)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			seed = seed*1664525 + 1013904223
			r := uint8(seed >> 24)
			seed = seed*1664525 + 1013904223
			g := uint8(seed >> 24)
			seed = seed*1664525 + 1013904223
			b := uint8(seed >> 24)
			img.Set(x, y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: 255,
			})
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
