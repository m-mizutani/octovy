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
}

type Branch struct {
	GitHubBranch
	LastScannedAt int64
	PkgTypes      []PkgType
	PkgCount      int64
	VulnCount     int64
}
