package types

import "errors"

var (
	// ErrInvalidOption is an error that indicates an invalid option is given by user via CLI or configuration
	ErrInvalidOption = errors.New("invalid option")

	// ErrInvalidRequest is an error that indicates an invalid HTTP request
	ErrInvalidRequest = errors.New("invalid request")

	// ErrInvalidResponse is an error that indicates a failure in data consistency in the application
	ErrValidationFailed = errors.New("validation failed")

	// ErrInvalidGitHubData is an error that indicates an invalid data provided by GitHub. Mainly used in GitHub API response
	ErrInvalidGitHubData = errors.New("invalid GitHub data")

	// ErrLogicError is an error that indicates a logic error in the application
	ErrLogicError = errors.New("logic error")
)
