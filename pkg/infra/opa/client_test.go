package opa_test

import (
	"context"
	"os"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/opa"
	"github.com/stretchr/testify/require"
)

func TestOPAClient(t *testing.T) {
	url, ok := os.LookupEnv("OPA_SERVER_URL")
	if !ok {
		t.Skip("OPA_SERVER_URL is not set")
	}

	obj := map[string]string{
		"user": "blue",
	}
	var resp model.GitHubCheckResult

	client, err := opa.New(&opa.Config{
		BaseURL: url,
		Path:    "octovy/testing",
	})

	require.NoError(t, err)
	require.NoError(t, client.Data(context.Background(), obj, resp))
}
