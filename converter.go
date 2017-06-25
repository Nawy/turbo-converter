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

func resizeImage(sourceImage image.Image, maxWidth, maxHeight int, imageType Type) *image.NRGBA {
	width := sourceImage.Bounds().Dx()
	height := sourceImage.Bounds().Dy()
	resizeValue := resizePoint{0, 0}

	if width > height {
		if width > maxWidth {
			resizeValue.Width = maxWidth
			resizeValue.Height = 0
		}
	} else {
		if height > maxHeight {
			resizeValue.Width = 0
			resizeValue.Height = maxHeight
		}
	}

	var filter imaging.ResampleFilter

	if imageType == IMAGE {
		if width > maxWidth || height > maxHeight {
			filter = getFilterByType(conf.Image.Downscale)
		}
	} else {
		if width > maxWidth || height > maxHeight {
			filter = getFilterByType(conf.Thumbnail.Resize.Downscale)
		}
	}

	return imaging.Resize(sourceImage, resizeValue.Width, resizeValue.Height, filter)
}

func convertImage(inputImage image.Image, output string) (image.Image, *Error) {

	log.Infof(
		"IMAGE: Sharpen: %f; Brightness: %f; Contrast: %f; Gamma %f;",
		conf.Image.PostProcessing.Sharpen,
		conf.Image.PostProcessing.Brightness,
		conf.Image.PostProcessing.Contrast,
		conf.Image.PostProcessing.Gamma,
	)

	processImage := inputImage
	if isNeedScale(processImage, conf.Image.MaxWidth, conf.Image.MaxHeight) {
		processImage = resizeImage(inputImage, conf.Image.MaxWidth, conf.Image.MaxHeight, IMAGE)
	}

	processImage = imaging.Sharpen(processImage, conf.Image.PostProcessing.Sharpen)
	processImage = imaging.AdjustBrightness(processImage, conf.Image.PostProcessing.Brightness)
	processImage = imaging.AdjustContrast(processImage, conf.Image.PostProcessing.Contrast)
	processImage = imaging.AdjustGamma(processImage, conf.Image.PostProcessing.Gamma)

	outputImage, err := os.Create(output)
	defer outputImage.Close()
	if err != nil {
		return nil, &Error{"Cannot open file " + output}
	}

	err = jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: conf.Image.Quality})
	if err != nil {
		return nil, &Error{"Cannot save image with name "}
	}

	return processImage, nil
}

func convertThumbnail(inputImage image.Image, output string) (image.Image, *Error) {

	log.Infof(
		"THUMBNAIL: Sharpen: %f; Brightness: %f; Contrast: %f; Gamma %f;",
		conf.Thumbnail.PostProcessing.Sharpen,
		conf.Thumbnail.PostProcessing.Brightness,
		conf.Thumbnail.PostProcessing.Contrast,
		conf.Thumbnail.PostProcessing.Gamma,
	)

	processImage := inputImage
	if isNeedScale(processImage, conf.Thumbnail.MaxWidth, conf.Thumbnail.MaxHeight) {
		processImage = resizeImage(
			inputImage,
			conf.Thumbnail.MaxWidth,
			conf.Thumbnail.MaxHeight,
			THUMBNAIL,
		)
	}

	processImage = imaging.Fill(
		processImage,
		conf.Thumbnail.MaxWidth,
		conf.Thumbnail.MaxHeight,
		imaging.Center,
		getFilterByType(conf.Thumbnail.Resize.Upscale),
	)

	processImage = imaging.Sharpen(processImage, conf.Thumbnail.PostProcessing.Sharpen)
	processImage = imaging.AdjustBrightness(processImage, conf.Thumbnail.PostProcessing.Brightness)
	processImage = imaging.AdjustContrast(processImage, conf.Thumbnail.PostProcessing.Contrast)
	processImage = imaging.AdjustGamma(processImage, conf.Thumbnail.PostProcessing.Gamma)

	outputImage, err := os.Create(output)
	defer outputImage.Close()
	if err != nil {
		return nil, &Error{"Cannot open file " + output}
	}

	err = jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: conf.Thumbnail.Quality})
	if err != nil {
		return nil, &Error{"Cannot save thumbnail with name " + output}
	}

	return processImage, nil
}

func getFilterByType(value string) imaging.ResampleFilter {
	switch value {
	case "BSpline":
		return imaging.BSpline
	case "Bartlett":
		return imaging.Bartlett
	case "Blackman":
		return imaging.Blackman
	case "Box":
		return imaging.Box
	case "CatmullRom":
		return imaging.CatmullRom
	case "Cosine":
		return imaging.Cosine
	case "Gaussian":
		return imaging.Gaussian
	case "Hamming":
		return imaging.Hamming
	case "Hann":
		return imaging.Hann
	case "Hermite":
		return imaging.Hermite
	case "Linear":
		return imaging.Linear
	case "MitchellNetravali":
		return imaging.MitchellNetravali
	case "NearestNeighbor":
		return imaging.NearestNeighbor
	case "Welch":
		return imaging.Welch
	default:
		return imaging.Lanczos
	}
}

func isNeedScale(img image.Image, maxWidth, maxHeight int) bool {
	if img.Bounds().Dx() > maxWidth || img.Bounds().Dy() > maxHeight {
		return true
	}
	return false
}
