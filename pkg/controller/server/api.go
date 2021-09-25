package server

import (
	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/pkg/errors"
)

var globalLogger = utils.Logger

const (
	contextUsecase      = "usecase"
	contextRequestIDKey = "requestID"
	contextLogger       = "logger"
	contextSession      = "session"

	cookieSessionID     = "session_id"
	cookieSessionSecret = "secret"
	cookieReferrerName  = "referrer"
)

type baseResponse struct {
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Error  string                 `json:"error"`
	Values map[string]interface{} `json:"values"`
}

func errResp(c *gin.Context, code int, err error) {
	var gErr *goerr.Error

	switch {
	case errors.As(err, &gErr):
		c.JSON(code, &errorResponse{
			Error:  gErr.Error(),
			Values: gErr.Values(),
		})
	default:
		c.JSON(code, &errorResponse{
			Error: err.Error(),
		})
	}
}

func New(uc usecase.Interface) *gin.Engine {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		c.Set(contextUsecase, uc)
	})
	engine.Use(requestLogging)
	engine.Use(authControl)
	engine.Use(getStaticFile)
	engine.Use(errorHandler)

	engine.GET("/auth/github", getAuthGitHub)
	engine.GET("/auth/github/callback", getAuthGitHubCallback)
	engine.GET("/auth/logout", getLogout)
	engine.POST("/webhook/github", postWebhookGitHub)

	r := engine.Group("/api/v1")
	r.POST("/webhook/github", postWebhookGitHub)
	r.GET("/scan/:scan_id", getScanReport)

	r.POST("/status/:owner/:repo_name", postVulnStatus)
	// r.GET("/user", getUser)

	return engine
}

func getUsecase(c *gin.Context) usecase.Interface {
	v, ok := c.Get(contextUsecase)
	if !ok {
		panic("No config in contextUsecase")
	}
	uc, ok := v.(usecase.Interface)
	if !ok {
		panic("Type mismatch for contextUsecase")
	}
	return uc
}

func isAuthenticated(c *gin.Context) (*ent.Session, error) {
	ssnID, err := c.Cookie(cookieSessionID)
	if err != nil || ssnID == "" {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "No session ID in cookie")
	}
	secret, err := c.Cookie(cookieSessionSecret)
	if err != nil || secret == "" {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "No session secret in cookie")
	}

	uc := getUsecase(c)
	ssn, err := uc.ValidateSession(c, ssnID)
	if err != nil {
		return nil, err
	}
	if ssn == nil {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "session not found")
	}
	if ssn.Token != secret {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "invalid session secret")
	}

	ssn.Token = "" // Erase token
	return ssn, nil
}
