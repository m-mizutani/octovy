package api

import "github.com/m-mizutani/goerr"

var (
	// API error
	errAPIInvalidParameter = goerr.New("Invalid API parameter")
	errInvalidWebhookData  = goerr.New("Invalid webhook data")
)
