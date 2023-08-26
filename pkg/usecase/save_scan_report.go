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
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/tabbed/pqtype"

	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	ttypes "github.com/aquasecurity/trivy/pkg/types"
)

func saveScanReport(ctx *model.Context, dbClient *sql.DB, report *ttypes.Report) error {
	for _, result := range report.Results {
		if err := saveVulnerabilities(ctx, dbClient, result.Vulnerabilities); err != nil {
			return err
		}

		if err := savePackages(ctx, dbClient, result.Type, result.Packages); err != nil {
			return err
		}
	}

	if err := saveScan(ctx, dbClient, report); err != nil {
		return err
	}

	return nil
}

func saveScan(ctx *model.Context, dbClient *sql.DB, report *ttypes.Report) error {
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

	for _, result := range report.Results {
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
			return goerr.Wrap(err, "saving package")
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

	vulnSet := map[string]*ttypes.DetectedVulnerability{}
	vulnIDs := []string{}
	for i, vuln := range vulns {
		vulnSet[vuln.VulnerabilityID] = &vulns[i]
		vulnIDs = append(vulnIDs, vuln.VulnerabilityID)
	}

	exists, err := q.GetVulnerabilities(ctx, vulnIDs)
	if err != nil {
		return goerr.Wrap(err, "getting vulnerabilities")
	}

	var updated []*ttypes.DetectedVulnerability
	for _, vuln := range exists {
		v, ok := vulnSet[vuln.ID]
		if !ok {
			return goerr.Wrap(types.ErrLogicError, "vulnerability ID is not found in vulnSet").With("vulnID", vuln.ID)
		}

		if v.LastModifiedDate.Before(vuln.LastModifiedAt.Time) {
			updated = append(updated, v)
		} else {
			delete(vulnSet, vuln.ID)
		}
	}

	// Update existing vulnerabilities
	for _, v := range updated {
		cvss, err := json.Marshal(v.CVSS)
		if err != nil {
			return goerr.Wrap(err, "marshaling CVSS").With("vulnID", v.VulnerabilityID).With("cvss", v.CVSS)
		}

		param := db.UpdateVulnerabilityParams{
			ID:          v.VulnerabilityID,
			Title:       v.Title,
			Description: v.Description,
			Severity:    v.Severity,
			CweIds:      v.CweIDs,
			Cvss:        pqtype.NullRawMessage{Valid: true, RawMessage: cvss},
			Reference:   v.References,
		}
		if v.PublishedDate != nil {
			param.PublishedAt = sql.NullTime{Valid: true, Time: *v.PublishedDate}
		}
		if v.LastModifiedDate != nil {
			param.LastModifiedAt = sql.NullTime{Valid: true, Time: *v.LastModifiedDate}
		}

		if err := q.UpdateVulnerability(ctx, param); err != nil {
			return goerr.Wrap(err, "updating vulnerability").With("vulnID", v.VulnerabilityID)
		}
	}

	// Insert new vulnerabilities
	for _, v := range vulnSet {
		cvss, err := json.Marshal(v.CVSS)
		if err != nil {
			return goerr.Wrap(err, "marshaling CVSS").With("vulnID", v.VulnerabilityID).With("cvss", v.CVSS)
		}

		param := db.SaveVulnerabilityParams{
			ID:          v.VulnerabilityID,
			Title:       v.Title,
			Description: v.Description,
			Severity:    v.Severity,
			CweIds:      v.CweIDs,
			Cvss:        pqtype.NullRawMessage{Valid: true, RawMessage: cvss},
			Reference:   v.References,
		}
		if v.PublishedDate != nil {
			param.PublishedAt = sql.NullTime{Valid: true, Time: *v.PublishedDate}
		}
		if v.LastModifiedDate != nil {
			param.LastModifiedAt = sql.NullTime{Valid: true, Time: *v.LastModifiedDate}
		}

		if err := q.SaveVulnerability(ctx, param); err != nil {
			return goerr.Wrap(err, "saving vulnerability").With("vulnID", v.VulnerabilityID)
		}
	}

	if err := tx.Commit(); err != nil {
		return goerr.Wrap(err, "commit transaction")
	}

	return nil
}
