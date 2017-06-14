package main

import (
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

type ResizePoint struct {
	width int
	heigh int
}

const MAX_WIDTH int = 2000
const MAX_HEIGHT int = 2000

func convertImage(input, output string, quality int) {
	inputImage, _ := imaging.Open(input)
	var resizeValue ResizePoint = ResizePoint{23, 23}
	if inputImage.Bounds().Dx() > inputImage.Bounds().Dy() {
		if inputImage.Bounds().Dx() > MAX_WIDTH {
			resizeValue.width = MAX_WIDTH
		}
	} else {
		if inputImage.Bounds().Dy() > MAX_HEIGHT {
			resizeValue.heigh = MAX_HEIGHT
		}
	}

	m := imaging.Resize(inputImage, resizeValue.width, resizeValue.heigh, imaging.Lanczos)
	m = imaging.Sharpen(m, 0.2)
	m = imaging.AdjustBrightness(m, 1)
	m = imaging.AdjustContrast(m, 1)

	outputImage, _ := os.Create(output)
	defer outputImage.Close()

	//imaging.Save(m, output)
	jpeg.Encode(outputImage, m, &jpeg.Options{Quality: quality})
}

func convertSmallImage(input, output string, quality int) {
	inputImage, _ := imaging.Open(input)

	m := imaging.Resize(inputImage, 300, 0, imaging.Lanczos)
	m = imaging.Fill(m, 200, 200, imaging.Center, imaging.Lanczos)
	m = imaging.Sharpen(m, 1.2)
	m = imaging.AdjustBrightness(m, 3)
	m = imaging.AdjustContrast(m, 1)
	m = imaging.AdjustGamma(m, 1.0)

	outputImage, _ := os.Create(output)
	defer outputImage.Close()

	//imaging.Save(m, output)
	jpeg.Encode(outputImage, m, &jpeg.Options{Quality: quality})
}
