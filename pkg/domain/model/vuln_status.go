package model

import (
	"strings"

	"github.com/m-mizutani/goerr"
)

type VulnStatusType string

const (
	StatusNone       VulnStatusType = "none"
	StatusSnoozed    VulnStatusType = "snoozed"
	StatusMitigated  VulnStatusType = "mitigated"
	StatusUnaffected VulnStatusType = "unaffected"
	StatusFixed      VulnStatusType = "fixed"
)

func (x VulnStatusType) IsValid() error {
	switch x {
	case StatusNone, StatusSnoozed, StatusFixed, StatusMitigated, StatusUnaffected:
		return nil
	}

	return goerr.Wrap(ErrInvalidValue, "Unsupported VulnStatusType").With("type", x)
}

func (x VulnStatusType) IsUpdatable() error {
	switch x {
	case StatusNone, StatusSnoozed, StatusMitigated, StatusUnaffected:
		return nil
	}

	return goerr.Wrap(ErrInvalidValue, "Only snoozed, mitigated and none are allowed to update").With("type", x)
}

type VulnPackageKey struct {
	Source  string
	PkgType PkgType
	PkgName string
	VulnID  string
}

func (x *VulnPackageKey) Key() string {
	return strings.Join([]string{x.Source, x.PkgName, x.VulnID}, "|")
}

type VulnStatus struct {
	GitHubRepo
	VulnPackageKey

	ID        string
	Status    VulnStatusType
	ExpiresAt int64
	CreatedAt int64
	UserID    string
	Comment   string
}

func (x *VulnStatus) IsValid() error {
	if err := x.GitHubRepo.IsValid(); err != nil {
		return err
	}
	if x.Source == "" {
		return goerr.Wrap(ErrInvalidValue, "Source must not be empty")
	}
	if x.PkgName == "" {
		return goerr.Wrap(ErrInvalidValue, "PkgName must not be empty")
	}
	if x.PkgType == "" {
		return goerr.Wrap(ErrInvalidValue, "PkgType must not be empty")
	}
	if x.VulnID == "" {
		return goerr.Wrap(ErrInvalidValue, "VulnID must not be empty")
	}
	if err := x.Status.IsValid(); err != nil {
		return err
	}

	if x.ExpiresAt < 0 {
		return goerr.Wrap(ErrInvalidValue, "ExpiresAt must be >= 0")
	}
	if x.CreatedAt <= 0 {
		return goerr.Wrap(ErrInvalidValue, "CreatedAt must be > 0")
	}

	if x.Status == StatusSnoozed && x.ExpiresAt == 0 {
		return goerr.Wrap(ErrInvalidValue, "Snoozed status must have ExpiresAt over 0")
	}

	return nil
}
