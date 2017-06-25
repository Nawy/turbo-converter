package main

import (
	"os"
	"strings"
)

// path from config
// hash generated hash code
func getPathWithHash(path, hash string) string {
	resultFolder := path + "/" + hash[0:1]
	os.MkdirAll(resultFolder, os.ModePerm)

	resultPath := resultFolder + "/" + hash[1:] + ".jpg"
	return resultPath
}

func getResponsePathWithHash(path, hash string) string {
	resultFolder := path + "/" + hash[0:1]

	resultPath := resultFolder + "/" + hash[1:] + ".jpg"
	return resultPath
}

func isFileExists(path string) (bool, error) {
	log.Infof("Check file %s", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getStorageImagePath(responsePath string) (string, *Error) {
	isHasPrefix := strings.HasPrefix(responsePath, conf.Image.ResponsePath)

	if !isHasPrefix {
		return "", &Error{"Wrong path for image " + responsePath}
	}

	resultString := strings.Replace(responsePath, conf.Image.ResponsePath, conf.Image.StoragePath, 1)

	return resultString, nil
}

func getStorageThumbnailPath(responsePath string) (string, *Error) {
	isHasPrefix := strings.HasPrefix(responsePath, conf.Thumbnail.ResponsePath)

	if !isHasPrefix {
		return "", &Error{"Wrong path for thumbnail " + responsePath}
	}

	resultString := strings.Replace(responsePath, conf.Thumbnail.ResponsePath, conf.Thumbnail.StoragePath, 1)

	return resultString, nil
}

func deleteFile(path string) (bool, *Error) {

	isFileExists, _ := isFileExists(path)
	if !isFileExists {
		return false, &Error{"File isn't exists " + path}
	}

	err := os.Remove(path)
	if err != nil {
		return false, &Error{"Cannot remove file " + path}
	}

	return true, nil
}
