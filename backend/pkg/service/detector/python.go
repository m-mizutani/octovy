package detector

import (
	pipver "github.com/aquasecurity/go-pep440-version"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/ghsa"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/python"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func matchPython(constraints, version string) (bool, error) {
	c, err := pipver.NewSpecifiers(constraints)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid patched version").With("ver", constraints)
	}
	v, err := pipver.Parse(version)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid version").With("ver", version)
	}

	return c.Check(v), nil
}

func isVulnerablePython(data *model.AdvisoryData, version string) (bool, error) {
	var adv python.Advisory
	if err := data.Unmarshal(&adv); err != nil {
		return false, err
	}

	for _, spec := range adv.Specs {
		if matched, err := matchPython(spec, version); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}

	return false, nil
}

func isVulnerablePythonGHSA(data *model.AdvisoryData, version string) (bool, error) {
	var adv ghsa.Advisory
	if err := data.Unmarshal(&adv); err != nil {
		return false, err
	}

	for _, vulnerable := range adv.VulnerableVersions {
		if matched, err := matchSemVer(vulnerable, version); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}

	return false, nil
}
