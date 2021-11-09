package db

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

// Rule
func (x *Client) GetCheckRules(ctx *model.Context) ([]*ent.CheckRule, error) {
	rules, err := x.client.CheckRule.Query().WithSeverity().All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return rules, nil
}

func (x *Client) CreateCheckRule(ctx *model.Context, req *model.RequestRule) (*ent.CheckRule, error) {
	rule, err := x.client.CheckRule.Create().
		SetCheckResult(req.Result).
		SetSeverityID(req.SeverityID).
		Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return rule, nil
}

func (x *Client) DeleteCheckRule(ctx *model.Context, id int) error {
	if err := x.client.CheckRule.DeleteOneID(id).Exec(ctx); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
