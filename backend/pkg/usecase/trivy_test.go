package usecase_test

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type githubMock struct {
	ListReleasesMock         func(owner string, repo string) ([]*github.RepositoryRelease, error)
	DownloadReleaseAssetMock func(owner string, repo string, assetID int64) (io.ReadCloser, error)
}

func (x *githubMock) ListReleases(owner string, repo string) ([]*github.RepositoryRelease, error) {
	return x.ListReleasesMock(owner, repo)
}

func (x *githubMock) DownloadReleaseAsset(owner string, repo string, assetID int64) (io.ReadCloser, error) {
	return x.DownloadReleaseAssetMock(owner, repo, assetID)
}

func TestUpdateTrivyDB(t *testing.T) {
	svc := service.New(&service.Config{
		S3Region: "ap-northeast-0",
		S3Bucket: "blue-bucket",
		S3Prefix: "five/",
	})

	newS3Mock, s3Mock := aws.NewMockS3()
	svc.NewS3 = newS3Mock

	calledListReleasesMock, calledDownloadReleaseAssetMock := 0, 0

	ghMock := &githubMock{
		ListReleasesMock: func(owner, repo string) ([]*github.RepositoryRelease, error) {
			calledListReleasesMock++
			assert.Equal(t, "aquasecurity", owner)
			assert.Equal(t, "trivy-db", repo)

			return []*github.RepositoryRelease{
				{
					Name: github.String("v1-20000000"),
					Assets: []github.ReleaseAsset{
						{
							Name: github.String("other.db.gz"),
							ID:   github.Int64(2345),
						},
						{
							Name: github.String("trivy.db.gz"),
							ID:   github.Int64(3456),
						},
					},
				},
			}, nil
		},

		DownloadReleaseAssetMock: func(owner, repo string, assetID int64) (io.ReadCloser, error) {
			calledDownloadReleaseAssetMock++
			assert.Equal(t, "aquasecurity", owner)
			assert.Equal(t, "trivy-db", repo)
			assert.Equal(t, int64(3456), assetID)

			return io.NopCloser(strings.NewReader("boom!")), nil
		},
	}

	svc.NewGitHub = func() infra.GitHubClient { return ghMock }

	err := usecase.UpdateTrivyDB(svc)
	require.NoError(t, err)

	assert.Equal(t, 1, calledListReleasesMock)
	assert.Equal(t, 1, calledDownloadReleaseAssetMock)

	bucket, ok := s3Mock.Objects["blue-bucket"]
	require.True(t, ok)
	obj, ok := bucket["five/db/trivy.db.gz"]
	require.True(t, ok)
	assert.NotNil(t, obj)
}
