package common

import (
	"io"

	"github.com/sirupsen/logrus"
)

type LogOption struct {
	LogLevel logrus.Level
	Logger   *logrus.Logger
}

func NewLogger(opt ...LogOption) *logrus.Logger {
	logger := logrus.New()
	if len(opt) == 0 {
		logger.Out = io.Discard
		return logger
	}
	if opt[0].Logger != nil {
		return opt[0].Logger
	}
	logger.SetLevel(opt[0].LogLevel)
	return logger
}
