package infra

import (
	"archive/zip"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

type Interfaces struct {
	// Factories
	NewDB            NewDB
	NewSecretManager NewSecretManager
	NewSQS           NewSQS
	NewHTTP          NewHTTPClient // Interface set
	FS               FS
}

// AWS
type SecretsManagerClient interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

type NewSecretManager func(region string) (SecretsManagerClient, error)

type SQSClient interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

type NewSQS func(region string) (SQSClient, error)

// DB
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

type NewDB func(region, tableName string) (DBClient, error)

// FileSystem
type FS interface {
	WriteFile(r io.Reader, path string) error
	OpenZip(path string) (*zip.ReadCloser, error)
	TempFile(dir, pattern string) (f *os.File, err error)
	Remove(name string) error
}

// HTTP
type NewHTTPClient func(http.RoundTripper) *http.Client
