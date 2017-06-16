package main

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

type resizePoint struct {
	Width  int
	Height int
}

func resizeImage(sourceImage image.Image, maxWidth, maxHeight int) *image.NRGBA {
	resizeValue := resizePoint{0, 0}

	if sourceImage.Bounds().Dx() > sourceImage.Bounds().Dy() {
		if sourceImage.Bounds().Dx() > conf.Image.MaxWidth {
			resizeValue.Width = maxWidth
		}
	} else {
		if sourceImage.Bounds().Dy() > conf.Image.MaxHeight {
			resizeValue.Height = maxHeight
		}
	}

	return imaging.Resize(sourceImage, resizeValue.Width, resizeValue.Height, imaging.Lanczos)
}

func convertImage(inputImage image.Image, output string) image.Image {

	processImage := resizeImage(inputImage, conf.Image.MaxWidth, conf.Image.MaxHeight)
	processImage = imaging.Sharpen(processImage, 0.2)
	processImage = imaging.AdjustBrightness(processImage, 1)
	processImage = imaging.AdjustContrast(processImage, 1)

	outputImage, _ := os.Create(output)
	defer outputImage.Close()

	jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: conf.Image.Quality})

	return processImage
}

func convertTumbnail(inputImage image.Image, output string) image.Image {

	processImage := imaging.Resize(inputImage, 300, 0, imaging.Lanczos)
	processImage = imaging.Fill(processImage, 200, 200, imaging.Center, imaging.Lanczos)

	processImage = imaging.Sharpen(processImage, 1.2)
	processImage = imaging.AdjustBrightness(processImage, 3)
	processImage = imaging.AdjustContrast(processImage, 1)
	processImage = imaging.AdjustGamma(processImage, 1.0)

	outputImage, _ := os.Create(output)
	defer outputImage.Close()

	//imaging.Save(m, output)
	jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: conf.Thumbnail.Quality})

	return processImage
}
