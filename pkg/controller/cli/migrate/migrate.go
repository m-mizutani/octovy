package migrate

import (
	"github.com/m-mizutani/gots/slice"
	"github.com/urfave/cli/v2"

	"github.com/m-mizutani/octovy/pkg/controller/cli/config"
)

func New() *cli.Command {
	var (
		dryRun   bool
		dbConfig config.DB
	)

	baseFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:        "dry-run",
			Usage:       "Dry run mode",
			Destination: &dryRun,
		},
	}

	return &cli.Command{
		Name:    "migrate",
		Aliases: []string{"m"},
		Usage:   "Migrate database",
		Flags: slice.Flatten(
			baseFlags,
			dbConfig.Flags(),
		),
		Action: func(c *cli.Context) error {
			return dbConfig.Migrate(dryRun)
		},
	}
}
