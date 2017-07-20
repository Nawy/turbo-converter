package service

import (
	"log"
	"os"

	config "turbo-converter/config"
)

type Error struct {
	msg string
}

func InitLogging() *os.File {
	f, err := os.OpenFile(config.GlobalConfig().Logging.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("Log started")
	return f
}
