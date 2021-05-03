package main

import (
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/controller"
)

var logger = golambda.Logger

func main() {
	golambda.Start(func(_ golambda.Event) (interface{}, error) {
		return controller.New().LambdaUpdateDB()
	})
}
