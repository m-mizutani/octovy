package main

import (
	"os"

	"github.com/m-mizutani/octovy/pkg/controller/cli"
)

func main() {
	if err := cli.New().Run(os.Args); err != nil {
		os.Exit(1)
	}
}
