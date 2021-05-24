package detector

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

var logger = golambda.Logger

type Detector struct {
	trivyDB interfaces.TrivyDBClient
}

func New(db interfaces.TrivyDBClient) *Detector {
	return &Detector{
		trivyDB: db,
	}
}

type isVulnerableFunc func(adv *model.AdvisoryData, version string) (bool, error)

var pkgTypeSourceMap map[model.PkgType]map[string]isVulnerableFunc

func init() {
	pkgTypeSourceMap = make(map[model.PkgType]map[string]isVulnerableFunc)

	pkgTypeSourceMap[model.PkgRubyGems] = map[string]isVulnerableFunc{
		"ruby-advisory-db":                  isVulnerableBundlerRubyAdv,
		"GitHub Security Advisory Rubygems": isVulnerableBundlerGHSA,
	}

	pkgTypeSourceMap[model.PkgGoModule] = map[string]isVulnerableFunc{
		"go::GitLab Advisory Database": isVulnerableGoMod,
	}

	pkgTypeSourceMap[model.PkgNPM] = map[string]isVulnerableFunc{
		"nodejs-security-wg":           isVulnerableNodeSecurityWG,
		"GitHub Security Advisory Npm": isVulnerableNodeGHSA,
	}

	pkgTypeSourceMap[model.PkgPyPI] = map[string]isVulnerableFunc{
		"GitHub Security Advisory Pip": isVulnerablePythonGHSA,
		"python-safety-db":             isVulnerablePython,
	}
}

func (x *Detector) Detect(pkgType model.PkgType, pkgName, version string) ([]*model.Vulnerability, error) {
	options, ok := pkgTypeSourceMap[pkgType]
	if !ok {
		logger.With("pkgType", pkgType).Warn("Unsupported pkgType")
		return nil, nil
	}

	var affected []*model.AdvisoryData
	for src, isVulnerable := range options {
		advisories, err := x.trivyDB.GetAdvisories(src, pkgName)
		if err != nil {
			return nil, err
		}

		for _, adv := range advisories {
			if vulnerable, err := isVulnerable(adv, version); err != nil {
				return nil, err
			} else if vulnerable {
				affected = append(affected, adv)
			}
		}
	}

	var vulnerabilities []*model.Vulnerability
	for _, adv := range affected {
		vuln, err := x.trivyDB.GetVulnerability(adv.VulnID)
		if err != nil {
			return nil, err
		}
		vulnerabilities = append(vulnerabilities, &model.Vulnerability{
			Detail: *vuln,
			VulnID: adv.VulnID,
		})
	}

	return vulnerabilities, nil
}

func (x Detector) TrivyDBMeta() (*model.TrivyDBMeta, error) {
	return x.trivyDB.GetDBMeta()
}
