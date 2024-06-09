package interfaces

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/url"

	"github.com/google/go-github/v53/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type StorageMock struct {
	Data map[string][]byte
}

var _ Storage = (*StorageMock)(nil)

func NewStorageMock() *StorageMock {
	return &StorageMock{
		Data: make(map[string][]byte),
	}
}

// Get implements Storage.
func (s *StorageMock) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if data, ok := s.Data[key]; ok {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return nil, nil
}

// Put implements Storage.
func (s *StorageMock) Put(ctx context.Context, key string, r io.ReadCloser) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.Data[key] = data
	return nil
}

func (s *StorageMock) Unmarshal(key string, v interface{}) error {
	data, ok := s.Data[key]
	if !ok {
		return io.EOF
	}

	if err := json.Unmarshal(data, v); err != nil {
		return goerr.Wrap(err, "Failed to unmarshal data")
	}

	return nil
}

type GitHubMock struct {
	MockGetArchiveURL      func(ctx context.Context, input *GetArchiveURLInput) (*url.URL, error)
	MockCreateIssueComment func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error
	MockListIssueComments  func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error)
	MockMinimizeComment    func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error
	MockCreateCheckRun     func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error)
	MockUpdateCheckRun     func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error
}

func (x *GitHubMock) GetArchiveURL(ctx context.Context, input *GetArchiveURLInput) (*url.URL, error) {
	return x.MockGetArchiveURL(ctx, input)
}

func (x *GitHubMock) CreateIssueComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error {
	return x.MockCreateIssueComment(ctx, repo, id, prID, body)
}

func (x *GitHubMock) ListIssueComments(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error) {
	return x.MockListIssueComments(ctx, repo, id, prID)
}

func (x *GitHubMock) MinimizeComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error {
	return x.MockMinimizeComment(ctx, repo, id, subjectID)
}

func (x *GitHubMock) CreateCheckRun(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error) {
	return x.MockCreateCheckRun(ctx, id, repo, commit)
}

func (x *GitHubMock) UpdateCheckRun(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
	return x.MockUpdateCheckRun(ctx, id, repo, checkID, opt)
}
