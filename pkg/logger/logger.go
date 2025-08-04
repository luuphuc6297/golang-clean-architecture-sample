package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Warn(args ...any)
	Debug(args ...any)
	WithField(key string, value any) Logger
	WithError(err error) Logger
}

type logger struct {
	logrus *logrus.Logger
}

func NewLogger() Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return &logger{logrus: log}
}

func (l *logger) Info(args ...any) {
	l.logrus.Info(args...)
}

func (l *logger) Error(args ...any) {
	l.logrus.Error(args...)
}

func (l *logger) Fatal(args ...any) {
	l.logrus.Fatal(args...)
}

func (l *logger) Warn(args ...any) {
	l.logrus.Warn(args...)
}

func (l *logger) Debug(args ...any) {
	l.logrus.Debug(args...)
}

func (l *logger) WithField(key string, value any) Logger {
	return &logger{logrus: l.logrus.WithField(key, value).Logger}
}

func (l *logger) WithError(err error) Logger {
	return &logger{logrus: l.logrus.WithError(err).Logger}
}
