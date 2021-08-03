package main

import (
	"os"

	"github.com/m-mizutani/octovy/pkg/controller"
)

func main() {
	if err := controller.New().RunCmd(os.Args, os.Environ()); err != nil {
		os.Exit(1)
	}
}
