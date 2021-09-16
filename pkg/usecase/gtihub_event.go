package usecase

import "github.com/google/go-github/v39/github"

func (x *usecase) HandleGitHubPushEvent(event *github.PushEvent) error {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) HandleGitHubPullReqEvent(event *github.PullRequestEvent) error {
	panic("not implemented") // TODO: Implement
}

func (x *usecase) HandleGitHubInstallationEvent(event *github.InstallationEvent) error {
	panic("not implemented") // TODO: Implement
}
