package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func getUser(c *gin.Context) {
	ssn := getSession(c)
	if ssn == nil {
		c.Error(goerr.Wrap(model.ErrAuthenticationFailed))
		return
	}

	uc := getUsecase(c)
	user, err := uc.LookupUser(model.NewContextWith(c), ssn.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: user})
}
