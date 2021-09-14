package infra

import (
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/trivydb"
)

type Factories struct {
	NewDB        db.Factory
	NewGitHub    github.Factory
	NewGitHubApp githubapp.Factory
	NewTrivyDB   trivydb.Factory
}

func New() Factories {
	return Factories{
		NewDB:        db.New,
		NewGitHub:    github.New,
		NewGitHubApp: githubapp.New,
		NewTrivyDB:   trivydb.New,
	}
}
