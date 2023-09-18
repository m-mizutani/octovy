package types

import "github.com/m-mizutani/goerr"

var (
	ErrInvalidOption = goerr.New("invalid option")

	ErrInvalidRequest = goerr.New("invalid request")

	ErrInvalidGitHubData = goerr.New("invalid GitHub data")

	ErrLogicError = goerr.New("logic error")
)
