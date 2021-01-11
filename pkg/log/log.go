package log

// Fields : Log fields
type Fields map[string]interface{}

// Options contains logging options which are
// common to any logger
type LoggerOptions struct {
	Prefix string
	Level  string
}

// Logger : Logger interface
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(fields Fields) Logger
	WithName(name string) Logger
}
