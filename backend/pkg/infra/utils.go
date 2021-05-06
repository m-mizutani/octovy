package infra

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"time"
)

type Utils struct {
	TimeNow  TimeNow
	TempFile TempFile
	OpenZip  OpenZip
	Remove   Remove
}

func DefaultUtils() Utils {
	return Utils{
		TimeNow: func() time.Time {
			return time.Now().UTC()
		},
		TempFile: ioutil.TempFile,
		OpenZip:  zip.OpenReader,
		Remove:   os.Remove,
	}
}
