package service

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/infra/aws"
	"github.com/m-mizutani/octovy/backend/pkg/infra/db"
	"github.com/m-mizutani/octovy/backend/pkg/infra/fs"
	"github.com/m-mizutani/octovy/backend/pkg/infra/net"
)

var logger = golambda.Logger

type Service struct {
	config *Config
	infra.Interfaces

	dbClient  infra.DBClient
	smClient  infra.SecretsManagerClient
	sqsClient infra.SQSClient
}

var defaultInfra = infra.Interfaces{
	NewDB:            db.NewDynamoClient,
	NewSecretManager: aws.NewSecretsManager,
	NewSQS:           aws.NewSQS,
	NewHTTP:          net.NewHTTP,
	FS:               &fs.FS{},
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
