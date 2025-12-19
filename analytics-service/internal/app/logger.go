package app

import (
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

func initLogger(level string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		log.SetLevel(log.DEBUG)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		log.SetLevel(log.INFO)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		log.SetLevel(log.WARN)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		log.SetLevel(log.ERROR)
	default:
		logrus.SetLevel(logrus.InfoLevel)
		log.SetLevel(log.INFO)
	}
}

