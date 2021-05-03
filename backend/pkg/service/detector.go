package service

import (
	"io/ioutil"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/service/detector"
)

func (x *Service) Detector() (*detector.Detector, error) {
	if x.trivyDBPath == "" {
		tmp, err := ioutil.TempFile("", "*.db")
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		if err := x.downloadTrivyDB(tmp); err != nil {
			return nil, err
		}
		x.trivyDBPath = tmp.Name()
	}

	db, err := x.NewTrivyDB(x.trivyDBPath)
	if err != nil {
		return nil, err
	}

	return detector.New(db), nil
}
