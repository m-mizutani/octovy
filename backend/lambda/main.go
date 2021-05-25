package main

import (
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/controller"
)

func main() {
	funcID := os.Getenv("LAMBDA_FUNC_ID")

	golambda.Start(func(event golambda.Event) (interface{}, error) {
		ctrl := controller.New()

		switch funcID {
		case "apiHandler":
			return ctrl.LambdaAPIHandler(event)
		case "updateDB":
			return ctrl.LambdaUpdateDB()
		case "feedback":
			return ctrl.LambdaFeedback(event)
		case "scanRepo":
			return ctrl.LambdaScanRepo(event)

		default:
			return nil, goerr.New("Unregistered LAMBDA_FUNC_ID").With("ID", funcID)
		}
	})
}
