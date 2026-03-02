package image

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
)

// 判断是否接近白色（允许一点误差）
func isWhite(c color.Color, threshold uint8) bool {
	r, g, b, _ := c.RGBA()

	R := uint8(r >> 8)
	G := uint8(g >> 8)
	B := uint8(b >> 8)

	return R > threshold && G > threshold && B > threshold
}

// 裁剪白边
func CropWhitespace(img image.Image, threshold uint8) image.Image {

	bounds := img.Bounds()

	minX := bounds.Max.X
	minY := bounds.Max.Y
	maxX := bounds.Min.X
	maxY := bounds.Min.Y

	found := false

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			if !isWhite(img.At(x, y), threshold) {

				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}

				found = true
			}
		}
	}

	if !found {
		return img
	}

	rect := image.Rect(minX, minY, maxX+1, maxY+1)

	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)
}

// []byte → 裁剪 → base64
func CropAndBase64(input []byte) ([]byte, error) {

	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	cropped := CropWhitespace(img, 240) // 240 = 接近白色

	var buf bytes.Buffer

	err = png.Encode(&buf, cropped)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
