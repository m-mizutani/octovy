package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) {
	ssn, err := isAuthenticated(c)
	if err != nil {
		c.Error(err)
		return
	}

	cfg := getConfig(c)

	user, err := cfg.Usecase.LookupUser(ssn.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: user})
}
