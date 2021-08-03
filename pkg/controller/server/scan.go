package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
)

func getScanReport(c *gin.Context) {
	cfg := getConfig(c)
	reportID := c.Param("report_id")

	report, err := cfg.Usecase.LookupScanReport(reportID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if report == nil {
		_ = c.Error(goerr.Wrap(errResourceNotFound, "No such report"))
	}

	c.JSON(http.StatusOK, baseResponse{Data: report})
}
