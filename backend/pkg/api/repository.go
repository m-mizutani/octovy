package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/backend/pkg/model"
)

func getRepos(c *gin.Context) {
	cfg := getConfig(c)
	repos, err := cfg.Service.DB().FindRepo()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: repos})
}

func getReposByOwner(c *gin.Context) {
	cfg := getConfig(c)
	owner := c.Param("owner")
	repos, err := cfg.Service.DB().FindRepoByOwner(owner)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: repos})
}

func getRepoInfo(c *gin.Context) {
	cfg := getConfig(c)
	owner := c.Param("owner")
	name := c.Param("name")
	repo, err := cfg.Service.DB().FindRepoByFullName(owner, name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: repo})
}

func getPackagesByRepoBranch(c *gin.Context) {
	cfg := getConfig(c)
	branch := &model.GitHubBranch{
		GitHubRepo: model.GitHubRepo{
			Owner:    c.Param("owner"),
			RepoName: c.Param("name"),
		},
		Branch: c.Param("branch"),
	}

	packages, err := cfg.Service.DB().FindPackageRecordsByBranch(branch)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: packages})
}

func getPackage(c *gin.Context) {
	cfg := getConfig(c)
	pkgType := c.Query("type")
	pkgName := c.Query("name")
	packages, err := cfg.Service.DB().FindPackageRecordsByName(model.PkgType(pkgType), pkgName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: packages})
}
