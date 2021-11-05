package db

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/enttest"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

type ScanResult struct {
	Repo            *ent.Repository
	Scan            *ent.Scan
	Packages        []*ent.PackageRecord
	Vulnerabilities []*ent.Vulnerability
}

type Interface interface {
	Open(dbType, dbConfig string) error
	Close() error

	CreateRepo(ctx *model.Context, repo *ent.Repository) (*ent.Repository, error)

	// Vulnerability
	PutVulnerabilities(ctx *model.Context, vulnerabilities []*ent.Vulnerability) error
	GetVulnerability(ctx *model.Context, id string) (*ent.Vulnerability, error)
	GetLatestVulnerabilities(ctx *model.Context, offset, limit int) ([]*ent.Vulnerability, error)
	GetVulnerabilityCount(ctx *model.Context) (int, error)

	// Severity
	CreateSeverity(ctx *model.Context, req *model.RequestSeverity) (*ent.Severity, error)
	DeleteSeverity(ctx *model.Context, id int) error
	GetSeverities(ctx *model.Context) ([]*ent.Severity, error)
	UpdateSeverity(ctx *model.Context, id int, req *model.RequestSeverity) error
	AssignSeverity(ctx *model.Context, vulnID string, id int) error

	PutPackages(ctx *model.Context, packages []*ent.PackageRecord) ([]*ent.PackageRecord, error)
	PutScan(ctx *model.Context, scan *ent.Scan, repo *ent.Repository, packages []*ent.PackageRecord) (*ent.Scan, error)
	PutVulnStatus(ctx *model.Context, repo *ent.Repository, status *ent.VulnStatus, userID int) (*ent.VulnStatus, error)
	GetVulnStatus(ctx *model.Context, repo *model.GitHubRepo) ([]*ent.VulnStatus, error)

	GetScan(ctx *model.Context, id string) (*ent.Scan, error)
	GetLatestScan(ctx *model.Context, branch model.GitHubBranch) (*ent.Scan, error)
	GetLatestScans(ctx *model.Context) ([]*ent.Scan, error)
	GetRepositories(ctx *model.Context) ([]*ent.Repository, error)
	GetRepositoriesWithVuln(ctx *model.Context, vulnID string) ([]*ent.Repository, error)

	// Rule
	GetRules(ctx *model.Context) ([]*ent.Rule, error)
	CreateRule(ctx *model.Context, req *model.RequestRule) (*ent.Rule, error)
	DeleteRule(ctx *model.Context, id int) error

	// Auth
	SaveAuthState(ctx *model.Context, state string, expiresAt int64) error
	HasAuthState(ctx *model.Context, state string, now int64) (bool, error)
	GetUser(ctx *model.Context, userID int) (*ent.User, error)
	PutUser(ctx *model.Context, user *ent.User) (int, error)
	PutSession(ctx *model.Context, ssn *ent.Session) error
	GetSession(ctx *model.Context, ssnID string, now int64) (*ent.Session, error)
	DeleteSession(ctx *model.Context, ssnID string) error
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
	dsn := "file:" + uuid.NewString() + "?mode=memory&cache=shared&_fk=1"
	db.client = enttest.Open(t, "sqlite3", dsn)
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

	if err := client.Schema.Create(model.NewContext()); err != nil {
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
