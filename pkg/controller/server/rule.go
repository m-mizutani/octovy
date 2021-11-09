package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func GetCheckRules(c *gin.Context) {
	uc := getUsecase(c)

	resp, err := uc.GetCheckRules(model.NewContextWith(c))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}

func createCheckRule(c *gin.Context) {
	uc := getUsecase(c)

	var req model.RequestCheckRule
	if err := c.BindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	resp, err := uc.CreateRule(model.NewContextWith(c), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, baseResponse{Data: resp})
}

func deleteCheckRule(c *gin.Context) {
	uc := getUsecase(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	if err := uc.DeleteRule(model.NewContextWith(c), int(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{})
}
