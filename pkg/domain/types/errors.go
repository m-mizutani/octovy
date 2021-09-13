package types

import "github.com/m-mizutani/goerr"

var (
	ErrDatabaseUnexpected   = goerr.New("database failure")
	ErrDatabaseInvalidInput = goerr.New("invalid input for database")
	ErrItemNotFound         = goerr.New("item not found")
	ErrInvalidChain         = goerr.New("invalid chain plugin")
	ErrInvalidInput         = goerr.New("invalid input data")
)
