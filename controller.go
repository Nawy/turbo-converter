package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/disintegration/imaging"
)

// ImageJSON response for upload request in convertation process
type ImageJSON struct {
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
}

// StatusResponseJSON status response, contains time and available space on the disk
type StatusResponseJSON struct {
	Time  string `json:"time"`
	Space string `json:"space"`
}

// ErrorResponseJSON error response contains only error description in JSON format
type ErrorResponseJSON struct {
	Description string `json:"description"`
}

// SimpleResponseJSON is simple one line response
type SimpleResponseJSON struct {
	Response string `json:"response"`
}

var semaphore Semaphore = make(Semaphore, 3)

// handler for request upload image(http)
func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if !isMethod(w, r, "POST") {
		return
	}

	semaphore.Lock()
	defer semaphore.Unlock()

	inputImage, err := imaging.Decode(r.Body)

	if err != nil {
		log.Error("Cannot decode uploaded image")
		jsonResponse(w, ErrorResponseJSON{"Body is empty"}, 400)
		return
	}

	imageHash := getHash()
	thumbnailHash := getHash()

	storageImagePath := getPathWithHash(conf.Image.StoragePath, imageHash)
	storageThumbnailPath := getPathWithHash(conf.Thumbnail.StoragePath, thumbnailHash)

	responseImagePath := getResponsePathWithHash(conf.Image.ResponsePath, imageHash)
	responseThumbnailPath := getResponsePathWithHash(conf.Thumbnail.ResponsePath, thumbnailHash)

	optimalImage, errImage := convertImage(inputImage, storageImagePath)
	if isError(errImage, w, "Cannot convert image", 500) {
		return
	}
	_, errThumbnail := convertThumbnail(optimalImage, storageThumbnailPath)
	if isError(errThumbnail, w, "Cannot convert thumbnail", 500) {
		return
	}

	log.Infof("Uploaded image=[%s] with thumbnail=[%s]", storageImagePath, storageThumbnailPath)

	response := ImageJSON{responseImagePath, responseThumbnailPath}
	jsonResponse(w, response, 200)
}

// handler for request delete image by image path(http)
func deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	if !isMethod(w, r, "DELETE") {
		return
	}

	decoder := json.NewDecoder(r.Body)
	requestJson := ImageJSON{}

	err := decoder.Decode(&requestJson)
	if err != nil {
		log.Errorf("\nError: %s", err)
		jsonResponse(w, ErrorResponseJSON{"Wrong JSON request!"}, 400)
		return
	}

	imagePath, error := getStorageImagePath(requestJson.Image)
	if isError(error, w, "Wrong image path", 400) {
		return
	}

	thumbnailPath, error := getStorageThumbnailPath(requestJson.Thumbnail)
	if isError(error, w, "Wrong thumbnail path", 400) {
		return
	}

	_, error = deleteFile(imagePath)
	if isError(error, w, "Wrong image path", 400) {
		return
	}

	_, error = deleteFile(thumbnailPath)
	if isError(error, w, "Wrong thumbnail path", 400) {
		return
	}
	codeResponse(w, 200)
}

// handler for request status of server
func statusHandler(w http.ResponseWriter, r *http.Request) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()

	if err != nil {
		log.Error("Internal error of reading free space")
		return
	}

	syscall.Statfs(wd, &stat)

	freeSpace := stat.Bavail * uint64(stat.Bsize)
	response := StatusResponseJSON{time.Now().String(), convertMemValue(freeSpace)}

	jsonResponse(w, response, 200)
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

func isError(err *Error, w http.ResponseWriter, message string, status int) bool {
	if err != nil {
		log.Errorf("Error: %s", err.msg)
		jsonResponse(w, ErrorResponseJSON{message}, status)
		return true
	}
	return false
}

func isMethod(w http.ResponseWriter, r *http.Request, methodName string) bool {
	if r.Method != methodName {
		jsonResponse(w, ErrorResponseJSON{"You should use " + methodName + " http method"}, 400)
		return false
	}
	return true
}

func codeResponse(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func jsonResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	responseJSON, err := json.Marshal(response)

	if err != nil {
		panic("Cannot create json response")
	}

	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(statusCode)
	w.Write(responseJSON)
}
