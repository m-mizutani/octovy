package githubapp_test

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/Netflix/go-env"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
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

	appID, err := strconv.ParseInt(props.GITHUB_APP_ID, 10, 64)
	require.NoError(t, err)

	app := githubapp.New(appID, installID, privateKey, props.GITHUB_ENDPOINT)

	buf := &WriteBuffer{}
	repo := &model.GitHubRepo{
		Owner:    props.GITHUB_ORG,
		RepoName: props.GITHUB_REPO_NAME,
	}
	require.NoError(t, app.GetCodeZip(repo, props.GITHUB_COMMIT, buf))
}
