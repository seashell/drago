package logrus

import (
	"os"

	log "github.com/seashell/drago/pkg/log"
	logrus "github.com/sirupsen/logrus"
)

const (
	Info  = "INFO"
	Warn  = "WARN"
	Debug = "DEBUG"
	Error = "ERROR"
	Fatal = "FATAL"
)

type logEntry struct {
	entry *logrus.Entry
}

type logger struct {
	config Config
	logger *logrus.Logger
}

// Config contains the configuration for the logger adapter
type Config struct {
	log.LoggerOptions
}

// NewLoggerAdapter creates a new Logger adapter
func NewLoggerAdapter(config Config) (log.Logger, error) {
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	formatter := &logrus.TextFormatter{}

	l := &logrus.Logger{
		Out:       os.Stdout,
		Level:     level,
		Formatter: formatter,
	}

	return &logger{config: config, logger: l}, nil
}

// Debugf :
func (l *logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(l.config.Prefix+format, args...)
}

// Infof :
func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(l.config.Prefix+format, args...)
}

// Warnf :
func (l *logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(l.config.Prefix+format, args...)
}

// Errorf :
func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(l.config.Prefix+format, args...)
}

// Fatalf :
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(l.config.Prefix+format, args...)
}

// Panicf :
func (l *logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(l.config.Prefix+format, args...)
}

// WithFields :
func (l *logger) WithFields(fields log.Fields) log.Logger {
	return &logEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

// WithName :
func (l *logger) WithName(name string) log.Logger {
	return &logEntry{
		entry: l.logger.WithField("name", name),
	}
}

// Debugf :
func (l *logEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Infof :
func (l *logEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warnf :
func (l *logEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Errorf :
func (l *logEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatalf :
func (l *logEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// Panicf :
func (l *logEntry) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

// WithFields :
func (l *logEntry) WithFields(fields log.Fields) log.Logger {
	return &logEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

// WithName :
func (l *logEntry) WithName(name string) log.Logger {
	return &logEntry{
		entry: l.entry.WithField("name", name),
	}
}

func convertToLogrusFields(fields log.Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}

	return logrusFields
}
