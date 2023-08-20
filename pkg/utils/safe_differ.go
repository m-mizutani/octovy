package utils

import (
	"io"
	"os"

	"log/slog"
)

func SafeClose(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			logger.Warn("Fail to close resource", slog.Any("error", err))
		}
	}
}

func SafeRemove(path string) {
	if err := os.Remove(path); err != nil {
		logger.Warn("Fail to remove file", slog.Any("error", err))
	}
}

func SafeRemoveAll(path string) {
	if err := os.RemoveAll(path); err != nil {
		logger.Warn("Fail to remove file", slog.Any("error", err))
	}
}
