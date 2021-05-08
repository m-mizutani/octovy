package trivydb

import (
	"encoding/json"
	"time"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	bolt "go.etcd.io/bbolt"
)

type TrivyDB struct {
	db *bolt.DB
}

const (
	vulnerabilityBucket = "vulnerability"
	trivyBucketName     = "trivy"
)

func New(dbPath string) (infra.TrivyDBClient, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		ReadOnly: true,
	})
	if err != nil {
		return nil, err
	}

	// TODO: Validate trivy DB

	return &TrivyDB{db: db}, nil
}

func (x *TrivyDB) GetAdvisories(source string, pkgName string) ([]*model.AdvisoryData, error) {
	var advisories []*model.AdvisoryData
	view := func(t *bolt.Tx) error {
		pkgBucket := t.Bucket([]byte(source))
		if pkgBucket == nil {
			return goerr.New("Invalid package source name for trivy DB").With("source", source)
		}

		// find nested bucket
		advBucket := pkgBucket.Bucket([]byte(pkgName))
		if advBucket == nil {
			return nil
		}

		err := advBucket.ForEach(func(k, v []byte) error {
			advisories = append(advisories, &model.AdvisoryData{
				VulnID: string(k),
				Data:   v,
			})
			return nil
		})
		if err != nil {
			return goerr.Wrap(err)
		}

		return nil
	}
	if err := x.db.View(view); err != nil {
		return nil, goerr.Wrap(err)
	}

	return advisories, nil
}

func (x *TrivyDB) GetVulnerability(vulnID string) (*types.Vulnerability, error) {
	var vuln *types.Vulnerability
	view := func(t *bolt.Tx) error {
		vulnBucket := t.Bucket([]byte(vulnerabilityBucket))
		if vulnBucket == nil {
			return goerr.New("Invalid trivy DB, no vulnerability bucket")
		}

		raw := vulnBucket.Get([]byte(vulnID))
		if raw == nil {
			return nil // VulnID is not found
		}

		if err := json.Unmarshal(raw, &vuln); err != nil {
			return goerr.Wrap(err, "Failed to unmarshal vulnerability in trivy DB").With("data", string(raw))
		}

		return nil
	}
	if err := x.db.View(view); err != nil {
		return nil, goerr.Wrap(err)
	}

	return vuln, nil
}

type metadata struct {
	DownloadedAt string
	NextUpdate   string
	Type         int64
	UpdatedAt    string
	Version      int64
}

func (x *TrivyDB) GetDBMeta() (*model.TrivyDBMeta, error) {
	var meta *model.TrivyDBMeta
	view := func(t *bolt.Tx) error {
		trivyBucket := t.Bucket([]byte(trivyBucketName))
		if trivyBucket == nil {
			return goerr.New("trivy bucket is not found in trivy DB")
		}

		metadataBucket := trivyBucket.Bucket([]byte("metadata"))
		if metadataBucket == nil {
			return goerr.New("metadata bucket is not found in trivy DB")
		}

		data := metadataBucket.Get([]byte("data"))
		if data == nil {
			return goerr.New("metadata is not found in trivy DB")
		}

		var v metadata
		if err := json.Unmarshal(data, &v); err != nil {
			return goerr.Wrap(err, "Failed to unmarshal trivy DB metadata").With("data", string(data))
		}

		ts, err := time.Parse("2006-01-02T15:04:05.00000000Z", v.UpdatedAt)
		if err != nil {
			return goerr.Wrap(err, "Can not parse updatedAt in trivy DB metadata").With("meta", v)
		}

		meta = &model.TrivyDBMeta{
			Version:   int(v.Version),
			Type:      int(v.Type),
			UpdatedAt: ts.Unix(),
		}

		return nil
	}
	if err := x.db.View(view); err != nil {
		return nil, goerr.Wrap(err)
	}

	return meta, nil
}

type TrivyDBMock struct {
	DBPath           string
	AdvisoryMap      map[string]map[string][]*model.AdvisoryData
	VulnerabilityMap map[string]*types.Vulnerability
	DBMeta           *model.TrivyDBMeta
}

func NewMock() (infra.NewTrivyDB, *TrivyDBMock) {
	mock := &TrivyDBMock{
		AdvisoryMap: map[string]map[string][]*model.AdvisoryData{
			"GitHub Security Advisory Rubygems": make(map[string][]*model.AdvisoryData),
			"GitHub Security Advisory Npm":      make(map[string][]*model.AdvisoryData),
			"GitHub Security Advisory Pip":      make(map[string][]*model.AdvisoryData),
			"go::GitLab Advisory Database":      make(map[string][]*model.AdvisoryData),
			"nodejs-security-wg":                make(map[string][]*model.AdvisoryData),
			"python-safety-db":                  make(map[string][]*model.AdvisoryData),
			"ruby-advisory-db":                  make(map[string][]*model.AdvisoryData),
		},
		VulnerabilityMap: make(map[string]*types.Vulnerability),
	}

	return func(dbPath string) (infra.TrivyDBClient, error) {
		mock.DBPath = dbPath
		return mock, nil
	}, mock
}

func (x *TrivyDBMock) GetAdvisories(source string, pkgName string) ([]*model.AdvisoryData, error) {
	pkgBucket, ok := x.AdvisoryMap[source]
	if !ok {
		return nil, goerr.New("Invalid package source name for trivy DB").With("source", source)
	}

	return pkgBucket[pkgName], nil
}

func (x *TrivyDBMock) GetVulnerability(vulnID string) (*types.Vulnerability, error) {
	return x.VulnerabilityMap[vulnID], nil
}

func (x *TrivyDBMock) GetDBMeta() (*model.TrivyDBMeta, error) {
	if x.DBMeta != nil {
		return x.DBMeta, nil
	} else {
		return &model.TrivyDBMeta{
			Version:   1,
			Type:      1,
			UpdatedAt: 12345,
		}, nil
	}
}
