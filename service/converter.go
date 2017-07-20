package service

import (
	"image"
	"image/jpeg"
	"log"
	"os"

	config "turbo-converter/config"
	model "turbo-converter/model"

	"github.com/disintegration/imaging"
)

type ResizePoint struct {
	Width  int
	Height int
}

func resizeImage(sourceImage image.Image, maxWidth, maxHeight int, imageType model.Type) *image.NRGBA {
	width := sourceImage.Bounds().Dx()
	height := sourceImage.Bounds().Dy()
	resizeValue := ResizePoint{0, 0}

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

	if imageType == model.IMAGE {
		if width > maxWidth || height > maxHeight {
			filter = getFilterByType(config.GlobalConfig().Image.Downscale)
		}
	} else {
		if width > maxWidth || height > maxHeight {
			filter = getFilterByType(config.GlobalConfig().Thumbnail.Resize.Downscale)
		}
	}

	return imaging.Resize(sourceImage, resizeValue.Width, resizeValue.Height, filter)
}

func ConvertImage(inputImage image.Image, output string) (image.Image, *Error) {

	log.Printf(
		"IMAGE: Sharpen: %f; Brightness: %f; Contrast: %f; Gamma %f;",
		config.GlobalConfig().Image.PostProcessing.Sharpen,
		config.GlobalConfig().Image.PostProcessing.Brightness,
		config.GlobalConfig().Image.PostProcessing.Contrast,
		config.GlobalConfig().Image.PostProcessing.Gamma,
	)

	processImage := inputImage
	if isNeedScale(processImage, config.GlobalConfig().Image.MaxWidth, config.GlobalConfig().Image.MaxHeight) {
		processImage = resizeImage(inputImage, config.GlobalConfig().Image.MaxWidth, config.GlobalConfig().Image.MaxHeight, model.IMAGE)
	}

	processImage = imaging.Sharpen(processImage, config.GlobalConfig().Image.PostProcessing.Sharpen)
	processImage = imaging.AdjustBrightness(processImage, config.GlobalConfig().Image.PostProcessing.Brightness)
	processImage = imaging.AdjustContrast(processImage, config.GlobalConfig().Image.PostProcessing.Contrast)
	processImage = imaging.AdjustGamma(processImage, config.GlobalConfig().Image.PostProcessing.Gamma)

	outputImage, err := os.Create(output)
	defer outputImage.Close()
	if err != nil {
		return nil, &Error{"Cannot open file " + output}
	}

	err = jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: config.GlobalConfig().Image.Quality})
	if err != nil {
		return nil, &Error{"Cannot save image with name "}
	}

	return processImage, nil
}

func ConvertThumbnail(inputImage image.Image, output string) (image.Image, *Error) {

	log.Printf(
		"THUMBNAIL: Sharpen: %f; Brightness: %f; Contrast: %f; Gamma %f;",
		config.GlobalConfig().Thumbnail.PostProcessing.Sharpen,
		config.GlobalConfig().Thumbnail.PostProcessing.Brightness,
		config.GlobalConfig().Thumbnail.PostProcessing.Contrast,
		config.GlobalConfig().Thumbnail.PostProcessing.Gamma,
	)

	processImage := inputImage
	if isNeedScale(processImage, config.GlobalConfig().Thumbnail.MaxWidth, config.GlobalConfig().Thumbnail.MaxHeight) {
		processImage = resizeImage(
			inputImage,
			config.GlobalConfig().Thumbnail.MaxWidth,
			config.GlobalConfig().Thumbnail.MaxHeight,
			model.THUMBNAIL,
		)
	}

	processImage = imaging.Fill(
		processImage,
		config.GlobalConfig().Thumbnail.MaxWidth,
		config.GlobalConfig().Thumbnail.MaxHeight,
		imaging.Center,
		getFilterByType(config.GlobalConfig().Thumbnail.Resize.Upscale),
	)

	processImage = imaging.Sharpen(processImage, config.GlobalConfig().Thumbnail.PostProcessing.Sharpen)
	processImage = imaging.AdjustBrightness(processImage, config.GlobalConfig().Thumbnail.PostProcessing.Brightness)
	processImage = imaging.AdjustContrast(processImage, config.GlobalConfig().Thumbnail.PostProcessing.Contrast)
	processImage = imaging.AdjustGamma(processImage, config.GlobalConfig().Thumbnail.PostProcessing.Gamma)

	outputImage, err := os.Create(output)
	defer outputImage.Close()
	if err != nil {
		return nil, &Error{"Cannot open file " + output}
	}

	err = jpeg.Encode(outputImage, processImage, &jpeg.Options{Quality: config.GlobalConfig().Thumbnail.Quality})
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

func isGraphics(inputImage image.Image) bool {

	pixels := make(map[uint32]bool)
	for x := 0; x <= inputImage.Bounds().Dx(); x++ {
		for y := 0; y < inputImage.Bounds().Dy(); y++ {
			r, g, b, _ := inputImage.At(x, y).RGBA()
			key := ((r*1000)+g)*1000 + b
			pixels[uint32(key)] = true
		}
	}

	length := len(pixels)
	size := inputImage.Bounds().Dx() * inputImage.Bounds().Dy()
	ratio := float32(length) / (float32(size) / 100.0)

	if ratio <= 1.5 {
		return true
	}
	return false
}

func isNeedScale(img image.Image, maxWidth, maxHeight int) bool {
	if img.Bounds().Dx() > maxWidth || img.Bounds().Dy() > maxHeight {
		return true
	}
	return false
}
