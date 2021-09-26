//go:build github
// +build github

package assets

import "embed"

var assets embed.FS

func Assets() *embed.FS {
	return &assets
}
