package config

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
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
