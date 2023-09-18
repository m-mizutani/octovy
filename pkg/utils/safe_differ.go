package utils

import (
	"database/sql"
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

func SafeRollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		logger.Warn("Fail to rollback transaction", slog.Any("error", err))
	}
}
