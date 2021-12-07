package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/m-mizutani/zlog"
	"github.com/pkg/errors"
)

func setSession(c *gin.Context, ssn *ent.Session) {
	c.Set(contextSession, ssn)
}

func getSession(c *gin.Context) *ent.Session {
	obj, ok := c.Get(contextSession)
	if !ok {
		return nil
	}
	ssn, ok := obj.(*ent.Session)
	if !ok {
		return nil
	}
	return ssn
}

type httpLog struct {
	Method    string
	Path      string
	Params    gin.Params
	RequestID string
	Remote    string
}

func setLog(c *gin.Context, key string, value interface{}) {
	log := getLog(c)
	if log == nil {
		log = utils.Logger.Log()
	}
	c.Set(model.ContextKeyLogger, log.With(key, value))
}

func getLog(c *gin.Context) *zlog.LogEntity {
	obj, ok := c.Get(model.ContextKeyLogger)
	if !ok {
		return nil
	}
	log, ok := obj.(*zlog.LogEntity)
	if !ok {
		panic("not matched with zlog.LogEntity")
	}
	return log
}

func requestLogging(c *gin.Context) {
	reqID := uuid.New().String()
	reqLog := httpLog{
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Params:    c.Params,
		RequestID: reqID,
		Remote:    c.ClientIP(),
	}
	utils.Logger.
		With("http", reqLog).
		With("user-agent", c.Request.UserAgent()).
		Info("HTTP Request")
	c.Set(contextRequestIDKey, reqID)
	setLog(c, "http", reqLog)
	c.Next()
}

func authControl(c *gin.Context) {
	loginURL := "/login"
	ctx := model.NewContextWith(c)

	if c.Request.URL.Path == "/auth/logout" {
		if ssnID, err := c.Cookie(cookieSessionID); err == nil {
			// Revoke session if session ID exists,
			uc := getUsecase(c)
			if err := uc.RevokeSession(ctx, ssnID); err != nil {
				c.Error(err)
				return
			}
		}
		c.SetCookie(cookieSessionID, "", 0, "", "/", true, true)
		c.SetCookie(cookieSessionSecret, "", 0, "", "/", true, true)
		c.Redirect(http.StatusFound, loginURL)
		c.Abort()
		return
	}

	// Pass through settings
	if strings.HasPrefix(c.Request.URL.Path, "/auth/") ||
		strings.HasPrefix(c.Request.URL.Path, "/_next/") ||
		strings.HasPrefix(c.Request.URL.Path, "/webhook/") ||
		c.Request.URL.Path == loginURL {
		c.Next()
		return
	}

	notAuthResp := func() {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			errResp(c, http.StatusUnauthorized, goerr.Wrap(model.ErrNotAuthenticated))
		} else {
			c.Redirect(http.StatusFound, loginURL)
		}
		c.Abort()
	}

	ssnID, err := c.Cookie(cookieSessionID)
	if err != nil || ssnID == "" {
		notAuthResp()
		return
	}
	secret, err := c.Cookie(cookieSessionSecret)
	if err != nil || secret == "" {
		notAuthResp()
		return
	}

	uc := getUsecase(c)
	ssn, err := uc.ValidateSession(ctx, ssnID)
	if err != nil {
		notAuthResp()
		return
	}

	// Not authenticated
	if ssn == nil || ssn.Token != secret {
		notAuthResp()
		return
	}

	ssn.Token = "" // Erase token
	setSession(c, ssn)
	setLog(c, "session_id", ssn.ID)
	c.Next()
}

func errorHandler(c *gin.Context) {
	c.Next()

	if ginError := c.Errors.Last(); ginError != nil {
		getUsecase(c).HandleError(model.NewContextWith(c), ginError)

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
			case errors.Is(err, model.ErrVulnerabilityNotFound):
				errResp(c, http.StatusNotFound, err)
			default:
				errResp(c, http.StatusInternalServerError, err)
			}
		} else {
			getLog(c).Err(ginError).Error("ginError")
			errResp(c, http.StatusInternalServerError, ginError)
		}
	}

}
