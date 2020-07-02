package logger

import (
	"os"

	"github.com/seashell/drago/pkg/logger"
	"github.com/sirupsen/logrus"
)

const (
	Info  = "INFO"
	Warn  = "WARN"
	Debug = "DEBUG"
	Error = "ERROR"
	Fatal = "FATAL"
)

type logrusLogEntry struct {
	entry *logrus.Entry
}

type logrusLogger struct {
	logger *logrus.Logger
}

// Configuration : Logger configuration
type Configuration struct {
	Level string
}

// New : Create a new Logger
func New(c Configuration) (logger.Logger, error) {
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		return nil, err
	}

	formatter := &logrus.TextFormatter{}

	l := &logrus.Logger{
		Out:       os.Stdout,
		Level:     level,
		Formatter: formatter,
	}

	return &logrusLogger{logger: l}, nil
}

// Debugf :
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof :
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warnf :
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Errorf :
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatalf :
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Panicf :
func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

// WithFields :
func (l *logrusLogger) WithFields(fields logger.Fields) logger.Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

// Debugf :
func (l *logrusLogEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Infof :
func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warnf :
func (l *logrusLogEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Errorf :
func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatalf :
func (l *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// Panicf :
func (l *logrusLogEntry) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

// WithFields :
func (l *logrusLogEntry) WithFields(fields logger.Fields) logger.Logger {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func convertToLogrusFields(fields logger.Fields) logrus.Fields {
	logrusFields := logrus.Fields{}

	for index, val := range fields {
		logrusFields[index] = val
	}

	return logrusFields
}
