package logging

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	Prefix string
	Level  logrus.Level
}

func GetLogger(config *LoggerConfig) *logrus.Logger {
	logger := logrus.New()
	fields := logrus.Fields{}

	if len(config.Prefix) > 0 {
		logger.SetFormatter(&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{config.Prefix},
		})
	}

	if config.Level > 0 {
		logger.Level = config.Level
	}

	logger.WithFields(fields)
	return logger
}
