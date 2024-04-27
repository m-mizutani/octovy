package usecase_test

import (
	"context"
	_ "embed"
	"os"
	"strconv"

	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/bq"
	"github.com/m-mizutani/octovy/pkg/infra/gh"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/m-mizutani/octovy/pkg/utils"
)

//go:embed testdata/octovy-test-code-main.zip
var testCodeZip []byte

//go:embed testdata/trivy-result.json
var testTrivyResult []byte

func TestScanGitHubRepo(t *testing.T) {
	mockGH := &ghMock{}
	mockHTTP := &httpMock{}
	mockTrivy := &trivyMock{}
	mockBQ := &bq.Mock{}

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(mockGH),
		infra.WithHTTPClient(mockHTTP),
		infra.WithTrivy(mockTrivy),
		infra.WithBigQuery(mockBQ),
	))

	ctx := context.Background()

	mockGH.mockGetArchiveURL = func(ctx context.Context, input *gh.GetArchiveURLInput) (*url.URL, error) {
		gt.V(t, input.Owner).Equal("m-mizutani")
		gt.V(t, input.Repo).Equal("octovy")
		gt.V(t, input.CommitID).Equal("f7c8851da7c7fcc46212fccfb6c9c4bda520f1ca")
		gt.V(t, input.InstallID).Equal(12345)

		resp := gt.R1(url.Parse("https://example.com/some/url.zip")).NoError(t)
		return resp, nil
	}

	mockHTTP.mockDo = func(req *http.Request) (*http.Response, error) {
		gt.V(t, req.URL.String()).Equal("https://example.com/some/url.zip")

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(testCodeZip)),
		}
		return resp, nil
	}

	mockTrivy.mockRun = func(ctx context.Context, args []string) error {
		gt.A(t, args).
			Contain([]string{"--format", "json"}).
			Contain([]string{"--list-all-pkgs"})

		for i := range args {
			if args[i] == "--output" {
				fd := gt.R1(os.Create(args[i+1])).NoError(t)
				gt.R1(fd.Write(testTrivyResult)).NoError(t)
				gt.NoError(t, fd.Close())
				return nil
			}
		}

		t.Error("no --output option")
		return nil
	}

	var calledBQCreateTable int
	mockBQ.FnCreateTable = func(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error {
		calledBQCreateTable++
		gt.Equal(t, table, "scans")
		return nil
	}

	mockBQ.FnGetMetadata = func(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error) {
		return nil, nil
	}

	var calledBQInsert int
	mockBQ.FnInsert = func(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error {
		calledBQInsert++
		return nil
	}

	gt.NoError(t, uc.ScanGitHubRepo(ctx, &model.ScanGitHubRepoInput{
		GitHubMetadata: model.GitHubMetadata{
			GitHubCommit: model.GitHubCommit{
				GitHubRepo: model.GitHubRepo{
					RepoID:   12345,
					Owner:    "m-mizutani",
					RepoName: "octovy",
				},
				CommitID: "f7c8851da7c7fcc46212fccfb6c9c4bda520f1ca",
			},
		},
		InstallID: 12345,
	}))
	gt.Equal(t, calledBQCreateTable, 1)
	gt.Equal(t, calledBQInsert, 1)
}

type ghMock struct {
	mockGetArchiveURL func(ctx context.Context, input *gh.GetArchiveURLInput) (*url.URL, error)
}

func (x *ghMock) GetArchiveURL(ctx context.Context, input *gh.GetArchiveURLInput) (*url.URL, error) {
	return x.mockGetArchiveURL(ctx, input)
}

type trivyMock struct {
	mockRun func(ctx context.Context, args []string) error
}

func (x *trivyMock) Run(ctx context.Context, args []string) error {
	return x.mockRun(ctx, args)
}

type httpMock struct {
	mockDo func(req *http.Request) (*http.Response, error)
}

func (x *httpMock) Do(req *http.Request) (*http.Response, error) {
	return x.mockDo(req)
}

func TestScanGitHubRepoWithData(t *testing.T) {

	if _, ok := os.LookupEnv("TEST_SCAN_GITHUB_REPO"); !ok {
		t.Skip("TEST_SCAN_GITHUB_REPO is not set")
	}

	// Setting up GitHub App
	strAppID := utils.LoadEnv(t, "TEST_OCTOVY_GITHUB_APP_ID")
	privateKey := utils.LoadEnv(t, "TEST_OCTOVY_GITHUB_APP_PRIVATE_KEY")

	appID := gt.R1(strconv.ParseInt(strAppID, 10, 64)).NoError(t)
	ghApp := gt.R1(gh.New(types.GitHubAppID(appID), types.GitHubAppPrivateKey(privateKey))).NoError(t)

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(ghApp),
	))

	ctx := context.Background()

	gt.NoError(t, uc.ScanGitHubRepo(ctx, &model.ScanGitHubRepoInput{
		GitHubMetadata: model.GitHubMetadata{
			GitHubCommit: model.GitHubCommit{
				GitHubRepo: model.GitHubRepo{
					RepoID:   41633205,
					Owner:    "m-mizutani",
					RepoName: "octovy",
				},
				CommitID: "6581604ef668e77a178e18dbc56e898f5fd87014",
			},
		},
		InstallID: 41633205,
	}))
}
