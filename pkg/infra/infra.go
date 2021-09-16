package infra

import (
	"io/ioutil"
	"time"

	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/github"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/queue"
	"github.com/m-mizutani/octovy/pkg/infra/trivydb"
)

type Interfaces struct {
	NewDB        db.Factory
	NewGitHub    github.Factory
	NewGitHubApp githubapp.Factory
	NewTrivyDB   trivydb.Factory
	ScanQueue    queue.Interface
	Utils        Utils
}

type Utils struct {
	Now      func() time.Time
	ReadFile func(fname string) ([]byte, error)
}

func New() Interfaces {
	return Interfaces{
		NewDB:        db.New,
		NewGitHub:    github.New,
		NewGitHubApp: githubapp.New,
		NewTrivyDB:   trivydb.New,
		ScanQueue:    queue.New(),
		Utils: Utils{
			Now:      time.Now,
			ReadFile: ioutil.ReadFile,
		},
	}
}
