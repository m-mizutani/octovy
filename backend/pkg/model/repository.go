package model

type GitHubRepo struct {
	Owner    string
	RepoName string
}

type GitHubBranch struct {
	GitHubRepo
	Branch string
}

type GitHubCommit struct {
	GitHubRepo
	CommitID string
}

type Repository struct {
	GitHubRepo
	URL           string
	Branches      []string
	DefaultBranch string
	InstallID     int64
	VulnerablePkg int64
	TotalPkg      int64
}
