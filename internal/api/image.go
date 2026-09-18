package api

import (
	"bytes"
	"encoding/binary"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

func imageDimensions(pathOrURL string, raw []byte, contentType string) (int, int) {
	if !strings.HasPrefix(contentType, "image/") {
		return 0, 0
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(raw))
	if err == nil {
		return cfg.Width, cfg.Height
	}
	if width, height, ok := webpDimensions(raw); ok {
		return width, height
	}
	return 0, 0
}

// webpDimensions reads dimensions from the WebP container without requiring a
// full WebP decoder. This keeps upload validation working for WebP covers even
// though the standard library does not register a WebP image decoder.
func webpDimensions(raw []byte) (int, int, bool) {
	if len(raw) < 16 || string(raw[:4]) != "RIFF" || string(raw[8:12]) != "WEBP" {
		return 0, 0, false
	}

	for offset := 12; offset+8 <= len(raw); {
		chunkType := string(raw[offset : offset+4])
		chunkSize := int(binary.LittleEndian.Uint32(raw[offset+4 : offset+8]))
		dataStart := offset + 8
		if chunkSize < 0 || dataStart > len(raw) || chunkSize > len(raw)-dataStart {
			return 0, 0, false
		}
		dataEnd := dataStart + chunkSize

		switch chunkType {
		case "VP8 ":
			// VP8 frame header: three-byte frame tag, three-byte start
			// code, then 16-bit little-endian width and height.
			if chunkSize >= 10 && raw[dataStart+3] == 0x9d && raw[dataStart+4] == 0x01 && raw[dataStart+5] == 0x2a {
				width := int(binary.LittleEndian.Uint16(raw[dataStart+6:dataStart+8])) & 0x3fff
				height := int(binary.LittleEndian.Uint16(raw[dataStart+8:dataStart+10])) & 0x3fff
				if width > 0 && height > 0 {
					return width, height, true
				}
			}
		case "VP8L":
			if chunkSize >= 5 && raw[dataStart] == 0x2f {
				width := 1 + (int(raw[dataStart+1]) & 0x3f) + (int(raw[dataStart+2])&0x3f)<<8
				height := 1 + (int(raw[dataStart+2]) >> 6) + int(raw[dataStart+3])<<2 + (int(raw[dataStart+4])&0x0f)<<10
				if width > 0 && height > 0 {
					return width, height, true
				}
			}
		case "VP8X":
			if chunkSize >= 10 {
				width := 1 + int(raw[dataStart+4]) + int(raw[dataStart+5])<<8 + int(raw[dataStart+6])<<16
				height := 1 + int(raw[dataStart+7]) + int(raw[dataStart+8])<<8 + int(raw[dataStart+9])<<16
				if width > 0 && height > 0 {
					return width, height, true
				}
			}
		}

		offset = dataEnd
		if chunkSize%2 == 1 {
			offset++
		}
	}
	return 0, 0, false
}
