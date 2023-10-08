package usecase_test

import (
	"encoding/json"
	"testing"
	"time"

	ptypes "github.com/aquasecurity/trivy-db/pkg/types"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	ttypes "github.com/aquasecurity/trivy/pkg/types"
	"github.com/google/uuid"
	"github.com/m-mizutani/gots/ptr"
	"github.com/m-mizutani/gots/rands"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/usecase"
)

func deepCopy[T any](t *testing.T, src T) T {
	var dst T
	data := gt.R1(json.Marshal(src)).NoError(t)
	gt.NoError(t, json.Unmarshal(data, &dst))
	return dst
}

func TestGetVulnDiffForGitHubRepo(t *testing.T) {
	dbClient := newTestDB(t)
	ctx := model.NewContext()
	meta := &usecase.GitHubRepoMetadata{
		GitHubCommit: usecase.GitHubCommit{
			GitHubRepo: usecase.GitHubRepo{
				Owner: "m-mizutani",
				Repo:  "octovy",
			},
			CommitID: uuid.NewString(),
		},
	}

	now := time.Now()
	salt := rands.AlphaNum(10)
	baseReport := ttypes.Report{
		ArtifactName: "github.com/m-mizutani/octovy",
		ArtifactType: ttypes.ClassLangPkg,
		Results: ttypes.Results{
			{
				Target: "Gemfile.lock",
				Class:  ttypes.ClassLangPkg,
				Type:   "bundler",
				Packages: []ftypes.Package{
					{
						Name:    "octokit_" + salt,
						Version: "4.18.0",
					},
				},
				Vulnerabilities: []ttypes.DetectedVulnerability{
					{
						VulnerabilityID:  "CVE-2020-1234-" + salt,
						PkgName:          "octokit_" + salt,
						InstalledVersion: "4.18.0",
						Vulnerability: ptypes.Vulnerability{
							Title:            "CVE-2020-1234",
							Description:      "test",
							Severity:         "MIDDLE",
							References:       []string{"https://example.com"},
							PublishedDate:    ptr.To(now),
							LastModifiedDate: ptr.To(now),
						},
					},
				},
			},
		},
	}

	gt.NoError(t, usecase.SaveScanReportGitHubRepo(ctx, dbClient, &baseReport, meta))

	t.Run("add vuln", func(t *testing.T) {
		report := deepCopy(t, baseReport)
		report.Results[0].Vulnerabilities = append(report.Results[0].Vulnerabilities, ttypes.DetectedVulnerability{
			VulnerabilityID:  "CVE-2020-5678-" + salt,
			PkgName:          "octokit_" + salt,
			InstalledVersion: "4.18.0",
			Vulnerability: ptypes.Vulnerability{
				Title:            "CVE-2020-5678",
				Description:      "test",
				Severity:         "MIDDLE",
				References:       []string{"https://example.com"},
				PublishedDate:    ptr.To(now),
				LastModifiedDate: ptr.To(now),
			},
		})

		report.Results[0].Packages = append(report.Results[0].Packages, ftypes.Package{
			Name:    "octokit_" + salt,
			Version: "4.18.0",
		})

		diffResults, err := usecase.GetVulnDiffForGitHubRepo(ctx, dbClient, &report, &meta.GitHubCommit)
		gt.NoError(t, err)
		gt.A(t, diffResults).Must().Length(1).At(0, func(t testing.TB, v usecase.DiffResult) {
			gt.Equal(t, v.Status, usecase.DiffStatusMod)
			gt.A(t, v.Add).Must().Length(1).At(0, func(t testing.TB, v ttypes.DetectedVulnerability) {
				gt.Equal(t, "CVE-2020-5678-"+salt, v.VulnerabilityID)
			})
			gt.A(t, v.Del).Must().Length(0)
		})
	})

	t.Run("del vuln", func(t *testing.T) {
		report := deepCopy(t, baseReport)
		report.Results[0].Vulnerabilities = []ttypes.DetectedVulnerability{}

		diffResults, err := usecase.GetVulnDiffForGitHubRepo(ctx, dbClient, &report, &meta.GitHubCommit)
		gt.NoError(t, err)
		gt.A(t, diffResults).Must().Length(1).At(0, func(t testing.TB, v usecase.DiffResult) {
			gt.Equal(t, v.Status, usecase.DiffStatusMod)
			gt.A(t, v.Add).Must().Length(0)
			gt.A(t, v.Del).Must().Length(1).At(0, func(t testing.TB, v ttypes.DetectedVulnerability) {
				gt.Equal(t, "CVE-2020-1234-"+salt, v.VulnerabilityID)
			})
		})
	})

	t.Run("add result", func(t *testing.T) {
		report := deepCopy(t, baseReport)
		report.Results = append(report.Results, ttypes.Result{
			Target: "go.mod",
			Class:  ttypes.ClassLangPkg,
			Type:   "go",
			Vulnerabilities: []ttypes.DetectedVulnerability{
				{
					VulnerabilityID:  "CVE-2020-3456-" + salt,
					PkgName:          "github.com/m-mizutani/octovy",
					InstalledVersion: "v1.0.0",
					Vulnerability: ptypes.Vulnerability{
						Title:       "CVE-2020-3456",
						Description: "test",
						Severity:    "MIDDLE",
					},
				},
			},
		})

		diffResults, err := usecase.GetVulnDiffForGitHubRepo(ctx, dbClient, &report, &meta.GitHubCommit)
		gt.NoError(t, err)
		gt.A(t, diffResults).Must().Length(1).At(0, func(t testing.TB, v usecase.DiffResult) {
			gt.Equal(t, v.Status, usecase.DiffStatusAdd)
			gt.A(t, v.Add).Must().Length(1).At(0, func(t testing.TB, v ttypes.DetectedVulnerability) {
				gt.Equal(t, "CVE-2020-3456-"+salt, v.VulnerabilityID)
			})
			gt.A(t, v.Del).Must().Length(0)
		})
	})

	t.Run("del result", func(t *testing.T) {
		report := deepCopy(t, baseReport)
		report.Results = []ttypes.Result{}

		diffResults, err := usecase.GetVulnDiffForGitHubRepo(ctx, dbClient, &report, &meta.GitHubCommit)
		gt.NoError(t, err)
		gt.A(t, diffResults).Must().Length(1).At(0, func(t testing.TB, v usecase.DiffResult) {
			gt.Equal(t, v.Status, usecase.DiffStatusDel)
			gt.A(t, v.Add).Must().Length(0)
			gt.A(t, v.Del).Must().Length(1).At(0, func(t testing.TB, v ttypes.DetectedVulnerability) {
				gt.Equal(t, "CVE-2020-1234-"+salt, v.VulnerabilityID)
			})
		})
	})
}
