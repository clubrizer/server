package pkg

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

var (
	log *logrus.Logger
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)
}

func Debug(format string, v ...interface{}) {
	log.WithField(getCaller()).Debugf(format, v...)
}

func Info(format string, v ...interface{}) {

	log.WithField(getCaller()).Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.WithField(getCaller()).Warnf(format, v...)
}

func Error(err error, format string, v ...interface{}) {
	log.WithField(getCaller()).WithField("error", err).Errorf(format, v...)
}

func Fatal(err error, format string, v ...interface{}) {
	log.WithField(getCaller()).WithField("error", err).Fatalf(format, v...)
}

func getCaller() (string, string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		log.Error("Failed to get log caller")
		return "", ""
	}

	return "caller", fmt.Sprintf("%s:%d", file, line)
}
