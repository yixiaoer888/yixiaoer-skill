package docx

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseExtractsTextAndEmbeddedImages(t *testing.T) {
	docxPath := writeTestDOCX(t)
	document, err := Parse(docxPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(document.Paragraphs) != 2 {
		t.Fatalf("expected two non-empty paragraphs, got %+v", document.Paragraphs)
	}
	if document.Paragraphs[0].Text != "文章标题" {
		t.Fatalf("unexpected title text: %+v", document.Paragraphs[0])
	}
	if len(document.Images) != 1 {
		t.Fatalf("expected one embedded image, got %+v", document.Images)
	}
	if !strings.Contains(document.HTML, `<img src="image-001.jpeg"`) {
		t.Fatalf("expected HTML to reference extracted image, got %s", document.HTML)
	}
	if _, err := os.Stat(document.Images[0].SourcePath); err != nil {
		t.Fatalf("expected extracted image to exist: %v", err)
	}
	baseDir := document.BaseDir
	if err := document.Cleanup(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		t.Fatalf("expected DOCX workspace to be removed, got %v", err)
	}
}

func writeTestDOCX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "article.docx")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	archive := zip.NewWriter(file)
	entries := map[string]string{
		"word/document.xml": `<?xml version="1.0"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><w:body><w:p><w:r><w:t>文章标题</w:t></w:r></w:p><w:p><w:r><w:t>正文</w:t></w:r><w:r><w:drawing><wp:inline><a:blip r:embed="rId1"/></wp:inline></w:drawing></w:r></w:p></w:body></w:document>`,
		"word/_rels/document.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/image1.jpeg"/></Relationships>`,
	}
	for name, body := range entries {
		writer, writerErr := archive.Create(name)
		if writerErr != nil {
			t.Fatal(writerErr)
		}
		if _, writerErr = writer.Write([]byte(body)); writerErr != nil {
			t.Fatal(writerErr)
		}
	}
	writer, err := archive.Create("word/media/image1.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("test-image")); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}
