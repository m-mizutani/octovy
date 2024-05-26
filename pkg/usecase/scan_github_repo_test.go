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
	"github.com/google/go-github/v53/github"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
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
	mockGH := &interfaces.GitHubMock{}
	mockHTTP := &httpMock{}
	mockTrivy := &trivyMock{}
	mockBQ := &bq.Mock{}
	mockStorage := interfaces.NewStorageMock()

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(mockGH),
		infra.WithHTTPClient(mockHTTP),
		infra.WithTrivy(mockTrivy),
		infra.WithBigQuery(mockBQ),
		infra.WithStorage(mockStorage),
	))

	ctx := context.Background()

	mockGH.MockGetArchiveURL = func(ctx context.Context, input *interfaces.GetArchiveURLInput) (*url.URL, error) {
		gt.V(t, input.Owner).Equal("m-mizutani")
		gt.V(t, input.Repo).Equal("octovy")
		gt.V(t, input.CommitID).Equal("f7c8851da7c7fcc46212fccfb6c9c4bda520f1ca")
		gt.V(t, input.InstallID).Equal(12345)

		resp := gt.R1(url.Parse("https://example.com/some/url.zip")).NoError(t)
		return resp, nil
	}
	mockGH.MockCreateCheckRun = func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error) {
		return 0, nil
	}
	mockGH.MockUpdateCheckRun = func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
		return nil
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
				Branch:   "main",
			},
		},
		InstallID: 12345,
	}))
	gt.Equal(t, calledBQCreateTable, 1)
	gt.Equal(t, calledBQInsert, 1)

	var commitScan *model.Scan
	gt.NoError(t, mockStorage.Unmarshal("m-mizutani/octovy/commit/f7c8851da7c7fcc46212fccfb6c9c4bda520f1ca/scan.json.gz", &commitScan))
	gt.Equal(t, commitScan.GitHub.Owner, "m-mizutani")

	var branchScan *model.Scan
	gt.NoError(t, mockStorage.Unmarshal("m-mizutani/octovy/branch/main/scan.json.gz", &branchScan))
	gt.Equal(t, branchScan.GitHub.Owner, "m-mizutani")
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

func TestScanGitHubRepoWithPR(t *testing.T) {
	mockGH := &interfaces.GitHubMock{}
	mockHTTP := &httpMock{}
	mockTrivy := &trivyMock{}
	mockBQ := &bq.Mock{}
	mockStorage := interfaces.NewStorageMock()

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(mockGH),
		infra.WithHTTPClient(mockHTTP),
		infra.WithTrivy(mockTrivy),
		infra.WithBigQuery(mockBQ),
		infra.WithStorage(mockStorage),
	))

	ctx := context.Background()

	mockHTTP.mockDo = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(testCodeZip)),
		}, nil
	}

	mockTrivy.mockRun = func(ctx context.Context, args []string) error {
		for i := range args {
			if args[i] == "--output" {
				fd := gt.R1(os.Create(args[i+1])).NoError(t)
				gt.R1(fd.Write(testTrivyResult)).NoError(t)
				gt.NoError(t, fd.Close())
				return nil
			}
		}
		t.Fatalf("no --output option")
		return nil
	}

	var calledBQCreateTable int
	mockBQ.FnCreateTable = func(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error {
		calledBQCreateTable++
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

	mockGH.MockGetArchiveURL = func(ctx context.Context, input *interfaces.GetArchiveURLInput) (*url.URL, error) {
		u := gt.R1(url.Parse("https://example.com/some/url.zip")).NoError(t)
		return u, nil
	}
	var calledMockListIssueComments int
	mockGH.MockListIssueComments = func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error) {
		calledMockListIssueComments++
		return nil, nil
	}
	var calledMockCreateIssueComment int
	mockGH.MockCreateIssueComment = func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error {
		calledMockCreateIssueComment++
		return nil
	}
	var calledMockGHCreateCheckRun int
	mockGH.MockCreateCheckRun = func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error) {
		calledMockGHCreateCheckRun++
		return 5, nil
	}
	var calledMockGHUpdateCheckRun int
	mockGH.MockUpdateCheckRun = func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
		gt.Equal(t, checkID, 5)
		gt.Equal(t, *opt.Status, "completed")
		gt.Equal(t, *opt.Conclusion, "success")
		calledMockGHUpdateCheckRun++
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
				Branch:   "main",
			},
			PullRequest: &model.GitHubPullRequest{
				Number:       123,
				ID:           12345,
				BaseBranch:   "main",
				BaseCommitID: "0f2324c367815ec3d928d21b892ce0ed9963aef3",
			},
		},
		InstallID: 12345,
	}))

	gt.Equal(t, calledBQCreateTable, 1)
	gt.Equal(t, calledBQInsert, 1)
	gt.Equal(t, calledMockListIssueComments, 1)
	gt.Equal(t, calledMockCreateIssueComment, 1)
	gt.Equal(t, calledMockGHCreateCheckRun, 1)
	gt.Equal(t, calledMockGHUpdateCheckRun, 1)

	var commitScan *model.Scan
	gt.NoError(t, mockStorage.Unmarshal("m-mizutani/octovy/commit/f7c8851da7c7fcc46212fccfb6c9c4bda520f1ca/scan.json.gz", &commitScan))
	gt.Equal(t, commitScan.GitHub.Owner, "m-mizutani")

	var branchScan *model.Scan
	gt.NoError(t, mockStorage.Unmarshal("m-mizutani/octovy/branch/main/scan.json.gz", &branchScan))
	gt.Equal(t, branchScan.GitHub.Owner, "m-mizutani")
}
