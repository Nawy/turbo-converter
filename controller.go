package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		fmt.Fprintf(w, "Error %d", 1)
	}

	imagePath := getPathWithHash(conf.Image.Path, getHash())
	thumbnailPath := getPathWithHash(conf.Thumbnail.Path, getHash())

	optimalImage := convertImage(inputImage, imagePath)
	convertTumbnail(optimalImage, thumbnailPath)

	response := ImageResponseJSON{imagePath, thumbnailPath}
	jsonResponse(w, response)
}

// handler for request delete image by image path(http)
func deleteImageHandler(w http.ResponseWriter, r *http.Request) {

}

// handler for request status of server
func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Result: %s", getHash())
}

func jsonResponse(w http.ResponseWriter, response interface{}) {
	responseJSON, err := json.Marshal(response)

	if err != nil {
		panic("Canno create json response")
	}

	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	w.Write(responseJSON)
}
