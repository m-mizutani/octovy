package db

import (
	"fmt"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

func packagePK(branch *model.GitHubBranch) string {
	return fmt.Sprintf("pkg:%s/%s@%s", branch.Owner, branch.RepoName, branch.Branch)
}
func packageSK(src string, pkgType model.PkgType, pkgName string, pkgVer string) string {
	return fmt.Sprintf("%s|%s:%s@%s", src, pkgType, pkgName, pkgVer)
}
func packagePK2(pkgType model.PkgType, pkgName string) string {
	return fmt.Sprintf("pkg:%s:%s", pkgType, pkgName)
}
func packageSK2(branch *model.GitHubBranch, pkgVer string) string {
	return fmt.Sprintf("%s/%s@%s|%s", branch.Owner, branch.RepoName, branch.Branch, pkgVer)
}

func (x *DynamoClient) InsertPackage(pkg *model.Package) error {
	record := &dynamoRecord{
		PK:  packagePK(&pkg.GitHubBranch),
		SK:  packageSK(pkg.Source, pkg.PkgType, pkg.PkgName, pkg.Version),
		PK2: packagePK2(pkg.PkgType, pkg.PkgName),
		SK2: packageSK2(&pkg.GitHubBranch, pkg.Version),
		Doc: pkg,
	}
	q := x.table.Put(record).If("attribute_not_exists(pk) AND attribute_not_exists(sk)")
	if err := q.Run(); err != nil {
		if !isConditionalCheckErr(err) {
			return wrapErr(err).With("record", record)
		}
	}

	return nil
}

func (x *DynamoClient) DeletePackage(pkg *model.Package) error {
	pk := packagePK(&pkg.GitHubBranch)
	sk := packageSK(pkg.Source, pkg.PkgType, pkg.PkgName, pkg.Version)

	if err := x.table.Delete("pk", pk).Range("sk", sk).Run(); err != nil {
		if !isNotFoundErr(err) {
			return wrapErr(err).With("pkg", pkg).With("pk", pk).With("sk", sk)
		}
	}
	return nil
}

func (x *DynamoClient) FindPackagesByName(pkgType model.PkgType, pkgName string) ([]*model.Package, error) {
	var records []*dynamoRecord
	pk2 := packagePK2(pkgType, pkgName)
	if err := x.table.Get("pk2", pk2).Index(dynamoGSIName).All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err).With("pk2", pk2)
		}
	}

	packages := make([]*model.Package, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&packages[i]); err != nil {
			return nil, err
		}
	}
	return packages, nil
}

func (x *DynamoClient) FindPackagesByBranch(branch *model.GitHubBranch) ([]*model.Package, error) {
	var records []*dynamoRecord
	pk := packagePK(branch)
	if err := x.table.Get("pk", pk).All(&records); err != nil {
		if !isNotFoundErr(err) {
			return nil, goerr.Wrap(err).With("pk", pk)
		}
	}

	packages := make([]*model.Package, len(records))
	for i := range records {
		if err := records[i].Unmarshal(&packages[i]); err != nil {
			return nil, err
		}
	}
	return packages, nil
}
