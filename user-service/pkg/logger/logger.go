package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger represents a wrapper over logrus.Logger for centralized logging.
type Logger struct {
	log *logrus.Logger
}

// NewLogger creates a new logger instance and configures it according to the provided logging level.
func NewLogger(level string) *Logger {
	l := logrus.New()

	// Output logs to standard output.
	l.Out = os.Stdout

	// Set logging level
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	l.SetLevel(parsedLevel)

	// Configure log format (can be changed to JSON or other if needed)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	return &Logger{l}
}

// Info writes an informational message.
func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

// Infof writes a formatted informational message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

// Debug writes a debug message.
func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

// Debugf writes a formatted debug message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

// Error writes an error message.
func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

// Errorf writes a formatted error message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
