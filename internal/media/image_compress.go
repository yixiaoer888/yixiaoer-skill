package media

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"math"
)

type CompressedImage struct {
	Data        []byte
	ContentType string
	Format      string
	Width       int
	Height      int
}

func CompressImageToMaxBytes(raw []byte, maxBytes int64) (CompressedImage, bool, error) {
	if maxBytes <= 0 || int64(len(raw)) <= maxBytes {
		return CompressedImage{}, false, nil
	}
	src, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return CompressedImage{}, false, err
	}
	base := flattenImage(src)
	bounds := base.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return CompressedImage{}, false, nil
	}

	for {
		for quality := 90; quality >= 35; quality -= 5 {
			encoded, err := encodeJPEG(base, quality)
			if err != nil {
				return CompressedImage{}, false, err
			}
			if int64(len(encoded)) <= maxBytes {
				return CompressedImage{
					Data:        encoded,
					ContentType: "image/jpeg",
					Format:      "jpg",
					Width:       base.Bounds().Dx(),
					Height:      base.Bounds().Dy(),
				}, true, nil
			}
		}

		nextWidth := int(math.Floor(float64(base.Bounds().Dx()) * 0.85))
		nextHeight := int(math.Floor(float64(base.Bounds().Dy()) * 0.85))
		if nextWidth < 64 || nextHeight < 64 {
			encoded, err := encodeJPEG(base, 30)
			if err != nil {
				return CompressedImage{}, false, err
			}
			if int64(len(encoded)) <= maxBytes {
				return CompressedImage{
					Data:        encoded,
					ContentType: "image/jpeg",
					Format:      "jpg",
					Width:       base.Bounds().Dx(),
					Height:      base.Bounds().Dy(),
				}, true, nil
			}
			return CompressedImage{}, false, nil
		}
		base = resizeNearest(base, nextWidth, nextHeight)
	}
}

func flattenImage(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), src, bounds.Min, draw.Over)
	return dst
}

func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	var buffer bytes.Buffer
	err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: quality})
	return buffer.Bytes(), err
}

func resizeNearest(src *image.RGBA, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()
	for y := 0; y < height; y++ {
		sy := srcBounds.Min.Y + y*srcHeight/height
		for x := 0; x < width; x++ {
			sx := srcBounds.Min.X + x*srcWidth/width
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}
