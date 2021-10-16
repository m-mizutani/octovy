package db

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/infra/ent/session"
)

// Auth
func (x *Client) SaveAuthState(ctx *model.Context, state string, expiresAt int64) error {
	_, err := x.client.AuthStateCache.Create().
		SetID(state).
		SetExpiresAt(expiresAt).
		Save(ctx)
	if err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) HasAuthState(ctx *model.Context, state string, now int64) (bool, error) {
	cache, err := x.client.AuthStateCache.Get(ctx, state)
	switch {
	case err == nil:
		return now < cache.ExpiresAt, nil
	case ent.IsNotFound(err):
		return false, nil
	default:
		return false, goerr.Wrap(err)
	}
}

func (x *Client) GetUser(ctx *model.Context, userID int) (*ent.User, error) {
	got, err := x.client.User.Get(ctx, userID)
	if ent.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, goerr.Wrap(err)
	}

	return got, nil
}

func (x *Client) PutUser(ctx *model.Context, user *ent.User) (int, error) {
	userID, err := x.client.User.Create().
		SetGithubID(user.GithubID).
		SetLogin(user.Login).
		SetName(user.Name).
		SetURL(user.URL).
		SetAvatarURL(user.AvatarURL).
		OnConflictColumns("github_id").
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return 0, goerr.Wrap(err)
	}
	return userID, nil
}

func (x *Client) PutSession(ctx *model.Context, ssn *ent.Session) error {
	_, err := x.client.Session.Create().
		SetID(ssn.ID).
		SetUserID(ssn.UserID).
		SetToken(ssn.Token).
		SetCreatedAt(ssn.CreatedAt).
		SetExpiresAt(ssn.ExpiresAt).
		Save(ctx)

	if err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) GetSession(ctx *model.Context, ssnID string, now int64) (*ent.Session, error) {
	ssn, err := x.client.Session.Query().Where(session.ID(ssnID)).WithLogin().First(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, goerr.Wrap(err)
	}

	return ssn, nil
}

func (x *Client) DeleteSession(ctx *model.Context, ssnID string) error {
	if err := x.client.Session.DeleteOneID(ssnID).Exec(ctx); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}
