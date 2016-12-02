package log

// Fields for logging
type Fields map[string]interface{}

// Interface for logging in TTN
type Interface interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	Fatalf(msg string, v ...interface{})
	WithField(string, interface{}) Interface
	WithFields(Fields) Interface
	WithError(error) Interface
}

var defaultLogger Interface = noopLogger{}

// Get returns the defaultLogger logger
func Get() Interface {
	return defaultLogger
}

// Set sets the default logger
func Set(log Interface) {
	defaultLogger = log
}

// noopLogger just does nothing
type noopLogger struct{}

func (l noopLogger) Debug(msg string)                        {}
func (l noopLogger) Info(msg string)                         {}
func (l noopLogger) Warn(msg string)                         {}
func (l noopLogger) Error(msg string)                        {}
func (l noopLogger) Fatal(msg string)                        {}
func (l noopLogger) Debugf(msg string, v ...interface{})     {}
func (l noopLogger) Infof(msg string, v ...interface{})      {}
func (l noopLogger) Warnf(msg string, v ...interface{})      {}
func (l noopLogger) Errorf(msg string, v ...interface{})     {}
func (l noopLogger) Fatalf(msg string, v ...interface{})     {}
func (l noopLogger) WithField(string, interface{}) Interface { return l }
func (l noopLogger) WithFields(Fields) Interface             { return l }
func (l noopLogger) WithError(error) Interface               { return l }
