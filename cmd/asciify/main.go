package main

import (
	"fmt"
	"image"
	_ "image/jpeg" // Import necessary image formats
	_ "image/png"
	"math"
	"os"

	"golang.org/x/term"
)

const shade = " .:coPO?@#"

func main() {
	file, err := os.Open(os.Args[1]) // Replace with your image file path
	if err != nil {
		fmt.Println("Error opening image:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	tWidth, tHeight, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Println(tHeight, tWidth)

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
					totalLum += int(math.Floor(((rf + gf + bf) / 0xffff) * float64(len(shade))))
				}
			}
			if sampleCount == 0 {
				continue
			}
			fmt.Print(string(shade[totalLum/sampleCount]), string(shade[totalLum/sampleCount]))
		}
		fmt.Print("\n")
	}
}
