package trivydb

import (
	"encoding/json"

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

type TrivyDBMock struct {
	DBPath           string
	AdvisoryMap      map[string]map[string][]*model.AdvisoryData
	VulnerabilityMap map[string]*types.Vulnerability
}

func NewMock() (infra.NewTrivyDB, *TrivyDBMock) {
	mock := &TrivyDBMock{
		AdvisoryMap:      make(map[string]map[string][]*model.AdvisoryData),
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
