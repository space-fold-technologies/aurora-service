package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *LoggerProvider

type LoggerProvider struct {
	log *logrus.Logger
}

func (l *LoggerProvider) init() {
	l.log = logrus.New()
	Formatter := new(logrus.TextFormatter)
	l.log.SetFormatter(Formatter)
	l.log.Out = os.Stdout
	Formatter.FullTimestamp = true
}

func GetInstance() *logrus.Logger {
	if logger == nil {
		logger = &LoggerProvider{}
		logger.init()
	}
	return logger.log
}
