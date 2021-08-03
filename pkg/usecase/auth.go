package usecase

import (
	"github.com/gofrs/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

const authStateTimeoutSecond = 600
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
	if len(state) < 32 {
		return nil, goerr.Wrap(model.ErrAuthenticationFailed, "Auth state is empty or not enough length")
	}

	now := x.svc.Infra.Utils.TimeNow().UTC().Unix()
	logger.With("now", now).With("state", state[:4]).Info("Looking up state")
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

	authClient := x.svc.Infra.NewGitHubAuth(x.config.GitHubEndpoint, x.config.GitHubWebURL)

	token, err := authClient.Authenticate(secrets.GitHubClientID, secrets.GitHubClientSecret, code)
	if err != nil {
		return nil, err
	}

	user, err := authClient.GetUser()
	if err != nil {
		return nil, err
	}
	token.UserID = user.UserID

	/*
		installations, err := authClient.GetInstallations()
		if err != nil {
			return nil, err
		}

		userPerm := &model.UserPermissions{
			UserID:      user.UserID,
			Permissions: make(map[string][]*model.RepoPermission),
		}

		for _, install := range installations {
			if install.ID == nil {
				logger.With("install", install).Error("No installID in Installation")
				continue
			}

			repositories, err := authClient.GetInstalledRepositories(*install.ID)
			if err != nil {
				return nil, err
			}

			for _, repo := range repositories {
				if repo.Owner == nil ||
					repo.Owner.Login == nil ||
					repo.Name == nil ||
					repo.Permissions == nil {
					logger.With("repo", repo).Error("Invalid repositroy data")
					continue
				}

				perm := &model.RepoPermission{
					GitHubRepo: model.GitHubRepo{
						Owner:    *repo.Owner.Login,
						RepoName: *repo.Name,
					},
					Permissions: *repo.Permissions,
				}
				userPerm.Permissions[perm.Owner] = append(userPerm.Permissions[perm.Owner], perm)
			}
		}
	*/

	if err := x.svc.DB().PutGitHubToken(token); err != nil {
		return nil, goerr.Wrap(err)
	}
	if err := x.svc.DB().PutUser(user); err != nil {
		return nil, goerr.Wrap(err)
	}
	/*
		if err := x.svc.DB().PutUserPermissions(userPerm); err != nil {
			return nil, goerr.Wrap(err)
		}
	*/
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

func (x *Default) RevokeSession(token string) error {
	return x.svc.DB().DeleteSession(token)
}
