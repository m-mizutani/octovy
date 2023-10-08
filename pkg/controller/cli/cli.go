package cli

import (
	"github.com/m-mizutani/octovy/pkg/controller/cli/migrate"
	"github.com/m-mizutani/octovy/pkg/controller/cli/scan"
	"github.com/m-mizutani/octovy/pkg/controller/cli/serve"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/urfave/cli/v2"
)

type CLI struct {
}

func New() *CLI {
	return &CLI{}
}

func (x *CLI) Run(argv []string) error {
	var (
		logLevel  string
		logFormat string
		logOutput string
	)

	app := &cli.App{
		Name:  "octovy",
		Usage: "Vulnerability management system with Trivy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "Log level [trace|debug|info|warn|error]",
				Aliases:     []string{"l"},
				EnvVars:     []string{"OCTOVY_LOG_LEVEL"},
				Destination: &logLevel,
				Value:       "info",
			},
			&cli.StringFlag{
				Name:        "log-format",
				Usage:       "Log format [text|json]",
				Aliases:     []string{"f"},
				EnvVars:     []string{"OCTOVY_LOG_FORMAT"},
				Destination: &logFormat,
				Value:       "text",
			},
			&cli.StringFlag{
				Name:        "log-output",
				Usage:       "Log output [-|stdout|stderr|<file>]",
				Aliases:     []string{"o"},
				EnvVars:     []string{"OCTOVY_LOG_OUTPUT"},
				Destination: &logOutput,
				Value:       "-",
			},
		},
		Commands: []*cli.Command{
			serve.New(),
			scan.New(),
			migrate.New(),
		},
		Before: func(ctx *cli.Context) error {
			if err := utils.ReconfigureLogger(logFormat, logLevel, logOutput); err != nil {
				return err
			}
			return nil
		},
	}

	if err := app.Run(argv); err != nil {
		utils.Logger().Error("fatal error", "error", err)
		return err
	}

	return nil
}
