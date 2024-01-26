package log

import "github.com/sirupsen/logrus"

// This Example Of Usage Logger
// =============================
// logger.Info("This is log info")
// logger.Warn("This is log warn")
// logger.Error("This is log error")
// logger.WithFields(logrus.Fields{
// "param": "field",
// }).Info("With Params")

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
}

func New() Logger {
	l := logrus.New()
	return l
}
