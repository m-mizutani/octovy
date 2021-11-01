package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func getSeverities(c *gin.Context) {
	uc := getUsecase(c)

	resp, err := uc.GetSeverities(model.NewContextWith(c))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}

type severityRequest struct {
	Label string
}

func createSeverity(c *gin.Context) {
	uc := getUsecase(c)

	var req severityRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	resp, err := uc.CreateSeverity(model.NewContextWith(c), req.Label)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, baseResponse{Data: resp})
}

func updateSeverity(c *gin.Context) {
	uc := getUsecase(c)

	var req severityRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	if err := uc.UpdateSeverity(model.NewContextWith(c), int(id), req.Label); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{})
}

func assignSeverity(c *gin.Context) {
	uc := getUsecase(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	vulnID := c.Param("vuln_id")
	if err := uc.AssignSeverity(model.NewContextWith(c), vulnID, int(id)); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{})
}

func deleteSeverity(c *gin.Context) {
	uc := getUsecase(c)

	resp, err := uc.GetRepositories(model.NewContextWith(c))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}
