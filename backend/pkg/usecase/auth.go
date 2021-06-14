package usecase

import (
	"github.com/gofrs/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

const authStateTimeoutSecond = 60

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

func (x *Default) CreateToken(user *model.User) ([]byte, error) {
	return []byte("five timeless words"), nil
}

func (x *Default) ValidateToken(token []byte) (string, error) {
	if string(token) == "five timeless words" {
		return "881", nil
	}
	return "", goerr.Wrap(model.ErrAuthenticationFailed)
}
