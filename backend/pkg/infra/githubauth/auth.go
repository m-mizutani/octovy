package githubauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

type GitHubAuthClient struct {
	ClientID     string
	ClientSecret string
	APIEndpoint  string
	WebEndpoint  string
}

func New(clientID, clientSecret, apiEndpoint, webEndpoint string) interfaces.GitHubAuth {
	if clientID == "" {
		golambda.EmitError(goerr.New("clientID is empty"))
		panic("clientID is empty")
	}
	if clientSecret == "" {
		golambda.EmitError(goerr.New("clientSecret is empty"))
		panic("clientSecret is empty")
	}
	if apiEndpoint == "" {
		apiEndpoint = "https://api.github.com"
	}
	if webEndpoint == "" {
		webEndpoint = "https://github.com"
	}

	return &GitHubAuthClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		APIEndpoint:  apiEndpoint,
		WebEndpoint:  webEndpoint,
	}
}

func (x *GitHubAuthClient) GetAccessToken(code string) (*model.User, *model.GitHubToken, error) {
	authReq := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}{
		ClientID:     x.ClientID,
		ClientSecret: x.ClientSecret,
		Code:         code,
	}
	authReqBody, err := json.Marshal(authReq)
	if err != nil {
		return nil, nil, goerr.Wrap(err, "Failed to encode authReq")
	}
	fmt.Println("auth", string(authReqBody))
	url := fmt.Sprintf("%s/login/oauth/access_token", strings.TrimSuffix(x.WebEndpoint, "/"))
	req, err := http.NewRequest("POST", url, bytes.NewReader(authReqBody))
	if err != nil {
		return nil, nil, goerr.Wrap(err, "Failed to create a new auth request").With("url", url)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, goerr.Wrap(err, "Failed to post access_token").With("url", url)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, goerr.Wrap(err, "Failed to post access_token").With("body", string(body)).With("code", resp.StatusCode).With("url", url)
	}

	var token model.GitHubToken
	// if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
	fmt.Println("body", string(body))
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, nil, goerr.Wrap(err, "Failed to parse GitHub access token").With("url", url)
	}

	apiURL := strings.TrimSuffix(x.APIEndpoint, "/") + "/user"
	userReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, nil, goerr.Wrap(err)
	}
	userReq.Header.Add("Authorization", "token "+token.AccessToken)
	req.Header.Add("Accept", "application/json")

	userResp, err := client.Do(userReq)
	if err != nil {
		return nil, nil, goerr.Wrap(err)
	}

	body, err = ioutil.ReadAll(userResp.Body)
	if err != nil {
		return nil, nil, goerr.Wrap(err)
	}
	if userResp.StatusCode != http.StatusOK {
		return nil, nil, goerr.Wrap(err, "Failed to get user info").With("body", string(body)).With("code", userResp.StatusCode)
	}

	var user github.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, nil, goerr.Wrap(err, "Failed to parse github user").With("body", string(body)).With("url", apiURL)
	}

	str := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	if user.ID == nil {
		return nil, nil, goerr.New("No GitHub user ID").With("user", user)
	}
	token.UserID = fmt.Sprintf("%d", *user.ID)
	return &model.User{
		UserID:    token.UserID,
		Login:     str(user.Login),
		Name:      str(user.Name),
		AvatarURL: str(user.AvatarURL),
		URL:       str(user.URL),
	}, &token, nil
}
