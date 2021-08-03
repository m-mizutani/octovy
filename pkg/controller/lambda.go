package controller

import (
	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/pkg/api"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func (x *Controller) LambdaAPIHandler(event golambda.Event) (interface{}, error) {
	var req events.APIGatewayProxyRequest
	if err := event.Bind(&req); err != nil {
		return nil, golambda.WrapError(err).With("event", event)
	}

	gin.SetMode(gin.ReleaseMode)
	engine := api.New(&api.Config{
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

	x.Config.TrivyDBPath = "/tmp/trivy.db"

	for _, record := range records {
		var req model.ScanRepositoryRequest
		if err := record.Bind(&req); err != nil {
			return nil, err
		}

		if err := x.Usecase.ScanRepository(&req); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (x *Controller) LambdaUpdateDB() (interface{}, error) {
	if err := x.Usecase.UpdateTrivyDB(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (x *Controller) LambdaFeedback(event golambda.Event) (interface{}, error) {
	records, err := event.DecapSQSBody()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		var req model.FeedbackRequest
		if err := record.Bind(&req); err != nil {
			return nil, err
		}

		if err := x.Usecase.FeedbackScanResult(&req); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
