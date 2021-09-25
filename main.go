package main

import (
	"os"

	"github.com/m-mizutani/octovy/pkg/controller"
)

func main() {
	_ = controller.New().RunCmd(os.Args)
}
