package model

import (
	"math/rand"
	"regexp"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

type RespVulnerability struct {
	Vulnerability *ent.Vulnerability `json:"vulnerability"`
	Affected      []*ent.Repository  `json:"affected"`
}

type RequestSeverity struct {
	Label string
	Color string
}

var (
	colorRegex = regexp.MustCompile("^#[0-9A-Fa-f]{6}$")
)

func randomColor() string {
	s := "#"
	chars := "0123456789ABCDEF"
	for i := 0; i < 6; i++ {
		s += string(chars[int(rand.Uint32())%len(chars)])
	}

	return s
}

func (x *RequestSeverity) IsValid() error {
	if x.Label == "" {
		return goerr.Wrap(ErrInvalidInput, "empty severity name is not allowed")
	}

	// Fill random color if empty
	if x.Color == "" {
		x.Color = randomColor()
	}

	if !colorRegex.MatchString(x.Color) {
		return goerr.Wrap(ErrInvalidInput, "invalid color schema")
	}

	return nil
}

type RequestCheckRule struct {
	SeverityID int
	Result     types.GitHubCheckResult
}

func (x *RequestCheckRule) IsValid() error {
	if err := x.Result.IsValid(); err != nil {
		return goerr.Wrap(err, "unsupported result")
	}

	return nil
}
