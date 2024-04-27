package trivy

type MisconfSummary struct {
	Successes  int
	Failures   int
	Exceptions int
}

type MisconfStatus string

// DetectedMisconfiguration holds detected misconfigurations
type DetectedMisconfiguration struct {
	Type          string        `json:",omitempty"`
	ID            string        `json:",omitempty"`
	AVDID         string        `json:",omitempty"`
	Title         string        `json:",omitempty"`
	Description   string        `json:",omitempty"`
	Message       string        `json:",omitempty"`
	Namespace     string        `json:",omitempty"`
	Query         string        `json:",omitempty"`
	Resolution    string        `json:",omitempty"`
	Severity      string        `json:",omitempty"`
	PrimaryURL    string        `json:",omitempty"`
	References    []string      `json:",omitempty"`
	Status        MisconfStatus `json:",omitempty"`
	Layer         Layer         `json:",omitempty"`
	CauseMetadata CauseMetadata `json:",omitempty"`

	// For debugging
	Traces []string `json:",omitempty"`
}

type CauseMetadata struct {
	Resource    string       `json:",omitempty"`
	Provider    string       `json:",omitempty"`
	Service     string       `json:",omitempty"`
	StartLine   int          `json:",omitempty"`
	EndLine     int          `json:",omitempty"`
	Code        Code         `json:",omitempty"`
	Occurrences []Occurrence `json:",omitempty"`
}

type Occurrence struct {
	Resource string `json:",omitempty"`
	Filename string `json:",omitempty"`
	Location Location
}

type Code struct {
	Lines []Line
}

type Line struct {
	Number      int    `json:"Number"`
	Content     string `json:"Content"`
	IsCause     bool   `json:"IsCause"`
	Annotation  string `json:"Annotation"`
	Truncated   bool   `json:"Truncated"`
	Highlighted string `json:"Highlighted,omitempty"`
	FirstCause  bool   `json:"FirstCause"`
	LastCause   bool   `json:"LastCause"`
}
