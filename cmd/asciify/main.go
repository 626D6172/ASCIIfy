package main

import (
	"fmt"
	"image"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/626d6172/ascii-render/pkg/asciify"
	"github.com/gosuri/uilive"
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
