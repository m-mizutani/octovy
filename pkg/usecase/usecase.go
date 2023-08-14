package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/infra"

	trivy_types "github.com/aquasecurity/trivy/pkg/types"
)

type UseCase struct{}

func New(clients *infra.Clients) *UseCase {
	return &UseCase{}
}

func (x *UseCase) ScanRepository(dir string) error {
	tmp, err := os.CreateTemp("", "trivy-scan-*.json")
	if err != nil {
		return goerr.Wrap(err, "creating trivy tmp output file")
	}
	if err := tmp.Close(); err != nil {
		return goerr.Wrap(err, "closing trivy tmp output file")
	}

	trivyOptions := []string{
		"fs",
		"--format", "json",
		"-o", tmp.Name(),
		"--list-all-pkgs",
		dir,
	}
	cmd := exec.Command("trivy", trivyOptions...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return goerr.Wrap(err, "retrieving stderr pipe")
	}
	if err := cmd.Run(); err != nil {
		msg, _ := io.ReadAll(stderr)
		return goerr.Wrap(err, "executing trivy").With("stderr", msg)
	}

	fmt.Println(tmp.Name())

	fd, err := os.Open(tmp.Name())
	if err != nil {
		return goerr.Wrap(err, "opening trivy tmp output file")
	}

	var report trivy_types.Report
	if err := json.NewDecoder(fd).Decode(&report); err != nil {
		return goerr.Wrap(err, "decoding trivy report")
	}

	return nil
}
