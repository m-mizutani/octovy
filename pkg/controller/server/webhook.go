package server

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/goerr"
)

func postWebhookGitHub(c *gin.Context) {
	uc := getUsecase(c)
	logger := getLogger(c)

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

	switch event := raw.(type) {
	case *github.PushEvent:
		if err := uc.HandleGitHubPushEvent(c, event); err != nil {
			_ = c.Error(err)
			return
		}

	case *github.PullRequestEvent:
		if err := uc.HandleGitHubPullReqEvent(c, event); err != nil {
			_ = c.Error(err)
			return
		}

	case *github.InstallationEvent:
		if err := uc.HandleGitHubInstallationEvent(c, event); err != nil {
			_ = c.Error(err)
			return
		}

	default:
		logger.Warn().Interface("event", event).Interface("type", githubEventType).Msg("Unsupported event")
	}

	c.JSON(200, baseResponse{Data: "OK"})
}
