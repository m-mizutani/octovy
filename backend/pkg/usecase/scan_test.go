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
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanRepository(t *testing.T) {
	t.Run("Append a new vulnerability", func(t *testing.T) {
		now := time.Unix(10000, 0)
		svc, dbClient, trivyDBMock := setupScanRepositoryService(t, "../testdata/src/bundler.zip")
		svc.Utils.TimeNow = func() time.Time { return now }

		trivyDBMock.AdvisoryMap["ruby-advisory-db"] = map[string][]*model.AdvisoryData{
			"rack": {
				{
					VulnID: "CVE-2020-8161",
					Data:   []byte(`{"PatchedVersions":["~\u003e 2.3.3","\u003e= 2.4.0"]}`),
				},
			},
		}
		trivyDBMock.VulnerabilityMap["CVE-2020-8161"] = &types.Vulnerability{
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

		uc := usecase.New()
		require.NoError(t, uc.ScanRepository(svc, req))
		rackPkgs1, err := dbClient.FindPackageRecordsByName(model.PkgRubyGems, "rack")
		require.NoError(t, err)
		require.Equal(t, 1, len(rackPkgs1))
		assert.Equal(t, 1, len(rackPkgs1[0].Package.Vulnerabilities))
		assert.Contains(t, rackPkgs1[0].Package.Vulnerabilities, "CVE-2020-8161")

		// Add new one
		trivyDBMock.AdvisoryMap["ruby-advisory-db"] = map[string][]*model.AdvisoryData{
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
		trivyDBMock.VulnerabilityMap["CVE-2020-9999"] = &types.Vulnerability{
			Title: "test vuln 2",
		}

		svc.Utils.TimeNow = func() time.Time { return now.Add(time.Second) }

		// Scan again
		require.NoError(t, uc.ScanRepository(svc, req))
		rackPkgs2, err := dbClient.FindPackageRecordsByName(model.PkgRubyGems, "rack")
		require.NoError(t, err)
		require.Equal(t, 1, len(rackPkgs2))
		assert.Equal(t, 2, len(rackPkgs2[0].Package.Vulnerabilities))
		assert.Contains(t, rackPkgs2[0].Package.Vulnerabilities, "CVE-2020-8161")
		assert.Contains(t, rackPkgs2[0].Package.Vulnerabilities, "CVE-2020-9999")

		vulns, err := dbClient.FindLatestVulnerabilities(10)
		require.NoError(t, err)
		require.Equal(t, 2, len(vulns))
		assert.Contains(t, []string{vulns[0].VulnID, vulns[1].VulnID}, "CVE-2020-8161")
		assert.Contains(t, []string{vulns[0].VulnID, vulns[1].VulnID}, "CVE-2020-9999")
	})
}

func TestScanBundler(t *testing.T) {
	svc, dbClient, trivyDBMock := setupScanRepositoryService(t, "../testdata/src/bundler.zip")

	trivyDBMock.AdvisoryMap["ruby-advisory-db"] = map[string][]*model.AdvisoryData{
		"rack": {
			{
				VulnID: "CVE-2020-8161",
				// patched version is modified for test
				Data: []byte(`{"PatchedVersions":["~\u003e 2.3.3","\u003e= 2.4.0"]}`),
			},
		},
	}
	trivyDBMock.VulnerabilityMap["CVE-2020-8161"] = &types.Vulnerability{
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

	uc := usecase.New()
	require.NoError(t, uc.ScanRepository(svc, req))

	packages, err := dbClient.FindPackageRecordsByBranch(&req.GitHubBranch)
	require.NoError(t, err)
	assert.Equal(t, 41, len(packages))

	vulns, err := dbClient.FindLatestVulnerabilities(10)
	require.NoError(t, err)
	require.Equal(t, 1, len(vulns))
	assert.Equal(t, "CVE-2020-8161", vulns[0].VulnID)

	rackPkgs, err := dbClient.FindPackageRecordsByName(model.PkgRubyGems, "rack")
	require.NoError(t, err)
	require.Equal(t, 1, len(rackPkgs))
	assert.Equal(t, "rack", rackPkgs[0].Package.Name)
	assert.Contains(t, rackPkgs[0].Package.Vulnerabilities, "CVE-2020-8161")
}

func TestScanGoModule(t *testing.T) {
	svc, dbClient, _ := setupScanRepositoryService(t, "../testdata/src/go_mod.zip")

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

	uc := usecase.New()
	require.NoError(t, uc.ScanRepository(svc, req))

	packages, err := dbClient.FindPackageRecordsByBranch(&req.GitHubBranch)
	require.NoError(t, err)
	assert.Equal(t, 147, len(packages))
}

func TestScanNPM(t *testing.T) {
	svc, dbClient, _ := setupScanRepositoryService(t, "../testdata/src/npm.zip")

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

	uc := usecase.New()
	require.NoError(t, uc.ScanRepository(svc, req))

	packages, err := dbClient.FindPackageRecordsByBranch(&req.GitHubBranch)
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

type httpMock struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (x *httpMock) RoundTrip(req *http.Request) (*http.Response, error) {
	return x.handler(req)
}

func newHTTPMockFactory(handler func(req *http.Request) (*http.Response, error)) infra.NewHTTPClient {
	return func(rt http.RoundTripper) *http.Client {
		return &http.Client{Transport: &httpMock{
			handler: handler,
		}}
	}
}

func newResp(code int, body interface{}) *http.Response {
	var rc io.ReadCloser
	if r, ok := body.(io.ReadCloser); ok {
		rc = r
	} else {
		msg, ok := body.([]byte)
		if !ok {
			msg, _ = json.Marshal(body)
		}
		rc = ioutil.NopCloser(bytes.NewBuffer(msg))
	}

	return &http.Response{
		StatusCode: code,
		Body:       rc,
	}
}

func setupScanRepositoryService(t *testing.T, scannedArchivePath string) (*service.Service, infra.DBClient, *trivydb.TrivyDBMock) {
	pem := genRSAKey(t)
	base64PEM := base64.StdEncoding.EncodeToString(pem)
	const (
		secretsARN       = "arn:aws:secretsmanager:us-east-0:123456789012:secret:tutorials/MyFirstSecret-jiObOV"
		installID  int64 = 234
	)

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

	cfg := service.NewConfig()
	cfg.SecretsARN = secretsARN
	cfg.GitHubEndpoint = "https://ghe.example.org/api/v3"
	cfg.TableName = dbClient.TableName()
	cfg.S3Region = "ap-northeast-0"
	cfg.S3Bucket = "my-db-bucket"
	cfg.S3Prefix = "test-prefix/"

	// Build service and injects mocks
	svc := service.New(cfg)
	svc.NewSecretManager = newSM
	svc.NewDB = func(region, tableName string) (infra.DBClient, error) {
		return dbClient, nil
	}

	// Build trivy mock
	newTrivyDBMock, trivyDBMock := trivydb.NewMock()
	svc.NewTrivyDB = newTrivyDBMock

	// Setup S3 mock
	newS3Mock, s3Mock := aws.NewMockS3()
	svc.NewS3 = newS3Mock
	s3Mock.Objects["my-db-bucket"] = map[string][]byte{
		"test-prefix/db/trivy.db.gz": []byte("boom!"),
	}

	// Setup trivy DB

	var calledGetArchiveLink, calledDownloadArchive int
	svc.NewHTTP = newHTTPMockFactory(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "ghe.example.org", req.URL.Host)
		switch req.URL.Path {
		case "/api/v3/repos/five/blue/zipball/beefcafe":
			calledGetArchiveLink++
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, "application/vnd.github.v3+json", req.Header.Get("Accept"))

			// In test, used stub and ghinstallation is not working
			assert.False(t, strings.HasPrefix(req.Header.Get("Authorization"), "token "))
			hdr := http.Header{}
			hdr.Set("Location", "https://ghe.example.org/_codeload/five/blue/legacy.zip/master?token=hogehoge")
			return &http.Response{
				StatusCode: http.StatusFound,
				Header:     hdr,
				Body:       ioutil.NopCloser(bytes.NewBuffer(nil)),
			}, nil

		case "/_codeload/five/blue/legacy.zip/master":
			calledDownloadArchive++
			assert.Equal(t, "GET", req.Method)
			assert.False(t, strings.HasPrefix(req.Header.Get("Authorization"), "token "))

			r, err := os.Open(scannedArchivePath)
			require.NoError(t, err)
			return newResp(http.StatusOK, r), nil

		default:
			require.Fail(t, "Invalid path", req.URL.Path)
			return &http.Response{}, nil
		}
	})

	t.Cleanup(func() {
		assert.Equal(t, "ap-northeast-0", s3Mock.Region)
		assert.Less(t, 0, calledDownloadArchive)
		assert.Less(t, 0, calledGetArchiveLink)
		require.Less(t, 0, len(s3Mock.GetInput))
		assert.Equal(t, "test-prefix/db/trivy.db.gz", *s3Mock.GetInput[0].Key)
	})

	return svc, dbClient, trivyDBMock
}
