package gsdiff

import (
	"image"
	_ "image/png"
	"os"
	"testing"
)

func BenchmarkCreateDiff(b *testing.B) {
	openImage := func(filename string) (image.Image, error) {
		f, err := os.Open(filename)
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
	im1, err := openImage("./1.png")
	if err != nil {
		b.Fatal(err)
	}
	im2, err := openImage("./2.png")
	if err != nil {
		b.Fatal(err)
	}
	if _, err := CreateDiff(im1, im2); err != nil {
		b.Error(err)
	}
}
