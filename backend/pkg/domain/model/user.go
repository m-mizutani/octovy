package model

import "github.com/m-mizutani/goerr"

type GitHubToken struct {
	UserID                string
	AccessToken           string `json:"access_token"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
}

type Session struct {
	UserID    string
	Token     string
	CreatedAt int64
	ExpiresAt int64
}

func (x *Session) IsValid() error {
	if x.UserID == "" {
		return goerr.Wrap(ErrInvalidInputValues, "UserID must not be empty")
	}
	if x.Token == "" {
		return goerr.Wrap(ErrInvalidInputValues, "Token must not be empty")
	}
	if x.CreatedAt <= 0 {
		return goerr.Wrap(ErrInvalidInputValues, "CreatedAt must not be > 0")
	}
	if x.ExpiresAt <= 0 {
		return goerr.Wrap(ErrInvalidInputValues, "ExpiresAt must not be > 0")
	}

	return nil
}

type User struct {
	UserID    string
	Login     string
	Name      string
	AvatarURL string
	URL       string
}
