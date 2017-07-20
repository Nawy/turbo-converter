package service

import (
	"log"
	"os"
	"strings"
	config "turbo-converter/config"
)

// path from config
// hash generated hash code
func GetPathWithHash(path, hash string) string {
	resultFolder := path + "/" + hash[0:1]
	os.MkdirAll(resultFolder, os.ModePerm)

	resultPath := resultFolder + "/" + hash[1:] + ".jpg"
	return resultPath
}

func GetResponsePathWithHash(path, hash string) string {
	resultFolder := path + "/" + hash[0:1]

	resultPath := resultFolder + "/" + hash[1:] + ".jpg"
	return resultPath
}

func isFileExists(path string) (bool, error) {
	log.Printf("Check file %s", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func GetStorageImagePath(responsePath string) (string, *Error) {
	isHasPrefix := strings.HasPrefix(responsePath, config.GlobalConfig().Image.ResponsePath)

	if !isHasPrefix {
		return "", &Error{"Wrong path for image " + responsePath}
	}

	resultString := strings.Replace(responsePath, config.GlobalConfig().Image.ResponsePath, config.GlobalConfig().Image.StoragePath, 1)

	return resultString, nil
}

func GetStorageThumbnailPath(responsePath string) (string, *Error) {
	isHasPrefix := strings.HasPrefix(responsePath, config.GlobalConfig().Thumbnail.ResponsePath)

	if !isHasPrefix {
		return "", &Error{"Wrong path for thumbnail " + responsePath}
	}

	resultString := strings.Replace(responsePath, config.GlobalConfig().Thumbnail.ResponsePath, config.GlobalConfig().Thumbnail.StoragePath, 1)

	return resultString, nil
}

func DeleteFile(path string) (bool, *Error) {

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
