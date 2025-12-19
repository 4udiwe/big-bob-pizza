package app

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func initLogger(level string) {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(logLevel)
	}

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}
