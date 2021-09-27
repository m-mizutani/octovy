package detector

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

type Detector struct {
	trivyDB trivydb.Interface
}

func New(db trivydb.Interface) *Detector {
	return &Detector{
		trivyDB: db,
	}
}

type isVulnerableFunc func(adv *model.AdvisoryData, version string) (bool, error)

var pkgTypeSourceMap map[types.PkgType]map[string]isVulnerableFunc

func init() {
	pkgTypeSourceMap = make(map[types.PkgType]map[string]isVulnerableFunc)

	pkgTypeSourceMap[types.PkgRubyGems] = map[string]isVulnerableFunc{
		"ruby-advisory-db":                  isVulnerableBundlerRubyAdv,
		"GitHub Security Advisory Rubygems": isVulnerableBundlerGHSA,
	}

	pkgTypeSourceMap[types.PkgGoModule] = map[string]isVulnerableFunc{
		"go::GitLab Advisory Database": isVulnerableGoMod,
	}

	pkgTypeSourceMap[types.PkgNPM] = map[string]isVulnerableFunc{
		"nodejs-security-wg":           isVulnerableNodeSecurityWG,
		"GitHub Security Advisory Npm": isVulnerableNodeGHSA,
	}

	pkgTypeSourceMap[types.PkgPyPI] = map[string]isVulnerableFunc{
		"GitHub Security Advisory Pip": isVulnerablePythonGHSA,
		"python-safety-db":             isVulnerablePython,
	}
}

func (x *Detector) Detect(pkgType types.PkgType, pkgName, version string) ([]*model.Vulnerability, error) {
	options, ok := pkgTypeSourceMap[pkgType]
	if !ok {
		logger.Warn().Interface("pkgType", pkgType).Msg("Unsupported pkgType")
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
				logger.Warn().Err(err).Interface("pkg", pkgName).Interface("adv", adv).Err(err).Send()
				continue
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
