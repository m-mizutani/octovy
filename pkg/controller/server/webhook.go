package server

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func postWebhookGitHub(c *gin.Context) {
	uc := getUsecase(c)

	githubEventType := c.Request.Header.Get("X-GitHub-Event")
	if githubEventType == "" {
		_ = c.Error(goerr.Wrap(errAPIInvalidParameter, "No X-GitHub-Event"))
		return
	}

	eventBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.Error(goerr.Wrap(err, "Failed to read github webhook event body"))
		return
	}

	if err := uc.VerifyGitHubSecret(c.GetHeader("X-Hub-Signature-256"), eventBody); err != nil {
		_ = c.Error(err)
		return
	}

	// github.com/google/go-github/v29/github have not support integration_installation
	if githubEventType == "integration_installation" {
		return
	}

	raw, err := github.ParseWebHook(githubEventType, eventBody)
	if err != nil {
		_ = c.Error(goerr.Wrap(err, "Failed to parse github webhook event body").With("body", string(eventBody)))
		return
	}

	ctx := model.NewContextWith(c)
	switch event := raw.(type) {
	case *github.PushEvent:
		if err := uc.HandleGitHubPushEvent(ctx, event); err != nil {
			_ = c.Error(err)
			return
		}

	case *github.PullRequestEvent:
		if err := uc.HandleGitHubPullReqEvent(ctx, event); err != nil {
			_ = c.Error(err)
			return
		}

	case *github.InstallationEvent:
		if err := uc.HandleGitHubInstallationEvent(ctx, event); err != nil {
			_ = c.Error(err)
			return
		}

	default:
		getLog(c).With("event", event).With("type", githubEventType).Warn("Unsupported event")
	}

	c.JSON(200, baseResponse{Data: "OK"})
}
