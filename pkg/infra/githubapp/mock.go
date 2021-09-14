package githubapp

import (
	"io"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type Mock struct {
	GetCodeZipMock         func(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error
	CreateIssueCommentMock func(repo *model.GitHubRepo, prID int, body string) error
	CreateCheckRunMock     func(repo *model.GitHubRepo, commit string) (int64, error)
	UpdateCheckRunMock     func(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error

	AppID     int64
	InstallID int64
	PEM       []byte
	Endpoint  string
}

func NewMock() (Factory, *Mock) {
	mock := &Mock{}
	return func(appID, installID int64, pem []byte, endpoint string) Interface {
		mock.AppID = appID
		mock.InstallID = installID
		mock.PEM = pem
		mock.Endpoint = endpoint

		return mock
	}, mock
}

func (x *Mock) GetCodeZip(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error {
	return x.GetCodeZipMock(repo, commitID, w)
}
func (x *Mock) CreateIssueComment(repo *model.GitHubRepo, prID int, body string) error {
	return x.CreateIssueCommentMock(repo, prID, body)
}
func (x *Mock) CreateCheckRun(repo *model.GitHubRepo, commit string) (int64, error) {
	return x.CreateCheckRunMock(repo, commit)
}
func (x *Mock) UpdateCheckRun(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
	return x.UpdateCheckRunMock(repo, checkID, opt)
}
