package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getOctovyMetadata(c *gin.Context) {
	cfg := getConfig(c)
	meta := cfg.Usecase.GetOctovyMetadata()

	c.JSON(http.StatusOK, baseResponse{Data: meta})
}
