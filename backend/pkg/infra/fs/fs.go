package fs

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"

	"github.com/m-mizutani/goerr"
)

type FS struct{}

func (x *FS) WriteFile(r io.Reader, path string) error {
	w, err := os.Create(path)
	if err != nil {
		return goerr.Wrap(err).With("path", path)
	}
	if _, err := io.Copy(w, r); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *FS) OpenZip(path string) (*zip.ReadCloser, error) {
	return zip.OpenReader(path)
}

func (x *FS) TempFile(dir, pattern string) (f *os.File, err error) {
	return ioutil.TempFile(dir, pattern)
}

func (x *FS) Remove(name string) error {
	return os.Remove(name)
}
