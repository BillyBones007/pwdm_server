package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Log *logrus.Logger

// NewLogger - returns logger.
func NewLogger() Log {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Out = os.Stdout
		log.Info("Failed to log to file, using default stdout")
	}
	return log
}
