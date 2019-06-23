package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/rocc0/gs-diff"
)

func main() {

	var output string

	flag.StringVar(&output, "output", "diff.png", "output filename")
	flag.StringVar(&output, "o", "diff.png", "output filename")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("usage: diff [<option>...] <image1> <image2>")
		os.Exit(1)
	}

	img1, err := mustLoadImage(args[0])
	if err != nil {
		log.Fatal(err)
	}
	img2, err := mustLoadImage(args[1])
	if err != nil {
		log.Fatal(err)
	}
	result, err := gsdiff.CreateDiff(img1, img2)
	if err != nil {
		log.Fatal(err)
	}
	if err := mustSaveImage(result, output); err != nil {
		log.Fatal(err)
	}
}

func mustLoadImage(filename string) (image.Image, error) {
	f, err := mustOpen(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func mustOpen(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func mustSaveImage(img image.Image, output string) error {
	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
