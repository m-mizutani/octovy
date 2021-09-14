package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func getAuthGitHub(c *gin.Context) {
	uc := getUsecase(c)

	state, err := uc.CreateAuthState()
	if err != nil {
		c.Error(err)
		return
	}

	clientID, err := uc.GetGitHubAppClientID()
	if err != nil {
		c.Error(err)
		return
	}

	v := url.Values{}
	v.Set("client_id", clientID)
	v.Set("response_type", "code")
	v.Set("state", state)

	redirectTo := "https://github.com/login/oauth/authorize?" + v.Encode()

	c.SetCookie(cookieReferrerName, c.Query("callback"), 60, "", "", true, true)
	c.Redirect(http.StatusFound, redirectTo)
}

func getAuthGitHubCallback(c *gin.Context) {
	uc := getUsecase(c)

	code := c.Query("code")
	state := c.Query("state")

	user, err := uc.AuthGitHubUser(code, state)
	if err != nil {
		errMsg := "Authentication failed in GitHub OAuth procedure, requestID: "
		if id, ok := c.Get(contextRequestIDKey); ok {
			errMsg += fmt.Sprintf("%s", id)
		} else {
			errMsg += "unknown"
		}

		v := url.Values{}
		v.Set("login_error", errMsg)
		c.Redirect(http.StatusFound, uc.FrontendURL()+"/#/repository?"+v.Encode())
		golambda.EmitError(err)
		return
	}

	ssn, err := uc.CreateSession(user)
	if err != nil {
		v := url.Values{}
		v.Set("login_error", "Failed to issue session token")
		c.Redirect(http.StatusFound, uc.FrontendURL()+"/#/repository?"+v.Encode())
		golambda.EmitError(err)
	}

	c.SetCookie(cookieTokenName, ssn.Token, 86400*7, "", "", true, true)
	redirectTo := uc.FrontendURL()
	if v, err := c.Cookie(cookieReferrerName); err == nil {
		redirectTo = strings.TrimSuffix(redirectTo, "/") + "/#/" + v
	}

	c.Redirect(http.StatusFound, redirectTo)
}

func getLogout(c *gin.Context) {
	uc := getUsecase(c)

	cookie, err := c.Cookie(cookieTokenName)
	if err != nil {
		c.Error(goerr.Wrap(model.ErrAuthenticationFailed, "No valid cookie"))
		return
	}

	if err := uc.RevokeSession(cookie); err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(cookieTokenName, "", 0, "", "/", true, true)
	c.Redirect(http.StatusFound, uc.FrontendURL())
}
