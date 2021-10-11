package model

import "github.com/m-mizutani/octovy/pkg/infra/ent"

type RespVulnerability struct {
	Vulnerability *ent.Vulnerability `json:"vulnerability"`
	Affected      []*ent.Repository  `json:"affected"`
}
