package usecase_test

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strconv"

	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
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
	testDB := newTestDB(t)

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(mockGH),
		infra.WithHTTPClient(mockHTTP),
		infra.WithTrivy(mockTrivy),
		infra.WithDB(testDB),
	))

	ctx := model.NewContext()

	mockGH.mockGetArchiveURL = func(ctx *model.Context, input *gh.GetArchiveURLInput) (*url.URL, error) {
		gt.V(t, input.Owner).Equal("m-mizutani")
		gt.V(t, input.Repo).Equal("octovy")
		gt.V(t, input.CommitID).Equal("1234567890")
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

	mockTrivy.mockRun = func(ctx *model.Context, args []string) error {
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

	gt.NoError(t, uc.ScanGitHubRepo(ctx, &usecase.ScanGitHubRepoInput{
		GitHubRepoMetadata: usecase.GitHubRepoMetadata{
			GitHubCommit: usecase.GitHubCommit{
				GitHubRepo: usecase.GitHubRepo{
					Owner: "m-mizutani",
					Repo:  "octovy",
				},
				CommitID: "1234567890",
			},
		},
		InstallID: 12345,
	}))
}

type ghMock struct {
	mockGetArchiveURL func(ctx *model.Context, input *gh.GetArchiveURLInput) (*url.URL, error)
}

func (x *ghMock) GetArchiveURL(ctx *model.Context, input *gh.GetArchiveURLInput) (*url.URL, error) {
	return x.mockGetArchiveURL(ctx, input)
}

type trivyMock struct {
	mockRun func(ctx *model.Context, args []string) error
}

func (x *trivyMock) Run(ctx *model.Context, args []string) error {
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
	strAppID, ok := os.LookupEnv("OCTOVY_GITHUB_APP_ID")
	if !ok {
		t.Error("OCTOVY_GITHUB_APP_ID is not set")
	}
	privateKey, ok := os.LookupEnv("OCTOVY_GITHUB_APP_PRIVATE_KEY")
	if !ok {
		t.Error("OCTOVY_GITHUB_APP_PRIVATE_KEY is not set")
	}
	appID := gt.R1(strconv.ParseInt(strAppID, 10, 64)).NoError(t)
	ghApp := gt.R1(gh.New(types.GitHubAppID(appID), types.GitHubAppPrivateKey(privateKey))).NoError(t)

	// Setting up database
	dbUser, ok := os.LookupEnv("OCTOVY_DB_USER")
	if !ok {
		t.Error("OCTOVY_DB_USER is not set")
	}
	dbPass, ok := os.LookupEnv("OCTOVY_DB_PASSWORD")
	if !ok {
		t.Error("OCTOVY_DB_PASS is not set")
	}
	dbName, ok := os.LookupEnv("OCTOVY_DB_NAME")
	if !ok {
		t.Error("OCTOVY_DB_NAME is not set")
	}
	dbPort, ok := os.LookupEnv("OCTOVY_DB_PORT")
	if !ok {
		t.Error("OCTOVY_DB_PORT is not set")
	}
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbPort)

	dbClient := gt.R1(sql.Open("postgres", dsn)).NoError(t)
	defer utils.SafeClose(dbClient)

	if t.Failed() {
		t.FailNow()
	}

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(ghApp),
		infra.WithDB(dbClient),
	))

	ctx := model.NewContext()

	gt.NoError(t, uc.ScanGitHubRepo(ctx, &usecase.ScanGitHubRepoInput{
		GitHubRepoMetadata: usecase.GitHubRepoMetadata{
			GitHubCommit: usecase.GitHubCommit{
				GitHubRepo: usecase.GitHubRepo{
					Owner: "m-mizutani",
					Repo:  "octovy",
				},
				CommitID: "6581604ef668e77a178e18dbc56e898f5fd87014",
			},
		},
		InstallID: 41633205,
	}))
}
