package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackage(t *testing.T) {
	t.Run("Insert and find package record", func(t *testing.T) {
		client := newTestTable(t)

		pkgs := []*model.PackageRecord{
			{
				Detected: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "five",
							RepoName: "blue",
						},
						Branch: "ao",
					},
					CommitID: "aaaaaaaa",
				},
				Source: "go.mod",
				Package: model.Package{
					Type:            model.PkgGoModule,
					Name:            "timeless",
					Version:         "1.2.3",
					Vulnerabilities: []string{},
				},
				ScannedAt: 1234,
			},
			{
				Detected: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "five",
							RepoName: "blue",
						},
						Branch: "ao",
					},
					CommitID: "aaaaaaaa",
				},
				Source: "go.mod",
				Package: model.Package{
					Type:            model.PkgGoModule,
					Name:            "words",
					Version:         "5.0.1",
					Vulnerabilities: []string{},
				},
				ScannedAt: 1234,
			},
			{
				Detected: model.ScanTarget{
					GitHubBranch: model.GitHubBranch{
						GitHubRepo: model.GitHubRepo{
							Owner:    "three",
							RepoName: "heaven",
						},
						Branch: "feel",
					},
					CommitID: "aaaaaaaa",
				},
				Source: "go.mod",
				Package: model.Package{
					Type:            model.PkgGoModule,
					Name:            "timeless",
					Version:         "1.2.4",
					Vulnerabilities: []string{},
				},
				ScannedAt: 1234,
			},
		}

		for _, pkg := range pkgs {
			inserted, err := client.InsertPackageRecord(pkg)
			require.True(t, inserted)
			require.NoError(t, err)
		}

		t.Run("Lookup branch", func(t *testing.T) {
			r, err := client.FindPackageRecordsByBranch(&model.GitHubBranch{
				GitHubRepo: model.GitHubRepo{
					Owner:    "five",
					RepoName: "blue",
				},
				Branch: "ao",
			})
			require.NoError(t, err)
			assert.Contains(t, r, pkgs[0])
			assert.Contains(t, r, pkgs[1])
			assert.NotContains(t, r, pkgs[2])
		})

		t.Run("Lookup package name", func(t *testing.T) {
			r, err := client.FindPackageRecordsByName(model.PkgGoModule, "timeless")
			require.NoError(t, err)
			assert.Contains(t, r, pkgs[0])
			assert.NotContains(t, r, pkgs[1])
			assert.Contains(t, r, pkgs[2])
		})
	})
}
