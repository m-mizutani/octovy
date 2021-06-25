package interfaces

import (
	"archive/zip"
	"io"
	"os"
	"time"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

type Infra struct {
	// Factories
	NewDB            NewDB
	NewTrivyDB       NewTrivyDB
	NewSecretManager NewSecretManager
	NewSQS           NewSQS
	NewS3            NewS3
	NewGitHubCom     NewGitHubCom
	NewGitHubApp     NewGitHubApp
	NewGitHubAuth    NewGitHubAuth
	Utils            *Utils
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
	FindOwners() ([]*model.Owner, error)

	UpdateBranch(branch *model.Branch) error
	LookupBranch(branch *model.GitHubBranch) (*model.Branch, error)
	FindLatestScannedBranch(repo *model.GitHubRepo, n int) ([]*model.Branch, error)

	InsertVulnerability(vuln *model.Vulnerability) error
	FindVulnerability(vulnID string) (*model.Vulnerability, error)
	FindLatestVulnerabilities(n int) ([]*model.Vulnerability, error)
	GetVulnerabilities(vulnIDs []string) ([]*model.Vulnerability, error)

	PutVulnStatus(status *model.VulnStatus) error
	GetVulnStatus(repo *model.GitHubRepo, now int64) ([]*model.VulnStatus, error)
	GetVulnStatusLogs(repo *model.GitHubRepo, key *model.VulnPackageKey) ([]*model.VulnStatus, error)

	SaveAuthState(state string, expiresAt int64) error
	HasAuthState(state string, now int64) (bool, error)
	PutUser(user *model.User) error
	GetUser(userID string) (*model.User, error)
	PutUserPermissions(perm *model.UserPermissions) error
	GetUserPermissions(userID string) (*model.UserPermissions, error)

	PutGitHubToken(token *model.GitHubToken) error
	GetGitHubToken(userID string) (*model.GitHubToken, error)
	PutSession(ssn *model.Session) error
	GetSession(token string, now int64) (*model.Session, error)
	DeleteSession(token string) error

	TableName() string
	Close() error
}

// GitHubCom accesses only github.com to download trivy DB. It does not require API endpoint configuration and credentials
type GitHubCom interface {
	ListReleases(owner, repo string) ([]*github.RepositoryRelease, error)
	DownloadReleaseAsset(owner, repo string, assetID int64) (io.ReadCloser, error)
}
type NewGitHubCom func() GitHubCom

// GitHubApp is GitHub App interface that requires both of App ID and Install ID. Additionally it needs to change API endpoint for GitHub Enterprise
type GitHubApp interface {
	GetCodeZip(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error
	CreateIssueComment(repo *model.GitHubRepo, prID int, body string) error
	CreateCheckRun(repo *model.GitHubRepo, commit string) (int64, error)
	UpdateCheckRun(repo *model.GitHubRepo, checkID int64, opt *github.UpdateCheckRunOptions) error
}
type NewGitHubApp func(appID, installID int64, pem []byte, endpoint string) GitHubApp

// GitHubAuth is for authentication of GitHub user. It does not require App ID and Install ID, but requires API endpoint configuration for GitHub Enterprise
type GitHubAuth interface {
	SetToken(token *model.GitHubToken)
	Authenticate(clientID, clientSecret, code string) (*model.GitHubToken, error)

	GetUser() (*model.User, error)
	GetInstallations() ([]*github.Installation, error)
	GetInstalledRepositories(installID int64) ([]*github.Repository, error)
}
type NewGitHubAuth func(apiEndpoint, webEndpoint string) GitHubAuth

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
type GenerateToken func(n int) string

type Utils struct {
	TimeNow       TimeNow
	TempFile      TempFile
	OpenZip       OpenZip
	Remove        Remove
	GenerateToken GenerateToken
}
