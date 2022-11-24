package utils

import (
	"os"

	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout))

func Logger() *slog.Logger {
	return logger
}

func ReconfigureLogger(logFormat, logLevel string) error {
	switch logFormat {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout))
	default:
		return goerr.New("invalid log format, should be 'text' or 'json': %s", logFormat)
	}

	switch logLevel {
	case "debug":
		logger.Enabled(slog.DebugLevel)
		fallthrough
	case "info":
		logger.Enabled(slog.InfoLevel)
		fallthrough
	case "warn":
		logger.Enabled(slog.WarnLevel)
		fallthrough
	case "error":
		logger.Enabled(slog.ErrorLevel)
	default:
		return goerr.New("invalid log format, should be 'debug', 'info', 'warn' or 'error': %s", logLevel)
	}

	return nil
}
