package infra

import (
	"archive/zip"
	"io"
	"net/http"
	"os"

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
	FS               FS
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
	InsertPackage(*model.Package) error
	DeletePackage(*model.Package) error
	FindPackagesByName(pkgType model.PkgType, pkgName string) ([]*model.Package, error)
	FindPackagesByBranch(*model.GitHubBranch) ([]*model.Package, error)

	InsertRepo(*model.Repository) (bool, error)
	SetRepoBranches(*model.GitHubRepo, []string) error
	SetRepoDefaultBranch(*model.GitHubRepo, string) error
	FindRepo() ([]*model.Repository, error)
	FindRepoByOwner(owner string) ([]*model.Repository, error)
	FindRepoByFullName(owner, name string) (*model.Repository, error)
	TableName() string
	Close() error
}

// FileSystem
type FS interface {
	WriteFile(r io.Reader, path string) error
	OpenZip(path string) (*zip.ReadCloser, error)
	TempFile(dir, pattern string) (f *os.File, err error)
	Remove(name string) error
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
}
