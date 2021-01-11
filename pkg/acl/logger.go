package acl

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/seashell/drago/pkg/log"
)

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
	levelFatal = "FATAL"
	levelPanic = "PANIC"
)

var levels = map[string]int{
	levelDebug: 5,
	levelInfo:  4,
	levelWarn:  3,
	levelError: 2,
	levelFatal: 1,
	levelPanic: 0,
}

type simpleLogger struct {
	name    string
	fields  log.Fields
	options log.LoggerOptions
}

func (l simpleLogger) Logf(level string, format string, args ...interface{}) {
	if l.isLevelEnabled(level) {
		fmt.Printf("%s%s", l.options.Prefix, fmt.Sprintf(format, args...))
	}
}

func (l simpleLogger) Debugf(format string, args ...interface{}) {
	l.Logf(levelDebug, format, args...)
}

func (l simpleLogger) Infof(format string, args ...interface{}) {
	l.Logf(levelInfo, format, args...)
}

func (l simpleLogger) Warnf(format string, args ...interface{}) {
	l.Logf(levelWarn, format, args...)
}

func (l simpleLogger) Errorf(format string, args ...interface{}) {
	l.Logf(levelError, format, args...)
}

func (l simpleLogger) Fatalf(format string, args ...interface{}) {
	l.Logf(levelFatal, format, args...)
}

func (l simpleLogger) Panicf(format string, args ...interface{}) {
	l.Logf(levelPanic, format, args...)
}

func (l simpleLogger) WithFields(fields log.Fields) log.Logger {
	return &simpleLogger{
		name:    l.name,
		fields:  fields,
		options: l.options,
	}
}

func (l simpleLogger) WithName(name string) log.Logger {
	return &simpleLogger{
		name:    name,
		fields:  l.fields,
		options: l.options,
	}
}

func (l *simpleLogger) isLevelEnabled(level string) bool {
	return levels[l.options.Level] >= levels[level]
}

func mergeFields(a log.Fields, b log.Fields) log.Fields {

	res := log.Fields{}

	for k, v := range a {
		res[k] = v
	}

	if err := mergo.Merge(&a, b, mergo.WithOverride); err != nil {
		return res
	}

	return res
}
