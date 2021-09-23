package model

import "github.com/m-mizutani/goerr"

var (
	ErrDatabaseUnexpected   = goerr.New("database failure")
	ErrDatabaseInvalidInput = goerr.New("invalid input for database")
	ErrItemNotFound         = goerr.New("item not found")
	ErrInvalidSystemValue   = goerr.New("invalid system value")

	ErrInvalidInput         = goerr.New("invalid input data")
	ErrAuthenticationFailed = goerr.New("authentication failed")
	ErrUserNotFound         = goerr.New("user not found")
	ErrInvalidWebhookData   = goerr.New("invalid webhook data")
	ErrGitHubAPI            = goerr.New("github API returns unexpected response")
)
