package zap

import (
	"fmt"

	log "github.com/seashell/drago/pkg/log"
	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

const (
	Info  = "INFO"
	Warn  = "WARN"
	Debug = "DEBUG"
	Error = "ERROR"
	Fatal = "FATAL"
)

type logger struct {
	config Config
	logger *zap.Logger
}

// Config contains the configuration for the logger adapter
type Config struct {
	log.LoggerOptions
}

// NewLoggerAdapter creates a new zap Logger adapter
func NewLoggerAdapter(config Config) (log.Logger, error) {

	level, err := parseZapLevel(config.Level)
	if err != nil {
		return nil, err
	}

	c := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := c.Build()
	if err != nil {
		return nil, err
	}

	return &logger{config: config, logger: l.WithOptions(zap.AddCallerSkip(1))}, nil
}

// Debugf :
func (l *logger) Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	l.logger.Debug(s)
}

// Infof :
func (l *logger) Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(l.config.Prefix+format, args...)
	l.logger.Info(s)
}

// Warnf :
func (l *logger) Warnf(format string, args ...interface{}) {
	s := fmt.Sprintf(l.config.Prefix+format, args...)
	l.logger.Warn(s)
}

// Errorf :
func (l *logger) Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(l.config.Prefix+format, args...)
	l.logger.Error(s)
}

// Fatalf :
func (l *logger) Fatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(l.config.Prefix+format, args...)
	l.logger.Fatal(s)
}

// Panicf :
func (l *logger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(l.config.Prefix+format, args...)
	l.logger.Panic(s)
}

// WithFields :
func (l *logger) WithFields(fields log.Fields) log.Logger {
	return &logger{
		logger: l.logger.With(convertToZapFields(fields)...),
	}
}

func parseZapLevel(l string) (zapcore.Level, error) {
	switch l {
	case Info:
		return zap.InfoLevel, nil
	case Warn:
		return zap.WarnLevel, nil
	case Debug:
		return zap.DebugLevel, nil
	case Error:
		return zap.ErrorLevel, nil
	case Fatal:
		return zap.FatalLevel, nil
	default:
		return 0, fmt.Errorf("unknown logging level: %s", l)
	}
}

func convertToZapFields(fields log.Fields) []zap.Field {
	zapFields := []zap.Field{}
	for index, val := range fields {
		zapFields = append(zapFields, zap.Any(index, val))
	}
	return zapFields
}
