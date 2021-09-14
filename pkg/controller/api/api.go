package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/pkg/errors"
)

var logger = golambda.Logger

const (
	contextUsecase      = "usecase"
	contextRequestIDKey = "requestID"
	cookieTokenName     = "token"
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

func New(uc usecase.Interface) *gin.Engine {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		reqID := uuid.New().String()
		logger.
			With("path", c.Request.URL.Path).
			With("params", c.Params).
			With("request_id", reqID).
			With("remote", c.ClientIP()).
			With("ua", c.Request.UserAgent()).
			Info("API request")

		c.Set(contextRequestIDKey, reqID)
		c.Set(contextUsecase, uc)
		c.Next()
	})

	engine.Use(func(c *gin.Context) {
		c.Next()

		if ginError := c.Errors.Last(); ginError != nil {
			if err := errors.Cause(ginError); err != nil {
				switch {
				case errors.Is(err, model.ErrInvalidInput):
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
	r.GET("/scan/report/:report_id", getScanReport)

	// r.GET("/vuln/:vuln_id", getVulnerability)
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
	cookie, err := c.Cookie(cookieTokenName)
	if err != nil || cookie == "" {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "No valid cookie")
	}

	uc := getUsecase(c)
	ssn, err := uc.ValidateSession(cookie)
	if err != nil {
		return nil, err
	}
	if ssn == nil {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Invalid user in token")
	}
	ssn.Token = "" // Erase token
	return ssn, nil
}
