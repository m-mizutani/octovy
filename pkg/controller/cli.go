package controller

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/urfave/cli/v2"
)

var logger = utils.Logger

func (x *Controller) RunCmd(args []string) error {
	app := &cli.App{
		Name:        "octovy",
		Version:     model.Version,
		Description: "Vulnerability management service for GitHub repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				EnvVars: []string{"OCTOVY_LOG_LEVEL"},
				Value:   "info",
				Usage:   "LogLevel [trace|debug|info|warn|error]",
			},
		},
		Commands: []*cli.Command{
			newServeCommand(x),
		},
		Before: func(c *cli.Context) error {
			if err := globalSetup(c); err != nil {
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error().Interface("config", x.Config).Err(err).Msg("Failed")
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
				Name:        "addr",
				Usage:       "server binding address",
				Aliases:     []string{"a"},
				EnvVars:     []string{"OCTOVY_ADDR"},
				Destination: &ctrl.Config.ServerAddr,
				Value:       "127.0.0.1",
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Port number",
				Aliases:     []string{"p"},
				EnvVars:     []string{"OCTOVY_PORT"},
				Destination: &ctrl.Config.ServerPort,
				Value:       9080,
			},

			&cli.StringFlag{
				Name:        "db-type",
				Usage:       "Database type [postgres|sqlite3]",
				EnvVars:     []string{"OCTOVY_DB_TYPE"},
				Destination: &ctrl.Config.DBType,
				Value:       "sqlite3",
			},
			&cli.StringFlag{
				Name:        "db-config",
				Usage:       "Database config as DSN",
				EnvVars:     []string{"OCTOVY_DB_CONFIG"},
				Destination: &ctrl.Config.DBConfig,
				Value:       "file:ent?mode=memory&cache=shared&_fk=1",
			},

			&cli.StringFlag{
				Name:        "frontend-url",
				EnvVars:     []string{"OCTOVY_FRONTEND_URL"},
				Destination: &ctrl.Config.FrontendURL,
				Required:    true,
			},

			&cli.BoolFlag{
				Name:        "webhook-only",
				Usage:       "Enable only webhook from GitHub. Frontend and API will be disabled",
				Destination: &ctrl.Config.WebhookOnly,
				EnvVars:     []string{"OCTOVY_WEBHOOK_ONLY"},
			},

			&cli.Int64Flag{
				Name:        "github-app-id",
				EnvVars:     []string{"OCTOVY_GITHUB_APP_ID"},
				Destination: &ctrl.Config.GitHubAppID,
				Required:    true,
			},
			&cli.PathFlag{
				Name:        "github-app-private-key",
				EnvVars:     []string{"OCTOVY_GITHUB_APP_PRIVATE_KEY"},
				Usage:       "GitHub App private key data (not file path)",
				Destination: &ctrl.Config.GitHubAppPrivateKey,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-app-client-id",
				EnvVars:     []string{"OCTOVY_GITHUB_CLIENT_ID"},
				Destination: &ctrl.Config.GitHubAppClientID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-app-client-secret",
				EnvVars:     []string{"OCTOVY_GITHUB_SECRET"},
				Destination: &ctrl.Config.GitHubAppSecret,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-webhook-secret",
				EnvVars:     []string{"OCTOVY_GITHUB_WEBHOOK_SECRET"},
				Destination: &ctrl.Config.GitHubWebhookSecret,
				Usage:       "Verify webhook request with the secret",
			},

			&cli.StringFlag{
				Name:        "trivy-db-path",
				EnvVars:     []string{"OCTOVY_TRIVY_DB_PATH"},
				Destination: &ctrl.Config.TrivyDBPath,
				Value:       "/tmp/trivy.db",
				Required:    true,
			},

			&cli.StringFlag{
				Name:        "sentry-dsn",
				EnvVars:     []string{"SENTRY_DSN"},
				Destination: &ctrl.Config.SentryDSN,
			},
			&cli.StringFlag{
				Name:        "sentry-env",
				EnvVars:     []string{"SENTRY_ENV"},
				Destination: &ctrl.Config.SentryEnv,
			},
		},
		Action: func(c *cli.Context) error {
			if err := ctrl.usecase.Init(); err != nil {
				return err
			}

			return serveCommand(c, ctrl)
		},
		After: func(c *cli.Context) error {
			ctrl.usecase.Shutdown()
			return nil
		},
	}
}

func serveCommand(c *cli.Context, ctrl *Controller) error {
	serverAddr := fmt.Sprintf("%s:%d", ctrl.Config.ServerAddr, ctrl.Config.ServerPort)

	engine := server.New(ctrl.usecase)

	gin.SetMode(gin.DebugMode)
	copiedConfig := *ctrl.Config
	// Removing sensitive data
	if copiedConfig.GitHubAppPrivateKey != "" {
		copiedConfig.GitHubAppPrivateKey = "[Removed]"
	}
	if copiedConfig.GitHubAppSecret != "" {
		copiedConfig.GitHubAppSecret = "[Removed]"
	}
	if copiedConfig.DBConfig != "" {
		copiedConfig.DBConfig = "[Removed]"
	}
	if copiedConfig.GitHubWebhookSecret != "" {
		copiedConfig.GitHubWebhookSecret = "[Removed]"
	}

	logger.Info().Interface("config", copiedConfig).Msg("Starting server...")
	if err := engine.Run(serverAddr); err != nil {
		logger.Error().Err(err).Interface("config", ctrl.Config).Msg("Server error")
	}

	return nil
}
