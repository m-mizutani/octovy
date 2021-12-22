package model

import (
	"math/rand"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/m-mizutani/goerr"
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

type RequestRepoLabel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (x *RequestRepoLabel) IsValid() error {
	if err := validation.Validate(x.Name,
		validation.Required,
		validation.Length(1, 64),
		is.ASCII,
	); err != nil {
		return ErrInvalidInput.Wrap(err).With("field", "name")
	}

	if err := validation.Validate(x.Description,
		validation.Length(0, 256),
	); err != nil {
		return ErrInvalidInput.Wrap(err).With("field", "description")
	}

	if err := validation.Validate(x.Color,
		is.HexColor,
		validation.Length(4, 7),
	); err != nil {
		return ErrInvalidInput.Wrap(err).With("field", "color")
	}

	return nil

}

type GetRepoScanRequest struct {
	GitHubRepo
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PushTrivyResultRequest struct {
	Target ScanTarget
	Report TrivyReport
}

func (x *PushTrivyResultRequest) IsValid() error {
	const ghNameRegex = `^[a-z0-9]+(?:(?:(?:[._]|__|[-]*)[a-z0-9]+)+)?$`
	ghNamePtn := regexp.MustCompile(ghNameRegex)
	commitPtn := regexp.MustCompile(`^[0-9a-f]{40}$`)

	if err := validation.ValidateStruct(&x.Target,
		validation.Field(&x.Target.Owner, validation.Required, validation.Match(ghNamePtn)),
		validation.Field(&x.Target.Name, validation.Required, validation.Match(ghNamePtn)),
		validation.Field(&x.Target.Branch, validation.Required, validation.Match(ghNamePtn)),
		validation.Field(&x.Target.Branch, validation.Required, validation.Match(ghNamePtn)),
		validation.Field(&x.Target.CommitID, validation.Length(40, 40), validation.Required, validation.Match(commitPtn)),
	); err != nil {
		return ErrInvalidInput.Wrap(err).With("req", x)
	}

	return nil
}
