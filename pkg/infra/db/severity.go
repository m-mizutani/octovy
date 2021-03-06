package db

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/severity"
)

// Severity
func (x *Client) CreateSeverity(ctx *model.Context, req *model.RequestSeverity) (*ent.Severity, error) {
	added, err := x.client.Severity.Create().
		SetLabel(req.Label).
		SetColor(req.Color).
		Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return added, nil
}

func (x *Client) DeleteSeverity(ctx *model.Context, id int) error {
	if err := x.client.Severity.DeleteOneID(id).Exec(ctx); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) GetSeverities(ctx *model.Context) ([]*ent.Severity, error) {
	got, err := x.client.Severity.Query().Order(ent.Asc(severity.FieldID)).All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	return got, nil
}

func (x *Client) UpdateSeverity(ctx *model.Context, id int, req *model.RequestSeverity) error {
	if err := x.client.Severity.UpdateOneID(id).SetLabel(req.Label).SetColor(req.Color).Exec(ctx); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) AssignSeverity(ctx *model.Context, vulnID string, id int) error {
	if err := x.client.Vulnerability.UpdateOneID(vulnID).SetCustomSeverityID(id).Exec(ctx); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}
