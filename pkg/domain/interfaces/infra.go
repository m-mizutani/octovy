package interfaces

import (
	"context"
	"io"
	"net/url"

	"cloud.google.com/go/bigquery"

	"github.com/google/go-github/v53/github"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type BigQuery interface {
	Insert(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error

	GetMetadata(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error)
	UpdateTable(ctx context.Context, table types.BQTableID, md bigquery.TableMetadataToUpdate, eTag string) error
	CreateTable(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error
}

type Storage interface {
	Put(ctx context.Context, key string, r io.ReadCloser) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
}

type GitHub interface {
	GetArchiveURL(ctx context.Context, input *GetArchiveURLInput) (*url.URL, error)
	CreateIssueComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error
	ListIssueComments(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error)
	MinimizeComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error
	CreateCheckRun(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error)
	UpdateCheckRun(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error
}

type GetArchiveURLInput struct {
	Owner     string
	Repo      string
	CommitID  string
	InstallID types.GitHubAppInstallID
}
