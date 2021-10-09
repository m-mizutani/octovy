package trivy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

const (
	DefaultName = "trivy"
)

type Interface interface {
	SetPath(path string)
	Scan(dir string) (*report.Report, error)
}

type Trivy struct {
	path string
}

func New() *Trivy {
	return &Trivy{
		path: DefaultName,
	}
}

func (x *Trivy) SetPath(path string) {
	x.path = path
}

func (x *Trivy) Scan(dir string) (*report.Report, error) {
	temp, err := ioutil.TempFile("", "*.json")
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if err := temp.Close(); err != nil {
		return nil, goerr.Wrap(err)
	}
	defer func() {
		if err := os.Remove(temp.Name()); err != nil {
			logger.Error().Err(err).Msg("Failed to remove temp file")
		}
	}()

	cmd := exec.Command(x.path, "fs", "--list-all-pkgs", "-f", "json", "-o", temp.Name(), dir)
	cmd.Env = os.Environ()
	// https://github.com/aquasecurity/trivy/discussions/1050
	cmd.Env = append(cmd.Env, "TRIVY_NEW_JSON_SCHEMA=true")

	out, err := cmd.CombinedOutput()
	if err != nil {
		cwd, _ := os.Getwd()
		logger.Error().Err(err).Str("out", string(out)).Str("cwd", cwd).Msg("failed")
		return nil, goerr.Wrap(err).With("path", x.path).With("out", string(out))
	}

	raw, err := ioutil.ReadFile(temp.Name())
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	var result report.Report
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, goerr.Wrap(err)
	}

	return &result, nil
}
