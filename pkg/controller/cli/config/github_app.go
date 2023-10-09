package config

import (
	"encoding/base64"

	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/urfave/cli/v2"
)

type GitHubApp struct {
	ID         types.GitHubAppID
	Secret     types.GitHubAppSecret
	privateKey types.GitHubAppPrivateKey
}

func (x *GitHubApp) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.Int64Flag{
			Name:        "github-app-id",
			Usage:       "GitHub App ID",
			Category:    "GitHub App",
			Destination: (*int64)(&x.ID),
			EnvVars:     []string{"OCTOVY_GITHUB_APP_ID"},
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "github-app-secret",
			Usage:       "GitHub App Secret",
			Category:    "GitHub App",
			Destination: (*string)(&x.Secret),
			EnvVars:     []string{"OCTOVY_GITHUB_APP_SECRET"},
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "github-app-private-key",
			Usage:       "GitHub App Private Key",
			Category:    "GitHub App",
			Destination: (*string)(&x.privateKey),
			EnvVars:     []string{"OCTOVY_GITHUB_APP_PRIVATE_KEY"},
			Required:    true,
		},
	}
}

func (x *GitHubApp) PrivateKey() types.GitHubAppPrivateKey {
	if raw, err := base64.StdEncoding.DecodeString(string(x.privateKey)); err == nil {
		return types.GitHubAppPrivateKey(raw)
	}
	return x.privateKey
}
