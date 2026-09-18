package docx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"os"
	pathpkg "path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const (
	maxDocumentXMLBytes   = 32 * 1024 * 1024
	maxEmbeddedImageBytes = 32 * 1024 * 1024
)

var numberedHeadingPattern = regexp.MustCompile(`^\s*\d+\s*[|｜]`)

// Document is the useful portion of a Word document for a WeChat article.
// BaseDir contains extracted embedded images and is valid until Cleanup is
// called. The generated HTML references those images by relative file name so
// the normal article-image materialization path can upload them and replace
// them with stable Yixiaoer URLs.
type Document struct {
	SourcePath string      `json:"sourcePath"`
	BaseDir    string      `json:"-"`
	HTML       string      `json:"html"`
	Paragraphs []Paragraph `json:"paragraphs"`
	Images     []Image     `json:"images"`
}

type Paragraph struct {
	Text     string `json:"text"`
	HTML     string `json:"html"`
	Style    string `json:"style,omitempty"`
	HasImage bool   `json:"hasImage,omitempty"`
}

type Image struct {
	Name       string `json:"name"`
	Target     string `json:"target"`
	SourcePath string `json:"sourcePath"`
	RelID      string `json:"relId"`
}

type segment struct {
	kind  string
	text  string
	relID string
	style runStyle
}

type rawParagraph struct {
	Style    string
	Segments []segment
}

type runStyle struct {
	Bold      bool
	Italic    bool
	Underline bool
	Strike    bool
}

type relationship struct {
	Target     string
	TargetMode string
}

// Parse reads an OOXML .docx file and extracts its text plus embedded images.
// It deliberately handles the WordprocessingML subset needed by the
// publisher rather than depending on a GUI or on a Python runtime.
func Parse(sourcePath string) (Document, error) {
	abs, err := filepath.Abs(strings.TrimSpace(sourcePath))
	if err != nil || strings.TrimSpace(sourcePath) == "" {
		return Document{}, docxUsage("invalid DOCX file path", map[string]interface{}{"path": sourcePath})
	}
	if !strings.EqualFold(filepath.Ext(abs), ".docx") {
		return Document{}, docxUsage("DOCX import requires a .docx file", map[string]interface{}{"path": abs})
	}
	if info, statErr := os.Stat(abs); statErr != nil {
		return Document{}, docxUsage("DOCX file not found", map[string]interface{}{"path": abs, "cause": statErr.Error()}).
			WithCategory("file_not_found").
			WithHint("请确认文章路径指向可读取的 .docx 文件。")
	} else if info.IsDir() {
		return Document{}, docxUsage("DOCX path is a directory", map[string]interface{}{"path": abs})
	}

	reader, err := zip.OpenReader(abs)
	if err != nil {
		return Document{}, docxUsage("failed to open DOCX archive", map[string]interface{}{"path": abs, "cause": err.Error()}).
			WithCategory("docx_archive").
			WithHint("请确认文件没有损坏，并且确实是 Office Open XML 的 .docx 文件。")
	}
	defer reader.Close()

	entries := make(map[string]*zip.File, len(reader.File))
	for _, entry := range reader.File {
		entries[entry.Name] = entry
	}
	documentEntry := entries["word/document.xml"]
	if documentEntry == nil {
		return Document{}, docxUsage("DOCX archive is missing word/document.xml", map[string]interface{}{"path": abs}).
			WithCategory("docx_archive").
			WithHint("请用 Word 或 WPS 重新另存为标准 .docx 文件后再导入。")
	}
	documentXML, err := readZipEntry(documentEntry, maxDocumentXMLBytes)
	if err != nil {
		return Document{}, err
	}
	relationships, err := readRelationships(entries["word/_rels/document.xml.rels"])
	if err != nil {
		return Document{}, err
	}
	rawParagraphs, err := parseDocument(documentXML)
	if err != nil {
		return Document{}, err
	}

	tempDir, err := os.MkdirTemp("", "yxer-docx-")
	if err != nil {
		return Document{}, docxUsage("failed to create DOCX image workspace", map[string]interface{}{"cause": err.Error()}).
			WithCategory("docx_workspace")
	}
	doc := Document{SourcePath: abs, BaseDir: tempDir}
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.RemoveAll(tempDir)
		}
	}()

	mediaByRelID := map[string]Image{}
	for _, paragraph := range rawParagraphs {
		for _, item := range paragraph.Segments {
			if item.kind != "image" {
				continue
			}
			if _, exists := mediaByRelID[item.relID]; exists {
				continue
			}
			relation, ok := relationships[item.relID]
			if !ok || relation.TargetMode != "" && !strings.EqualFold(relation.TargetMode, "Internal") {
				return Document{}, docxUsage("DOCX image relationship is missing or external", map[string]interface{}{"relId": item.relID}).
					WithCategory("docx_image")
			}
			zipName := resolveRelationshipTarget("word/document.xml", relation.Target)
			entry := entries[zipName]
			if entry == nil {
				return Document{}, docxUsage("DOCX embedded image is missing", map[string]interface{}{"relId": item.relID, "target": relation.Target, "zipPath": zipName}).
					WithCategory("docx_image")
			}
			raw, readErr := readZipEntry(entry, maxEmbeddedImageBytes)
			if readErr != nil {
				return Document{}, readErr
			}
			ext := strings.ToLower(filepath.Ext(zipName))
			if ext == "" {
				ext = ".bin"
			}
			name := fmt.Sprintf("image-%03d%s", len(doc.Images)+1, ext)
			imagePath := filepath.Join(tempDir, name)
			if writeErr := os.WriteFile(imagePath, raw, 0o600); writeErr != nil {
				return Document{}, docxUsage("failed to extract DOCX embedded image", map[string]interface{}{"target": relation.Target, "path": imagePath, "cause": writeErr.Error()}).
					WithCategory("docx_image")
			}
			mediaByRelID[item.relID] = Image{Name: name, Target: relation.Target, SourcePath: imagePath, RelID: item.relID}
			doc.Images = append(doc.Images, mediaByRelID[item.relID])
		}
	}

	var htmlBuilder strings.Builder
	for _, rawParagraph := range rawParagraphs {
		paragraph := renderParagraph(rawParagraph, mediaByRelID)
		if paragraph.Text == "" && !paragraph.HasImage {
			continue
		}
		doc.Paragraphs = append(doc.Paragraphs, paragraph)
		tag := "p"
		if isHeading(paragraph) {
			tag = "h2"
		}
		htmlBuilder.WriteString("<")
		htmlBuilder.WriteString(tag)
		htmlBuilder.WriteString(">")
		htmlBuilder.WriteString(paragraph.HTML)
		htmlBuilder.WriteString("</")
		htmlBuilder.WriteString(tag)
		htmlBuilder.WriteString(">")
	}
	doc.HTML = htmlBuilder.String()
	cleanup = false
	return doc, nil
}

// Cleanup removes the temporary extraction directory. It is safe to call
// multiple times and is intended to be deferred by the workflow using Parse.
func (d *Document) Cleanup() error {
	if d == nil || strings.TrimSpace(d.BaseDir) == "" {
		return nil
	}
	err := os.RemoveAll(d.BaseDir)
	d.BaseDir = ""
	return err
}

func readZipEntry(entry *zip.File, limit int64) ([]byte, error) {
	if entry == nil {
		return nil, docxUsage("DOCX archive entry is missing", nil).WithCategory("docx_archive")
	}
	if entry.UncompressedSize64 > uint64(limit) {
		return nil, docxUsage("DOCX archive entry exceeds the import limit", map[string]interface{}{
			"entry":        entry.Name,
			"limitBytes":   limit,
			"contentBytes": entry.UncompressedSize64,
		}).WithCategory("docx_size").WithHint("请压缩文档中的图片或拆分文档后再导入。")
	}
	reader, err := entry.Open()
	if err != nil {
		return nil, docxUsage("failed to read DOCX archive entry", map[string]interface{}{"entry": entry.Name, "cause": err.Error()}).WithCategory("docx_archive")
	}
	defer reader.Close()
	raw, err := io.ReadAll(io.LimitReader(reader, limit+1))
	if err != nil {
		return nil, docxUsage("failed to read DOCX archive entry", map[string]interface{}{"entry": entry.Name, "cause": err.Error()}).WithCategory("docx_archive")
	}
	if int64(len(raw)) > limit {
		return nil, docxUsage("DOCX archive entry exceeds the import limit", map[string]interface{}{"entry": entry.Name, "limitBytes": limit}).WithCategory("docx_size")
	}
	return raw, nil
}

func readRelationships(entry *zip.File) (map[string]relationship, error) {
	if entry == nil {
		return map[string]relationship{}, nil
	}
	raw, err := readZipEntry(entry, 4*1024*1024)
	if err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(strings.NewReader(string(raw)))
	result := map[string]relationship{}
	for {
		token, tokenErr := decoder.Token()
		if tokenErr == io.EOF {
			return result, nil
		}
		if tokenErr != nil {
			return nil, docxUsage("failed to parse DOCX relationships", map[string]interface{}{"cause": tokenErr.Error()}).WithCategory("docx_xml")
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Relationship" {
			continue
		}
		id := attr(start, "Id")
		target := attr(start, "Target")
		if id != "" && target != "" {
			result[id] = relationship{Target: target, TargetMode: attr(start, "TargetMode")}
		}
	}
}

func parseDocument(raw []byte) ([]rawParagraph, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(raw)))
	var paragraphs []rawParagraph
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return paragraphs, nil
		}
		if err != nil {
			return nil, docxUsage("failed to parse DOCX document.xml", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml").WithHint("请用 Word 或 WPS 重新另存为标准 .docx 文件后再导入。")
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "p" {
			continue
		}
		paragraph, paragraphErr := readParagraph(decoder)
		if paragraphErr != nil {
			return nil, paragraphErr
		}
		paragraphs = append(paragraphs, paragraph)
	}
}

func readParagraph(decoder *xml.Decoder) (rawParagraph, error) {
	var paragraph rawParagraph
	for {
		token, err := decoder.Token()
		if err != nil {
			return paragraph, docxUsage("failed to parse DOCX paragraph", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch typed := token.(type) {
		case xml.StartElement:
			switch typed.Name.Local {
			case "pPr":
				style, styleErr := readParagraphProperties(decoder)
				if styleErr != nil {
					return paragraph, styleErr
				}
				paragraph.Style = style
			case "r":
				segments, runErr := readRun(decoder)
				if runErr != nil {
					return paragraph, runErr
				}
				paragraph.Segments = append(paragraph.Segments, segments...)
			case "hyperlink", "smartTag", "sdtContent":
				segments, nestedErr := readNestedParagraphContent(decoder, typed.Name.Local)
				if nestedErr != nil {
					return paragraph, nestedErr
				}
				paragraph.Segments = append(paragraph.Segments, segments...)
			default:
				if skipErr := skipElement(decoder); skipErr != nil {
					return paragraph, skipErr
				}
			}
		case xml.EndElement:
			if typed.Name.Local == "p" {
				return paragraph, nil
			}
		}
	}
}

func readParagraphProperties(decoder *xml.Decoder) (string, error) {
	style := ""
	for {
		token, err := decoder.Token()
		if err != nil {
			return style, docxUsage("failed to parse DOCX paragraph properties", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch typed := token.(type) {
		case xml.StartElement:
			if typed.Name.Local == "pStyle" {
				if value := attr(typed, "val"); value != "" {
					style = value
				}
			}
			if skipErr := skipElement(decoder); skipErr != nil {
				return style, skipErr
			}
		case xml.EndElement:
			if typed.Name.Local == "pPr" {
				return style, nil
			}
		}
	}
}

func readNestedParagraphContent(decoder *xml.Decoder, endName string) ([]segment, error) {
	var segments []segment
	for {
		token, err := decoder.Token()
		if err != nil {
			return segments, docxUsage("failed to parse nested DOCX paragraph content", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch typed := token.(type) {
		case xml.StartElement:
			switch typed.Name.Local {
			case "r":
				run, runErr := readRun(decoder)
				if runErr != nil {
					return segments, runErr
				}
				segments = append(segments, run...)
			default:
				if skipErr := skipElement(decoder); skipErr != nil {
					return segments, skipErr
				}
			}
		case xml.EndElement:
			if typed.Name.Local == endName {
				return segments, nil
			}
		}
	}
}

func readRun(decoder *xml.Decoder) ([]segment, error) {
	var segments []segment
	style := runStyle{}
	for {
		token, err := decoder.Token()
		if err != nil {
			return segments, docxUsage("failed to parse DOCX run", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch typed := token.(type) {
		case xml.StartElement:
			switch typed.Name.Local {
			case "rPr":
				var styleErr error
				style, styleErr = readRunProperties(decoder, style)
				if styleErr != nil {
					return segments, styleErr
				}
			case "t", "instrText":
				text, textErr := readElementText(decoder, typed.Name.Local)
				if textErr != nil {
					return segments, textErr
				}
				if typed.Name.Local == "t" && text != "" {
					segments = append(segments, segment{kind: "text", text: text, style: style})
				}
			case "tab":
				if skipErr := skipElement(decoder); skipErr != nil {
					return segments, skipErr
				}
				segments = append(segments, segment{kind: "text", text: "\t", style: style})
			case "br", "cr":
				if skipErr := skipElement(decoder); skipErr != nil {
					return segments, skipErr
				}
				segments = append(segments, segment{kind: "text", text: "\n", style: style})
			case "drawing", "pict":
				relID, drawingErr := readDrawingRelationship(decoder, typed.Name.Local)
				if drawingErr != nil {
					return segments, drawingErr
				}
				if relID != "" {
					segments = append(segments, segment{kind: "image", relID: relID})
				}
			default:
				if skipErr := skipElement(decoder); skipErr != nil {
					return segments, skipErr
				}
			}
		case xml.EndElement:
			if typed.Name.Local == "r" {
				return segments, nil
			}
		}
	}
}

func readRunProperties(decoder *xml.Decoder, style runStyle) (runStyle, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			return style, docxUsage("failed to parse DOCX run properties", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch typed := token.(type) {
		case xml.StartElement:
			switch typed.Name.Local {
			case "b", "bCs":
				style.Bold = true
			case "i", "iCs":
				style.Italic = true
			case "u":
				style.Underline = true
			case "strike", "dstrike":
				style.Strike = true
			}
			if skipErr := skipElement(decoder); skipErr != nil {
				return style, skipErr
			}
		case xml.EndElement:
			if typed.Name.Local == "rPr" {
				return style, nil
			}
		}
	}
}

func readDrawingRelationship(decoder *xml.Decoder, endName string) (string, error) {
	relID := ""
	depth := 1
	for depth > 0 {
		token, err := decoder.Token()
		if err != nil {
			return relID, docxUsage("failed to parse DOCX image drawing", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch typed := token.(type) {
		case xml.StartElement:
			depth++
			if typed.Name.Local == "blip" || typed.Name.Local == "imagedata" {
				if value := firstAttr(typed, "embed", "id"); value != "" {
					relID = value
				}
			}
		case xml.EndElement:
			depth--
			if depth == 0 && typed.Name.Local != endName {
				return relID, nil
			}
		}
	}
	return relID, nil
}

func readElementText(decoder *xml.Decoder, endName string) (string, error) {
	var builder strings.Builder
	depth := 1
	for depth > 0 {
		token, err := decoder.Token()
		if err != nil {
			return builder.String(), docxUsage("failed to parse DOCX text", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch typed := token.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		case xml.CharData:
			builder.Write([]byte(typed))
		}
	}
	return builder.String(), nil
}

func skipElement(decoder *xml.Decoder) error {
	depth := 1
	for depth > 0 {
		token, err := decoder.Token()
		if err != nil {
			return docxUsage("failed to skip unsupported DOCX element", map[string]interface{}{"cause": err.Error()}).WithCategory("docx_xml")
		}
		switch token.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}
	return nil
}

func renderParagraph(raw rawParagraph, images map[string]Image) Paragraph {
	var textBuilder strings.Builder
	var htmlBuilder strings.Builder
	hasImage := false
	for _, item := range raw.Segments {
		if item.kind == "image" {
			image, ok := images[item.relID]
			if !ok {
				continue
			}
			hasImage = true
			htmlBuilder.WriteString(`<img src="`)
			htmlBuilder.WriteString(html.EscapeString(image.Name))
			htmlBuilder.WriteString(`" alt="" />`)
			continue
		}
		textBuilder.WriteString(item.text)
		value := html.EscapeString(item.text)
		value = strings.ReplaceAll(value, "\n", "<br>")
		value = strings.ReplaceAll(value, "\t", "&emsp;")
		if item.style.Bold {
			value = "<strong>" + value + "</strong>"
		}
		if item.style.Italic {
			value = "<em>" + value + "</em>"
		}
		if item.style.Underline {
			value = "<u>" + value + "</u>"
		}
		if item.style.Strike {
			value = "<s>" + value + "</s>"
		}
		htmlBuilder.WriteString(value)
	}
	return Paragraph{Text: strings.TrimSpace(textBuilder.String()), HTML: htmlBuilder.String(), Style: raw.Style, HasImage: hasImage}
}

func isHeading(paragraph Paragraph) bool {
	style := strings.ToLower(strings.TrimSpace(paragraph.Style))
	return strings.Contains(style, "heading") || numberedHeadingPattern.MatchString(paragraph.Text)
}

func resolveRelationshipTarget(source, target string) string {
	target = strings.ReplaceAll(strings.TrimSpace(target), "\\", "/")
	if strings.HasPrefix(target, "/") {
		target = strings.TrimPrefix(target, "/")
	}
	return pathpkg.Clean(pathpkg.Join(pathpkg.Dir(source), target))
}

func attr(element xml.StartElement, name string) string {
	for _, item := range element.Attr {
		if item.Name.Local == name {
			return strings.TrimSpace(item.Value)
		}
	}
	return ""
}

func firstAttr(element xml.StartElement, names ...string) string {
	for _, name := range names {
		if value := attr(element, name); value != "" {
			return value
		}
	}
	return ""
}

func docxUsage(message string, details interface{}) *yxerrors.Error {
	return yxerrors.Usage(message, details).
		WithCategory("docx_import").
		WithHint("请确认输入是可读取的标准 .docx 文件；导入微信公众号文章时还需要单独提供横版封面。")
}
