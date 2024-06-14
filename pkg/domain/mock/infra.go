// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/google/go-github/v53/github"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"net/url"
	"sync"
)

// Ensure, that GitHubMock does implement interfaces.GitHub.
// If this is not the case, regenerate this file with moq.
var _ interfaces.GitHub = &GitHubMock{}

// GitHubMock is a mock implementation of interfaces.GitHub.
//
//	func TestSomethingThatUsesGitHub(t *testing.T) {
//
//		// make and configure a mocked interfaces.GitHub
//		mockedGitHub := &GitHubMock{
//			CreateCheckRunFunc: func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error) {
//				panic("mock out the CreateCheckRun method")
//			},
//			CreateIssueCommentFunc: func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error {
//				panic("mock out the CreateIssueComment method")
//			},
//			GetArchiveURLFunc: func(ctx context.Context, input *interfaces.GetArchiveURLInput) (*url.URL, error) {
//				panic("mock out the GetArchiveURL method")
//			},
//			ListIssueCommentsFunc: func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error) {
//				panic("mock out the ListIssueComments method")
//			},
//			MinimizeCommentFunc: func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error {
//				panic("mock out the MinimizeComment method")
//			},
//			UpdateCheckRunFunc: func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
//				panic("mock out the UpdateCheckRun method")
//			},
//		}
//
//		// use mockedGitHub in code that requires interfaces.GitHub
//		// and then make assertions.
//
//	}
type GitHubMock struct {
	// CreateCheckRunFunc mocks the CreateCheckRun method.
	CreateCheckRunFunc func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error)

	// CreateIssueCommentFunc mocks the CreateIssueComment method.
	CreateIssueCommentFunc func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error

	// GetArchiveURLFunc mocks the GetArchiveURL method.
	GetArchiveURLFunc func(ctx context.Context, input *interfaces.GetArchiveURLInput) (*url.URL, error)

	// ListIssueCommentsFunc mocks the ListIssueComments method.
	ListIssueCommentsFunc func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error)

	// MinimizeCommentFunc mocks the MinimizeComment method.
	MinimizeCommentFunc func(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error

	// UpdateCheckRunFunc mocks the UpdateCheckRun method.
	UpdateCheckRunFunc func(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateCheckRun holds details about calls to the CreateCheckRun method.
		CreateCheckRun []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID types.GitHubAppInstallID
			// Repo is the repo argument value.
			Repo *model.GitHubRepo
			// Commit is the commit argument value.
			Commit string
		}
		// CreateIssueComment holds details about calls to the CreateIssueComment method.
		CreateIssueComment []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Repo is the repo argument value.
			Repo *model.GitHubRepo
			// ID is the id argument value.
			ID types.GitHubAppInstallID
			// PrID is the prID argument value.
			PrID int
			// Body is the body argument value.
			Body string
		}
		// GetArchiveURL holds details about calls to the GetArchiveURL method.
		GetArchiveURL []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input *interfaces.GetArchiveURLInput
		}
		// ListIssueComments holds details about calls to the ListIssueComments method.
		ListIssueComments []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Repo is the repo argument value.
			Repo *model.GitHubRepo
			// ID is the id argument value.
			ID types.GitHubAppInstallID
			// PrID is the prID argument value.
			PrID int
		}
		// MinimizeComment holds details about calls to the MinimizeComment method.
		MinimizeComment []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Repo is the repo argument value.
			Repo *model.GitHubRepo
			// ID is the id argument value.
			ID types.GitHubAppInstallID
			// SubjectID is the subjectID argument value.
			SubjectID string
		}
		// UpdateCheckRun holds details about calls to the UpdateCheckRun method.
		UpdateCheckRun []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID types.GitHubAppInstallID
			// Repo is the repo argument value.
			Repo *model.GitHubRepo
			// CheckID is the checkID argument value.
			CheckID int64
			// Opt is the opt argument value.
			Opt *github.UpdateCheckRunOptions
		}
	}
	lockCreateCheckRun     sync.RWMutex
	lockCreateIssueComment sync.RWMutex
	lockGetArchiveURL      sync.RWMutex
	lockListIssueComments  sync.RWMutex
	lockMinimizeComment    sync.RWMutex
	lockUpdateCheckRun     sync.RWMutex
}

// CreateCheckRun calls CreateCheckRunFunc.
func (mock *GitHubMock) CreateCheckRun(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, commit string) (int64, error) {
	if mock.CreateCheckRunFunc == nil {
		panic("GitHubMock.CreateCheckRunFunc: method is nil but GitHub.CreateCheckRun was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		ID     types.GitHubAppInstallID
		Repo   *model.GitHubRepo
		Commit string
	}{
		Ctx:    ctx,
		ID:     id,
		Repo:   repo,
		Commit: commit,
	}
	mock.lockCreateCheckRun.Lock()
	mock.calls.CreateCheckRun = append(mock.calls.CreateCheckRun, callInfo)
	mock.lockCreateCheckRun.Unlock()
	return mock.CreateCheckRunFunc(ctx, id, repo, commit)
}

// CreateCheckRunCalls gets all the calls that were made to CreateCheckRun.
// Check the length with:
//
//	len(mockedGitHub.CreateCheckRunCalls())
func (mock *GitHubMock) CreateCheckRunCalls() []struct {
	Ctx    context.Context
	ID     types.GitHubAppInstallID
	Repo   *model.GitHubRepo
	Commit string
} {
	var calls []struct {
		Ctx    context.Context
		ID     types.GitHubAppInstallID
		Repo   *model.GitHubRepo
		Commit string
	}
	mock.lockCreateCheckRun.RLock()
	calls = mock.calls.CreateCheckRun
	mock.lockCreateCheckRun.RUnlock()
	return calls
}

// CreateIssueComment calls CreateIssueCommentFunc.
func (mock *GitHubMock) CreateIssueComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int, body string) error {
	if mock.CreateIssueCommentFunc == nil {
		panic("GitHubMock.CreateIssueCommentFunc: method is nil but GitHub.CreateIssueComment was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Repo *model.GitHubRepo
		ID   types.GitHubAppInstallID
		PrID int
		Body string
	}{
		Ctx:  ctx,
		Repo: repo,
		ID:   id,
		PrID: prID,
		Body: body,
	}
	mock.lockCreateIssueComment.Lock()
	mock.calls.CreateIssueComment = append(mock.calls.CreateIssueComment, callInfo)
	mock.lockCreateIssueComment.Unlock()
	return mock.CreateIssueCommentFunc(ctx, repo, id, prID, body)
}

// CreateIssueCommentCalls gets all the calls that were made to CreateIssueComment.
// Check the length with:
//
//	len(mockedGitHub.CreateIssueCommentCalls())
func (mock *GitHubMock) CreateIssueCommentCalls() []struct {
	Ctx  context.Context
	Repo *model.GitHubRepo
	ID   types.GitHubAppInstallID
	PrID int
	Body string
} {
	var calls []struct {
		Ctx  context.Context
		Repo *model.GitHubRepo
		ID   types.GitHubAppInstallID
		PrID int
		Body string
	}
	mock.lockCreateIssueComment.RLock()
	calls = mock.calls.CreateIssueComment
	mock.lockCreateIssueComment.RUnlock()
	return calls
}

// GetArchiveURL calls GetArchiveURLFunc.
func (mock *GitHubMock) GetArchiveURL(ctx context.Context, input *interfaces.GetArchiveURLInput) (*url.URL, error) {
	if mock.GetArchiveURLFunc == nil {
		panic("GitHubMock.GetArchiveURLFunc: method is nil but GitHub.GetArchiveURL was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Input *interfaces.GetArchiveURLInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockGetArchiveURL.Lock()
	mock.calls.GetArchiveURL = append(mock.calls.GetArchiveURL, callInfo)
	mock.lockGetArchiveURL.Unlock()
	return mock.GetArchiveURLFunc(ctx, input)
}

// GetArchiveURLCalls gets all the calls that were made to GetArchiveURL.
// Check the length with:
//
//	len(mockedGitHub.GetArchiveURLCalls())
func (mock *GitHubMock) GetArchiveURLCalls() []struct {
	Ctx   context.Context
	Input *interfaces.GetArchiveURLInput
} {
	var calls []struct {
		Ctx   context.Context
		Input *interfaces.GetArchiveURLInput
	}
	mock.lockGetArchiveURL.RLock()
	calls = mock.calls.GetArchiveURL
	mock.lockGetArchiveURL.RUnlock()
	return calls
}

// ListIssueComments calls ListIssueCommentsFunc.
func (mock *GitHubMock) ListIssueComments(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, prID int) ([]*model.GitHubIssueComment, error) {
	if mock.ListIssueCommentsFunc == nil {
		panic("GitHubMock.ListIssueCommentsFunc: method is nil but GitHub.ListIssueComments was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Repo *model.GitHubRepo
		ID   types.GitHubAppInstallID
		PrID int
	}{
		Ctx:  ctx,
		Repo: repo,
		ID:   id,
		PrID: prID,
	}
	mock.lockListIssueComments.Lock()
	mock.calls.ListIssueComments = append(mock.calls.ListIssueComments, callInfo)
	mock.lockListIssueComments.Unlock()
	return mock.ListIssueCommentsFunc(ctx, repo, id, prID)
}

// ListIssueCommentsCalls gets all the calls that were made to ListIssueComments.
// Check the length with:
//
//	len(mockedGitHub.ListIssueCommentsCalls())
func (mock *GitHubMock) ListIssueCommentsCalls() []struct {
	Ctx  context.Context
	Repo *model.GitHubRepo
	ID   types.GitHubAppInstallID
	PrID int
} {
	var calls []struct {
		Ctx  context.Context
		Repo *model.GitHubRepo
		ID   types.GitHubAppInstallID
		PrID int
	}
	mock.lockListIssueComments.RLock()
	calls = mock.calls.ListIssueComments
	mock.lockListIssueComments.RUnlock()
	return calls
}

// MinimizeComment calls MinimizeCommentFunc.
func (mock *GitHubMock) MinimizeComment(ctx context.Context, repo *model.GitHubRepo, id types.GitHubAppInstallID, subjectID string) error {
	if mock.MinimizeCommentFunc == nil {
		panic("GitHubMock.MinimizeCommentFunc: method is nil but GitHub.MinimizeComment was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Repo      *model.GitHubRepo
		ID        types.GitHubAppInstallID
		SubjectID string
	}{
		Ctx:       ctx,
		Repo:      repo,
		ID:        id,
		SubjectID: subjectID,
	}
	mock.lockMinimizeComment.Lock()
	mock.calls.MinimizeComment = append(mock.calls.MinimizeComment, callInfo)
	mock.lockMinimizeComment.Unlock()
	return mock.MinimizeCommentFunc(ctx, repo, id, subjectID)
}

// MinimizeCommentCalls gets all the calls that were made to MinimizeComment.
// Check the length with:
//
//	len(mockedGitHub.MinimizeCommentCalls())
func (mock *GitHubMock) MinimizeCommentCalls() []struct {
	Ctx       context.Context
	Repo      *model.GitHubRepo
	ID        types.GitHubAppInstallID
	SubjectID string
} {
	var calls []struct {
		Ctx       context.Context
		Repo      *model.GitHubRepo
		ID        types.GitHubAppInstallID
		SubjectID string
	}
	mock.lockMinimizeComment.RLock()
	calls = mock.calls.MinimizeComment
	mock.lockMinimizeComment.RUnlock()
	return calls
}

// UpdateCheckRun calls UpdateCheckRunFunc.
func (mock *GitHubMock) UpdateCheckRun(ctx context.Context, id types.GitHubAppInstallID, repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error {
	if mock.UpdateCheckRunFunc == nil {
		panic("GitHubMock.UpdateCheckRunFunc: method is nil but GitHub.UpdateCheckRun was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		ID      types.GitHubAppInstallID
		Repo    *model.GitHubRepo
		CheckID int64
		Opt     *github.UpdateCheckRunOptions
	}{
		Ctx:     ctx,
		ID:      id,
		Repo:    repo,
		CheckID: checkID,
		Opt:     opt,
	}
	mock.lockUpdateCheckRun.Lock()
	mock.calls.UpdateCheckRun = append(mock.calls.UpdateCheckRun, callInfo)
	mock.lockUpdateCheckRun.Unlock()
	return mock.UpdateCheckRunFunc(ctx, id, repo, checkID, opt)
}

// UpdateCheckRunCalls gets all the calls that were made to UpdateCheckRun.
// Check the length with:
//
//	len(mockedGitHub.UpdateCheckRunCalls())
func (mock *GitHubMock) UpdateCheckRunCalls() []struct {
	Ctx     context.Context
	ID      types.GitHubAppInstallID
	Repo    *model.GitHubRepo
	CheckID int64
	Opt     *github.UpdateCheckRunOptions
} {
	var calls []struct {
		Ctx     context.Context
		ID      types.GitHubAppInstallID
		Repo    *model.GitHubRepo
		CheckID int64
		Opt     *github.UpdateCheckRunOptions
	}
	mock.lockUpdateCheckRun.RLock()
	calls = mock.calls.UpdateCheckRun
	mock.lockUpdateCheckRun.RUnlock()
	return calls
}

// Ensure, that BigQueryMock does implement interfaces.BigQuery.
// If this is not the case, regenerate this file with moq.
var _ interfaces.BigQuery = &BigQueryMock{}

// BigQueryMock is a mock implementation of interfaces.BigQuery.
//
//	func TestSomethingThatUsesBigQuery(t *testing.T) {
//
//		// make and configure a mocked interfaces.BigQuery
//		mockedBigQuery := &BigQueryMock{
//			CreateTableFunc: func(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error {
//				panic("mock out the CreateTable method")
//			},
//			GetMetadataFunc: func(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error) {
//				panic("mock out the GetMetadata method")
//			},
//			InsertFunc: func(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error {
//				panic("mock out the Insert method")
//			},
//			UpdateTableFunc: func(ctx context.Context, table types.BQTableID, md bigquery.TableMetadataToUpdate, eTag string) error {
//				panic("mock out the UpdateTable method")
//			},
//		}
//
//		// use mockedBigQuery in code that requires interfaces.BigQuery
//		// and then make assertions.
//
//	}
type BigQueryMock struct {
	// CreateTableFunc mocks the CreateTable method.
	CreateTableFunc func(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error

	// GetMetadataFunc mocks the GetMetadata method.
	GetMetadataFunc func(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error)

	// InsertFunc mocks the Insert method.
	InsertFunc func(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error

	// UpdateTableFunc mocks the UpdateTable method.
	UpdateTableFunc func(ctx context.Context, table types.BQTableID, md bigquery.TableMetadataToUpdate, eTag string) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateTable holds details about calls to the CreateTable method.
		CreateTable []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Table is the table argument value.
			Table types.BQTableID
			// Md is the md argument value.
			Md *bigquery.TableMetadata
		}
		// GetMetadata holds details about calls to the GetMetadata method.
		GetMetadata []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Table is the table argument value.
			Table types.BQTableID
		}
		// Insert holds details about calls to the Insert method.
		Insert []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// TableID is the tableID argument value.
			TableID types.BQTableID
			// Schema is the schema argument value.
			Schema bigquery.Schema
			// Data is the data argument value.
			Data any
		}
		// UpdateTable holds details about calls to the UpdateTable method.
		UpdateTable []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Table is the table argument value.
			Table types.BQTableID
			// Md is the md argument value.
			Md bigquery.TableMetadataToUpdate
			// ETag is the eTag argument value.
			ETag string
		}
	}
	lockCreateTable sync.RWMutex
	lockGetMetadata sync.RWMutex
	lockInsert      sync.RWMutex
	lockUpdateTable sync.RWMutex
}

// CreateTable calls CreateTableFunc.
func (mock *BigQueryMock) CreateTable(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error {
	if mock.CreateTableFunc == nil {
		panic("BigQueryMock.CreateTableFunc: method is nil but BigQuery.CreateTable was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Table types.BQTableID
		Md    *bigquery.TableMetadata
	}{
		Ctx:   ctx,
		Table: table,
		Md:    md,
	}
	mock.lockCreateTable.Lock()
	mock.calls.CreateTable = append(mock.calls.CreateTable, callInfo)
	mock.lockCreateTable.Unlock()
	return mock.CreateTableFunc(ctx, table, md)
}

// CreateTableCalls gets all the calls that were made to CreateTable.
// Check the length with:
//
//	len(mockedBigQuery.CreateTableCalls())
func (mock *BigQueryMock) CreateTableCalls() []struct {
	Ctx   context.Context
	Table types.BQTableID
	Md    *bigquery.TableMetadata
} {
	var calls []struct {
		Ctx   context.Context
		Table types.BQTableID
		Md    *bigquery.TableMetadata
	}
	mock.lockCreateTable.RLock()
	calls = mock.calls.CreateTable
	mock.lockCreateTable.RUnlock()
	return calls
}

// GetMetadata calls GetMetadataFunc.
func (mock *BigQueryMock) GetMetadata(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error) {
	if mock.GetMetadataFunc == nil {
		panic("BigQueryMock.GetMetadataFunc: method is nil but BigQuery.GetMetadata was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Table types.BQTableID
	}{
		Ctx:   ctx,
		Table: table,
	}
	mock.lockGetMetadata.Lock()
	mock.calls.GetMetadata = append(mock.calls.GetMetadata, callInfo)
	mock.lockGetMetadata.Unlock()
	return mock.GetMetadataFunc(ctx, table)
}

// GetMetadataCalls gets all the calls that were made to GetMetadata.
// Check the length with:
//
//	len(mockedBigQuery.GetMetadataCalls())
func (mock *BigQueryMock) GetMetadataCalls() []struct {
	Ctx   context.Context
	Table types.BQTableID
} {
	var calls []struct {
		Ctx   context.Context
		Table types.BQTableID
	}
	mock.lockGetMetadata.RLock()
	calls = mock.calls.GetMetadata
	mock.lockGetMetadata.RUnlock()
	return calls
}

// Insert calls InsertFunc.
func (mock *BigQueryMock) Insert(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error {
	if mock.InsertFunc == nil {
		panic("BigQueryMock.InsertFunc: method is nil but BigQuery.Insert was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		TableID types.BQTableID
		Schema  bigquery.Schema
		Data    any
	}{
		Ctx:     ctx,
		TableID: tableID,
		Schema:  schema,
		Data:    data,
	}
	mock.lockInsert.Lock()
	mock.calls.Insert = append(mock.calls.Insert, callInfo)
	mock.lockInsert.Unlock()
	return mock.InsertFunc(ctx, tableID, schema, data)
}

// InsertCalls gets all the calls that were made to Insert.
// Check the length with:
//
//	len(mockedBigQuery.InsertCalls())
func (mock *BigQueryMock) InsertCalls() []struct {
	Ctx     context.Context
	TableID types.BQTableID
	Schema  bigquery.Schema
	Data    any
} {
	var calls []struct {
		Ctx     context.Context
		TableID types.BQTableID
		Schema  bigquery.Schema
		Data    any
	}
	mock.lockInsert.RLock()
	calls = mock.calls.Insert
	mock.lockInsert.RUnlock()
	return calls
}

// UpdateTable calls UpdateTableFunc.
func (mock *BigQueryMock) UpdateTable(ctx context.Context, table types.BQTableID, md bigquery.TableMetadataToUpdate, eTag string) error {
	if mock.UpdateTableFunc == nil {
		panic("BigQueryMock.UpdateTableFunc: method is nil but BigQuery.UpdateTable was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Table types.BQTableID
		Md    bigquery.TableMetadataToUpdate
		ETag  string
	}{
		Ctx:   ctx,
		Table: table,
		Md:    md,
		ETag:  eTag,
	}
	mock.lockUpdateTable.Lock()
	mock.calls.UpdateTable = append(mock.calls.UpdateTable, callInfo)
	mock.lockUpdateTable.Unlock()
	return mock.UpdateTableFunc(ctx, table, md, eTag)
}

// UpdateTableCalls gets all the calls that were made to UpdateTable.
// Check the length with:
//
//	len(mockedBigQuery.UpdateTableCalls())
func (mock *BigQueryMock) UpdateTableCalls() []struct {
	Ctx   context.Context
	Table types.BQTableID
	Md    bigquery.TableMetadataToUpdate
	ETag  string
} {
	var calls []struct {
		Ctx   context.Context
		Table types.BQTableID
		Md    bigquery.TableMetadataToUpdate
		ETag  string
	}
	mock.lockUpdateTable.RLock()
	calls = mock.calls.UpdateTable
	mock.lockUpdateTable.RUnlock()
	return calls
}