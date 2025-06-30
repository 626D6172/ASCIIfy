package main

import (
	"fmt"
	"os"

	"github.com/626d6172/ascii-render/pkg/asciify"
)

const shade = " .:coPO?@#"

func main() {
	file, err := os.Open(os.Args[1]) // Replace with your image file path
	if err != nil {
		fmt.Println("Error opening image:", err)
		return
	}
	defer file.Close()

	err = asciify.SoftASCII(file)
	if err != nil {
		panic(err)
	}
}
