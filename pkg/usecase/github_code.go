package usecase

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/githubapp"
)

type codes struct {
	Path string
}

func (x *codes) RemoveAll() {
	/*
		if err := os.RemoveAll(x.Path); err != nil {
			logger.Error().Interface("dirname", x.Path).Msg("Failed to remove src files")
		}
	*/
}

func extractCode(f *zip.File, dst string) error {
	if f.FileInfo().IsDir() {
		return nil
	}

	target := stepDownDirectory(f.Name)
	if target == "" {
		return nil
	}

	fpath := filepath.Join(dst, target)
	if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) {
		return goerr.Wrap(model.ErrInvalidGitHubData, "illegal file path of zip").With("path", fpath)
	}

	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return goerr.Wrap(err).With("fpath", fpath)
	}

	// #nosec
	out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return goerr.Wrap(err).With("fpath", fpath)
	}
	// #nosec, avoiding false positive
	defer func() {
		if err := out.Close(); err != nil {
			logger.With("err", err).Error("Close zip output file")
		}
	}()

	rc, err := f.Open()
	if err != nil {
		return goerr.Wrap(err)
	}
	defer func() {
		if err := rc.Close(); err != nil {
			logger.With("err", err).Error("Close zip input file")
		}
	}()

	// #nosec
	_, err = io.Copy(out, rc)
	if err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

func setupGitHubCodes(ctx *model.Context, req *model.ScanRepositoryRequest, app githubapp.Interface) (*codes, error) {
	tmp, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			ctx.Log().With("filename", tmp.Name()).Error("Failed to remove zip file")
		}
	}()

	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	resp := &codes{Path: tmpdir}

	if err := app.GetCodeZip(&req.GitHubRepo, req.CommitID, tmp); err != nil {
		return resp, err
	}

	zipFile, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return resp, goerr.Wrap(err).With("file", tmp.Name())
	}
	defer func() {
		if err := zipFile.Close(); err != nil {
			ctx.Log().With("zip", zipFile).With("err", err).Error("Failed to close zip file")
		}
	}()

	// Extract a source code zip file
	for _, f := range zipFile.File {
		if err := extractCode(f, tmpdir); err != nil {
			return resp, err
		}
	}

	return resp, nil
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
