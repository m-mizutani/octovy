package usecase

import (
	"errors"

	"github.com/m-mizutani/goerr"
)

func (x *usecase) handleError(err error) {
	// Logging
	entry := logger.Error()
	var gerr *goerr.Error
	if errors.As(err, &gerr) {
		for key, value := range gerr.Values() {
			entry = entry.Interface(key, value)
		}
		entry = entry.Interface("stacktrace", gerr.Stacks())
	}

	entry.Msg(err.Error())

	if x.testErrorHandler != nil {
		x.testErrorHandler(err)
	}
}
