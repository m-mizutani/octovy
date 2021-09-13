package db

import (
	"context"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/enttest"
)

type ScanResult struct {
	Branch          *ent.Branch
	Scan            *ent.Scan
	Packages        []*ent.PackageRecord
	Vulnerabilities []*ent.Vulnerability
}

type BranchKey struct {
	Owner    string
	RepoName string
	Branch   string
}

type Interface interface {
	GetBranch(ctx context.Context, branch *BranchKey) (*ent.Branch, error)

	PutVulnerabilities(ctx context.Context, vulnerabilities []*ent.Vulnerability) error
	PutPackages(ctx context.Context, packages []*ent.PackageRecord) ([]*ent.PackageRecord, error)
	PutScan(ctx context.Context, scan *ent.Scan, branch *ent.Branch, packages []*ent.PackageRecord) (*ent.Scan, error)

	GetScan(ctx context.Context, id int) (*ent.Scan, error)
	GetLatestScan(ctx context.Context, owner, repoName, branch string) (*ent.Scan, error)

	Close() error
}

type Client struct {
	client *ent.Client

	lock  bool
	mutex sync.Mutex
}

func newClient() *Client {
	return &Client{}
}

func New(dbType, dbConfig string) (Interface, error) {
	client := newClient()
	if err := client.init(dbType, dbConfig); err != nil {
		return nil, err
	}
	return client, nil
}

func NewDBMock(t *testing.T) Interface {
	db := newClient()
	db.client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	db.lock = true
	return db
}

func (x *Client) init(dbType, dbConfig string) error {
	client, err := ent.Open(dbType, dbConfig)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	x.client = client

	if err := client.Schema.Create(context.Background()); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) Close() error {
	if err := x.client.Close(); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	return nil
}
