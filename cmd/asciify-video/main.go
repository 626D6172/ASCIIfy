package main

import (
	"image"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/626d6172/ascii-render/pkg/asciify"
	vidio "github.com/AlexEidt/Vidio"
	"github.com/gosuri/uilive"
	"golang.org/x/term"
)

const shade = " .:coPO?@#"

func main() {
	video, err := vidio.NewVideo(os.Args[1])
	if err != nil {
		panic(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
	video.SetFrameBuffer(img.Pix)

	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	screen := uilive.New()
	screen.RefreshInterval = time.Millisecond * 50
	screen.Start()
	c := make(chan []byte, 5)
	go func() {
		buf := asciify.CreateScreenBuffer(w, h)

		screen.Write(buf)

		for video.Read() {
			asciify.ImageToASCIIToBuf(img, w, h, buf)
			c <- buf
		}
		close(c)
	}()

	ticker := time.NewTicker(time.Millisecond * time.Duration(1000/video.FPS()))
	for {
		b, ok := <-c
		if !ok {
			break
		}
		<-ticker.C
		screen.Write(b)
		screen.Flush()
	}
	screen.Stop()
}
