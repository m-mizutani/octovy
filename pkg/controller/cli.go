package controller

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/controller/api"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/urfave/cli/v2"
)

var logger = utils.Logger

func (x *Controller) RunCmd(args []string) error {
	app := &cli.App{
		Name:        "octovy",
		Description: "Utility command of octovy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				EnvVars: []string{"OCTOVY_LOG_LEVEL"},
				Usage:   "LogLevel [trace|debug|info|warn|error]",
			},
		},
		Commands: []*cli.Command{
			newServeCommand(x),
		},
		Before: globalSetup,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error().Err(err).Msg("Failed")
		return err
	}

	return nil
}

func globalSetup(c *cli.Context) error {
	// Setup logger
	if err := utils.SetLogLevel(c.String("log-level")); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

func newServeCommand(ctrl *Controller) *cli.Command {

	return &cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "Addr",
				Usage:       "server binding address",
				Aliases:     []string{"a"},
				EnvVars:     []string{"OCTOVY_ADDR"},
				Destination: &ctrl.Config.ServerAddr,
				Value:       "127.0.0.1",
			},
			&cli.IntFlag{
				Name:        "Port",
				Usage:       "Port number",
				Aliases:     []string{"p"},
				EnvVars:     []string{"OCTOVY_PORT"},
				Destination: &ctrl.Config.ServerPort,
				Value:       9080,
			},

			&cli.StringFlag{
				Name:        "frontend-url",
				EnvVars:     []string{"OCTOVY_FRONTEND_URL"},
				Destination: &ctrl.Config.FrontendURL,
			},

			&cli.Int64Flag{
				Name:        "github-app-id",
				EnvVars:     []string{"OCTOVY_GITHUB_APP_ID"},
				Destination: &ctrl.Config.GitHubAppID,
			},
			&cli.PathFlag{
				Name:        "github-app-pem",
				EnvVars:     []string{"OCTOVY_GITHUB_APP_PRIVATE_KEY"},
				Usage:       "GitHub App private key file path",
				Destination: &ctrl.Config.GitHubAppPrivateKeyPath,
			},
			&cli.Int64Flag{
				Name:        "github-app-client-id",
				EnvVars:     []string{"OCTOVY_GITHUB_CLIENT_ID"},
				Destination: &ctrl.Config.GitHubAppClientID,
			},
			&cli.StringFlag{
				Name:        "github-app-client-id",
				EnvVars:     []string{"OCTOVY_GITHUB_SECRET"},
				Destination: &ctrl.Config.GitHubAppSecret,
			},
		},
		Action: func(c *cli.Context) error {
			return apiCommand(c, ctrl)
		},
	}
}

func apiCommand(c *cli.Context, ctrl *Controller) error {
	serverAddr := fmt.Sprintf("%s:%d", ctrl.Config.ServerAddr, ctrl.Config.ServerPort)

	engine := api.New(ctrl.usecase)

	gin.SetMode(gin.DebugMode)
	logger.Info().Interface("config", ctrl.Config).Msg("Starting server...")
	if err := engine.Run(serverAddr); err != nil {
		logger.Error().Err(err).Interface("config", ctrl.Config).Msg("Server error")
	}

	return nil
}
