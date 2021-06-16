package infra

import (
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/infra/db"
	"github.com/m-mizutani/octovy/backend/pkg/infra/github"
	"github.com/m-mizutani/octovy/backend/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/backend/pkg/infra/githubauth"
	"github.com/m-mizutani/octovy/backend/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/backend/pkg/infra/utils"
)

func New() *interfaces.Infra {
	var defaultInfra = &interfaces.Infra{
		NewDB:            db.NewDynamoClient,
		NewTrivyDB:       trivydb.New,
		NewSecretManager: aws.NewSecretsManager,
		NewSQS:           aws.NewSQS,
		NewS3:            aws.NewS3,
		NewGitHubApp:     githubapp.New,
		NewGitHub:        github.New,
		NewGitHubAuth:    githubauth.New,
		Utils:            utils.DefaultUtils(),
	}
	return defaultInfra
}
