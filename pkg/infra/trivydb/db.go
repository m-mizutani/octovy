package trivydb

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	bolt "go.etcd.io/bbolt"
)

type Interface interface {
	GetAdvisories(source, pkgName string) ([]*model.AdvisoryData, error)
	GetVulnerability(vulnID string) (*types.Vulnerability, error)
	GetDBMeta() (*model.TrivyDBMeta, error)
}

type TrivyDB struct {
	db *bolt.DB
}

const (
	vulnerabilityBucket = "vulnerability"
	trivyBucketName     = "trivy"
)

type Factory func(dbPath string) (Interface, error)

func New(dbPath string) (Interface, error) {
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

		updatedAtParts := strings.Split(v.UpdatedAt, ".")
		ts, err := time.Parse("2006-01-02T15:04:05", updatedAtParts[0])
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
