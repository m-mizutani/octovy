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

	// Pass through settings
	if strings.HasPrefix(c.Request.URL.Path, "/auth/") ||
		strings.HasPrefix(c.Request.URL.Path, "/_next/") ||
		c.Request.URL.Path == loginURL {
		c.Next()
		return
	}

	ssnID, err := c.Cookie(cookieSessionID)
	if err != nil || ssnID == "" {
		logger.Warn().Err(err).Msg("No session ID, redirect to /login")
		c.Redirect(http.StatusFound, loginURL)
		return
	}
	secret, err := c.Cookie(cookieSessionSecret)
	if err != nil || secret == "" {
		logger.Warn().Err(err).Msg("No secret, redirect to /login")
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	uc := getUsecase(c)
	ssn, err := uc.ValidateSession(c, ssnID)
	if err != nil {
		if errors.Is(err, model.ErrAuthenticationFailed) {
			c.Redirect(http.StatusFound, loginURL)
		} else {
			c.Error(err)
		}
		return
	}

	// Not authenticated
	if ssn == nil || ssn.Token != secret {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			errResp(c, http.StatusUnauthorized, goerr.New("auth error"))
			return
		} else {
			c.Redirect(http.StatusFound, loginURL)
			return
		}
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
