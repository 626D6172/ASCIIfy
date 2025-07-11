package main

import (
	"image"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/bmarse/ascii-render/pkg/asciify"
	"github.com/gosuri/uilive"
	"golang.org/x/term"
)

const shade = " .:coPO?@#"

func main() {
	file, err := os.Open(os.Args[1]) // Replace with your image file path
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	screen := uilive.New()
	screen.RefreshInterval = time.Millisecond * 200
	screen.Start()
	buf := asciify.CreateScreenBuffer(w, h)

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	asciify.ImageToASCIIToBuf(img, w, h, buf)
	screen.Write(buf)
	screen.Stop()
}
