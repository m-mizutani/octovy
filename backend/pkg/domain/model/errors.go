package model

import "github.com/m-mizutani/goerr"

var (
	// Data validation
	ErrInvalidInputValues = goerr.New("Invalid input values")

	// System data validation
	ErrInvalidSecretValues = goerr.New("Unexpected values in SecretsManager")

	// Generic system error
	ErrSystem = goerr.New("System error")

	// API
	ErrInvalidWebhookData   = goerr.New("Invalid webhook data")
	ErrAuthenticationFailed = goerr.New("Authentication failed")
	ErrUserNotFound         = goerr.New("User not found")
)
