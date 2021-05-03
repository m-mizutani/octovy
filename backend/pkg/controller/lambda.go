package controller

import (
	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/api"
	"github.com/m-mizutani/octovy/backend/pkg/model"
	"github.com/m-mizutani/octovy/backend/pkg/service"
	"github.com/m-mizutani/octovy/backend/pkg/usecase"
)

func (x *Controller) LambdaAPIHandler(event golambda.Event) (interface{}, error) {
	var req events.APIGatewayProxyRequest
	if err := event.Bind(&req); err != nil {
		return nil, golambda.WrapError(err).With("event", event)
	}

	svc := service.New(x.Config)

	gin.SetMode(gin.ReleaseMode)
	engine := api.New(&api.Config{
		Service:  svc,
		Usecase:  x.Usecase,
		AssetDir: "assets",
	})

	ginLambda := ginadapter.New(engine)

	return ginLambda.Proxy(req)
}

func (x *Controller) LambdaScanRepo(event golambda.Event) (interface{}, error) {
	records, err := event.DecapSQSBody()
	if err != nil {
		return nil, goerr.Wrap(err).With("event", event)
	}

	svc := service.New(x.Config)
	for _, record := range records {
		var req model.ScanRepositoryRequest
		if err := record.Bind(&req); err != nil {
			return nil, err
		}

		if err := x.Usecase.ScanRepository(svc, &req); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (x *Controller) LambdaUpdateDB() (interface{}, error) {
	svc := service.New(x.Config)

	if err := usecase.UpdateTrivyDB(svc); err != nil {
		return nil, err
	}
	return nil, nil
}
