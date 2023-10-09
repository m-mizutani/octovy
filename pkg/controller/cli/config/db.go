package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/k0kubun/sqldef"
	"github.com/k0kubun/sqldef/database"
	"github.com/k0kubun/sqldef/database/postgres"
	"github.com/k0kubun/sqldef/schema"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/urfave/cli/v2"

	db "github.com/m-mizutani/octovy/database"
)

type DB struct {
	User     string
	Password string `masq:"secret"`
	Host     string
	Port     int
	DBName   string
	SSLMode  string
}

func (x *DB) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "db-user",
			Category:    "Database",
			Usage:       "database user",
			EnvVars:     []string{"OCTOVY_DB_USER"},
			Required:    true,
			Destination: &x.User,
		},
		&cli.StringFlag{
			Name:        "db-password",
			Category:    "Database",
			Usage:       "database password",
			EnvVars:     []string{"OCTOVY_DB_PASSWORD"},
			Destination: &x.Password,
		},
		&cli.StringFlag{
			Name:        "db-host",
			Category:    "Database",
			Usage:       "database host",
			EnvVars:     []string{"OCTOVY_DB_HOST"},
			Destination: &x.Host,
			Value:       "localhost",
		},
		&cli.IntFlag{
			Name:        "db-port",
			Category:    "Database",
			Usage:       "database port",
			EnvVars:     []string{"OCTOVY_DB_PORT"},
			Destination: &x.Port,
			Value:       5432,
		},
		&cli.StringFlag{
			Name:        "db-name",
			Category:    "Database",
			Usage:       "database name",
			EnvVars:     []string{"OCTOVY_DB_NAME"},
			Required:    true,
			Destination: &x.DBName,
		},
		&cli.StringFlag{
			Name:        "db-ssl-mode",
			Category:    "Database",
			Usage:       "database SSL mode",
			EnvVars:     []string{"OCTOVY_DB_SSL_MODE"},
			Destination: &x.SSLMode,
			Value:       "disable",
		},
	}
}

func (x *DB) DSN() string {
	var options []string
	if x.User != "" {
		options = append(options, "user="+x.User)
	}
	if x.Password != "" {
		options = append(options, "password="+x.Password)
	}
	if x.Host != "" {
		options = append(options, "host="+x.Host)
	}
	if x.Port != 0 {
		options = append(options, "port="+fmt.Sprintf("%d", x.Port))
	}
	if x.DBName != "" {
		options = append(options, "dbname="+x.DBName)
	}
	if x.SSLMode != "" {
		options = append(options, "sslmode="+x.SSLMode)
	}

	return strings.Join(options, " ")
}

func (x *DB) Connect(ctx context.Context) (*sql.DB, error) {
	dbClient, err := sql.Open("postgres", x.DSN())
	if err != nil {
		return nil, goerr.Wrap(err, "failed to open database")
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(20*time.Second))
	defer cancel()

	if err := waitDBReady(ctx, dbClient); err != nil {
		return nil, err
	}

	return dbClient, nil
}

func waitDBReady(ctx context.Context, db *sql.DB) error {
	utils.Logger().Info("waiting database ready")
	for {
		err := db.PingContext(ctx)
		if err == nil {
			utils.Logger().Info("database is ready")
			return nil
		} else if errors.Is(err, context.DeadlineExceeded) {
			return goerr.Wrap(err, "database is not ready, and timeout")
		}
		utils.Logger().Warn("database is not ready", slog.Any("error", err))
		time.Sleep(1 * time.Second)
	}
}

func (x *DB) Migrate(dryRun bool) error {
	utils.Logger().Info("migrating database", slog.Any("config.DB", x))

	options := &sqldef.Options{
		DesiredDDLs:     db.Schema(),
		DryRun:          dryRun,
		Export:          false,
		EnableDropTable: true,
		// BeforeApply:     opts.BeforeApply,
		// Config: database.ParseGeneratorConfig(opts.Config),
	}

	config := database.Config{
		DbName:   x.DBName,
		User:     x.User,
		Password: x.Password,
		Host:     x.Host,
		Port:     x.Port,
	}
	if err := os.Setenv("PGSSLMODE", x.SSLMode); err != nil {
		return goerr.Wrap(err, "failed to set PGSSLMODE")
	}

	db, err := postgres.NewDatabase(config)
	if err != nil {
		return goerr.Wrap(err, "failed to open database")
	}
	defer utils.SafeClose(db)

	sqlParser := postgres.NewParser()
	sqldef.Run(schema.GeneratorModePostgres, db, sqlParser, options)

	return nil
}
