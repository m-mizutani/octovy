package trivy

import (
	"io"
	"os/exec"

	"github.com/m-mizutani/goerr"
)

type Client interface {
	Run(args []string) error
}

type clientImpl struct {
	path string
}

func New(path string) Client {
	return &clientImpl{
		path: path,
	}
}

func (x *clientImpl) Run(args []string) error {
	cmd := exec.Command(x.path, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return goerr.Wrap(err, "retrieving stderr pipe")
	}
	if err := cmd.Run(); err != nil {
		msg, _ := io.ReadAll(stderr)
		return goerr.Wrap(err, "executing trivy").With("stderr", msg)
	}

	return nil
}
