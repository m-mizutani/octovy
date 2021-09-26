package utils

import (
	cryptoRand "crypto/rand"
	"math"
	"math/big"
	mathRand "math/rand"
)

func initGenerateToken() {
	v, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	mathRand.Seed(v.Int64())
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
