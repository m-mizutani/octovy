package cli

import (
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
	)

	app := &cli.App{
		Name:  "octovy",
		Usage: "Vulnerability management system with Trivy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "Log level [trace|debug|info|warn|error]",
				Destination: &logLevel,
				Value:       "info",
			},
			&cli.StringFlag{
				Name:        "log-format",
				Usage:       "Log format [text|json]",
				Destination: &logFormat,
				Value:       "text",
			},
		},
		Commands: []*cli.Command{
			serve.New(),
			scan.New(),
		},
		Before: func(ctx *cli.Context) error {
			if err := utils.ReconfigureLogger(logFormat, logLevel); err != nil {
				return err
			}
			return nil
		},
	}

	if err := app.Run(argv); err != nil {
		utils.Logger().Error("fatal error", err)
		return err
	}

	return nil
}
