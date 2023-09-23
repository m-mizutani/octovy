package usecase

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/sqlc-dev/pqtype"

	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	ttypes "github.com/aquasecurity/trivy/pkg/types"
)

func saveScanReportGitHubRepo(ctx *model.Context, dbClient *sql.DB, report *ttypes.Report, meta *GitHubRepoMetadata) error {
	for _, result := range report.Results {
		if err := saveVulnerabilities(ctx, dbClient, result.Vulnerabilities); err != nil {
			return err
		}

		if err := savePackages(ctx, dbClient, result.Type, result.Packages); err != nil {
			return err
		}
	}

	if err := saveScanGitHubRepo(ctx, dbClient, report, meta); err != nil {
		return err
	}

	return nil
}

func saveScanGitHubRepo(ctx *model.Context, dbClient *sql.DB, report *ttypes.Report, meta *GitHubRepoMetadata) error {
	tx, err := dbClient.Begin()
	if err != nil {
		return goerr.Wrap(err, "starting transaction")
	}
	defer utils.SafeRollback(tx)
	q := db.New(tx)

	scanID := uuid.New()
	if err := q.SaveScan(ctx, db.SaveScanParams{
		ID:           scanID,
		ArtifactName: report.ArtifactName,
		ArtifactType: string(report.ArtifactType),
	}); err != nil {
		return goerr.Wrap(err, "saving scan")
	}

	if err := q.SaveMetaGithubRepository(ctx, db.SaveMetaGithubRepositoryParams{
		ID:       uuid.New(),
		ScanID:   scanID,
		Owner:    meta.Owner,
		RepoName: meta.Repo,
		CommitID: meta.CommitID,
		Branch: sql.NullString{
			String: meta.Branch,
			Valid:  meta.Branch != "",
		},
		IsDefaultBranch: sql.NullBool{
			Bool:  meta.IsDefaultBranch,
			Valid: meta.Branch != "",
		},
		BaseCommitID: sql.NullString{
			String: meta.BaseCommitID,
			Valid:  meta.BaseCommitID != "",
		},
		PullRequestID: sql.NullInt32{
			Int32: int32(meta.PullRequestID),
			Valid: meta.PullRequestID != 0,
		},
	}); err != nil {
		return goerr.Wrap(err, "saving meta github repository")
	}

	for _, result := range report.Results {
		// @TODO: Support other types
		if result.Class != ttypes.ClassLangPkg && result.Class != ttypes.ClassOSPkg {
			continue
		}

		resultID := uuid.New()
		if err := q.SaveResult(ctx, db.SaveResultParams{
			ID:         resultID,
			ScanID:     scanID,
			Target:     result.Target,
			TargetType: string(result.Type),
			Class:      db.TargetClass(result.Class),
		}); err != nil {
			return goerr.Wrap(err, "saving result")
		}

		for _, vuln := range result.Vulnerabilities {
			if err := q.SaveResultVulnerability(ctx, db.SaveResultVulnerabilityParams{
				ID:       uuid.New(),
				ResultID: resultID,
				VulnID:   vuln.VulnerabilityID,
				PkgID:    calcPackageID(result.Type, vuln.PkgName, vuln.InstalledVersion),
				FixedVersion: sql.NullString{
					String: vuln.FixedVersion,
					Valid:  vuln.FixedVersion != "",
				},
				PrimaryUrl: sql.NullString{
					String: vuln.PrimaryURL,
					Valid:  vuln.PrimaryURL != "",
				},
			}); err != nil {
				return goerr.Wrap(err, "saving result vulnerability")
			}
		}

		for _, pkg := range result.Packages {
			if err := q.SaveResultPackage(ctx, db.SaveResultPackageParams{
				ID:       uuid.New(),
				ResultID: resultID,
				PkgID:    calcPackageID(result.Type, pkg.Name, pkg.Version),
			}); err != nil {
				return goerr.Wrap(err, "saving result package")
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return goerr.Wrap(err, "commit transaction")
	}

	return nil
}

// calcPackageID returns a unique ID for a package. It is calculated by sha512 from package name, version and type.
func calcPackageID(typ, name, version string) string {
	src := bytes.Join([][]byte{
		[]byte(typ),
		[]byte(name),
		[]byte(version),
	}, []byte{0})

	hash := sha256.New()
	hash.Write(src)
	hv := hash.Sum(nil)
	return hex.EncodeToString(hv)
}

func savePackages(ctx *model.Context, dbClient *sql.DB, typ string, packages []ftypes.Package) error {
	tx, err := dbClient.Begin()
	if err != nil {
		return goerr.Wrap(err, "starting transaction")
	}
	defer utils.SafeRollback(tx)
	q := db.New(tx)

	pkgSet := map[string]*ftypes.Package{}
	pkgIDs := []string{}
	for i, pkg := range packages {
		pkgID := calcPackageID(typ, pkg.Name, pkg.Version)
		pkgSet[pkgID] = &packages[i]
		pkgIDs = append(pkgIDs, calcPackageID(typ, pkg.Name, pkg.Version))
	}

	exists, err := q.GetPackages(ctx, pkgIDs)
	if err != nil {
		return goerr.Wrap(err, "getting packages")
	}

	for _, pkg := range exists {
		delete(pkgSet, pkg.ID)
	}

	for pkgID, pkg := range pkgSet {
		if err := q.SavePackage(ctx, db.SavePackageParams{
			ID:         pkgID,
			TargetType: typ,
			Name:       pkg.Name,
			Version:    pkg.Version,
		}); err != nil {
			return goerr.Wrap(err, "saving package").With("pkgID", pkgID).With("pkg", pkg)
		}
	}

	if err := tx.Commit(); err != nil {
		return goerr.Wrap(err, "commit transaction")
	}

	return nil
}

func saveVulnerabilities(ctx *model.Context, dbClient *sql.DB, vulns []ttypes.DetectedVulnerability) error {
	tx, err := dbClient.Begin()
	if err != nil {
		return goerr.Wrap(err, "starting transaction")
	}
	defer utils.SafeRollback(tx)
	q := db.New(tx)

	for _, vuln := range vulns {
		data, err := json.Marshal(vuln.Vulnerability)
		if err != nil {
			return goerr.Wrap(err, "marshaling vulnerability").With("vuln", vuln)
		}

		if err := q.SaveVulnerability(ctx, db.SaveVulnerabilityParams{
			ID:       vuln.VulnerabilityID,
			Title:    vuln.Title,
			Severity: vuln.Severity,
			PublishedAt: sql.NullTime{
				Valid: vuln.PublishedDate != nil,
				Time:  *vuln.PublishedDate,
			},
			LastModifiedAt: sql.NullTime{
				Valid: vuln.LastModifiedDate != nil,
				Time:  *vuln.LastModifiedDate,
			},
			Data: pqtype.NullRawMessage{Valid: true, RawMessage: data},
		}); err != nil {
			return goerr.Wrap(err, "saving vulnerability").With("vuln", vuln)
		}
	}

	if err := tx.Commit(); err != nil {
		return goerr.Wrap(err, "commit transaction")
	}

	return nil
}
