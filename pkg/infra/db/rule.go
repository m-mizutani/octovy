package db

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

// Rule
func (x *Client) GetRules(ctx *model.Context) ([]*ent.Rule, error) {
	rules, err := x.client.Rule.Query().WithSeverity().All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return rules, nil
}

func (x *Client) CreateRule(ctx *model.Context, req *model.RequestRule) (*ent.Rule, error) {
	rule, err := x.client.Rule.Create().
		SetAction(req.Action).
		SetSeverityID(req.SeverityID).
		Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return rule, nil
}

func (x *Client) DeleteRule(ctx *model.Context, id int) error {
	if err := x.client.Rule.DeleteOneID(id).Exec(ctx); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
