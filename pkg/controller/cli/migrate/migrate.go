package migrate

import (
	"log/slog"
	"os"

	"github.com/k0kubun/sqldef"
	"github.com/k0kubun/sqldef/database"
	"github.com/k0kubun/sqldef/database/postgres"
	"github.com/k0kubun/sqldef/schema"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gots/slice"
	"github.com/urfave/cli/v2"

	db "github.com/m-mizutani/octovy/database"
	"github.com/m-mizutani/octovy/pkg/controller/cli/config"
	"github.com/m-mizutani/octovy/pkg/utils"
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
		Name:  "migrate",
		Usage: "Migrate database",
		Flags: slice.Flatten(
			baseFlags,
			dbConfig.Flags(),
		),
		Action: func(c *cli.Context) error {
			utils.Logger().Info("migrating database", slog.Any("config.DB", dbConfig))

			options := &sqldef.Options{
				DesiredDDLs:     db.Schema(),
				DryRun:          dryRun,
				Export:          false,
				EnableDropTable: true,
				// BeforeApply:     opts.BeforeApply,
				// Config: database.ParseGeneratorConfig(opts.Config),
			}

			config := database.Config{
				DbName:   dbConfig.DBName,
				User:     dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     dbConfig.Port,
			}
			os.Setenv("PGSSLMODE", dbConfig.SSLMode)

			db, err := postgres.NewDatabase(config)
			if err != nil {
				return goerr.Wrap(err, "failed to open database")
			}
			defer utils.SafeClose(db)

			sqlParser := postgres.NewParser()
			sqldef.Run(schema.GeneratorModePostgres, db, sqlParser, options)

			return nil
		},
	}
}
