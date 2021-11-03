package utils

import (
	"os"

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
		Logger.Emitter = zlog.NewWriterWith(zlog.NewConsoleFormatter(), os.Stdout)
	case "json":
		Logger.Emitter = zlog.NewWriterWith(zlog.NewJsonFormatter(), os.Stdout)
	default:
		return goerr.New("invalid log format: " + logFormat)
	}
	return nil
}
