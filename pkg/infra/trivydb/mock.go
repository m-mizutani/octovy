package trivydb

import (
	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type Mock struct {
	DBPath           string
	AdvisoryMap      map[string]map[string][]*model.AdvisoryData
	VulnerabilityMap map[string]*types.Vulnerability
	DBMeta           *model.TrivyDBMeta
}

func NewMock() (Factory, *Mock) {
	mock := &Mock{
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

	return func(dbPath string) (Interface, error) {
		mock.DBPath = dbPath
		return mock, nil
	}, mock
}

func (x *Mock) GetAdvisories(source string, pkgName string) ([]*model.AdvisoryData, error) {
	pkgBucket, ok := x.AdvisoryMap[source]
	if !ok {
		return nil, goerr.New("Invalid package source name for trivy DB").With("source", source)
	}

	return pkgBucket[pkgName], nil
}

func (x *Mock) GetVulnerability(vulnID string) (*types.Vulnerability, error) {
	return x.VulnerabilityMap[vulnID], nil
}

func (x *Mock) GetDBMeta() (*model.TrivyDBMeta, error) {
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
