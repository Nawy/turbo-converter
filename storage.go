package main

import "os"

// path from config
// hash generated hash code
func getPathWithHash(path, hash string) string {
	resultFolder := path + "/" + hash[0:1]
	os.MkdirAll(resultFolder, os.ModePerm)

	resultPath := resultFolder + "/" + hash[1:] + ".jpg"
	return resultPath
}
