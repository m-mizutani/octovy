package db

import (
	"context"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/enttest"
)

type ScanResult struct {
	Repo            *ent.Repository
	Scan            *ent.Scan
	Packages        []*ent.PackageRecord
	Vulnerabilities []*ent.Vulnerability
}

type Interface interface {
	Open(dbType, dbConfig string) error
	CreateRepo(ctx context.Context, repo *ent.Repository) (*ent.Repository, error)

	PutVulnerabilities(ctx context.Context, vulnerabilities []*ent.Vulnerability) error
	PutPackages(ctx context.Context, packages []*ent.PackageRecord) ([]*ent.PackageRecord, error)
	PutScan(ctx context.Context, scan *ent.Scan, repo *ent.Repository, packages []*ent.PackageRecord) (*ent.Scan, error)

	GetScan(ctx context.Context, id string) (*ent.Scan, error)
	GetLatestScan(ctx context.Context, branch model.GitHubBranch) (*ent.Scan, error)

	Close() error
}

type Factory func(dbType, dbConfig string) (Interface, error)

type Client struct {
	client *ent.Client

	disableOpen bool
	lock        bool
	mutex       sync.Mutex
}

func newClient() *Client {
	return &Client{}
}

func New() *Client {
	return newClient()
}

func NewMock(t *testing.T) *Client {
	db := newClient()
	db.client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	db.disableOpen = true
	db.lock = true
	return db
}

func (x *Client) Open(dbType, dbConfig string) error {
	if x.disableOpen {
		return nil
	}

	client, err := ent.Open(dbType, dbConfig)
	if err != nil {
		return model.ErrDatabaseUnexpected.Wrap(err)
	}
	x.client = client

	if err := client.Schema.Create(context.Background()); err != nil {
		return model.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) Close() error {
	if err := x.client.Close(); err != nil {
		return model.ErrDatabaseUnexpected.Wrap(err)
	}
	return nil
}
