package usecase

import (
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"os"
	"strings"

	gh "github.com/google/go-github/v39/github"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/pkg/usecase/detector"
)

const (
	trivyDBOwner     = "aquasecurity"
	trivyDBRepo      = "trivy-db"
	trivyDBName      = "trivy.db.gz"
	trivyDBSchemaVer = "v1"
)

// vulnDetector manages trivy DB and detector
type vulnDetector struct {
	github     github.Interface
	newTrivyDB trivydb.Factory
	dbPath     string

	dt *detector.Detector
}

func newVulnDetector(client github.Interface, newTrivyDB trivydb.Factory, dbPath string) *vulnDetector {
	return &vulnDetector{
		github:     client,
		newTrivyDB: newTrivyDB,
		dbPath:     dbPath,
	}
}

func (x *vulnDetector) Detect(pkgType types.PkgType, pkgName, version string) ([]*model.Vulnerability, error) {
	if x.dt == nil {
		panic("trivy DB for detector is not set")
	}

	return x.dt.Detect(pkgType, pkgName, version)
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func (x *vulnDetector) initDetector() error {
	db, err := x.newTrivyDB(x.dbPath)
	if err != nil {
		return err
	}

	x.dt = detector.New(db)
	return nil
}

func (x *vulnDetector) RefreshDB() error {
	// TODO: improve trivy DB refresh logic
	if isFileExist(x.dbPath) {
		logger.Debug().Str("path", x.dbPath).Msg("Found existing DB file")
		return x.initDetector()
	}

	ctx := context.Background()
	releases, err := x.github.ListReleases(ctx, trivyDBOwner, trivyDBRepo)
	if err != nil {
		return goerr.Wrap(err, "ListRelease error")
	}

	dbReader, err := downloadLatestTrivyDB(ctx, x.github, releases)
	if err != nil {
		return err
	}
	defer dbReader.Close()

	tmp, err := ioutil.TempFile("", "*.db")
	if err != nil {
		return goerr.Wrap(err)
	}

	gz, err := gzip.NewReader(dbReader)
	if err != nil {
		return goerr.Wrap(err)
	}

	if _, err := io.Copy(tmp, gz); err != nil {
		return goerr.Wrap(err, "Failed to save trivyDB").With("dst", tmp.Name())
	}
	if err := tmp.Close(); err != nil {
		return goerr.Wrap(err).With("path", tmp.Name())
	}
	logger.Debug().Str("path", tmp.Name()).Msg("Saved trivy DB file")

	if x.dbPath == "" {
		x.dbPath = tmp.Name()
	} else {
		if err := os.Rename(tmp.Name(), x.dbPath); err != nil {
			return goerr.Wrap(err)
		}
		logger.Debug().Str("old", tmp.Name()).Str("new", x.dbPath).Msg("Rename trivy DB file")
	}

	return x.initDetector()
}

func downloadLatestTrivyDB(ctx context.Context, client github.Interface, releases []*gh.RepositoryRelease) (io.ReadCloser, error) {
	dbNamePrefix := trivyDBSchemaVer + "-"

	for _, release := range releases {
		if !strings.HasPrefix(release.GetName(), dbNamePrefix) {
			continue
		}

		for _, asset := range release.Assets {
			if asset.GetName() != trivyDBName {
				continue
			}

			rc, err := client.DownloadReleaseAsset(ctx, trivyDBOwner, trivyDBRepo, asset.GetID())
			if err != nil {
				logger.Warn().Err(err).Interface("asset", asset).Msg("Failed to download trivy DB")
				continue
			}

			logger.Debug().Interface("asset", asset).Msg("Downloading trivy DB file")

			return rc, nil
		}
	}

	return nil, goerr.New("No available trivy-db asset")
}
