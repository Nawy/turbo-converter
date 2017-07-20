package controller

import (
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	config "turbo-converter/config"
	model "turbo-converter/model"
	service "turbo-converter/service"

	httpm "github.com/Nawy/go-rest-server/model"
	"github.com/disintegration/imaging"
)

var semaphore service.Semaphore = make(service.Semaphore, 3)

// handler for request upload image(http)
func UploadImageHandler(context *httpm.HTTPContext) *httpm.ResponseEntity {

	semaphore.Lock()
	defer semaphore.Unlock()

	inputImage, err := imaging.Decode(context.Request.Body)

	if err != nil {
		log.Println("Cannot decode uploaded image")
		context.WriteError("Body is empty", httpm.HTTP_BAD_REQUEST)
		return &httpm.ResponseEntity{Body: nil, Headers: nil}
	}

	imageHash := service.GetHash()
	thumbnailHash := service.GetHash()

	storageImagePath := service.GetPathWithHash(config.GlobalConfig().Image.StoragePath, imageHash)
	storageThumbnailPath := service.GetPathWithHash(config.GlobalConfig().Thumbnail.StoragePath, thumbnailHash)

	responseImagePath := service.GetResponsePathWithHash(config.GlobalConfig().Image.ResponsePath, imageHash)
	responseThumbnailPath := service.GetResponsePathWithHash(config.GlobalConfig().Thumbnail.ResponsePath, thumbnailHash)

	optimalImage, errImage := service.ConvertImage(inputImage, storageImagePath)
	if errImage != nil {
		context.WriteError("Cannot convert image", httpm.HTTP_INTERNAL_SERVER_ERROR)
		return &httpm.ResponseEntity{Body: nil, Headers: nil}
	}

	_, errThumbnail := service.ConvertThumbnail(optimalImage, storageThumbnailPath)
	if errThumbnail != nil {
		context.WriteError("Cannot convert thumbnail", httpm.HTTP_INTERNAL_SERVER_ERROR)
		return &httpm.ResponseEntity{Body: nil, Headers: nil}
	}

	log.Printf("Uploaded image=[%s] with thumbnail=[%s]", storageImagePath, storageThumbnailPath)

	response := model.ImageJSON{Image: responseImagePath, Thumbnail: responseThumbnailPath}
	return &httpm.ResponseEntity{Body: response, Headers: nil}
}

// handler for request delete image by image path(http)
func DeleteImageHandler(context *httpm.HTTPContext) *httpm.ResponseEntity {

	requestJSON := model.ImageJSON{Image: "", Thumbnail: ""}
	context.GetBody(requestJSON)

	imagePath, error := service.GetStorageImagePath(requestJSON.Image)
	if error != nil {
		context.WriteError("Wrong image path", httpm.HTTP_BAD_REQUEST)
		return httpm.EmptyResponse()
	}

	thumbnailPath, error := service.GetStorageThumbnailPath(requestJSON.Thumbnail)
	if error != nil {
		context.WriteError("Wrong thumbnail path", httpm.HTTP_BAD_REQUEST)
		return httpm.EmptyResponse()
	}

	_, error = service.DeleteFile(imagePath)
	if error != nil {
		context.WriteError("Wrong image path", httpm.HTTP_BAD_REQUEST)
		return httpm.EmptyResponse()
	}

	_, error = service.DeleteFile(thumbnailPath)
	if error != nil {
		context.WriteError("Wrong thumbnail path", httpm.HTTP_BAD_REQUEST)
		return httpm.EmptyResponse()
	}
	return httpm.EmptyResponse()
}

// handler for request status of server
func StatusHandler(context *httpm.HTTPContext) *httpm.ResponseEntity {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()

	if err != nil {
		context.WriteError("Internal error of reading free space", httpm.HTTP_INTERNAL_SERVER_ERROR)
		return httpm.EmptyResponse()
	}

	syscall.Statfs(wd, &stat)

	freeSpace := stat.Bavail * uint64(stat.Bsize)
	response := model.StatusResponseJSON{Time: time.Now().String(), Space: convertMemValue(freeSpace)}

	return &httpm.ResponseEntity{Body: response, Headers: nil}
}

func convertMemValue(memBytes uint64) string {
	var result float64
	memDimension := "BYTE"
	if float64(memBytes)/1024 >= 1 {
		result = float64(memBytes) / 1024
		memDimension = "KB"
	}

	if result/1024 >= 1 {
		result = result / 1024
		memDimension = "MB"
	}

	if result/1024 >= 1 {
		result = result / 1024
		memDimension = "GB"
	}

	resultString := strconv.FormatFloat(result, 'f', 2, 64) + memDimension
	return resultString
}
