package usecase

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/utils"
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
		utils.Logger.With("dsn", x.config.SentryDSN).With("env", x.config.SentryEnv).Debug("sentry initialized")
	}
	return nil
}

func (x *usecase) flushError() {
	sentry.Flush(2 * time.Second)
	utils.Logger.Debug("sentry flushed")
}

// HandleError handles a notable error. Logging error and send it to sentry if configured. It should handle an error caused by system, not a user.
func (x *usecase) HandleError(ctx *model.Context, err error) {
	// Logging
	entry := ctx.Log()
	var gerr *goerr.Error
	if errors.As(err, &gerr) {
		for key, value := range gerr.Values() {
			entry = entry.With(key, value)
		}
		entry = entry.With("stacktrace", gerr.Stacks())
	}

	if x.config.SentryDSN != "" {
		evID := sentry.CaptureException(err)
		entry = entry.With("sentry.EventID", evID)
	}

	entry.Error(err.Error())

	if x.testErrorHandler != nil {
		x.testErrorHandler(err)
	}
}
