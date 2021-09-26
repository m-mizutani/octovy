package server

import "github.com/m-mizutani/goerr"

var (
	// API error
	errAPIInvalidParameter = goerr.New("Invalid API parameter")

	errResourceNotFound = goerr.New("Resource not found")
)
