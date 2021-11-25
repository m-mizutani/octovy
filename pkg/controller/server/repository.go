package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
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

func getRepository(c *gin.Context) {
	uc := getUsecase(c)

	resp, err := uc.GetRepository(model.NewContextWith(c), &model.GitHubRepo{
		Owner:    c.Param("owner"),
		RepoName: c.Param("repo"),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}

func getRepositoryScan(c *gin.Context) {
	uc := getUsecase(c)

	req := &model.GetRepoScanRequest{
		GitHubRepo: model.GitHubRepo{
			Owner:    c.Param("owner"),
			RepoName: c.Param("repo"),
		},
	}

	if s := c.Query("offset"); s != "" {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			c.Error(goerr.Wrap(err, "offset is not valid number"))
			return
		}
		req.Offset = int(n)
	}
	if s := c.Query("limit"); s != "" {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			c.Error(err)
			return
		}
		req.Limit = int(n)
	}

	resp, err := uc.GetRepositoryScan(model.NewContextWith(c), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}
