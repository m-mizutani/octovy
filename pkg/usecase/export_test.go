package usecase

import (
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/service"
)

func InjectDBClient(usecases interfaces.Usecases, dbClient interfaces.DBClient) {
	uc, ok := usecases.(*Default)
	if !ok {
		panic("Failed dynamic cast to usecase.Default")
	}
	uc.svc.Infra.NewDB = func(region, tableName string) (interfaces.DBClient, error) {
		return dbClient, nil
	}
}

func ExposeService(usecases interfaces.Usecases) *service.Service {
	uc, ok := usecases.(*Default)
	if !ok {
		panic("Failed dynamic cast to usecase.Default")
	}
	return uc.svc
}

func EjectDBClient(usecases interfaces.Usecases) interfaces.DBClient {
	uc, ok := usecases.(*Default)
	if !ok {
		panic("Failed dynamic cast to usecase.Default")
	}
	return uc.svc.DB()
}
