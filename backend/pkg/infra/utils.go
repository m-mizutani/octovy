package infra

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"time"

	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
)

func DefaultUtils() *interfaces.Utils {
	return &interfaces.Utils{
		TimeNow: func() time.Time {
			return time.Now().UTC()
		},
		TempFile: ioutil.TempFile,
		OpenZip:  zip.OpenReader,
		Remove:   os.Remove,
	}
}
