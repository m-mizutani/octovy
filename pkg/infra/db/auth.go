package db

import (
	"github.com/guregu/dynamo"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

func authStatePK(state string) string {
	return "github_auth_state:" + state
}
func authStateSK() string {
	return "*"
}

func (x *DynamoClient) SaveAuthState(state string, expiresAt int64) error {
	if expiresAt == 0 {
		return goerr.Wrap(model.ErrInvalidValue, "expiresAt must be > 0")
	}
	record := dynamoRecord{
		PK:        authStatePK(state),
		SK:        authStateSK(),
		ExpiresAt: &expiresAt,
	}
	if err := x.table.Put(record).Run(); err != nil {
		return err
	}

	return nil
}

func (x *DynamoClient) HasAuthState(state string, now int64) (bool, error) {
	var record dynamoRecord
	pk := authStatePK(state)
	sk := authStateSK()

	q := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).Filter("? < expires_at", now)
	if err := q.One(&record); err != nil {
		if isNotFoundErr(err) {
			return false, nil
		}
		return false, goerr.Wrap(err)
	}

	return true, nil
}

func userPK(userID string) string {
	return "user:" + userID
}
func userSK() string {
	return "*"
}

func (x *DynamoClient) PutUser(user *model.User) error {
	record := dynamoRecord{
		PK:  userPK(user.UserID),
		SK:  userSK(),
		Doc: user,
	}
	if err := x.table.Put(record).Run(); err != nil {
		return err
	}

	return nil
}

func (x *DynamoClient) GetUser(userID string) (*model.User, error) {
	if userID == "" {
		return nil, goerr.Wrap(model.ErrInvalidValue, "userID must not be empty")
	}
	var record dynamoRecord

	pk := userPK(userID)
	sk := userSK()

	q := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk)
	if err := q.One(&record); err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, goerr.Wrap(err)
	}

	var user *model.User
	if err := record.Unmarshal(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func userPermissionsPK(userID string) string {
	return "user_permissions:" + userID
}
func userPermissionsSK() string {
	return "*"
}

func (x *DynamoClient) PutUserPermissions(perm *model.UserPermissions) error {
	record := dynamoRecord{
		PK:  userPermissionsPK(perm.UserID),
		SK:  userPermissionsSK(),
		Doc: perm,
	}
	if err := x.table.Put(record).Run(); err != nil {
		return err
	}

	return nil

}
func (x *DynamoClient) GetUserPermissions(userID string) (*model.UserPermissions, error) {
	if userID == "" {
		return nil, goerr.Wrap(model.ErrInvalidValue, "userID must not be empty")
	}
	var record dynamoRecord

	pk := userPermissionsPK(userID)
	sk := userPermissionsSK()

	q := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk)
	if err := q.One(&record); err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, goerr.Wrap(err)
	}

	var perm *model.UserPermissions
	if err := record.Unmarshal(&perm); err != nil {
		return nil, err
	}

	return perm, nil
}

func gitHubTokenPK(userID string) string {
	return "github_token:" + userID
}
func gitHubTokenSK() string {
	return "*"
}

func (x *DynamoClient) PutGitHubToken(token *model.GitHubToken) error {
	record := dynamoRecord{
		PK:  gitHubTokenPK(token.UserID),
		SK:  gitHubTokenSK(),
		Doc: token,
	}
	if err := x.table.Put(record).Run(); err != nil {
		return err
	}

	return nil

}

func (x *DynamoClient) GetGitHubToken(userID string) (*model.GitHubToken, error) {
	if userID == "" {
		return nil, goerr.Wrap(model.ErrInvalidValue, "userID must not be empty")
	}
	var record dynamoRecord

	pk := gitHubTokenPK(userID)
	sk := gitHubTokenSK()

	q := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk)
	if err := q.One(&record); err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, goerr.Wrap(err)
	}

	var token *model.GitHubToken
	if err := record.Unmarshal(&token); err != nil {
		return nil, err
	}

	return token, nil
}

func sessionPK(token string) string {
	return "session:" + token
}
func sessionSK() string {
	return "*"
}

func (x *DynamoClient) PutSession(ssn *model.Session) error {
	if err := ssn.IsValid(); err != nil {
		return err
	}

	record := dynamoRecord{
		PK:        sessionPK(ssn.Token),
		SK:        sessionSK(),
		Doc:       ssn,
		ExpiresAt: &ssn.ExpiresAt,
	}
	if err := x.table.Put(record).Run(); err != nil {
		return err
	}

	return nil

}

func (x *DynamoClient) GetSession(token string, now int64) (*model.Session, error) {
	if token == "" {
		return nil, goerr.Wrap(model.ErrInvalidValue, "token must not be empty")
	}
	var record dynamoRecord

	pk := sessionPK(token)
	sk := sessionSK()

	q := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).Filter("? < expires_at", now)
	if err := q.One(&record); err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, goerr.Wrap(err)
	}

	var ssn *model.Session
	if err := record.Unmarshal(&ssn); err != nil {
		return nil, err
	}

	return ssn, nil
}

func (x DynamoClient) DeleteSession(token string) error {
	if token == "" {
		return goerr.Wrap(model.ErrInvalidValue, "token must not be empty")
	}

	pk := sessionPK(token)
	sk := sessionSK()

	if err := x.table.Delete("pk", pk).Range("sk", sk).Run(); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
