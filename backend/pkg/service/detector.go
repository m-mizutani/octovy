package service

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/service/detector"
)

func (x *Service) Detector() (*detector.Detector, error) {
	x.trivyDBPath = x.config.TrivyDBPath
	var dbOut io.WriteCloser

	if x.trivyDBPath == "" {
		tmp, err := ioutil.TempFile("", "*.db")
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		x.trivyDBPath = tmp.Name()
		dbOut = tmp
	} else {
		fs, err := os.Create(x.trivyDBPath)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		dbOut = fs
	}

	if err := x.downloadTrivyDB(dbOut); err != nil {
		return nil, err
	}
	if err := dbOut.Close(); err != nil {
		return nil, goerr.Wrap(err).With("path", x.trivyDBPath)
	}

	db, err := x.NewTrivyDB(x.trivyDBPath)
	if err != nil {
		return nil, err
	}

	return detector.New(db), nil
}
