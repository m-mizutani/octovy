package usecase

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/service"
)

const (
	trivyDBOwner     = "aquasecurity"
	trivyDBRepo      = "trivy-db"
	trivyDBName      = "trivy.db.gz"
	trivyDBSchemaVer = "v1"
)

func downloadLatestTrivyDB(svc *service.Service, releases []*github.RepositoryRelease) (io.ReadCloser, error) {
	client := svc.Infra.NewGitHub()
	dbNamePrefix := trivyDBSchemaVer + "-"

	for _, release := range releases {
		if !strings.HasPrefix(release.GetName(), dbNamePrefix) {
			continue
		}

		for _, asset := range release.Assets {
			if asset.GetName() != trivyDBName {
				continue
			}

			rc, err := client.DownloadReleaseAsset(trivyDBOwner, trivyDBRepo, asset.GetID())
			if err != nil {
				logger.With("err", err).With("asset", asset).Warn("Failed to download trivy DB")
				continue
			}

			return rc, nil
		}
	}

	return nil, goerr.New("No available trivy-db asset")
}

func (x *Default) UpdateTrivyDB() error {
	client := x.svc.Infra.NewGitHub()

	releases, err := client.ListReleases(trivyDBOwner, trivyDBRepo)
	if err != nil {
		return goerr.Wrap(err, "ListRelease error")
	}

	dbReader, err := downloadLatestTrivyDB(x.svc, releases)
	if err != nil {
		return err
	}
	defer dbReader.Close()

	temp, err := ioutil.TempFile("", "*.gz")
	if err != nil {
		return goerr.Wrap(err)
	}

	if _, err := io.Copy(temp, dbReader); err != nil {
		return goerr.Wrap(err, "Failed to save trivyDB").With("dst", temp.Name())
	}

	if _, err := temp.Seek(0, 0); err != nil {
		return goerr.Wrap(err)
	}

	if err := x.svc.UploadTrivyDB(temp); err != nil {
		return err
	}

	return nil
}
