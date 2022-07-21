package logging

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var writer = log.NewSyncWriter(os.Stderr)
var logLevel string = "info"

const (
	Error = "error"
	Warn  = "warn"
	Info  = "info"
	Debug = "debug"
)

type logger struct {
	logger log.Logger
}

func (lo *logger) Log(keyvals ...interface{}) error {
	return lo.logger.Log(keyvals...)
}

func (lo *logger) Debug(keyvals ...interface{}) error {
	return level.Debug(lo.logger).Log(keyvals...)
}

func (lo *logger) Info(keyvals ...interface{}) error {
	return level.Info(lo.logger).Log(keyvals...)
}

func (lo *logger) Warn(keyvals ...interface{}) error {
	return level.Warn(lo.logger).Log(keyvals...)
}

func (lo *logger) Error(keyvals ...interface{}) error {
	return level.Error(lo.logger).Log(keyvals...)
}

func (lo *logger) Fatal(keyvals ...interface{}) {
	log.WithPrefix(lo.logger, level.Key(), "fatal")
	os.Exit(1)
}

func (lo *logger) Panic(keyvals ...interface{}) {
	log.WithPrefix(lo.logger, level.Key(), "panic")
	os.Exit(1)
}

type Logger interface {
	Log(keyvals ...interface{}) error
	Debug(keyvals ...interface{}) error
	Info(keyvals ...interface{}) error
	Warn(keyvals ...interface{}) error
	Error(keyvals ...interface{}) error
	Fatal(keyvals ...interface{})
	Panic(keyvals ...interface{})
}

func Setup(level string) {
	logLevel = level
}

func NewLogger() Logger {
	l := log.NewLogfmtLogger(writer)
	l = level.NewFilter(l, level.Allow(level.ParseDefault(logLevel, level.InfoValue())))
	return &logger{logger: l}
}
