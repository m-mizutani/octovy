package model

import "github.com/m-mizutani/goerr"

type VulnRespType string

const (
	RespSnooze VulnRespType = "snooze"
	RespNever  VulnRespType = "never"
)

func (x VulnRespType) IsValid() error {
	switch x {
	case RespSnooze:
		return nil
	case RespNever:
		return nil
	}

	return goerr.Wrap(ErrInvalidInputValues, "Unsupported vulnRespType").With("type", x)
}

type VulnResponse struct {
	GitHubRepo
	PkgType   PkgType
	PkgName   string
	Type      VulnRespType
	VulnID    string
	Duration  int64
	CreatedAt int64
}

func (x *VulnResponse) IsValid() error {
	if err := x.GitHubRepo.IsValid(); err != nil {
		return err
	}
	if x.PkgName == "" {
		return goerr.Wrap(ErrInvalidInputValues, "PkgName must not be empty")
	}
	if x.PkgType == "" {
		return goerr.Wrap(ErrInvalidInputValues, "PkgType must not be empty")
	}
	if x.VulnID == "" {
		return goerr.Wrap(ErrInvalidInputValues, "VulnID must not be empty")
	}
	if err := x.Type.IsValid(); err != nil {
		return err
	}

	return nil
}
