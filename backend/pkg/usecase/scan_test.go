package usecase_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/backend/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSet struct {
	db        interfaces.DBClient
	sqs       *aws.MockSQS
	trivy     *trivydb.TrivyDBMock
	githubapp *githubapp.Mock
}

func TestScanRepository(t *testing.T) {
	t.Run("Append a new vulnerability", func(t *testing.T) {
		now := time.Unix(10000, 0)
		uc, mock := setupScanRepositoryService(t, "../testdata/src/bundler.zip")
		svc := usecase.ExposeService(uc)
		svc.Infra.Utils.TimeNow = func() time.Time { return now }

		mock.trivy.AdvisoryMap["ruby-advisory-db"] = map[string][]*model.AdvisoryData{
			"rack": {
				{
					VulnID: "CVE-2020-8161",
					Data:   []byte(`{"PatchedVersions":["~\u003e 2.3.3","\u003e= 2.4.0"]}`),
				},
			},
		}
		mock.trivy.VulnerabilityMap["CVE-2020-8161"] = &types.Vulnerability{
			Title: "test vuln",
		}

		req := &model.ScanRepositoryRequest{
			ScanTarget: model.ScanTarget{
				GitHubBranch: model.GitHubBranch{
					GitHubRepo: model.GitHubRepo{
						Owner:    "five",
						RepoName: "blue",
					},
					Branch: "master",
				},
				CommitID:       "beefcafe",
				UpdatedAt:      1234,
				IsTargetBranch: true,
			},
			InstallID: 999,
			Feedback: &model.FeedbackOptions{
				PullReqID: model.Int(456),
			},
		}

		require.NoError(t, uc.ScanRepository(req))
		rackPkgs1, err := mock.db.FindPackageRecordsByName(model.PkgRubyGems, "rack")
		require.NoError(t, err)
		require.Equal(t, 1, len(rackPkgs1))
		assert.Equal(t, 1, len(rackPkgs1[0].Package.Vulnerabilities))
		assert.Contains(t, rackPkgs1[0].Package.Vulnerabilities, "CVE-2020-8161")

		// Add new one
		mock.trivy.AdvisoryMap["ruby-advisory-db"] = map[string][]*model.AdvisoryData{
			"rack": {
				{
					VulnID: "CVE-2020-8161",
					Data:   []byte(`{"PatchedVersions":["~\u003e 2.3.3","\u003e= 2.4.0"]}`),
				},
				{
					VulnID: "CVE-2020-9999",
					Data:   []byte(`{"PatchedVersions":["~\u003e 2.3.3","\u003e= 2.4.0"]}`),
				},
			},
		}
		mock.trivy.VulnerabilityMap["CVE-2020-9999"] = &types.Vulnerability{
			Title: "test vuln 2",
		}

		svc.Infra.Utils.TimeNow = func() time.Time { return now.Add(time.Second) }

		// Scan again
		require.NoError(t, uc.ScanRepository(req))
		rackPkgs2, err := mock.db.FindPackageRecordsByName(model.PkgRubyGems, "rack")
		require.NoError(t, err)
		require.Equal(t, 1, len(rackPkgs2))
		assert.Equal(t, 2, len(rackPkgs2[0].Package.Vulnerabilities))
		assert.Contains(t, rackPkgs2[0].Package.Vulnerabilities, "CVE-2020-8161")
		assert.Contains(t, rackPkgs2[0].Package.Vulnerabilities, "CVE-2020-9999")

		vulns, err := mock.db.FindLatestVulnerabilities(10)
		require.NoError(t, err)
		require.Equal(t, 2, len(vulns))
		assert.Contains(t, []string{vulns[0].VulnID, vulns[1].VulnID}, "CVE-2020-8161")
		assert.Contains(t, []string{vulns[0].VulnID, vulns[1].VulnID}, "CVE-2020-9999")

		require.Equal(t, 2, len(mock.sqs.Input))
		assert.NotNil(t, mock.sqs.Input[0].QueueUrl)
		assert.Equal(t, "https://feedback.queue.url", *mock.sqs.Input[0].QueueUrl)

		var feedbackReq model.FeedbackRequest
		require.NoError(t, json.Unmarshal([]byte(*mock.sqs.Input[0].MessageBody), &feedbackReq))
		require.NotNil(t, feedbackReq.Options.PullReqID)
		require.Nil(t, feedbackReq.Options.CheckSuiteID)
		assert.NotEmpty(t, feedbackReq.ReportID)
		assert.Equal(t, int64(999), feedbackReq.InstallID)
		assert.Equal(t, 456, *feedbackReq.Options.PullReqID)
	})
}

func TestScanBundler(t *testing.T) {
	uc, mock := setupScanRepositoryService(t, "../testdata/src/bundler.zip")

	mock.trivy.AdvisoryMap["ruby-advisory-db"] = map[string][]*model.AdvisoryData{
		"rack": {
			{
				VulnID: "CVE-2020-8161",
				// patched version is modified for test
				Data: []byte(`{"PatchedVersions":["~\u003e 2.3.3","\u003e= 2.4.0"]}`),
			},
		},
	}
	mock.trivy.VulnerabilityMap["CVE-2020-8161"] = &types.Vulnerability{
		Title: "test vuln",
	}

	req := &model.ScanRepositoryRequest{
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "five",
					RepoName: "blue",
				},
				Branch: "master",
			},
			CommitID:       "beefcafe",
			UpdatedAt:      1234,
			IsTargetBranch: true,
		},
		InstallID: 999,
	}

	require.NoError(t, uc.ScanRepository(req))
	assert.Equal(t, int64(999), mock.githubapp.InstallID)

	packages, err := mock.db.FindPackageRecordsByBranch(&req.GitHubBranch)
	require.NoError(t, err)
	assert.Equal(t, 41, len(packages))

	vulns, err := mock.db.FindLatestVulnerabilities(10)
	require.NoError(t, err)
	require.Equal(t, 1, len(vulns))
	assert.Equal(t, "CVE-2020-8161", vulns[0].VulnID)

	rackPkgs, err := mock.db.FindPackageRecordsByName(model.PkgRubyGems, "rack")
	require.NoError(t, err)
	require.Equal(t, 1, len(rackPkgs))
	assert.Equal(t, "rack", rackPkgs[0].Package.Name)
	assert.Contains(t, rackPkgs[0].Package.Vulnerabilities, "CVE-2020-8161")
}

func TestScanGoModule(t *testing.T) {
	uc, mock := setupScanRepositoryService(t, "../testdata/src/go_mod.zip")

	req := &model.ScanRepositoryRequest{
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "five",
					RepoName: "blue",
				},
				Branch: "master",
			},
			CommitID:       "beefcafe",
			UpdatedAt:      1234,
			IsTargetBranch: true,
		},
		InstallID: 999,
	}

	require.NoError(t, uc.ScanRepository(req))

	packages, err := mock.db.FindPackageRecordsByBranch(&req.GitHubBranch)
	require.NoError(t, err)
	assert.Equal(t, 147, len(packages))
}

func TestScanNPM(t *testing.T) {
	uc, mock := setupScanRepositoryService(t, "../testdata/src/npm.zip")

	req := &model.ScanRepositoryRequest{
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "five",
					RepoName: "blue",
				},
				Branch: "master",
			},
			CommitID:       "beefcafe",
			UpdatedAt:      1234,
			IsTargetBranch: true,
		},
		InstallID: 999,
	}

	require.NoError(t, uc.ScanRepository(req))

	packages, err := mock.db.FindPackageRecordsByBranch(&req.GitHubBranch)
	require.NoError(t, err)
	assert.Equal(t, 50, len(packages))
}

func genRSAKey(t *testing.T) []byte {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func setupScanRepositoryService(t *testing.T, scannedArchivePath string) (interfaces.Usecases, *mockSet) {
	pem := genRSAKey(t)
	base64PEM := base64.StdEncoding.EncodeToString(pem)
	const secretsARN = "arn:aws:secretsmanager:us-east-0:123456789012:secret:tutorials/MyFirstSecret-jiObOV"

	// mocking DB
	dbClient := newTestTable(t)
	inserted, err := dbClient.InsertRepo(&model.Repository{
		GitHubRepo: model.GitHubRepo{
			Owner:    "five",
			RepoName: "blue",
		},
	})
	require.NoError(t, err)
	require.True(t, inserted)

	// mocking SecretsManager
	newSM, mockSM := aws.NewMockSecretsManagerSet()
	mockSM.OutData[secretsARN] = map[string]string{
		"github_app_private_key": base64PEM,
		"github_app_id":          "123",
	}

	// mocking SQS
	newSQS, mockSQS := aws.NewMockSQSSet()

	cfg := &model.Config{
		SecretsARN:           secretsARN,
		GitHubEndpoint:       "https://ghe.example.org/api/v3",
		FeedbackRequestQueue: "https://feedback.queue.url",
		TableName:            dbClient.TableName(),
		S3Region:             "ap-northeast-0",
		S3Bucket:             "my-db-bucket",
		S3Prefix:             "test-prefix/",
	}

	// Build service and injects mocks
	uc := usecase.New(cfg)
	svc := usecase.ExposeService(uc)
	svc.Infra.NewSecretManager = newSM
	svc.Infra.NewDB = func(region, tableName string) (interfaces.DBClient, error) {
		return dbClient, nil
	}
	svc.Infra.NewSQS = newSQS

	// Build trivy mock
	newTrivyDBMock, trivyDBMock := trivydb.NewMock()
	svc.Infra.NewTrivyDB = newTrivyDBMock

	// Setup S3 mock
	newS3Mock, s3Mock := aws.NewMockS3()
	svc.Infra.NewS3 = newS3Mock
	s3Mock.Objects["my-db-bucket"] = map[string][]byte{
		"test-prefix/db/trivy.db.gz": []byte("boom!"),
	}

	newGitHubAppMock, gitHubAppMock := githubapp.NewMock()
	gitHubAppMock.GetCodeZipMock = func(repo *model.GitHubRepo, commitID string, w io.WriteCloser) error {
		r, err := os.Open(scannedArchivePath)
		require.NoError(t, err)
		_, err = io.Copy(w, r)
		require.NoError(t, err)
		return nil
	}
	svc.Infra.NewGitHubApp = newGitHubAppMock
	// Setup trivy DB

	t.Cleanup(func() {
		assert.Equal(t, int64(123), gitHubAppMock.AppID)
		assert.Equal(t, "https://ghe.example.org/api/v3", gitHubAppMock.Endpoint)
		assert.Equal(t, "ap-northeast-0", s3Mock.Region)
		require.Less(t, 0, len(s3Mock.GetInput))
		assert.Equal(t, "test-prefix/db/trivy.db.gz", *s3Mock.GetInput[0].Key)
	})

	return uc, &mockSet{
		db:        dbClient,
		sqs:       mockSQS,
		trivy:     trivyDBMock,
		githubapp: gitHubAppMock,
	}
}
