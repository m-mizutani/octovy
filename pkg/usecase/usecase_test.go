package usecase_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/pkg/usecase"

	gh "github.com/google/go-github/v39/github"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSet struct {
	DB        *db.Client
	GitHub    *github.Mock
	GtiHubApp *githubapp.Mock
	TrivyDB   *trivydb.Mock
	Utils     *infra.Utils
}

func setupUsecase(t *testing.T) (usecase.Interface, *mockSet) {
	uc := usecase.NewUsecase(&model.Config{})

	dbClient := db.NewMock(t)
	ghClient := github.NewMock()
	newGitHubApp, ghApp := githubapp.NewMock()
	newTrivyDB, trivyDB := trivydb.NewMock()
	util := infra.NewUtils()

	uc.InjectInfra(&infra.Interfaces{
		DB:           dbClient,
		GitHub:       ghClient,
		NewGitHubApp: newGitHubApp,
		NewTrivyDB:   newTrivyDB,
		Utils:        util,
	})

	usecase.SetErrorHandler(uc, func(err error) {
		require.NoError(t, err)
	})

	return uc, &mockSet{
		DB:        dbClient,
		GitHub:    ghClient,
		GtiHubApp: ghApp,
		TrivyDB:   trivyDB,
		Utils:     util,
	}
}

func genRSAKey(t *testing.T) []byte {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func injectGitHubMock(t *testing.T, mock *mockSet) {
	var calledListReleasesMock,
		calledDownloadReleaseAssetMock,
		calledGetCodeZipMock,
		calledCreateCheckRunMock,
		calledUpdateCheckRunMock int

	mock.GitHub.ListReleasesMock = func(ctx context.Context, owner, repo string) ([]*gh.RepositoryRelease, error) {
		calledListReleasesMock++

		assert.Equal(t, "aquasecurity", owner)
		assert.Equal(t, "trivy-db", repo)

		return []*gh.RepositoryRelease{
			{
				Name: gh.String("v1-20000000"),
				Assets: []*gh.ReleaseAsset{
					{
						Name: gh.String("xxx.db.gz"),
						ID:   gh.Int64(2345),
					},
					{
						Name: gh.String("trivy.db.gz"),
						ID:   gh.Int64(3456),
					},
				},
			},
		}, nil
	}

	mock.GitHub.DownloadReleaseAssetMock = func(ctx context.Context, owner, repo string, assetID int64) (io.ReadCloser, error) {
		calledDownloadReleaseAssetMock++
		assert.Equal(t, "aquasecurity", owner)
		assert.Equal(t, "trivy-db", repo)
		assert.Equal(t, int64(3456), assetID)

		buf := &bytes.Buffer{}
		gz := gzip.NewWriter(buf)
		gz.Write([]byte("boom!"))
		require.NoError(t, gz.Close())
		return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
	}

	mock.GtiHubApp.GetCodeZipMock = func(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error {
		calledGetCodeZipMock++
		raw, err := os.ReadFile("./testdata/sample-repo.zip")
		require.NoError(t, err)
		w.Write(raw)
		w.Close()
		return nil
	}

	mock.Utils.ReadFile = func(fname string) ([]byte, error) {
		return genRSAKey(t), nil
	}

	const dummyCheckID int64 = 999
	mock.GtiHubApp.CreateCheckRunMock = func(repo *model.GitHubRepo, commit string) (int64, error) {
		calledCreateCheckRunMock++
		return dummyCheckID, nil
	}

	mock.GtiHubApp.UpdateCheckRunMock = func(repo *model.GitHubRepo, checkID int64, opt *gh.UpdateCheckRunOptions) error {
		calledUpdateCheckRunMock++
		assert.Equal(t, dummyCheckID, checkID)
		return nil
	}

	t.Cleanup(func() {
		assert.Equal(t, 1, calledListReleasesMock)
		assert.Equal(t, 1, calledDownloadReleaseAssetMock)
		assert.Equal(t, 1, calledGetCodeZipMock)
		assert.Equal(t, 1, calledCreateCheckRunMock)
		assert.Equal(t, 1, calledUpdateCheckRunMock)
	})
}
