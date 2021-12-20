package model

import (
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

type ScanTarget struct {
	GitHubBranch
	CommitID     string
	UpdatedAt    int64
	RequestedAt  int64
	URL          string
	TargetBranch string
}

type ScanRepositoryRequest struct {
	ScanTarget
	InstallID     int64
	PullReqNumber *int
	PullReqAction string
}

type UpdateVulnStatusRequest struct {
	GitHubRepo
	UserID int
	ent.VulnStatus
}

type ImportRequest struct {
	Target ScanTarget
	Report TrivyReport
}
