package serve

import (
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/service"

	"github.com/urfave/cli/v2"
)

func New() *cli.Command {
	var (
		addr                string
		port                int
		gitHubAppID         types.GitHubAppID
		gitHubAppSecret     types.GitHubAppSecret
		gitHubAppPrivateKey types.GitHubAppPrivateKey
	)
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Server mode",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				Usage:       "Binding address",
				Value:       "127.0.0.1",
				Destination: &addr,
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Binding port",
				Value:       5080,
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "github-app-id",
				Usage:       "GitHub App ID",
				Destination: (*string)(&gitHubAppID),
				EnvVars:     []string{"OCTOVY_GITHUB_APP_ID"},
			},
			&cli.StringFlag{
				Name:        "github-app-secret",
				Usage:       "GitHub App Secret",
				Destination: (*string)(&gitHubAppSecret),
				EnvVars:     []string{"OCTOVY_GITHUB_APP_SECRET"},
			},
			&cli.StringFlag{
				Name:        "github-app-private-key",
				Usage:       "GitHub App Private Key",
				Destination: (*string)(&gitHubAppPrivateKey),
				EnvVars:     []string{"OCTOVY_GITHUB_APP_PRIVATE_KEY"},
			},
		},
		Action: func(ctx *cli.Context) error {
			clients := infra.New()
			svc := service.New(clients)
			httpServer := server.New(svc)

			if err := httpServer.Listen(addr, port); err != nil {
				return err
			}
			return nil
		},
	}
}
