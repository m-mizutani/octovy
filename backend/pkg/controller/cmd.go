package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/backend/pkg/api"
	"github.com/m-mizutani/octovy/backend/pkg/service"
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
			newAPICommand(x),
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

type apiCommandConfig struct {
	AWSRegion string
	TableName string
	SecretARN string
	AssetDir  string
	Addr      string
	Port      int

	ctrl *Controller
}

func newAPICommand(ctrl *Controller) *cli.Command {
	config := &apiCommandConfig{
		ctrl: ctrl,
	}

	return &cli.Command{
		Name: "api",
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

			// Required to handle asset. Not necessary if testing with webpack server
			&cli.StringFlag{
				Name:        "asset-dir",
				Aliases:     []string{"d"},
				EnvVars:     []string{"OCTOVY_ASSET_DIR"},
				Destination: &config.AssetDir,
			},

			// Required to handle webhook. Normally not necessary
			&cli.StringFlag{
				Name:        "secret-arn",
				Aliases:     []string{"s"},
				EnvVars:     []string{"OCTOVY_SECRET_ARN"},
				Destination: &config.SecretARN,
			},
		},
		Action: func(c *cli.Context) error {
			return apiCommand(c, config)
		},
	}
}

func apiCommand(c *cli.Context, config *apiCommandConfig) error {
	svc := service.New(&service.Config{
		AwsRegion:  config.AWSRegion,
		TableName:  config.TableName,
		SecretsARN: config.SecretARN,
	})
	engine := api.New(&api.Config{
		Service:  svc,
		Usecase:  config.ctrl.Usecase,
		AssetDir: config.AssetDir,
	})

	gin.SetMode(gin.DebugMode)
	logger.Info().Interface("config", config).Msg("Starting server...")
	if err := engine.Run(fmt.Sprintf("%s:%d", config.Addr, config.Port)); err != nil {
		logger.Error().Err(err).Interface("config", config).Msg("Server error")
	}

	return nil
}
