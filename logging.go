package main

import (
	"os"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("mainLogger")

type Error struct {
	msg string
}

func initLogging() {
	format := logging.MustStringFormatter(conf.Logging.Format)
	logFile, err := os.OpenFile(conf.Logging.Path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic("Cannot open loggin file" + err.Error())
	}

	fileLogging := logging.NewLogBackend(logFile, "", 0)
	fileLoggingFormater := logging.NewBackendFormatter(fileLogging, format)

	logging.SetBackend(fileLoggingFormater)
}
