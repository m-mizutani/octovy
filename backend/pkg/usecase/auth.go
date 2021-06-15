package usecase

import (
	"github.com/gofrs/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

const authStateTimeoutSecond = 60
const sessionTokenTimeoutSecond = 24 * 60 * 60 * 7
const sessionTokenLength = 128

func (x *Default) CreateAuthState() (string, error) {
	v, err := uuid.NewV4()
	if err != nil {
		return "", goerr.Wrap(err)
	}
	state := v.String()

	now := x.svc.Infra.Utils.TimeNow().UTC().Unix()
	if err := x.svc.DB().SaveAuthState(state, now+authStateTimeoutSecond); err != nil {
		return "", err
	}

	return state, nil
}

func (x *Default) LookupUser(userID string) (*model.User, error) {
	user, err := x.svc.DB().GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, goerr.Wrap(model.ErrUserNotFound)
	}

	return user, nil
}

func (x *Default) GetGitHubAppClientID() (string, error) {
	secrets, err := x.svc.GetSecrets()
	if err != nil {
		return "", err
	}

	return secrets.GitHubClientID, nil
}

func (x *Default) AuthGitHubUser(code, state string) (*model.User, error) {
	if state == "" {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Auth state is empty")
	}

	now := x.svc.Infra.Utils.TimeNow().UTC().Unix()
	found, err := x.svc.DB().HasAuthState(state, now)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if !found {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Auth state is not found")
	}

	secrets, err := x.svc.GetSecrets()
	if err != nil {
		return nil, err
	}

	authClient := x.svc.Infra.NewGitHubAuth(secrets.GitHubClientID, secrets.GitHubClientSecret, x.config.GitHubEndpoint, x.config.GitHubWebURL)
	user, token, err := authClient.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	if err := x.svc.DB().PutGitHubToken(token); err != nil {
		return nil, goerr.Wrap(err)
	}
	if err := x.svc.DB().PutUser(user); err != nil {
		return nil, goerr.Wrap(err)
	}

	return user, nil
}

func (x *Default) CreateSession(user *model.User) (*model.Session, error) {
	token := x.svc.Infra.Utils.GenerateToken(sessionTokenLength)
	now := x.svc.Infra.Utils.TimeNow()
	ssn := &model.Session{
		UserID:    user.UserID,
		Token:     token,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Unix() + sessionTokenTimeoutSecond,
	}

	if err := x.svc.DB().PutSession(ssn); err != nil {
		return nil, err
	}

	return ssn, nil
}

func (x *Default) ValidateSession(token string) (*model.Session, error) {
	now := x.svc.Infra.Utils.TimeNow()
	ssn, err := x.svc.DB().GetSession(token, now.Unix())
	if err != nil {
		return nil, err
	}
	if ssn == nil {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Invalid token")
	}

	return ssn, nil
}
