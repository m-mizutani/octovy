package db

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
)

// Repository Label
func (x *Client) CreateRepoLabel(ctx *model.Context, req *model.RequestRepoLabel) (*ent.RepoLabel, error) {
	added, err := x.client.RepoLabel.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetColor(req.Color).
		Save(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return added, nil
}

func (x *Client) UpdateRepoLabel(ctx *model.Context, id int, req *model.RequestRepoLabel) error {
	_, err := x.client.RepoLabel.UpdateOneID(id).
		SetName(req.Name).
		SetDescription(req.Description).
		SetColor(req.Color).
		Save(ctx)

	if err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) DeleteRepoLabel(ctx *model.Context, id int) error {
	if err := x.client.RepoLabel.DeleteOneID(id).Exec(ctx); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) GetRepoLabels(ctx *model.Context) ([]*ent.RepoLabel, error) {
	resp, err := x.client.RepoLabel.Query().All(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	return resp, nil
}

func (x *Client) AssignRepoLabel(ctx *model.Context, repoID int, labelID int) error {
	_, err := x.client.Repository.UpdateOneID(repoID).AddLabelIDs(labelID).Save(ctx)
	if err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) UnassignRepoLabel(ctx *model.Context, repoID int, labelID int) error {
	_, err := x.client.Repository.UpdateOneID(repoID).RemoveLabelIDs(labelID).Save(ctx)
	if err != nil {
		return goerr.Wrap(err)
	}
	return nil
}
