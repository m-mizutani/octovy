package usecase_test

import (
	"io"
	"os"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/opa"
	"github.com/m-mizutani/octovy/pkg/infra/policy"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
	"github.com/m-mizutani/octovy/pkg/usecase"

	gh "github.com/google/go-github/v39/github"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSet struct {
	DB        *db.Client
	GitHub    *github.Mock
	GtiHubApp *githubapp.Mock
	Trivy     *trivy.Mock
	Utils     *infra.Utils
	OPA       *opa.Mock
}

type testOption func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet)

func optDBMock() testOption {
	return func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet) {
		dbClient := db.NewMock(t)
		infra.DB = dbClient
		mock.DB = dbClient
	}
}

func optGitHubMock() testOption {
	return func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet) {
		ghClient := github.NewMock()
		infra.GitHub = ghClient
		mock.GitHub = ghClient
	}
}

func optGitHubAppMock() testOption {
	return func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet) {
		newGitHubApp, ghApp := githubapp.NewMock()
		infra.NewGitHubApp = newGitHubApp
		mock.GtiHubApp = ghApp
	}
}

func optGitHubAppMockZip() testOption {
	return func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet) {
		if mock.GtiHubApp == nil {
			require.Fail(t, "optGitHubAppMock should be called at first")
		}

		var calledGetCodeZipMock int

		mock.GtiHubApp.GetCodeZipMock = func(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error {
			calledGetCodeZipMock++
			raw, err := os.ReadFile("./testdata/sample-repo.zip")
			require.NoError(t, err)
			w.Write(raw)
			w.Close()
			return nil
		}

		t.Cleanup(func() {
			assert.GreaterOrEqual(t, calledGetCodeZipMock, 1)
		})
	}
}

func optOPAServer() testOption {
	return func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet) {
		opaMock := &opa.Mock{}
		mock.OPA = opaMock
		infra.OPAClient = opaMock
	}
}

func optCheckRule(rule string, update func(repo *model.GitHubRepo, checkID int64, opt *gh.UpdateCheckRunOptions) error) testOption {
	return func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet) {
		if mock.GtiHubApp == nil {
			require.Fail(t, "optGitHubAppMock should be called at first")
		}

		check, err := policy.NewCheck(rule)
		require.NoError(t, err)
		infra.CheckPolicy = check

		var calledCreateCheckRunMock int

		const dummyCheckID int64 = 999
		mock.GtiHubApp.CreateCheckRunMock = func(repo *model.GitHubRepo, commit string) (int64, error) {
			calledCreateCheckRunMock++
			return dummyCheckID, nil
		}

		mock.GtiHubApp.UpdateCheckRunMock = update

		t.Cleanup(func() {
			assert.GreaterOrEqual(t, calledCreateCheckRunMock, 1)
		})
	}
}

func optTrivy() testOption {
	return func(t *testing.T, cfg *model.Config, infra *infra.Clients, mock *mockSet) {
		trivyClient := trivy.NewMock()
		infra.Trivy = trivyClient
		mock.Trivy = trivyClient
	}
}

func setupUsecase(t *testing.T, options ...testOption) (*usecase.Usecase, *mockSet) {
	utils := infra.NewUtils()
	var cfg model.Config

	mock := &mockSet{
		Utils: utils,
	}
	inf := &infra.Clients{
		Utils: utils,
	}

	for _, opt := range options {
		opt(t, &cfg, inf, mock)
	}

	uc, err := usecase.New(&cfg, inf)
	require.NoError(t, err)

	return uc, mock
}
