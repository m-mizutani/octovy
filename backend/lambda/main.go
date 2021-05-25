package main

import (
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/controller"
)

func cleanupTempDir() {
	files, err := filepath.Glob("/tmp/*.zip")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func main() {
	funcID := os.Getenv("LAMBDA_FUNC_ID")

	golambda.Start(func(event golambda.Event) (interface{}, error) {
		cleanupTempDir()
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
