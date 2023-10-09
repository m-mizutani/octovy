package model

type GitHubRepo struct {
	Owner string
	Repo  string
}

type GitHubCommit struct {
	GitHubRepo
	CommitID string
}
