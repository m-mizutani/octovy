package trivy

import (
	"bytes"
	"os/exec"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type Client interface {
	Run(ctx *model.Context, args []string) error
}

type clientImpl struct {
	path string
}

func New(path string) Client {
	return &clientImpl{
		path: path,
	}
}

func (x *clientImpl) Run(ctx *model.Context, args []string) error {
	cmd := exec.CommandContext(ctx, x.path, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		ctx.Logger().With("stderr", stderr.String()).With("stdout", stdout.String()).Error("trivy failed")
		return goerr.Wrap(err, "executing trivy").
			With("stderr", stderr.String()).
			With("stdout", stdout.String())
	}

	return nil
}
