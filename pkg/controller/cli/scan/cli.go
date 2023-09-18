package scan

import (
	"github.com/urfave/cli/v2"
)

func New() *cli.Command {
	var (
		dir string
	)
	return &cli.Command{
		Name:  "scan",
		Usage: "Local scan mode",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dir",
				Usage:       "Target directory",
				Value:       ".",
				Destination: &dir,
			},
		},
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}
