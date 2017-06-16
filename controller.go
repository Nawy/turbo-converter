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

type ImageResponseJSON struct {
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
}

type StatusResponseJSON struct {
	Time  string `json:"time"`
	Space string `json:"space"`
}

// handler for request upload image(http)
func uploadImageHandler(w http.ResponseWriter, r *http.Request) {

	inputImage, err := imaging.Decode(r.Body)

	if err != nil {
		log.Error("Cannot decode uploaded date")
		return
	}

	imagePath := getPathWithHash(conf.Image.Path, getHash())
	thumbnailPath := getPathWithHash(conf.Thumbnail.Path, getHash())

	optimalImage := convertImage(inputImage, imagePath)
	convertTumbnail(optimalImage, thumbnailPath)

	log.Infof("Uploaded image=[%s] with tumblnail=[%s]", imagePath, thumbnailPath)

	response := ImageResponseJSON{imagePath, thumbnailPath}
	jsonResponse(w, response)
}

// handler for request delete image by image path(http)
func deleteImageHandler(w http.ResponseWriter, r *http.Request) {

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

	jsonResponse(w, response)
}

func convertMemValue(memBytes uint64) string {
	var result float64 = float64(memBytes)
	memDimension := "BYTE"
	if result/1024 >= 1 {
		result = result / 1024
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

func jsonResponse(w http.ResponseWriter, response interface{}) {
	responseJSON, err := json.Marshal(response)

	if err != nil {
		panic("Canno create json response")
	}

	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	w.Write(responseJSON)
}
