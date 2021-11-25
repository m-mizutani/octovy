package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func getRepoLabels(c *gin.Context) {
	uc := getUsecase(c)

	resp, err := uc.GetRepoLabels(model.NewContextWith(c))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}

func createRepoLabel(c *gin.Context) {
	uc := getUsecase(c)

	var req model.RequestRepoLabel
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	resp, err := uc.CreateRepoLabel(model.NewContextWith(c), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, baseResponse{Data: resp})
}

func updateRepoLabel(c *gin.Context) {
	uc := getUsecase(c)

	var req model.RequestRepoLabel
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	if err := uc.UpdateRepoLabel(model.NewContextWith(c), int(id), &req); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{})
}

func assignRepoLabel(c *gin.Context) {
	uc := getUsecase(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}
	repoID, err := strconv.ParseInt(c.Param("repo_id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	if err := uc.AssignRepoLabel(model.NewContextWith(c), int(repoID), int(id)); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, baseResponse{})
}

func unassignRepoLabel(c *gin.Context) {
	uc := getUsecase(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}
	repoID, err := strconv.ParseInt(c.Param("repo_id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	if err := uc.UnassignRepoLabel(model.NewContextWith(c), int(repoID), int(id)); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{})
}

func deleteRepoLabel(c *gin.Context) {
	uc := getUsecase(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	if err := uc.DeleteRepoLabel(model.NewContextWith(c), int(id)); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{})
}
