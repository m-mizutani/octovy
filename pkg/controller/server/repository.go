package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getRepositories(c *gin.Context) {
	uc := getUsecase(c)

	resp, err := uc.GetRepositories(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}
