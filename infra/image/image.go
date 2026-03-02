package image

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"math"

	drawx "golang.org/x/image/draw"
)

// ensureNRGBA 将任意图片转为 *image.NRGBA 以支持直接像素访问。
// 若已是 *image.NRGBA 则零拷贝返回。
func ensureNRGBA(img image.Image) *image.NRGBA {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)
	return dst
}

// CropWhitespace 裁剪白边：直接像素访问，避免逐像素 img.At() 接口开销
func CropWhitespace(img image.Image, threshold uint8) image.Image {
	nrgba := ensureNRGBA(img)
	bounds := nrgba.Bounds()
	pix := nrgba.Pix
	stride := nrgba.Stride
	t := threshold
	w := bounds.Dx()

	minX, minY := bounds.Min.X, bounds.Min.Y
	maxX, maxY := bounds.Max.X-1, bounds.Max.Y-1

	// 1. 从上往下扫，确定上边界
	topFound := false
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rowOff := (y - bounds.Min.Y) * stride
		for x := 0; x < w; x++ {
			i := rowOff + x*4
			if pix[i] < t || pix[i+1] < t || pix[i+2] < t {
				minY = y
				topFound = true
				break
			}
		}
		if topFound {
			break
		}
	}

	// 如果全图都是白色，直接返回原图
	if !topFound {
		return img
	}

	// 2. 从下往上扫，确定下边界
	bottomFound := false
	for y := bounds.Max.Y - 1; y >= minY; y-- {
		rowOff := (y - bounds.Min.Y) * stride
		for x := 0; x < w; x++ {
			i := rowOff + x*4
			if pix[i] < t || pix[i+1] < t || pix[i+2] < t {
				maxY = y
				bottomFound = true
				break
			}
		}
		if bottomFound {
			break
		}
	}

	// 3. 从左往右扫，确定左边界（仅扫描 minY..maxY 范围）
	leftFound := false
	h := maxY - minY + 1
	for x := 0; x < w; x++ {
		for dy := 0; dy < h; dy++ {
			i := (minY-bounds.Min.Y+dy)*stride + x*4
			if pix[i] < t || pix[i+1] < t || pix[i+2] < t {
				minX = bounds.Min.X + x
				leftFound = true
				break
			}
		}
		if leftFound {
			break
		}
	}

	// 4. 从右往左扫，确定右边界
	rightFound := false
	for x := w - 1; x >= minX-bounds.Min.X; x-- {
		for dy := 0; dy < h; dy++ {
			i := (minY-bounds.Min.Y+dy)*stride + x*4
			if pix[i] < t || pix[i+1] < t || pix[i+2] < t {
				maxX = bounds.Min.X + x
				rightFound = true
				break
			}
		}
		if rightFound {
			break
		}
	}

	return nrgba.SubImage(image.Rect(minX, minY, maxX+1, maxY+1))
}

// CropCompress 完整流程
func CropCompress(input []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	// 1. 裁剪白边（阈值 240）
	cropped := CropWhitespace(img, 240)

	// 2. 缩放并压缩
	return ResizeAndCompress(cropped, 80, 0.6)
}

// ResizeAndCompress 缩放并转为 JPEG
func ResizeAndCompress(img image.Image, quality int, scale float64) ([]byte, error) {
	bounds := img.Bounds()
	newWidth := int(math.Max(1, float64(bounds.Dx())*scale))
	newHeight := int(math.Max(1, float64(bounds.Dy())*scale))

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 使用 drawx.BiLinear 兼顾速度与质量
	drawx.BiLinear.Scale(dst, dst.Bounds(), img, bounds, drawx.Src, nil)

	var buf bytes.Buffer
	// 预估大小，减少扩容次数
	buf.Grow(newWidth * newHeight / 4)

	err := jpeg.Encode(&buf, dst, &jpeg.Options{
		Quality: quality,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
