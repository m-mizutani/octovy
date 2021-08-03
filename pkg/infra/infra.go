package infra

import (
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/infra/aws"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
	"github.com/m-mizutani/octovy/pkg/infra/githubauth"
	"github.com/m-mizutani/octovy/pkg/infra/githubcom"
	"github.com/m-mizutani/octovy/pkg/infra/trivydb"
	"github.com/m-mizutani/octovy/pkg/infra/utils"
)

func New() *interfaces.Infra {
	var defaultInfra = &interfaces.Infra{
		NewDB:            db.NewDynamoClient,
		NewTrivyDB:       trivydb.New,
		NewSecretManager: aws.NewSecretsManager,
		NewSQS:           aws.NewSQS,
		NewS3:            aws.NewS3,
		NewGitHubApp:     githubapp.New,
		NewGitHubCom:     githubcom.New,
		NewGitHubAuth:    githubauth.New,
		Utils:            utils.DefaultUtils(),
	}
	return defaultInfra
}
