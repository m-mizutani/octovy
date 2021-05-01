package model

import "github.com/m-mizutani/goerr"

var (
	// Data validation
	ErrInvalidScanRequest = goerr.New("Invalid repository scan request")

	// System data validation
	ErrInvalidSecretValues = goerr.New("Unexpected values in SecretsManager")

	// Generic system error
	ErrSystem = goerr.New("System error")
)
