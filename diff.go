package gsdiff

import (
	"errors"
	"image"
	"image/color"
	"sync"
)

type coords struct {
	x, y int
}

const (
	red = iota
	green

	sumRGB           = 765
	sumRGBOnePercent = 7.6
)

/*
	CreateDiff creates new images from diff between two images.
	Red color set if the diff between colors of two pixels is negative, otherwise green color is set.
*/
func CreateDiff(img1, img2 image.Image) (image.Image, error) {
	//drawing background
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return nil, errors.New("images have different bounds")
	}

	diff := image.NewRGBA(image.Rect(0, 0, bounds1.Dx(), bounds1.Dy()))
	var wg sync.WaitGroup
	xLimit, yLimit := bounds1.Dx(), bounds1.Dy()
	//filling every row in different goroutine
	for i := 0; i <= xLimit; i++ {
		wg.Add(1)
		go fillPixels(img1, img2, diff, coords{i, yLimit}, &wg)
	}

	wg.Wait()
	return diff, nil
}

//filling range of pixels
func fillPixels(img1, img2 image.Image, diff *image.RGBA, crd coords, wg *sync.WaitGroup) {
	for y := 0; y < crd.y; y++ {
		diff.Set(crd.x, y, chooseColor(img1.At(crd.x, y), img2.At(crd.x, y)))
	}
	wg.Done()
}

/*
	chooseColor. Calculating difference by subtracting sum of r+g+b of one color from r+g+b of another
	if result is positive than the color is green, otherwise red.
*/
func chooseColor(clr1, clr2 color.Color) color.RGBA {
	c1, c2, c3, _ := clr1.RGBA()
	c4, c5, c6, _ := clr2.RGBA()
	colorDiff := int(c1+c2+c3) - int(c4+c5+c6)

	if colorDiff == 0 { //don't change color if both colors are the same
		r, g, b, a := clr2.RGBA()
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	} else if colorDiff > 0 {
		return calcCoefficient(colorDiff, green)
	}
	return calcCoefficient(colorDiff, red)
}

/*
	calcCoefficient
	765 - sum of max values of r+g+b
	coefficient calculated between 128 and 255 of red or green channel respectively
*/
func calcCoefficient(diff, clr int) color.RGBA {
	coefficient := uint8(128 + ((sumRGB-float64(diff))/sumRGBOnePercent)*1.28)
	if clr == red {
		return color.RGBA{255, coefficient, coefficient, 255} //red
	}
	return color.RGBA{coefficient, 255, coefficient, 255} //green
}
