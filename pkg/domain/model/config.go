package model

import (
	_ "embed"
	"time"

	"cuelang.org/go/cue/cuecontext"
	"github.com/m-mizutani/goerr"
)

type Config struct {
	IgnoreTargets []IgnoreTarget
}

//go:embed schema/ignore.cue
var ignoreCue []byte

type IgnoreTarget struct {
	File  string
	Vulns []IgnoreVuln
}

func (x *IgnoreTarget) Validate() error {
	for _, v := range x.Vulns {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type IgnoreVuln struct {
	ID          string
	Description string
	ExpiresAt   time.Time
}

func (x *IgnoreVuln) Validate() error {
	maxExpiresAt := time.Now().Add(time.Hour * 24 * 90)
	if x.ExpiresAt.After(maxExpiresAt) {
		return goerr.New("expiresAt is too far in the future, must be within 90 days from now")
	}

	return nil
}

func LoadConfig(configData ...[]byte) (*Config, error) {
	ctx := cuecontext.New()

	// Load the schema
	schemaInstance := ctx.CompileBytes(ignoreCue)
	if schemaInstance.Err() != nil {
		return nil, goerr.Wrap(schemaInstance.Err(), "failed to compile schema")
	}

	for _, data := range configData {
		// Load the configuration
		configInstance := ctx.CompileBytes(data)
		if configInstance.Err() != nil {
			return nil, goerr.Wrap(configInstance.Err(), "failed to compile configuration")
		}

		// Merge the schema and config
		mergedInstance := schemaInstance.Unify(configInstance)
		if mergedInstance.Err() != nil {
			return nil, goerr.Wrap(mergedInstance.Err(), "failed to unify schema and config")
		}

		schemaInstance = mergedInstance
	}

	// Extract the configuration into a Go struct
	var config Config
	if err := schemaInstance.Value().Decode(&config); err != nil {
		return nil, goerr.Wrap(err, "failed to decode configuration")
	}

	return &config, nil
}
