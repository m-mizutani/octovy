package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func getAuthGitHub(c *gin.Context) {
	cfg := getConfig(c)
	meta := cfg.Usecase.GetOctovyMetadata()

	state, err := cfg.Usecase.CreateAuthState()
	if err != nil {
		c.Error(err)
		return
	}

	clientID, err := cfg.Usecase.GetGitHubAppClientID()
	if err != nil {
		c.Error(err)
		return
	}

	v := url.Values{}
	v.Set("client_id", clientID)
	v.Set("response_type", "code")
	v.Set("state", state)

	redirectTo := strings.TrimSuffix(meta.GitHubWebURL, "/") + "/login/oauth/authorize?" + v.Encode()

	c.SetCookie(cookieReferrerName, c.Query("callback"), 60, "", "", true, true)
	c.Redirect(http.StatusFound, redirectTo)
}

func getAuthGitHubCallback(c *gin.Context) {
	cfg := getConfig(c)
	meta := cfg.Usecase.GetOctovyMetadata()

	code := c.Query("code")
	state := c.Query("state")

	user, err := cfg.Usecase.AuthGitHubUser(code, state)
	if err != nil {
		errMsg := "Authentication failed in GitHub OAuth procedure, requestID: "
		if id, ok := c.Get(contextRequestIDKey); ok {
			errMsg += fmt.Sprintf("%s", id)
		} else {
			errMsg += "unknown"
		}

		v := url.Values{}
		v.Set("login_error", errMsg)
		c.Redirect(http.StatusFound, meta.FrontendURL+"?"+v.Encode())
		golambda.EmitError(err)
		return
	}

	ssn, err := cfg.Usecase.CreateSession(user)
	if err != nil {
		v := url.Values{}
		v.Set("login_error", "Failed to issue session token")
		c.Redirect(http.StatusFound, meta.FrontendURL+"?"+v.Encode())
		golambda.EmitError(err)
	}

	c.SetCookie(cookieTokenName, ssn.Token, 86400*7, "", "", true, true)
	redirectTo := meta.FrontendURL
	if v, err := c.Cookie(cookieReferrerName); err == nil {
		redirectTo = strings.TrimSuffix(redirectTo, "/") + "/#/" + v
	}

	c.Redirect(http.StatusFound, redirectTo)
}

func getLogout(c *gin.Context) {
	cfg := getConfig(c)
	meta := cfg.Usecase.GetOctovyMetadata()

	cookie, err := c.Cookie(cookieTokenName)
	if err != nil {
		c.Error(goerr.Wrap(model.ErrAuthenticationFailed, "No valid cookie"))
		return
	}

	if err := cfg.Usecase.RevokeSession(cookie); err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(cookieTokenName, "", 0, "", "/", true, true)
	c.Redirect(http.StatusFound, meta.FrontendURL)
}
