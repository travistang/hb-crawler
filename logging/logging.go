package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type LoggerConfig struct {
	Prefix string
	Level  logrus.Level
}

func GetLogger(config *LoggerConfig) *logrus.Logger {
	logger := logrus.New()
	fields := logrus.Fields{}

	if len(config.Prefix) > 0 {
		logger.Formatter = new(prefixed.TextFormatter)
		fields["prefix"] = fmt.Sprintf("[%s]", config.Prefix)
	}

	if config.Level > 0 {
		logger.Level = config.Level
	}

	logger.WithFields(fields)
	return logger
}
