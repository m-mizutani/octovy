package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/pkg/errors"
)

var logger = golambda.Logger

const (
	contextConfig       = "config"
	contextRequestIDKey = "requestID"
	cookieTokenName     = "token"
	cookieReferrerName  = "referrer"
)

type Config struct {
	Usecase  interfaces.Usecases
	AssetDir string
}

type baseResponse struct {
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Error  string                 `json:"error"`
	Values map[string]interface{} `json:"values"`
}

func errResp(c *gin.Context, code int, err error) {
	var wErr *golambda.Error
	var gErr *goerr.Error

	switch {
	case errors.As(err, &wErr):
		logger.With("stack", wErr.Stacks()).With("values", wErr.Values()).With("msg", wErr.Error()).Error("Failed with golambda.Error")
		c.JSON(code, &errorResponse{
			Error:  wErr.Error(),
			Values: wErr.Values(),
		})
	case errors.As(err, &gErr):
		logger.With("stack", gErr.Stacks()).With("values", gErr.Values()).With("msg", gErr.Error()).Error("Failed with goerr.Error")
		c.JSON(code, &errorResponse{
			Error:  gErr.Error(),
			Values: gErr.Values(),
		})
	default:
		logger.With("error", wErr).Error("Failed with normal Error")
		c.JSON(code, &errorResponse{
			Error: err.Error(),
		})
	}
}

func New(cfg *Config) *gin.Engine {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		reqID := uuid.New().String()
		logger.
			With("path", c.FullPath()).
			With("params", c.Params).
			With("request_id", reqID).
			With("remote", c.ClientIP()).
			With("ua", c.Request.UserAgent()).
			Info("API request")

		c.Set(contextRequestIDKey, reqID)
		c.Set(contextConfig, cfg)
		c.Next()
	})

	engine.Use(func(c *gin.Context) {
		c.Next()

		if ginError := c.Errors.Last(); ginError != nil {
			if err := errors.Cause(ginError); err != nil {
				switch {
				case errors.Is(err, model.ErrInvalidValue):
					errResp(c, http.StatusNotAcceptable, err)
				case errors.Is(err, errResourceNotFound):
					errResp(c, http.StatusNotFound, err)
				case errors.Is(err, model.ErrAuthenticationFailed):
					errResp(c, http.StatusUnauthorized, err)
				case errors.Is(err, model.ErrUserNotFound):
					errResp(c, http.StatusNotFound, err)
				default:
					golambda.EmitError(err)
					errResp(c, http.StatusInternalServerError, err)
				}
			} else {
				golambda.EmitError(err)
				errResp(c, http.StatusInternalServerError, ginError)
			}
		}
	})

	engine.GET("/", getIndex)
	engine.GET("/bundle.js", getBundleJS)
	engine.GET("/auth/github", getAuthGitHub)
	engine.GET("/auth/github/callback", getAuthGitHubCallback)
	engine.GET("/auth/logout", getLogout)
	engine.POST("/webhook/github", postWebhookGitHub)

	r := engine.Group("/api/v1")
	r.POST("/webhook/github", postWebhookGitHub)
	r.GET("/repo", getOwners)
	r.GET("/repo/:owner", getReposByOwner)
	r.GET("/repo/:owner/:name", getRepoInfo)
	r.GET("/repo/:owner/:name/:branch", getBranchInfo)
	r.GET("/scan/report/:report_id", getScanReport)
	r.GET("/package", getPackage)
	r.GET("/vuln/:vuln_id", getVulnerability)
	r.POST("/status/:owner/:repo_name", postVulnStatus)
	r.GET("/meta/octovy", getOctovyMetadata)
	r.GET("/user", getUser)

	return engine
}

func getConfig(c *gin.Context) *Config {
	v, ok := c.Get(contextConfig)
	if !ok {
		panic("No config in contextConfig")
	}
	config, ok := v.(*Config)
	if !ok {
		panic("Type mismatch for contextConfig")
	}
	return config

}

func isAuthenticated(c *gin.Context) (*model.Session, error) {
	cookie, err := c.Cookie(cookieTokenName)
	if err != nil || cookie == "" {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "No valid cookie")
	}

	cfg := getConfig(c)
	ssn, err := cfg.Usecase.ValidateSession(cookie)
	if err != nil {
		return nil, err
	}
	if ssn == nil {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Invalid user in token")
	}
	ssn.Token = "" // Erase token
	return ssn, nil
}
