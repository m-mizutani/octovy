package service

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/infra/db"
	"github.com/m-mizutani/octovy/backend/pkg/infra/github"
	"github.com/m-mizutani/octovy/backend/pkg/infra/net"
	"github.com/m-mizutani/octovy/backend/pkg/infra/trivydb"
)

var logger = golambda.Logger

type Service struct {
	config *Config
	infra.Interfaces

	trivyDBPath string
	dbClient    infra.DBClient
}

var defaultInfra = infra.Interfaces{
	NewDB:            db.NewDynamoClient,
	NewTrivyDB:       trivydb.New,
	NewSecretManager: aws.NewSecretsManager,
	NewSQS:           aws.NewSQS,
	NewS3:            aws.NewS3,
	NewHTTP:          net.NewHTTP,
	NewGitHub:        github.New,
	Utils:            infra.DefaultUtils(),
}

func New(config *Config) *Service {
	return &Service{
		Interfaces: defaultInfra,
		config:     config,
	}
}

func (x *Service) DB() infra.DBClient {
	if x.dbClient == nil {
		client, err := x.NewDB(x.config.AwsRegion, x.config.TableName)
		if err != nil {
			panic("Failed NewDB: " + err.Error())
		}
		x.dbClient = client
	}
	return x.dbClient
}
