package usecase_test

import (
	_ "embed"

	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

//go:embed testdata/octovy-test-code-main.zip
var testCodeZip []byte

func TestScanGitHubRepo(t *testing.T) {
	mockGH := &ghAppMock{}
	mockHTTP := &httpMock{}

	uc := usecase.New(infra.New(
		infra.WithGitHubApp(mockGH),
		infra.WithHTTPClient(mockHTTP),
	))

	ctx := model.NewContext()

	mockGH.mockGetArchiveURL = func(ctx *model.Context, input *githubapp.GetArchiveURLInput) (*url.URL, error) {
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

	gt.NoError(t, uc.ScanGitHubRepo(ctx, &usecase.ScanGitHubRepoInput{
		Owner:     "m-mizutani",
		Repo:      "octovy",
		CommitID:  "1234567890",
		InstallID: 12345,
	}))
}

type ghAppMock struct {
	mockGetArchiveURL func(ctx *model.Context, input *githubapp.GetArchiveURLInput) (*url.URL, error)
}

func (x *ghAppMock) GetArchiveURL(ctx *model.Context, input *githubapp.GetArchiveURLInput) (*url.URL, error) {
	return x.mockGetArchiveURL(ctx, input)
}

type httpMock struct {
	mockDo func(req *http.Request) (*http.Response, error)
}

func (x *httpMock) Do(req *http.Request) (*http.Response, error) {
	return x.mockDo(req)
}
