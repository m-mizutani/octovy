package service

import (
	"encoding/base64"
	"strconv"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type secretValues struct {
	GitHubAppPrivateKey string `json:"github_app_private_key"`
	GitHubClientSecret  string `json:"github_client_secret"`

	// AppID and ClientID does not required to be hide, but save them with private key and client secret for life cycle management of the variables.
	GitHubAppID    string `json:"github_app_id"`
	GitHubClientID string `json:"github_client_id"`
}

func (x *secretValues) GithubAppPEM() ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(x.GitHubAppPrivateKey)
	if err != nil {
		return nil, goerr.Wrap(model.ErrInvalidSecretValues, err.Error()).With("key.length", len(x.GitHubAppPrivateKey))
	}
	return raw, nil
}

func (x *secretValues) GetGitHubAppID() (int64, error) {
	n, err := strconv.ParseInt(x.GitHubAppID, 10, 64)
	if err != nil {
		return 0, goerr.Wrap(model.ErrInvalidSecretValues, err.Error()).With("github_app_id", x.GitHubAppID)
	}
	return n, nil
}

func (x *Service) GetSecrets() (*secretValues, error) {
	var values secretValues
	if err := golambda.GetSecretValuesWithFactory(x.config.SecretsARN, &values, func(region string) (golambda.SecretsManagerClient, error) {
		return x.Infra.NewSecretManager(region)
	}); err != nil {
		return nil, err
	}

	return &values, nil
}
