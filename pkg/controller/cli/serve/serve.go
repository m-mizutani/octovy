package serve

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/m-mizutani/octovy/pkg/utils"

	"github.com/urfave/cli/v2"
)

func New() *cli.Command {
	var (
		addr                string
		trivyPath           string
		gitHubAppID         types.GitHubAppID
		gitHubAppSecret     types.GitHubAppSecret
		gitHubAppPrivateKey types.GitHubAppPrivateKey
	)
	return &cli.Command{
		Name:  "serve",
		Usage: "Server mode",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				Usage:       "Binding address",
				Value:       "127.0.0.1:8000",
				Destination: &addr,
			},
			&cli.StringFlag{
				Name:        "trivy-path",
				Usage:       "Path to trivy binary",
				Value:       "trivy",
				Destination: &trivyPath,
			},
			&cli.Int64Flag{
				Name:        "github-app-id",
				Usage:       "GitHub App ID",
				Destination: (*int64)(&gitHubAppID),
				EnvVars:     []string{"OCTOVY_GITHUB_APP_ID"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-app-secret",
				Usage:       "GitHub App Secret",
				Destination: (*string)(&gitHubAppSecret),
				EnvVars:     []string{"OCTOVY_GITHUB_APP_SECRET"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "github-app-private-key",
				Usage:       "GitHub App Private Key",
				Destination: (*string)(&gitHubAppPrivateKey),
				EnvVars:     []string{"OCTOVY_GITHUB_APP_PRIVATE_KEY"},
				Required:    true,
			},
		},
		Action: func(ctx *cli.Context) error {
			ghApp, err := githubapp.New(gitHubAppID, gitHubAppPrivateKey)
			if err != nil {
				return err
			}

			clients := infra.New(
				infra.WithGitHubApp(ghApp),
				infra.WithTrivy(trivy.New(trivyPath)),
			)
			uc := usecase.New(clients)
			s := server.New(uc, gitHubAppSecret)

			serverErr := make(chan error, 1)
			httpServer := &http.Server{
				Addr:    addr,
				Handler: s.Mux(),

				ReadHeaderTimeout: 10 * time.Second,
			}

			go func() {
				utils.Logger().Info("starting server", "addr", addr)
				if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
					serverErr <- goerr.Wrap(err, "failed to listen and serve")
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case err := <-serverErr:
				return err

			case sig := <-quit:
				utils.Logger().Info("shutting down server", "signal", sig)

				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := httpServer.Shutdown(ctx); err != nil {
					return goerr.Wrap(err, "failed to shutdown server")
				}
			}

			return nil
		},
	}
}
