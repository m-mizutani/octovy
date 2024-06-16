package model

import (
	_ "embed"
	"os"
	"path/filepath"
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

func BuildConfig(configData ...[]byte) (*Config, error) {
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

// LoadConfigsFromDir loads configuration files from the repository. The configuration files are used to scan the repository with Trivy. The configuration .cue files are read recursively from the root directory of the repository.
func LoadConfigsFromDir(path string) (*Config, error) {
	// If path does not exist, return an empty configuration
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{}, nil
	}

	// read .cue files recursively from the root directory
	var configData [][]byte

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(filePath) == ".cue" {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return goerr.Wrap(err, "failed to read file").With("path", filePath)
			}
			configData = append(configData, data)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// If no configuration files are found, return an empty configuration
	if len(configData) == 0 {
		return &Config{}, nil
	}

	cfg, err := BuildConfig(configData...)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to load config")
	}

	return cfg, nil
}
