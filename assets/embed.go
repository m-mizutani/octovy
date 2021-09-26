//go:build !github
// +build !github

package assets

import "embed"

//go:embed out/*
//go:embed out/_next/static/*/*
//go:embed out/_next/static/chunks/pages/*.js
//go:embed out/_next/static/chunks/pages/*/*.js
var assets embed.FS

func Assets() *embed.FS {
	return &assets
}
