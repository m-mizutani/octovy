package service

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
	"github.com/m-mizutani/octovy/backend/pkg/infra"
)

var logger = golambda.Logger

type Service struct {
	config *model.Config
	Infra  *interfaces.Infra

	trivyDBPath string
	dbClient    interfaces.DBClient
}

func New(config *model.Config) *Service {
	return &Service{
		Infra:  infra.New(),
		config: config,
	}
}

func (x *Service) DB() interfaces.DBClient {
	if x.dbClient == nil {
		client, err := x.Infra.NewDB(x.config.AwsRegion, x.config.TableName)
		if err != nil {
			panic("Failed NewDB: " + err.Error())
		}
		x.dbClient = client
	}
	return x.dbClient
}
