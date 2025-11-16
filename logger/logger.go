package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Log is the singleton logger instance
	Log *logrus.Logger

	// DefaultFields are fields that will be added to every log entry
	DefaultFields logrus.Fields
)

// Init initializes the logger with default settings
func Init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true, // Show full level names (DEBUG, ERROR instead of DEBU, ERRO)
	})
	DefaultFields = make(logrus.Fields)
}

// InitJSON initializes the logger with JSON formatter (useful for production)
func InitJSON() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	DefaultFields = make(logrus.Fields)
}

// SetLevel sets the log level (panic, fatal, error, warn, info, debug, trace)
func SetLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	if Log == nil {
		Init()
	}
	Log.SetLevel(logLevel)
	return nil
}

// SetDefaultFields sets default fields that will be added to every log entry
func SetDefaultFields(fields logrus.Fields) {
	if DefaultFields == nil {
		DefaultFields = make(logrus.Fields)
	}
	for k, v := range fields {
		DefaultFields[k] = v
	}
}

// WithFields creates an entry with the default fields plus the provided fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	if Log == nil {
		Init()
	}
	allFields := make(logrus.Fields)
	for k, v := range DefaultFields {
		allFields[k] = v
	}
	for k, v := range fields {
		allFields[k] = v
	}
	return Log.WithFields(allFields)
}

// WithField creates an entry with the default fields plus a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return WithFields(logrus.Fields{key: value})
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Debug(args...)
	} else {
		Log.Debug(args...)
	}
}

// Debugf logs a formatted message at level Debug
func Debugf(format string, args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Debugf(format, args...)
	} else {
		Log.Debugf(format, args...)
	}
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Info(args...)
	} else {
		Log.Info(args...)
	}
}

// Infof logs a formatted message at level Info
func Infof(format string, args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Infof(format, args...)
	} else {
		Log.Infof(format, args...)
	}
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Warn(args...)
	} else {
		Log.Warn(args...)
	}
}

// Warnf logs a formatted message at level Warn
func Warnf(format string, args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Warnf(format, args...)
	} else {
		Log.Warnf(format, args...)
	}
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Error(args...)
	} else {
		Log.Error(args...)
	}
}

// Errorf logs a formatted message at level Error
func Errorf(format string, args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Errorf(format, args...)
	} else {
		Log.Errorf(format, args...)
	}
}

// Fatal logs a message at level Fatal and exits
func Fatal(args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Fatal(args...)
	} else {
		Log.Fatal(args...)
	}
}

// Fatalf logs a formatted message at level Fatal and exits
func Fatalf(format string, args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Fatalf(format, args...)
	} else {
		Log.Fatalf(format, args...)
	}
}

// Panic logs a message at level Panic and panics
func Panic(args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Panic(args...)
	} else {
		Log.Panic(args...)
	}
}

// Panicf logs a formatted message at level Panic and panics
func Panicf(format string, args ...interface{}) {
	if Log == nil {
		Init()
	}
	if len(DefaultFields) > 0 {
		Log.WithFields(DefaultFields).Panicf(format, args...)
	} else {
		Log.Panicf(format, args...)
	}
}

