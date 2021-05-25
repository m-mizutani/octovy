package githubapp

import (
	"io"

	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

type Mock struct {
	GetCodeZipMock         func(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error
	CreateIssueCommentMock func(repo *model.GitHubRepo, prID int, body string) error
	AppID                  int64
	InstallID              int64
	PEM                    []byte
	Endpoint               string
}

func NewMock() (interfaces.NewGitHubApp, *Mock) {
	mock := &Mock{}
	return func(appID, installID int64, pem []byte, endpoint string) interfaces.GitHubApp {
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
