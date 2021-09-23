package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/m-mizutani/octovy/pkg/utils"
)

const authStateTimeoutSecond = 600
const sessionTokenTimeoutSecond = 24 * 60 * 60 * 7
const sessionTokenLength = 128

func (x *usecase) CreateAuthState(ctx context.Context) (string, error) {
	state := utils.GenerateToken(128)

	now := x.infra.Utils.Now().Unix()
	if err := x.infra.DB.SaveAuthState(ctx, state, now+authStateTimeoutSecond); err != nil {
		return "", err
	}

	return state, nil
}

func (x *usecase) LookupUser(ctx context.Context, userID int) (*ent.User, error) {
	user, err := x.infra.DB.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, goerr.Wrap(model.ErrUserNotFound)
	}

	return user, nil
}

func (x *usecase) AuthGitHubUser(ctx context.Context, code, state string) (*ent.User, error) {
	if len(state) < 32 {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Auth state is empty or not enough length")
	}

	now := x.infra.Utils.Now().Unix()
	logger.Debug().Int64("now", now).Str("state", state[:4]).Msg("Looking up state")
	found, err := x.infra.DB.HasAuthState(ctx, state, now)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if !found {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Auth state is not found")
	}

	token, err := x.infra.GitHub.Authenticate(ctx, x.config.GitHubAppClientID, x.config.GitHubAppSecret, code)
	if err != nil {
		return nil, err
	}

	githubUser, err := x.infra.GitHub.GetUser(ctx, token)
	if err != nil {
		return nil, err
	}
	_, err = x.infra.GitHub.GetOrgs(ctx, token)
	if err != nil {
		return nil, err
	}

	if githubUser.ID == nil {
		return nil, goerr.Wrap(model.ErrGitHubAPI, "user.ID is null")
	}
	if githubUser.Login == nil {
		return nil, goerr.Wrap(model.ErrGitHubAPI, "user.Login is null")
	}
	if githubUser.Name == nil {
		return nil, goerr.Wrap(model.ErrGitHubAPI, "user.Name is null")
	}
	if githubUser.URL == nil {
		return nil, goerr.Wrap(model.ErrGitHubAPI, "user.URL is null")
	}
	if githubUser.AvatarURL == nil {
		return nil, goerr.Wrap(model.ErrGitHubAPI, "user.AvatarURL is null")
	}
	user := &ent.User{
		GithubID:  *githubUser.ID,
		Login:     *githubUser.Login,
		Name:      *githubUser.Name,
		URL:       *githubUser.URL,
		AvatarURL: *githubUser.AvatarURL,
	}

	userID, err := x.infra.DB.PutUser(ctx, user)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	user.ID = userID
	return user, nil
}

func (x *usecase) CreateSession(ctx context.Context, user *ent.User) (*ent.Session, error) {
	token := utils.GenerateToken(sessionTokenLength)
	now := x.infra.Utils.Now()
	ssn := &ent.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Unix() + sessionTokenTimeoutSecond,
	}

	if err := x.infra.DB.PutSession(ctx, ssn); err != nil {
		return nil, err
	}

	return ssn, nil
}

func (x *usecase) ValidateSession(ctx context.Context, ssnID string) (*ent.Session, error) {
	ssn, err := x.infra.DB.GetSession(ctx, ssnID, x.infra.Utils.Now().Unix())
	if err != nil {
		return nil, err
	}
	if ssn == nil {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Invalid token")
	}

	return ssn, nil
}

func (x *usecase) RevokeSession(ctx context.Context, ssnID string) error {
	return x.infra.DB.DeleteSession(ctx, ssnID)
}
