package infra

import (
	"io/ioutil"
	"time"

	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/trivy"
)

type Interfaces struct {
	DB           db.Interface
	GitHub       github.Interface
	NewGitHubApp githubapp.Factory
	Trivy        trivy.Interface
	Utils        *Utils
}

type Utils struct {
	Now      func() time.Time
	ReadFile func(fname string) ([]byte, error)
}

func NewUtils() *Utils {
	return &Utils{
		Now:      time.Now,
		ReadFile: ioutil.ReadFile,
	}
}

func New() *Interfaces {
	return &Interfaces{
		DB:           db.New(),
		GitHub:       github.New(),
		NewGitHubApp: githubapp.New,
		Trivy:        trivy.New(),
		Utils:        NewUtils(),
	}
}
