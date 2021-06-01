package api

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
)

func postWebhookGitHub(c *gin.Context) {
	cfg := getConfig(c)

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

	// github.com/google/go-github/v29/github have not support integration_installation
	if githubEventType == "integration_installation" {
		return
	}

	raw, err := github.ParseWebHook(githubEventType, eventBody)
	if err != nil {
		_ = c.Error(goerr.Wrap(err, "Failed to parse github webhook event body").With("body", string(eventBody)))
		return
	}

	switch event := raw.(type) {
	case *github.PushEvent:
		if err := cfg.Usecase.HandleGitHubPushEvent(event); err != nil {
			_ = c.Error(err)
			return
		}

	case *github.PullRequestEvent:
		if err := cfg.Usecase.HandleGitHubPullReqEvent(event); err != nil {
			_ = c.Error(err)
			return
		}

	case *github.InstallationEvent:
		if err := cfg.Usecase.HandleGitHubInstallationEvent(event); err != nil {
			_ = c.Error(err)
			return
		}

	default:
		logger.With("event", event).With("type", githubEventType).Warn("Unsupported event")
	}

	c.JSON(200, baseResponse{Data: "OK"})
}
