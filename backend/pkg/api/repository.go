package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func getOwners(c *gin.Context) {
	cfg := getConfig(c)
	owners, err := cfg.Usecase.FindOwners()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: owners})
}

func getReposByOwner(c *gin.Context) {
	cfg := getConfig(c)
	owner := c.Param("owner")
	repos, err := cfg.Usecase.FindReposByOwner(owner)
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
	repo, err := cfg.Usecase.FindReposByFullName(owner, name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: repo})
}

func getBranchInfo(c *gin.Context) {
	cfg := getConfig(c)
	owner := c.Param("owner")
	name := c.Param("name")
	branch := c.Param("branch")

	resp, err := cfg.Usecase.LookupBranch(&model.GitHubBranch{
		GitHubRepo: model.GitHubRepo{
			Owner:    owner,
			RepoName: name,
		},
		Branch: branch,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	if resp == nil {
		_ = c.Error(goerr.Wrap(errResourceNotFound, "No such branch"))
	}

	c.JSON(http.StatusOK, baseResponse{Data: resp})
}

/*
func getLatestScanResult(c *gin.Context) {
	cfg := getConfig(c)
	owner := c.Param("owner")
	name := c.Param("name")
	ref := c.Param("ref")

	if isValidCommitID(ref) {
		commit := &model.GitHubCommit{
			GitHubRepo: model.GitHubRepo{
				Owner:    owner,
				RepoName: name,
			},
			CommitID: ref,
		}

		result, err := cfg.Usecase.FindScanResult(commit)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, baseResponse{Data: result})
	} else {
		branch := &model.GitHubBranch{
			GitHubRepo: model.GitHubRepo{
				Owner:    owner,
				RepoName: name,
			},
			Branch: ref,
		}
		results, err := cfg.Usecase.FindLatestScanResults(branch, 1)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(results) > 0 {
			c.JSON(http.StatusOK, baseResponse{Data: results[0]})
		} else {
			c.JSON(http.StatusOK, baseResponse{Data: nil})
		}
	}
}
*/

func getPackage(c *gin.Context) {
	cfg := getConfig(c)
	pkgType := c.Query("type")
	pkgName := c.Query("name")
	packages, err := cfg.Usecase.FindPackageRecordsByName(model.PkgType(pkgType), pkgName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, baseResponse{Data: packages})
}
