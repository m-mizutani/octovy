package model

type GitHubRepo struct {
	Owner    string
	RepoName string
}

type GitHubBranch struct {
	GitHubRepo
	Branch string
}

type Repository struct {
	GitHubRepo
	URL           string
	Branches      []string
	DefaultBranch string
}
