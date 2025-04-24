package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger представляет обёртку над logrus.Logger для централизованного логирования.
type Logger struct {
	log *logrus.Logger
}

// NewLogger создаёт новый экземпляр логгера и настраивает его по переданному уровню логирования.
func NewLogger(level string) *Logger {
	l := logrus.New()

	// Вывод логов в стандартный вывод.
	l.Out = os.Stdout

	// Устанавливаем уровень логирования
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	l.SetLevel(parsedLevel)

	// Настройка формата логов (можно изменить на JSON или другой, если требуется)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	return &Logger{l}
}

// Info записывает информационное сообщение.
func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

// Infof записывает информационное сообщение с форматированием.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

// Debug записывает отладочное сообщение.
func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

// Debugf записывает отладочное сообщение с форматированием.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

// Error записывает сообщение об ошибке.
func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

// Errorf записывает сообщение об ошибке с форматированием.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
