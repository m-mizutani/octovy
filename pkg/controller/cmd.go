package controller

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/api"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var logger zerolog.Logger

func (x *Controller) RunCmd(args []string, envVars []string) error {
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
	logLevel := c.String("log-level")
	var zeroLogLevel zerolog.Level

	switch strings.ToLower(logLevel) {
	case "trace":
		zeroLogLevel = zerolog.TraceLevel
	case "debug":
		zeroLogLevel = zerolog.DebugLevel
	case "info":
		zeroLogLevel = zerolog.InfoLevel
	case "warn":
		zeroLogLevel = zerolog.WarnLevel
	case "error":
		zeroLogLevel = zerolog.ErrorLevel
	default:
		zeroLogLevel = zerolog.InfoLevel
	}

	console := zerolog.ConsoleWriter{Out: os.Stdout}

	logger = zerolog.New(console).Level(zeroLogLevel).With().Timestamp().Logger()
	logger.Debug().Str("LogLevel", logLevel).Msg("Set log level")
	return nil
}

type serveCommandConfig struct {
	AWSRegion string
	TableName string
	AssetDir  string
	Addr      string
	Port      int

	FrontendURL    string
	GitHubWebURL   string
	GitHubEndpoint string
	SecretsARN     string

	DebugMode bool

	ctrl *Controller
}

func newServeCommand(ctrl *Controller) *cli.Command {
	config := &serveCommandConfig{
		ctrl: ctrl,
	}

	return &cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "aws-region",
				Aliases:     []string{"r"},
				EnvVars:     []string{"AWS_REGION"},
				Destination: &config.AWSRegion,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "table-name",
				Aliases:     []string{"t"},
				EnvVars:     []string{"OCTOVY_TABLE_NAME"},
				Destination: &config.TableName,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "Addr",
				Usage:       "server binding address",
				Aliases:     []string{"a"},
				EnvVars:     []string{"OCTOVY_ADDR"},
				Destination: &config.Addr,
				Value:       "127.0.0.1",
			},
			&cli.IntFlag{
				Name:        "Port",
				Usage:       "Port number",
				Aliases:     []string{"p"},
				EnvVars:     []string{"OCTOVY_PORT"},
				Destination: &config.Port,
				Value:       9080,
			},

			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug mode",
				Aliases:     []string{"d"},
				Destination: &config.DebugMode,
			},

			&cli.StringFlag{
				Name:        "secrets-arn",
				EnvVars:     []string{"OCTOVY_SECRETS_ARN"},
				Aliases:     []string{"s"},
				Destination: &config.SecretsARN,
			},

			&cli.StringFlag{
				Name:        "github-web-url",
				EnvVars:     []string{"OCTOVY_GITHUB_WEB_URL"},
				Destination: &config.GitHubWebURL,
				Value:       "https://github.com",
			},
			&cli.StringFlag{
				Name:        "github-endpoint",
				EnvVars:     []string{"OCTOVY_GITHUB_ENDPOINT"},
				Destination: &config.GitHubEndpoint,
				Value:       "https://api.github.com",
			},
			&cli.StringFlag{
				Name:        "frontend-url",
				EnvVars:     []string{"OCTOVY_FRONTEND_URL"},
				Destination: &config.FrontendURL,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			return apiCommand(c, config)
		},
	}
}

func apiCommand(c *cli.Context, config *serveCommandConfig) error {
	serverAddr := fmt.Sprintf("%s:%d", config.Addr, config.Port)

	config.ctrl.Config.AwsRegion = config.AWSRegion
	config.ctrl.Config.TableName = config.TableName
	config.ctrl.Config.SecretsARN = config.SecretsARN

	config.ctrl.Config.GitHubWebURL = config.GitHubWebURL
	config.ctrl.Config.GitHubEndpoint = config.GitHubEndpoint
	config.ctrl.Config.FrontendURL = config.FrontendURL

	errCh := make(chan error, 128)

	// Server thread
	go func() {
		engine := api.New(&api.Config{
			Usecase:  config.ctrl.Usecase,
			AssetDir: config.AssetDir,
		})

		gin.SetMode(gin.ReleaseMode)
		if config.DebugMode {
			gin.SetMode(gin.DebugMode)
		}

		logger.Info().Interface("config", config).Msg("Starting server...")
		if err := engine.Run(serverAddr); err != nil {
			errCh <- err
		}
	}()

	for err := range errCh {
		logger.Error().Err(err).Interface("config", config).Msg("Server error")
		if errors.Is(err, model.ErrFatal) {
			logger.Error().Err(err).Msg("shuting down by fatal error")
			return err
		}
	}

	return nil
}
