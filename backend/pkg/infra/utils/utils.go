package utils

import (
	"archive/zip"
	cryptoRand "crypto/rand"
	"io/ioutil"
	"math"
	"math/big"
	mathRand "math/rand"
	"os"
	"time"

	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
)

func init() {
	v, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	mathRand.Seed(v.Int64())
}

func DefaultUtils() *interfaces.Utils {
	return &interfaces.Utils{
		TimeNow: func() time.Time {
			return time.Now().UTC()
		},
		TempFile:      ioutil.TempFile,
		OpenZip:       zip.OpenReader,
		Remove:        os.Remove,
		GenerateToken: GenerateToken,
	}
}

var randomTokenChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateToken(n int) string {
	s := make([]rune, n)
	for i := 0; i < n; i++ {
		v := mathRand.Int63() % int64(len(randomTokenChars))
		s[i] = randomTokenChars[v]
	}
	return string(s)
}
