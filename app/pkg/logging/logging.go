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
	l log.Logger
}

func (lo *logger) Log(keyvals ...interface{}) error {
	return lo.l.Log(keyvals...)
}

func (lo *logger) Debug(keyvals ...interface{}) error {
	return level.Debug(lo.l).Log(keyvals...)
}

func (lo *logger) Info(keyvals ...interface{}) error {
	return level.Info(lo.l).Log(keyvals...)
}

func (lo *logger) Warn(keyvals ...interface{}) error {
	return level.Warn(lo.l).Log(keyvals...)
}

func (lo *logger) Error(keyvals ...interface{}) error {
	return level.Error(lo.l).Log(keyvals...)
}

type Logger interface {
	Log(keyvals ...interface{}) error
	Debug(keyvals ...interface{}) error
	Info(keyvals ...interface{}) error
	Warn(keyvals ...interface{}) error
	Error(keyvals ...interface{}) error
}

func Setup(level string) {
	logLevel = level
}

func GetLogger() Logger {
	l := log.NewLogfmtLogger(writer)
	l = level.NewFilter(l, level.Allow(level.ParseDefault(logLevel, level.InfoValue())))
	return &logger{l: l}
}
