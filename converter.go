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

	log.Infof(
		"IMAGE: Sharpen: %f\nBrightness: %f\nContrast: %f\nGamma %f\n",
		conf.Image.PostProcessing.Sharpen,
		conf.Image.PostProcessing.Brightness,
		conf.Image.PostProcessing.Contrast,
		conf.Image.PostProcessing.Gamma,
	)
	processImage := resizeImage(inputImage, conf.Image.MaxWidth, conf.Image.MaxHeight, imaging.Lanczos)
	processImage = imaging.Sharpen(processImage, conf.Image.PostProcessing.Sharpen)
	processImage = imaging.AdjustBrightness(processImage, conf.Image.PostProcessing.Brightness)
	processImage = imaging.AdjustContrast(processImage, conf.Image.PostProcessing.Contrast)
	processImage = imaging.AdjustGamma(processImage, conf.Image.PostProcessing.Gamma)

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

func convertThumbnail(inputImage image.Image, output string) (image.Image, error) {

	log.Infof(
		"THUMBNAIL: Sharpen: %f\nBrightness: %f\nContrast: %f\nGamma %f\n",
		conf.Thumbnail.PostProcessing.Sharpen,
		conf.Thumbnail.PostProcessing.Brightness,
		conf.Thumbnail.PostProcessing.Contrast,
		conf.Thumbnail.PostProcessing.Gamma,
	)

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

	processImage = imaging.Sharpen(processImage, conf.Thumbnail.PostProcessing.Sharpen)
	processImage = imaging.AdjustBrightness(processImage, conf.Thumbnail.PostProcessing.Brightness)
	processImage = imaging.AdjustContrast(processImage, conf.Thumbnail.PostProcessing.Contrast)
	processImage = imaging.AdjustGamma(processImage, conf.Thumbnail.PostProcessing.Gamma)

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
