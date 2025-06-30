// Package asciify holds all the logic for converting images into fun ascii art
package asciify

import (
	"image"
	"io"
	"math"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gosuri/uilive"
	"github.com/nfnt/resize"
	"golang.org/x/term"
)

const SHADER = " .:coPO?@#"

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

func SoftASCII(r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	img = ResizeImage(img, uint(w/2), uint(h))

	bounds := img.Bounds()
	buf := CreateScreenBuffer(w, h)
	screen := uilive.New()
	screen.RefreshInterval = time.Millisecond * 200
	screen.Start()

	for y := range h {
		for x := 0; x < w; x += 2 {
			c := img.At(x/2, y)
			if x/2 > bounds.Max.X {
				buf[(y*w)+x] = SHADER[0]
				buf[(y*w)+x+1] = SHADER[0]
				continue
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
	}

	screen.Write(buf)
	screen.Stop()

	return nil
}
