// Package asciify holds all the logic for converting images into fun ascii art
package asciify

import (
	"fmt"
	"image"
	"io"
	"math"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
	"golang.org/x/term"
)

const SHADER = " .:coPO?@#"

func ResizeImage(m image.Image, w uint, h uint) image.Image {
	bounds := m.Bounds()
	imageWidth := bounds.Max.X
	imageHeight := bounds.Max.Y

	var newImage image.Image
	if imageWidth > imageHeight {
		newImage = resize.Resize(w, 0, m, resize.Lanczos3)
	} else {
		newImage = resize.Resize(0, h, m, resize.Lanczos3)
	}

	return newImage
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

	m := ResizeImage(img, uint(w/2), uint(h))

	bounds := m.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	fmt.Println(width, height, w, h)

	for y := range h {
		for x := range w {
			if x >= width {
				continue
			}
			c := m.At(x, y)
			r, g, b, _ := c.RGBA()
			// Calc lum https://en.wikipedia.org/wiki/Relative_luminance
			rf := 0.2126 * float64(r)
			gf := 0.7152 * float64(g)
			bf := 0.0722 * float64(b)
			lum := int(math.Floor(((rf + gf + bf) / 0xffff) * float64(len(SHADER)-1)))
			fmt.Print(string(SHADER[lum]), string(SHADER[lum]))
		}
		fmt.Print("\n")
	}

	return nil
}
