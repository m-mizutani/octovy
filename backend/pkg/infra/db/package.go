package db

import (
	"fmt"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

func packageRecordPK(branch *model.GitHubBranch) string {
	return fmt.Sprintf("pkg:%s/%s@%s", branch.Owner, branch.RepoName, branch.Branch)
}
func packageRecordSK(src string, pkgType model.PkgType, pkgName string, pkgVer string) string {
	return fmt.Sprintf("%s|%s|%s|%s", src, pkgType, pkgName, pkgVer)
}
func packageRecordPK2(pkgType model.PkgType, pkgName string) string {
	return fmt.Sprintf("pkg:%s|%s", pkgType, pkgName)
}
func packageRecordSK2(branch *model.GitHubBranch, pkgVer string) string {
	return fmt.Sprintf("%s/%s@%s|%s", branch.Owner, branch.RepoName, branch.Branch, pkgVer)
}

func (x *DynamoClient) InsertPackageRecord(pkg *model.PackageRecord) (bool, error) {
	record := &dynamoRecord{
		PK:  packageRecordPK(&pkg.Detected.GitHubBranch),
		SK:  packageRecordSK(pkg.Source, pkg.Type, pkg.Name, pkg.Version),
		PK2: packageRecordPK2(pkg.Type, pkg.Name),
		SK2: packageRecordSK2(&pkg.Detected.GitHubBranch, pkg.Version),
		PK3: packageRecordPK(&pkg.Detected.GitHubBranch),
		SK3: packageRecordSK(pkg.Source, pkg.Type, pkg.Name, pkg.Version),
		Doc: pkg,
	}
	q := x.table.Put(record).If("attribute_not_exists(pk) AND attribute_not_exists(sk)")
	if err := q.Run(); err != nil {
		if !isConditionalCheckErr(err) {
			return false, goerr.Wrap(err).With("record", record)
		}

		return false, nil
	}

	return true, nil
}

func (x *DynamoClient) RemovePackageRecord(pkg *model.PackageRecord) error {
	pk := packageRecordPK(&pkg.Detected.GitHubBranch)
	sk := packageRecordSK(pkg.Source, pkg.Type, pkg.Name, pkg.Version)

	q := x.table.Update("pk", pk).
		Range("sk", sk).
		Set("doc.'Removed'", true).
		Remove("pk3", "sk3").
		Set("doc.'Removed'", true).
		Set("doc.'ScannedAt'", pkg.ScannedAt).
		If("doc.'ScannedAt' < ?", pkg.ScannedAt)

	if err := q.Run(); err != nil {
		if !isConditionalCheckErr(err) {
			return goerr.Wrap(err).With("pkg", pkg).With("pk", pk).With("sk", sk)
		}
	}

	return nil
}

func (x *DynamoClient) UpdatePackageRecord(pkg *model.PackageRecord) error {
	pk := packageRecordPK(&pkg.Detected.GitHubBranch)
	sk := packageRecordSK(pkg.Source, pkg.Type, pkg.Name, pkg.Version)

	// Record exists, then update vulnerability info
	update := x.table.Update("pk", pk).
		Range("sk", sk).
		Set("doc.'Vulnerabilities'", pkg.Vulnerabilities).
		Set("doc.'ScannedAt'", pkg.ScannedAt).
		If("doc.'ScannedAt' < ?", pkg.ScannedAt)

	if err := update.Run(); err != nil {
		if !isConditionalCheckErr(err) {
			return goerr.Wrap(err)
		}
	}

	return nil
}

func (x *DynamoClient) FindPackageRecordsByName(pkgType model.PkgType, pkgName string) ([]*model.PackageRecord, error) {
	var records []*dynamoRecord
	pk2 := packageRecordPK2(pkgType, pkgName)
	if err := x.table.Get("pk2", pk2).Index(dynamoGSIName2nd).All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err).With("pk2", pk2)
		}
	}

	packageRecords := make([]*model.PackageRecord, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&packageRecords[i]); err != nil {
			return nil, err
		}
	}
	return packageRecords, nil
}

func (x *DynamoClient) FindPackageRecordsByBranch(branch *model.GitHubBranch) ([]*model.PackageRecord, error) {
	var records []*dynamoRecord
	pk := packageRecordPK(branch)
	if err := x.table.Get("pk3", pk).Index(dynamoGSIName3rd).All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err).With("pk", pk)
		}
	}

	packageRecords := make([]*model.PackageRecord, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&packageRecords[i]); err != nil {
			return nil, err
		}
	}
	return packageRecords, nil
}
