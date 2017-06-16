package main

import (
	"os"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func initLogging() {
	logFile, err := os.OpenFile(conf.Logging.Path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic("Cannot open loggin file" + err.Error())
	}

	fileLogging := logging.NewLogBackend(logFile, "", 0)
	fileLoggingFormater := logging.NewBackendFormatter(fileLogging, format)

	logging.SetBackend(fileLoggingFormater)
}
