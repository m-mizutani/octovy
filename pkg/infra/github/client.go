package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v39/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/utils"
)

var logger = utils.Logger

// This package is used to download trivy database, not used by GitHub App.
type Interface interface {
	Authenticate(ctx context.Context, clientID, clientSecret, code string) (*model.GitHubToken, error)
	GetUser(ctx context.Context, token *model.GitHubToken) (*github.User, error)
}

type Client struct {
	client *github.Client
}

func New() *Client {
	return &Client{
		client: github.NewClient(&http.Client{}),
	}
}

func (x *Client) Authenticate(ctx context.Context, clientID, clientSecret, code string) (*model.GitHubToken, error) {
	if clientID == "" {
		return nil, goerr.Wrap(model.ErrInvalidSystemValue, "clientID is empty")
	}
	if clientSecret == "" {
		return nil, goerr.Wrap(model.ErrInvalidSystemValue, "clientSecret is empty")
	}
	if code == "" {
		return nil, goerr.Wrap(model.ErrInvalidSystemValue, "code is empty")
	}

	authReq := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Code:         code,
	}
	authReqBody, err := json.Marshal(authReq)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to encode authReq")
	}

	webURL := "https://github.com/login/oauth/access_token"
	req, err := http.NewRequestWithContext(ctx, "POST", webURL, bytes.NewReader(authReqBody))
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create a new auth request").With("url", webURL)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to post access_token").With("url", webURL)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, goerr.Wrap(err, "Failed to post access_token").With("body", string(body)).With("code", resp.StatusCode).With("url", webURL)
	}

	var token model.GitHubToken
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, goerr.Wrap(err, "Failed to parse GitHub access token").With("url", webURL)
	}

	return &token, nil
}

func (x *Client) GetUser(ctx context.Context, token *model.GitHubToken) (*github.User, error) {
	body, err := getRequest(ctx, token, "https://api.github.com/user")
	if err != nil {
		return nil, err
	}

	var user github.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, goerr.Wrap(err, "Failed to parse github user").With("body", string(body))
	}

	logger.Debug().Interface("user", user).Msg("Got github user")

	if user.ID == nil {
		return nil, goerr.New("No GitHub user ID").With("user", user)
	}

	return &user, nil
}

func getRequest(ctx context.Context, token *model.GitHubToken, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, goerr.Wrap(err, "Failed to get").With("body", string(body)).With("url", url).With("code", resp.StatusCode)
	}

	return body, nil
}
