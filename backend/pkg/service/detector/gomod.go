package detector

import (
	"github.com/aquasecurity/go-version/pkg/semver"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

type gomodAdvisory struct {
	PatchedVersions    []string `json:",omitempty"`
	VulnerableVersions []string `json:",omitempty"`
}

func matchSemVer(constraints, version string) (bool, error) {
	c, err := semver.NewConstraints(constraints)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid patched version").With("ver", constraints)
	}
	v, err := semver.Parse(version)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid version").With("ver", version)
	}

	return c.Check(v), nil
}

func isVulnerableGoMod(data *model.AdvisoryData, version string) (bool, error) {
	var adv gomodAdvisory
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
