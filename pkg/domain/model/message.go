package model

import (
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

type ScanTarget struct {
	GitHubBranch
	CommitID       string
	UpdatedAt      int64
	RequestedAt    int64
	URL            string
	IsPullRequest  bool
	IsTargetBranch bool
}

type ScanRepositoryRequest struct {
	ScanTarget
	InstallID int64
	Feedback  *FeedbackOptions
}

type FeedbackOptions struct {
	PullReqID     *int
	PullReqBranch string
	CheckID       *int64
}

type UpdateVulnStatusRequest struct {
	GitHubRepo
	UserID string
	ent.VulnStatus
}
