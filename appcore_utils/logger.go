package appcore_utils

import (
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

// NewLogger returns a new logger instance
func NewLogger(configs *Configurations) *logrus.Logger {
	l := logrus.New()

	if configs.ObserveIsActive {
		l.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
		)))
	}
	return l
}
