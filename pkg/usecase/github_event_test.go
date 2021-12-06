package usecase_test

import (
	"io/ioutil"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/require"
)

func TestVerityGitHubSecret(t *testing.T) {
	webhookBodySample, err := ioutil.ReadFile("./testdata/webhook-body-sample.json")
	secret := "my-github-secret"
	require.NoError(t, err)

	t.Run("pass when verifying valid signature", func(t *testing.T) {
		uc := usecase.New(&model.Config{
			GitHubWebhookSecret: secret,
		}, nil)
		require.NoError(t, uc.VerifyGitHubSecret("sha256=7b91b4881b9ad0ea39f0e335786e6242b97bb2d9038b25d358dae70a14424535", webhookBodySample))
	})

	t.Run("error when verifying invalid signature", func(t *testing.T) {
		uc := usecase.New(&model.Config{
			GitHubWebhookSecret: secret,
		}, nil)
		require.ErrorIs(t, model.ErrInvalidWebhookData, uc.VerifyGitHubSecret("sha256=7b91b4881b9ad0ea39f0e335786e6242b97bb2d9038b25d358dae70000000000", webhookBodySample))
	})

	t.Run("ignore when no secret even if invalid signature", func(t *testing.T) {
		uc := usecase.New(&model.Config{
			GitHubWebhookSecret: "",
		}, nil)
		require.NoError(t, uc.VerifyGitHubSecret("sha256=7b91b4881b9ad0ea39f0e335786e6242b97bb2d9038b25d358dae70000000000", webhookBodySample))
	})

	t.Run("error when no signature", func(t *testing.T) {
		uc := usecase.New(&model.Config{
			GitHubWebhookSecret: secret,
		}, nil)
		require.ErrorIs(t, model.ErrInvalidWebhookData, uc.VerifyGitHubSecret("", webhookBodySample))
	})

}
