package usecase

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr"
)

func (x *usecase) initErrorHandler() error {
	if x.config.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         x.config.SentryDSN,
			Environment: x.config.SentryEnv,
		})
		if err != nil {
			return goerr.Wrap(err)
		}
		logger.Debug().Str("dsn", x.config.SentryDSN).Str("env", x.config.SentryEnv).Msg("sentry initialized")
	}
	return nil
}

func (x *usecase) flushError() {
	sentry.Flush(2 * time.Second)
	logger.Debug().Msg("sentry flushed")
}

// HandleError handles a notable error. Logging error and send it to sentry if configured. It should handle an error caused by system, not a user.
func (x *usecase) HandleError(err error) {
	// Logging
	entry := logger.Error()
	var gerr *goerr.Error
	if errors.As(err, &gerr) {
		for key, value := range gerr.Values() {
			entry = entry.Interface(key, value)
		}
		entry = entry.Interface("stacktrace", gerr.Stacks())
	}

	if x.config.SentryDSN != "" {
		evID := sentry.CaptureException(err)
		entry = entry.Interface("sentry.EventID", evID)
	}

	entry.Msg(err.Error())

	if x.testErrorHandler != nil {
		x.testErrorHandler(err)
	}
}
