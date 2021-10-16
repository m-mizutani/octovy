package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func getRepositories(c *gin.Context) {
	uc := getUsecase(c)

	resp, err := uc.GetRepositories(model.NewContextWith(c))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}
