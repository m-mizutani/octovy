package service_test

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/Netflix/go-env"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/stretchr/testify/require"
)

type WriteBuffer struct {
	bytes.Buffer
}

func (x *WriteBuffer) Close() error { return nil }

func TestGitHubDownload(t *testing.T) {
	var props struct {
		GITHUB_APP_ID          string `env:"GITHUB_APP_ID"`
		GITHUB_APP_PRIVATE_KEY string `env:"GITHUB_APP_PRIVATE_KEY"`
		GITHUB_COMMIT          string `env:"GITHUB_COMMIT"`
		GITHUB_ENDPOINT        string `env:"GITHUB_ENDPOINT"`
		GITHUB_INSTALL_ID      string `env:"GITHUB_INSTALL_ID"`
		GITHUB_ORG             string `env:"GITHUB_ORG"`
		GITHUB_REPO_NAME       string `env:"GITHUB_REPO_NAME"`
	}
	_, err := env.UnmarshalFromEnviron(&props)
	require.NoError(t, err)

	if props.GITHUB_APP_ID == "" ||
		props.GITHUB_APP_PRIVATE_KEY == "" ||
		props.GITHUB_COMMIT == "" ||
		props.GITHUB_ENDPOINT == "" ||
		props.GITHUB_INSTALL_ID == "" ||
		props.GITHUB_ORG == "" ||
		props.GITHUB_REPO_NAME == "" {
		t.Logf("props => %+v\n", props)
		t.Skip("Not enough paramters")
	}

	installID, err := strconv.ParseInt(props.GITHUB_INSTALL_ID, 10, 64)
	require.NoError(t, err)

	privateKey, err := ioutil.ReadFile(props.GITHUB_APP_PRIVATE_KEY)
	require.NoError(t, err)

	newSM, mockSM := aws.NewMockSecretsManagerSet()
	mockSM.OutData["arn:aws:secretsmanager:us-east-0:123456789012:secret:tutorials/MyFirstSecret-jiObOV"] = map[string]string{
		"github_app_private_key": base64.StdEncoding.EncodeToString(privateKey),
		"github_app_id":          props.GITHUB_APP_ID,
	}

	cfg := model.NewConfig()
	cfg.SecretsARN = "arn:aws:secretsmanager:us-east-0:123456789012:secret:tutorials/MyFirstSecret-jiObOV"
	cfg.GitHubEndpoint = props.GITHUB_ENDPOINT
	svc := service.New(cfg)
	svc.Infra.NewSecretManager = newSM

	buf := &WriteBuffer{}
	repo := &model.GitHubRepo{
		Owner:    props.GITHUB_ORG,
		RepoName: props.GITHUB_REPO_NAME,
	}
	require.NoError(t, svc.GetCodeZip(repo, props.GITHUB_COMMIT, installID, buf))
}
