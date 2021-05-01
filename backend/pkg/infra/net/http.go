package net

import (
	"net/http"
)

func NewHTTP(itr http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: itr,
	}
}
