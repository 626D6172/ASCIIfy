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

	"golang.org/x/term"
)

const SHADER = " .:coPO?@#"

func SoftASCII(r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	tWidth, tHeight, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	widthRatio := width / tWidth
	heightRatio := height / tHeight

	ratio := max(heightRatio, widthRatio)

	for y := bounds.Min.Y; y < height; y += ratio {
		for x := bounds.Min.X; x < width; x += ratio {
			totalLum := 0
			sampleCount := 0
			for i := range ratio {
				for ii := range ratio {
					sampleCount++
					c := img.At(x+i, y+ii)
					r, g, b, _ := c.RGBA()
					// Calc lum https://en.wikipedia.org/wiki/Relative_luminance
					rf := 0.2126 * float64(r)
					gf := 0.7152 * float64(g)
					bf := 0.0722 * float64(b)
					totalLum += int(math.Floor(((rf + gf + bf) / 0xffff) * float64(len(SHADER)-1)))
				}
			}
			if sampleCount == 0 {
				continue
			}
			fmt.Print(string(SHADER[totalLum/sampleCount]), string(SHADER[totalLum/sampleCount]))
		}
		fmt.Print("\n")
	}

	return nil
}
