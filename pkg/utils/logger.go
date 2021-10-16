package utils

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
	"github.com/m-mizutani/zlog/filter"
)

var Logger *zlog.Logger

func initLogger() {
	Logger = zlog.New()
	Logger.Filters = zlog.Filters{
		filter.Tag(),
	}
}

func SetLogLevel(logLevel string) error {
	return Logger.SetLogLevel(logLevel)
}

func SetLogFormat(logFormat string) error {
	switch logFormat {
	case "console":
		Logger.Formatter = zlog.NewConsoleFormatter()
	case "json":
		Logger.Formatter = zlog.NewJsonFormatter()
	default:
		return goerr.New("invalid log format: " + logFormat)
	}
	return nil
}
