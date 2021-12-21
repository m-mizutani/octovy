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

func New(uc *usecase.Usecase, options ...*Option) *gin.Engine {
	engine := gin.Default()
	opt := mergeOption(options)

	engine.Use(func(c *gin.Context) {
		c.Set(contextUsecase, uc)
	})
	engine.Use(requestLogging)
	if !opt.DisableAuth {
		engine.Use(authControl)
	}
	if !uc.DisableFrontend() {
		engine.Use(getStaticFile)
	}
	engine.Use(errorHandler)

	if !uc.DisableWebhookGitHub() {
		engine.POST("/webhook/github", postWebhookGitHub)
	}
	if !uc.DisableWebhookTrivy() {
		engine.POST("/webhook/trivy", postWebhookTrivy)
	}

	if !uc.DisableFrontend() {
		engine.GET("/auth/github", getAuthGitHub)
		engine.GET("/auth/github/callback", getAuthGitHubCallback)

		r := engine.Group("/api/v1")

		r.GET("/repository", getRepositories)
		r.GET("/repository/:owner/:repo", getRepository)
		r.GET("/repository/:owner/:repo/scan", getRepositoryScan)
		r.GET("/vulnerability", getVulnerabilities)
		r.GET("/vulnerability/:vuln_id", getVulnerability)
		r.POST("/vulnerability", postVulnerability)
		r.GET("/scan/:scan_id", getScanReport)
		r.GET("/scan/:scan_id/report", getScanPackages)

		r.POST("/status/:owner/:repo_name", postVulnStatus)
		r.GET("/user", getUser)

		r.GET("/severity", getSeverities)
		r.POST("/severity", createSeverity)
		r.PUT("/severity/:id", updateSeverity)
		r.POST("/severity/:id/assign/:vuln_id", assignSeverity)
		r.DELETE("/severity/:id", deleteSeverity)

		r.GET("/repo-label", getRepoLabels)
		r.POST("/repo-label", createRepoLabel)
		r.PUT("/repo-label/:id", updateRepoLabel)
		r.POST("/repo-label/:id/assign/:repo_id", assignRepoLabel)
		r.DELETE("/repo-label/:id/assign/:repo_id", unassignRepoLabel)
		r.DELETE("/repo-label/:id", deleteRepoLabel)
	}

	return engine
}

func getUsecase(c *gin.Context) *usecase.Usecase {
	v, ok := c.Get(contextUsecase)
	if !ok {
		panic("No config in contextUsecase")
	}
	uc, ok := v.(*usecase.Usecase)
	if !ok {
		panic("Type mismatch for contextUsecase")
	}
	return uc
}
