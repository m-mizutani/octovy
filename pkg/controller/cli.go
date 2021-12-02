package controller

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/usecase"
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
			&cli.StringFlag{
				Name:    "log-format",
				EnvVars: []string{"OCTOVY_LOG_FORMAT"},
				Value:   "console",
				Usage:   "LogFormat [console|json]",
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
		logger.With("config", x.Config).With("err", err).Error("Failed")
		return err
	}

	return nil
}

func globalSetup(c *cli.Context) error {
	// Setup logger
	if err := utils.SetLogLevel(c.String("log-level")); err != nil {
		return goerr.Wrap(err)
	}
	if err := utils.SetLogFormat(c.String("log-format")); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

func newServeCommand(ctrl *Controller) *cli.Command {
	var checkRuleFile string
	var infraCfg infra.Config

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
				Destination: &infraCfg.DBType,
				Value:       "sqlite3",
			},
			&cli.StringFlag{
				Name:        "db-config",
				Usage:       "Database config as DSN",
				EnvVars:     []string{"OCTOVY_DB_CONFIG"},
				Destination: &infraCfg.DBConfig,
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
				Destination: &infraCfg.GitHubAppID,
				Required:    true,
			},
			&cli.PathFlag{
				Name:        "github-app-private-key",
				EnvVars:     []string{"OCTOVY_GITHUB_APP_PRIVATE_KEY"},
				Usage:       "GitHub App private key data (not file path)",
				Destination: &infraCfg.GitHubAppPrivateKey,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-app-client-id",
				EnvVars:     []string{"OCTOVY_GITHUB_CLIENT_ID"},
				Destination: &infraCfg.GitHubAppClientID,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-app-client-secret",
				EnvVars:     []string{"OCTOVY_GITHUB_SECRET"},
				Destination: &infraCfg.GitHubAppSecret,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-webhook-secret",
				EnvVars:     []string{"OCTOVY_GITHUB_WEBHOOK_SECRET"},
				Destination: &ctrl.Config.GitHubWebhookSecret,
				Usage:       "Verify webhook request with the secret",
			},

			&cli.StringFlag{
				Name:        "check-policy-data",
				EnvVars:     []string{"OCTOVY_CHECK_POLICY_DATA"},
				Destination: &infraCfg.CheckPolicyData,
				Usage:       "Check result policy in Rego (plain text)",
			},
			&cli.StringFlag{
				Name:        "check-policy-file",
				EnvVars:     []string{"OCTOVY_CHECK_POLICY_FILE"},
				Destination: &checkRuleFile,
				Usage:       "Check result policy in Rego (file path)",
			},

			&cli.StringFlag{
				Name:        "opa-url",
				EnvVars:     []string{"OCTOVY_OPA_URL"},
				Destination: &infraCfg.OPAServerURL,
			},
			&cli.BoolFlag{
				Name:        "opa-use-iap",
				EnvVars:     []string{"OCTOVY_OPA_IAP"},
				Destination: &infraCfg.OPAUseGoogleIAP,
			},

			&cli.StringFlag{
				Name:        "trivy-path",
				EnvVars:     []string{"OCTOVY_TRIVY_PATH"},
				Destination: &infraCfg.TrivyPath,
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
			if checkRuleFile != "" {
				raw, err := ioutil.ReadFile(checkRuleFile)
				if err != nil {
					return goerr.Wrap(err, "fail to read check rule file")
				}
				if infraCfg.CheckPolicyData != "" {
					logger.With("existed", infraCfg.CheckPolicyData).Warn("both of --check-rule-file and --check-rule-data are specified. check-rule-data will be overwritten")
				}
				infraCfg.CheckPolicyData = string(raw)
			}

			clients, err := infra.New(&infraCfg)
			if err != nil {
				return err
			}
			uc := usecase.New(ctrl.Config, clients)
			defer uc.Close()

			serverAddr := fmt.Sprintf("%s:%d", ctrl.Config.ServerAddr, ctrl.Config.ServerPort)

			engine := server.New(uc)

			gin.SetMode(gin.DebugMode)

			logger.With("config", ctrl.Config).Info("Starting server...")
			if err := engine.Run(serverAddr); err != nil {
				return err
			}
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
	}
}
