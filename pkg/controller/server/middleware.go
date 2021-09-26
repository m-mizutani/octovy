package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func setLogger(c *gin.Context, logger zerolog.Logger) {
	c.Set(contextLogger, logger)
}

func getLogger(c *gin.Context) zerolog.Logger {
	obj, ok := c.Get(contextLogger)
	if !ok {
		return globalLogger
	}
	logger, ok := obj.(zerolog.Logger)
	if !ok {
		return globalLogger
	}
	return logger
}

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

func requestLogging(c *gin.Context) {
	reqID := uuid.New().String()
	globalLogger.Info().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Interface("params", c.Params).
		Str("request_id", reqID).
		Str("remote", c.ClientIP()).
		Str("ua", c.Request.UserAgent()).
		Msg("Request")
	c.Set(contextRequestIDKey, reqID)
	c.Set(contextLogger, globalLogger.With().Str("request_id", reqID).Logger())
	c.Next()
}

func authControl(c *gin.Context) {
	logger := getLogger(c)
	loginURL := "/login"

	if c.Request.URL.Path == "/auth/logout" {
		if ssnID, err := c.Cookie(cookieSessionID); err == nil {
			// Revoke session if session ID exists,
			uc := getUsecase(c)
			if err := uc.RevokeSession(c, ssnID); err != nil {
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
		logger.Debug().Msg("notAuthResp")
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
	ssn, err := uc.ValidateSession(c, ssnID)
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
	setLogger(c, logger.With().Str("session_id", ssn.ID).Logger())
	c.Next()
}

func errorHandler(c *gin.Context) {
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

}
