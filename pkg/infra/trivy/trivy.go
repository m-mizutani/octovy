package trivy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

const (
	DefaultName = "trivy"
)

type Interface interface {
	SetPath(path string)
	Scan(dir string) (*model.TrivyReport, error)
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

func (x *Trivy) Scan(dir string) (*model.TrivyReport, error) {
	temp, err := ioutil.TempFile("", "*.json")
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if err := temp.Close(); err != nil {
		return nil, goerr.Wrap(err)
	}
	defer func() {
		if err := os.Remove(temp.Name()); err != nil {
			utils.Logger.With("err", err).Error("Failed to remove temp file")
		}
	}()

	// #nosec
	cmd := exec.Command(x.path, "fs", "--list-all-pkgs", "-f", "json", "-o", temp.Name(), dir)
	cmd.Env = os.Environ()
	// https://github.com/aquasecurity/trivy/discussions/1050
	cmd.Env = append(cmd.Env, "TRIVY_NEW_JSON_SCHEMA=true")

	out, err := cmd.CombinedOutput()
	if err != nil {
		utils.Logger.With("err", err).With("out", string(out)).Error("failed")
		return nil, goerr.Wrap(err).With("path", x.path).With("out", string(out))
	}

	raw, err := ioutil.ReadFile(temp.Name())
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	var result model.TrivyReport
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, goerr.Wrap(err)
	}

	return &result, nil
}
