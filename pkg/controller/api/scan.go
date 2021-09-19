package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
)

func getScanReport(c *gin.Context) {
	uc := getUsecase(c)
	scanID := c.Param("scan_id")

	report, err := uc.LookupScanReport(c, scanID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if report == nil {
		_ = c.Error(goerr.Wrap(errResourceNotFound, "No such report"))
	}

	c.JSON(http.StatusOK, baseResponse{Data: report})
}
