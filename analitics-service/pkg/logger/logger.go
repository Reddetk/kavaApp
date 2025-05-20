package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger represents a wrapper over logrus.Logger for centralized logging.
type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
	Debug(ctx context.Context, msg string, args ...interface{})
}

type loggerImpl struct {
	log *logrus.Logger
}

// NewLogger creates a new logger instance and configures it according to the provided logging level.
func NewLogger(level string) Logger {
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

	return &loggerImpl{l}
}

// Info writes an informational message.
func (l *loggerImpl) Info(ctx context.Context, msg string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[args[i].(string)] = args[i+1]
		}
	}
	l.log.WithFields(fields).Info(msg)
}

// Warn writes a warning message.
func (l *loggerImpl) Warn(ctx context.Context, msg string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[args[i].(string)] = args[i+1]
		}
	}
	l.log.WithFields(fields).Warn(msg)
}

// Error writes an error message.
func (l *loggerImpl) Error(ctx context.Context, msg string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[args[i].(string)] = args[i+1]
		}
	}
	l.log.WithFields(fields).Error(msg)
}

// Debug writes a debug message.
func (l *loggerImpl) Debug(ctx context.Context, msg string, args ...interface{}) {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[args[i].(string)] = args[i+1]
		}
	}
	l.log.WithFields(fields).Debug(msg)
}