package detector

import (
	"github.com/aquasecurity/go-gem-version"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/bundler"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/ghsa"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func matchBundler(constraints, version string) (bool, error) {
	c, err := gem.NewConstraints(constraints)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid patched version").With("ver", constraints)
	}
	v, err := gem.NewVersion(version)
	if err != nil {
		return false, goerr.Wrap(err, "Invalid version").With("ver", version)
	}

	return c.Check(v), nil
}

func isVulnerableBundlerRubyAdv(data *model.AdvisoryData, version string) (bool, error) {
	var adv bundler.Advisory
	if err := data.Unmarshal(&adv); err != nil {
		return false, err
	}

	nonVulnVersions := append(adv.PatchedVersions, adv.UnaffectedVersions...)
	for _, nonVulns := range nonVulnVersions {
		if matched, err := matchBundler(nonVulns, version); err != nil {
			return false, err
		} else if matched {
			return false, nil
		}
	}

	return true, nil
}

func isVulnerableBundlerGHSA(data *model.AdvisoryData, version string) (bool, error) {
	var adv ghsa.Advisory
	if err := data.Unmarshal(&adv); err != nil {
		return false, err
	}

	for _, vulnerable := range adv.VulnerableVersions {
		if matched, err := matchBundler(vulnerable, version); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}

	return false, nil
}
