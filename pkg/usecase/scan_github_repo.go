package usecase

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra"
	gh "github.com/m-mizutani/octovy/pkg/infra/gh"
	"github.com/m-mizutani/octovy/pkg/utils"

	ttype "github.com/aquasecurity/trivy/pkg/types"
)

type ScanGitHubRepoInput struct {
	Owner     string
	Repo      string
	CommitID  string
	InstallID types.GitHubAppInstallID
}

func (x *ScanGitHubRepoInput) Validate() error {
	if x.Owner == "" {
		return goerr.Wrap(types.ErrInvalidOption, "owner is empty")
	}
	if x.Repo == "" {
		return goerr.Wrap(types.ErrInvalidOption, "repo is empty")
	}
	if x.CommitID == "" {
		return goerr.Wrap(types.ErrInvalidOption, "commit ID is empty")
	}
	if x.InstallID == 0 {
		return goerr.Wrap(types.ErrInvalidOption, "install ID is empty")
	}

	return nil
}

// ScanGitHubRepo is a usecase to download a source code from GitHub and scan it with Trivy. Using GitHub App credentials to download a private repository, then the app should be installed to the repository and have read access.
// After scanning, the result is stored to the database. The temporary files are removed after the scan.
func (x *UseCase) ScanGitHubRepo(ctx *model.Context, input *ScanGitHubRepoInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	zipURL, err := x.clients.GitHubApp().GetArchiveURL(ctx, &gh.GetArchiveURLInput{
		Owner:     input.Owner,
		Repo:      input.Repo,
		CommitID:  input.CommitID,
		InstallID: input.InstallID,
	})
	if err != nil {
		return err
	}

	// Download zip file
	tmpZip, err := os.CreateTemp("", fmt.Sprintf("octovy_code.%s.%s.%s.*.zip",
		input.Owner, input.Repo, input.CommitID,
	))
	if err != nil {
		return goerr.Wrap(err, "failed to create temp file for zip file")
	}
	defer utils.SafeRemove(tmpZip.Name())

	if err := downloadZipFile(ctx, x.clients.HTTPClient(), zipURL, tmpZip); err != nil {
		return err
	}
	if err := tmpZip.Close(); err != nil {
		return goerr.Wrap(err, "failed to close temp file for zip file")
	}

	// Extract zip file to local temp directory
	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("octovy.%s.%s.%s.*", input.Owner, input.Repo, input.CommitID))
	if err != nil {
		return goerr.Wrap(err, "failed to create temp directory for zip file")
	}
	defer utils.SafeRemoveAll(tmpDir)

	if err := extractZipFile(ctx, tmpZip.Name(), tmpDir); err != nil {
		return err
	}

	// Scan local directory
	tmpResult, err := os.CreateTemp("", fmt.Sprintf("octovy_result.%s.%s.%s.*.json",
		input.Owner, input.Repo, input.CommitID,
	))
	if err != nil {
		return goerr.Wrap(err, "failed to create temp file for scan result")
	}
	defer utils.SafeRemove(tmpResult.Name())
	if err := tmpResult.Close(); err != nil {
		return goerr.Wrap(err, "failed to close temp file for scan result")
	}

	if err := x.clients.Trivy().Run(ctx, []string{
		"fs",
		"--exit-code", "0",
		"--no-progress",
		"--format", "json",
		"--output", tmpResult.Name(),
		"--list-all-pkgs",
		tmpDir,
	}); err != nil {
		return goerr.Wrap(err, "failed to scan local directory")
	}

	var report ttype.Report
	if err := unmarshalFile(tmpResult.Name(), &report); err != nil {
		return err
	}

	ctx.Logger().Info("Scan result", slog.Any("report", tmpResult.Name()))

	if err := saveScanReport(ctx, x.clients.DB(), &report); err != nil {
		return err
	}

	return nil
}

func unmarshalFile(path string, v any) error {
	fd, err := os.Open(path)
	if err != nil {
		return goerr.Wrap(err, "failed to open file").With("path", path)
	}
	defer utils.SafeClose(fd)

	if err := json.NewDecoder(fd).Decode(v); err != nil {
		return goerr.Wrap(err, "failed to decode json").With("path", path)
	}

	return nil
}

func downloadZipFile(ctx *model.Context, httpClient infra.HTTPClient, zipURL *url.URL, w io.Writer) error {
	zipReq, err := http.NewRequestWithContext(ctx, http.MethodGet, zipURL.String(), nil)
	if err != nil {
		return goerr.Wrap(err, "failed to create request for zip file").With("url", zipURL)
	}

	zipResp, err := httpClient.Do(zipReq)
	if err != nil {
		return goerr.Wrap(err, "failed to download zip file").With("url", zipURL)
	}
	defer zipResp.Body.Close()

	if zipResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(zipResp.Body)
		return goerr.Wrap(types.ErrInvalidGitHubData, "failed to download zip file").With("url", zipURL).With("resp", zipResp).With("body", body)
	}

	if _, err = io.Copy(w, zipResp.Body); err != nil {
		return goerr.Wrap(err, "failed to write zip file").With("url", zipURL).With("resp", zipResp)
	}

	return nil
}

func extractZipFile(ctx *model.Context, src, dst string) error {
	zipFile, err := zip.OpenReader(src)
	if err != nil {
		return goerr.Wrap(err).With("file", src)
	}
	defer utils.SafeClose(zipFile)

	// Extract a source code zip file
	for _, f := range zipFile.File {
		if err := extractCode(ctx, f, dst); err != nil {
			return err
		}
	}

	return nil
}

func extractCode(ctx *model.Context, f *zip.File, dst string) error {
	if f.FileInfo().IsDir() {
		return nil
	}

	target := stepDownDirectory(f.Name)
	if target == "" {
		return nil
	}

	fpath := filepath.Join(dst, target)
	if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) {
		return goerr.Wrap(types.ErrInvalidGitHubData, "illegal file path of zip").With("path", fpath)
	}

	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return goerr.Wrap(err, "failed to create directory").With("path", fpath)
	}

	// #nosec
	out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return goerr.Wrap(err).With("fpath", fpath)
	}
	defer utils.SafeClose(out)

	rc, err := f.Open()
	if err != nil {
		return goerr.Wrap(err)
	}
	defer utils.SafeClose(rc)

	// #nosec
	_, err = io.Copy(out, rc)
	if err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

func stepDownDirectory(fpath string) string {
	if len(fpath) > 0 && fpath[0] == filepath.Separator {
		fpath = fpath[1:]
	}

	p := fpath
	var arr []string
	for {
		d, f := filepath.Split(p)
		if d == "" {
			break
		}
		arr = append([]string{f}, arr...)
		p = filepath.Clean(d)
	}

	return filepath.Join(arr...)
}
