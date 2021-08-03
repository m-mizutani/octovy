package model

import "github.com/m-mizutani/goerr"

var (
	// Internal error
	ErrNoAuthenticatedClient = goerr.New("No authenticated client, token required")

	// Data validation
	// ErrInvalidValue: This is caused by user or system
	ErrInvalidValue = goerr.New("Invalid input value")
	// ErrInvalidSystemValue: This is caused by only system
	ErrInvalidSystemValue = goerr.New("Invalid system value")

	// System data validation
	ErrInvalidSecretValues = goerr.New("Unexpected values in SecretsManager")

	// Generic system error
	ErrSystem = goerr.New("System error")

	// API
	ErrInvalidWebhookData   = goerr.New("Invalid webhook data")
	ErrAuthenticationFailed = goerr.New("Authentication failed")
	ErrUserNotFound         = goerr.New("User not found")

	// Fatal error, system can not be recovered and should go to shutdown
	ErrFatal = goerr.New("Fatal error")
)
