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

func resizeImage(sourceImage image.Image, maxWidth, maxHeight int, filter imaging.ResampleFilter) *image.NRGBA {
	resizeValue := resizePoint{sourceImage.Bounds().Dx(), sourceImage.Bounds().Dy()}

	if sourceImage.Bounds().Dx() > sourceImage.Bounds().Dy() {
		if sourceImage.Bounds().Dx() > maxWidth {
			resizeValue.Width = maxWidth
			resizeValue.Height = 0
		}
	} else {
		if sourceImage.Bounds().Dy() > maxHeight {
			resizeValue.Width = 0
			resizeValue.Height = maxHeight
		}
	}

	return imaging.Resize(sourceImage, resizeValue.Width, resizeValue.Height, filter)
}

func convertImage(inputImage image.Image, output string) (image.Image, error) {

	processImage := resizeImage(inputImage, conf.Image.MaxWidth, conf.Image.MaxHeight, imaging.Lanczos)
	processImage = imaging.Sharpen(processImage, 0.2)
	processImage = imaging.AdjustBrightness(processImage, 1)
	processImage = imaging.AdjustContrast(processImage, 1)

	outputImage, err := os.Create(output)
	defer outputImage.Close()
	if err != nil {
		log.Errorf("Cannot open file %s", output)
		return nil, err
	}

	err = jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: conf.Image.Quality})
	if err != nil {
		log.Errorf("Cannot save image with name %s", output)
	}

	return processImage, err
}

func convertTumbnail(inputImage image.Image, output string) (image.Image, error) {

	processImage := resizeImage(
		inputImage,
		conf.Thumbnail.MaxWidth,
		conf.Thumbnail.MaxHeight,
		imaging.Blackman,
	)
	processImage = imaging.Fill(
		processImage,
		conf.Thumbnail.MaxWidth,
		conf.Thumbnail.MaxHeight,
		imaging.Center,
		imaging.Blackman,
	)

	processImage = imaging.Sharpen(processImage, 1.2)
	processImage = imaging.AdjustBrightness(processImage, 3)
	processImage = imaging.AdjustContrast(processImage, 1)
	processImage = imaging.AdjustGamma(processImage, 1.0)

	outputImage, err := os.Create(output)
	defer outputImage.Close()
	if err != nil {
		log.Errorf("Cannot open file %s", output)
		return nil, err
	}

	err = jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: conf.Thumbnail.Quality})
	if err != nil {
		log.Errorf("Cannot save thumbnail with name %s", output)
	}

	return processImage, err
}
