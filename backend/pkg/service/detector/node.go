package detector

import (
	npm "github.com/aquasecurity/go-npm-version/pkg"

	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/ghsa"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/node"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

func matchNode(constraints, version string) (bool, error) {
	c, err := npm.NewConstraints(constraints)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid patched version").With("ver", constraints)
	}
	v, err := npm.NewVersion(version)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid version").With("ver", version)
	}

	return c.Check(v), nil
}

func isVulnerableNodeSecurityWG(data *model.AdvisoryData, version string) (bool, error) {
	var adv node.Advisory
	if err := data.Unmarshal(&adv); err != nil {
		return false, err
	}

	if adv.VulnerableVersions != "" {
		if matched, err := matchNode(adv.VulnerableVersions, version); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}

	if adv.PatchedVersions != "" {
		if matched, err := matchNode(adv.PatchedVersions, version); err != nil {
			return false, err
		} else if matched {
			return false, nil
		}
	}

	return true, nil
}

func isVulnerableNodeGHSA(data *model.AdvisoryData, version string) (bool, error) {
	var adv ghsa.Advisory
	if err := data.Unmarshal(&adv); err != nil {
		return false, err
	}

	for _, vulnerable := range adv.VulnerableVersions {
		if matched, err := matchNode(vulnerable, version); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}

	return false, nil
}
