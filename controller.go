package main

import "net/http"

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	convertTumbnail(r.Body, "test2.jpg")
}
