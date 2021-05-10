package infra

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

type Interfaces struct {
	// Factories
	NewDB            NewDB
	NewTrivyDB       NewTrivyDB
	NewSecretManager NewSecretManager
	NewSQS           NewSQS
	NewS3            NewS3
	NewGitHub        NewGitHub
	NewHTTP          NewHTTPClient // Interface set
	Utils            Utils
}

// AWS
// SecretsManager
type NewSecretManager func(region string) (SecretsManagerClient, error)
type SecretsManagerClient interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

// SQS
type NewSQS func(region string) (SQSClient, error)
type SQSClient interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

// S3
type NewS3 func(region string) (S3Client, error)
type S3Client interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

// DB
type NewDB func(region, tableName string) (DBClient, error)
type DBClient interface {
	InsertPackageRecord(*model.PackageRecord) (bool, error)
	RemovePackageRecord(*model.PackageRecord) error
	UpdatePackageRecord(*model.PackageRecord) error
	FindPackageRecordsByName(pkgType model.PkgType, pkgName string) ([]*model.PackageRecord, error)
	FindPackageRecordsByBranch(*model.GitHubBranch) ([]*model.PackageRecord, error)

	InsertScanReport(*model.ScanReport) error
	LookupScanReport(reportID string) (*model.ScanReport, error)
	FindScanLogsByBranch(branch *model.GitHubBranch, n int) ([]*model.ScanLog, error)
	FindScanLogsByCommit(commit *model.GitHubCommit, n int) ([]*model.ScanLog, error)

	InsertRepo(*model.Repository) (bool, error)
	UpdateBranchIfDefault(*model.GitHubRepo, *model.Branch) error
	SetRepoDefaultBranchName(*model.GitHubRepo, string) error
	FindRepo() ([]*model.Repository, error)
	FindRepoByOwner(owner string) ([]*model.Repository, error)
	FindRepoByFullName(owner, name string) (*model.Repository, error)

	UpdateBranch(branch *model.Branch) error
	LookupBranch(branch *model.GitHubBranch) (*model.Branch, error)
	FindLatestScannedBranch(repo *model.GitHubRepo, n int) ([]*model.Branch, error)

	InsertVulnerability(vuln *model.Vulnerability) error
	FindVulnerability(vulnID string) (*model.Vulnerability, error)
	FindLatestVulnerabilities(n int) ([]*model.Vulnerability, error)

	TableName() string
	Close() error
}

// HTTP
type NewHTTPClient func(http.RoundTripper) *http.Client

// GitHub
type NewGitHub func() GitHubClient
type GitHubClient interface {
	ListReleases(owner, repo string) ([]*github.RepositoryRelease, error)
	DownloadReleaseAsset(owner, repo string, assetID int64) (io.ReadCloser, error)
}

// Trivy DB
type NewTrivyDB func(dbPath string) (TrivyDBClient, error)
type TrivyDBClient interface {
	GetAdvisories(source, pkgName string) ([]*model.AdvisoryData, error)
	GetVulnerability(vulnID string) (*types.Vulnerability, error)
	GetDBMeta() (*model.TrivyDBMeta, error)
}

// Utils
type TimeNow func() time.Time
type WriteFile func(r io.Reader, path string) error
type OpenZip func(path string) (*zip.ReadCloser, error)
type TempFile func(dir, pattern string) (f *os.File, err error)
type Remove func(name string) error
