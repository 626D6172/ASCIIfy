// Package asciify holds all the logic for converting images into fun ascii art
package asciify

import (
	"image"
	"math"

	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
)

const SHADER = "  .:coPO?@#"

func CreateScreenBuffer(w int, h int) []byte {
	buf := make([]byte, h*(w+1))
	for i := range buf {
		buf[i] = ' '
	}

	// line endings
	for i := 1; i <= h; i++ {
		buf[((w+1)*i)-1] = '\n'
	}

	return buf
}

func ResizeImage(m image.Image, w uint, h uint) image.Image {
	return resize.Thumbnail(w, h, m, resize.Lanczos3)
}

func ImageToASCIIToBuf(img image.Image, w int, h int, buf []byte) {
	img = ResizeImage(img, uint(w/2), uint(h))

	bounds := img.Bounds()

	for y := range h {
		go func(y int) {
			for x := 0; x < w; x += 2 {
				c := img.At(x/2, y)
				if x/2 > bounds.Max.X {
					buf[(y*w)+x] = SHADER[0]
					buf[(y*w)+x+1] = SHADER[0]
					return
				}
				r, g, b, _ := c.RGBA()
				// Calc lum https://en.wikipedia.org/wiki/Relative_luminance
				rf := 0.2126 * float64(r)
				gf := 0.7152 * float64(g)
				bf := 0.0722 * float64(b)
				lum := int(math.Floor(((rf + gf + bf) / 0xffff) * float64(len(SHADER)-1)))
				buf[(y*w)+x] = SHADER[lum]
				buf[(y*w)+x+1] = SHADER[lum]
			}
		}(y)
	}
}
