package model

import "github.com/aquasecurity/trivy-db/pkg/types"

type TrivyResults []TrivyResult

type TrivyReport struct {
	SchemaVersion int          `json:",omitempty"`
	ArtifactName  string       `json:",omitempty"`
	ArtifactType  string       `json:",omitempty"`
	Metadata      Metadata     `json:",omitempty"`
	Results       TrivyResults `json:",omitempty"`
}

// Metadata represents a metadata of artifact
type Metadata struct {
	Size int64 `json:",omitempty"`

	// Container image
	ImageID     string   `json:",omitempty"`
	DiffIDs     []string `json:",omitempty"`
	RepoTags    []string `json:",omitempty"`
	RepoDigests []string `json:",omitempty"`
}
type TrivyResult struct {
	Target          string                  `json:"Target"`
	Class           string                  `json:"Class,omitempty"`
	Type            string                  `json:"Type,omitempty"`
	Packages        []TrivyPackage          `json:"Packages,omitempty"`
	Vulnerabilities []DetectedVulnerability `json:"Vulnerabilities,omitempty"`
}

type TrivyPackage struct {
	Name            string `json:",omitempty"`
	Version         string `json:",omitempty"`
	Release         string `json:",omitempty"`
	Epoch           int    `json:",omitempty"`
	Arch            string `json:",omitempty"`
	SrcName         string `json:",omitempty"`
	SrcVersion      string `json:",omitempty"`
	SrcRelease      string `json:",omitempty"`
	SrcEpoch        int    `json:",omitempty"`
	Modularitylabel string `json:",omitempty"` // only for Red Hat based distributions
	License         string `json:",omitempty"`

	// Each package metadata have the file path, while the package from lock files does not have.
	FilePath string `json:",omitempty"`
}

type DetectedVulnerability struct {
	VulnerabilityID  string   `json:",omitempty"`
	VendorIDs        []string `json:",omitempty"`
	PkgName          string   `json:",omitempty"`
	PkgPath          string   `json:",omitempty"` // It will be filled in the case of language-specific packages such as egg/wheel and gemspec
	InstalledVersion string   `json:",omitempty"`
	FixedVersion     string   `json:",omitempty"`
	SeveritySource   string   `json:",omitempty"`
	PrimaryURL       string   `json:",omitempty"`

	types.Vulnerability
}
