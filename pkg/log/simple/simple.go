package simple

import (
	"fmt"

	log "github.com/seashell/drago/pkg/log"
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
	name   string
	fields map[string]interface{}
}

// Config contains the configuration for the logger adapter
type Config struct {
	log.LoggerOptions
}

// NewLoggerAdapter creates a new Logger adapter
func NewLoggerAdapter(config Config) (log.Logger, error) {
	return &logger{config: config}, nil
}

// Debugf :
func (l *logger) Debugf(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+l.name+": "+l.config.Prefix+format+"\n", args...)
}

// Infof :
func (l *logger) Infof(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+l.name+": "+l.config.Prefix+format+"\n", args...)
}

// Warnf :
func (l *logger) Warnf(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+l.name+": "+l.config.Prefix+format+"\n", args...)
}

// Errorf :
func (l *logger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+l.name+": "+l.config.Prefix+format+"\n", args...)
}

// Fatalf :
func (l *logger) Fatalf(format string, args ...interface{}) {
	fmt.Printf("[FATAL] "+l.name+": "+l.config.Prefix+format+"\n", args...)
}

// Panicf :
func (l *logger) Panicf(format string, args ...interface{}) {
	fmt.Printf("[PANIC] "+l.name+": "+l.config.Prefix+format+"\n", args...)
}

// WithFields :
func (l *logger) WithFields(fields log.Fields) log.Logger {

	nl := &logger{
		fields: l.fields,
	}

	for k, v := range fields {
		nl.fields[k] = v
	}

	return nl
}

// WithName :
func (l *logger) WithName(name string) log.Logger {
	return &logger{
		name: name,
	}
}
