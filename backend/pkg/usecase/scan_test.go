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

	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func genRSAKey(t *testing.T) []byte {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

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

func TestScanRepository(t *testing.T) {
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
		Branches: []string{"master"},
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

	// Build service and injects mocks
	svc := service.New(cfg)
	svc.NewSecretManager = newSM
	svc.NewDB = func(region, tableName string) (infra.DBClient, error) {
		return dbClient, nil
	}

	var calledGetArchiveLink, calledDownloadArchive bool
	svc.NewHTTP = newHTTPMockFactory(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "ghe.example.org", req.URL.Host)
		switch req.URL.Path {
		case "/api/v3/repos/five/blue/zipball/beefcafe":
			calledGetArchiveLink = true
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
			calledDownloadArchive = true
			assert.Equal(t, "GET", req.Method)
			assert.False(t, strings.HasPrefix(req.Header.Get("Authorization"), "token "))

			r, err := os.Open("../testdata/src/bundler.zip")
			require.NoError(t, err)
			return newResp(http.StatusOK, r), nil

		default:
			require.Fail(t, "Invalid path", req.URL.Path)
			return &http.Response{}, nil
		}
	})

	req := &model.ScanRepositoryRequest{
		ScanTarget: model.ScanTarget{
			GitHubBranch: model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "five",
					RepoName: "blue",
				},
				Branch: "master",
			},
			Ref:       "beefcafe",
			UpdatedAt: 1234,
		},
		InstallID: installID,
	}

	uc := usecase.New()
	require.NoError(t, uc.ScanRepository(svc, req))
	assert.True(t, calledGetArchiveLink)
	assert.True(t, calledDownloadArchive)

	packages, err := dbClient.FindPackagesByBranch(&req.GitHubBranch)
	require.NoError(t, err)
	assert.Equal(t, 41, len(packages))
}
