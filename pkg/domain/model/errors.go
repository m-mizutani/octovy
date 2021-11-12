package model

import "github.com/m-mizutani/goerr"

var (
	ErrDatabaseUnexpected   = goerr.New("database failure")
	ErrDatabaseInvalidInput = goerr.New("invalid input for database")
	ErrItemNotFound         = goerr.New("item not found")
	ErrInvalidSystemValue   = goerr.New("invalid system value")

	ErrInvalidGitHubData = goerr.New("invalid github data")

	ErrInvalidInput          = goerr.New("invalid input data")
	ErrAuthenticationFailed  = goerr.New("authentication failed")
	ErrNotAuthenticated      = goerr.New("not authenticated request")
	ErrNotAuthorized         = goerr.New("not authorized request")
	ErrUserNotFound          = goerr.New("user not found")
	ErrVulnerabilityNotFound = goerr.New("vulnerability not found")
	ErrInvalidWebhookData    = goerr.New("invalid webhook data")
	ErrGitHubAPI             = goerr.New("github API returns unexpected response")

	// Rule error
	ErrInvalidPolicyResult = goerr.New("invalid rule result")
)
