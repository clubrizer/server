// Package log provides a logger to be used within the Clubrizer server.
package log

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

// Debug logs the given inputs with the DEBUG log level.
func Debug(format string, v ...interface{}) {
	log.WithField(getCaller()).Debugf(format, v...)
}

// Info logs the given inputs with the INFO log level.
func Info(format string, v ...interface{}) {

	log.WithField(getCaller()).Infof(format, v...)
}

// Warn logs the given inputs with the WARN log level.
func Warn(format string, v ...interface{}) {
	log.WithField(getCaller()).Warnf(format, v...)
}

// Error logs the given inputs with the ERROR log level.
func Error(err error, format string, v ...interface{}) {
	log.WithField(getCaller()).WithField("error", err).Errorf(format, v...)
}

// Fatal logs the given inputs and then calls os.Exit(1).
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
