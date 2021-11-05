package server

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"

	"github.com/m-mizutani/octovy/pkg/usecase"
)

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

type Option struct {
	DisableAuth bool
}

func mergeOption(options []*Option) *Option {
	var merged Option
	for _, opt := range options {
		merged.DisableAuth = opt.DisableAuth
	}
	return &merged
}

func New(uc usecase.Interface, options ...*Option) *gin.Engine {
	engine := gin.Default()
	opt := mergeOption(options)

	engine.Use(func(c *gin.Context) {
		c.Set(contextUsecase, uc)
	})
	engine.Use(requestLogging)
	if !opt.DisableAuth {
		engine.Use(authControl)
	}
	engine.Use(getStaticFile)
	engine.Use(errorHandler)

	engine.POST("/webhook/github", postWebhookGitHub)
	if !uc.WebhookOnly() {
		engine.GET("/auth/github", getAuthGitHub)
		engine.GET("/auth/github/callback", getAuthGitHubCallback)

		r := engine.Group("/api/v1")

		r.GET("/repository", getRepositories)
		r.GET("/vulnerability", getVulnerabilities)
		r.GET("/vulnerability/:vuln_id", getVulnerability)
		r.POST("/vulnerability", postVulnerability)
		r.GET("/scan/:scan_id", getScanReport)

		r.POST("/status/:owner/:repo_name", postVulnStatus)
		r.GET("/user", getUser)

		r.GET("/severity", getSeverities)
		r.POST("/severity", createSeverity)
		r.PUT("/severity/:id", updateSeverity)
		r.POST("/severity/:id/assign/:vuln_id", assignSeverity)
		r.DELETE("/severity/:id", deleteSeverity)

		r.GET("/rule", getRules)
		r.POST("/rule", createRule)
		r.DELETE("/rule/:id", deleteRule)
	}

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
