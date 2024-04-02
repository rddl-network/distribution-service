package service

import (
	"fmt"
	"os"
	"slices"

	stdLog "log"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type AppLogger struct {
	logger log.Logger
}

func getLogger(logLevel string) AppLogger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// logLevel should be set to "debug", "info", "warn", or "error"
	if !slices.Contains([]string{"debug", "info", "warn", "error"}, logLevel) {
		stdLog.Panicln("logLevel should be set to debug, info, warn or error")
	}

	// Set log level
	switch logLevel {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	default:
		logger = level.NewFilter(logger, level.AllowError())
	}

	return AppLogger{logger}
}

func (al AppLogger) Error(kevals ...interface{}) {
	err := level.Error(al.logger).Log(kevals...)
	if err != nil {
		fmt.Println(err)
	}
}

func (al AppLogger) Warn(kevals ...interface{}) {
	err := level.Warn(al.logger).Log(kevals...)
	if err != nil {
		fmt.Println(err)
	}
}

func (al AppLogger) Info(kevals ...interface{}) {
	err := level.Info(al.logger).Log(kevals...)
	if err != nil {
		fmt.Println(err)
	}
}

func (al AppLogger) Debug(kevals ...interface{}) {
	err := level.Debug(al.logger).Log(kevals...)
	if err != nil {
		fmt.Println(err)
	}
}
