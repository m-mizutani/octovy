package githubauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

var logger = golambda.Logger

type GitHubAuthClient struct {
	APIEndpoint string
	WebEndpoint string
	httpClient  *http.Client
}

func (x *GitHubAuthClient) apiURL(path string, values url.Values) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	apiURL := strings.TrimRight(x.APIEndpoint, "/") + path
	if values != nil {
		apiURL += "?" + values.Encode()
	}
	return apiURL
}

func New(apiEndpoint, webEndpoint string) interfaces.GitHubAuth {
	if apiEndpoint == "" {
		apiEndpoint = "https://api.github.com"
	}
	if webEndpoint == "" {
		webEndpoint = "https://github.com"
	}

	return &GitHubAuthClient{
		APIEndpoint: apiEndpoint,
		WebEndpoint: webEndpoint,
	}
}

type authTransport struct {
	token *model.GitHubToken
}

func (x *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "token "+x.token.AccessToken)
	client := &http.Client{}
	return client.Do(req)
}

func (x *GitHubAuthClient) SetToken(token *model.GitHubToken) {
	x.httpClient = &http.Client{
		Transport: &authTransport{
			token: token,
		},
	}
}

func (x *GitHubAuthClient) Authenticate(clientID, clientSecret, code string) (*model.GitHubToken, error) {
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

	webURL := strings.TrimRight(x.WebEndpoint, "/") + "/login/oauth/access_token"
	req, err := http.NewRequest("POST", webURL, bytes.NewReader(authReqBody))
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
	// if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {

	if err := json.Unmarshal(body, &token); err != nil {
		return nil, goerr.Wrap(err, "Failed to parse GitHub access token").With("url", webURL)
	}

	x.SetToken(&token)

	return &token, nil
}

func (x *GitHubAuthClient) GetUser() (*model.User, error) {
	if x.httpClient == nil {
		return nil, goerr.Wrap(model.ErrNoAuthenticatedClient)
	}

	url := x.apiURL("/user", nil)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := x.httpClient.Do(req)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, goerr.Wrap(err, "Failed to get user info").With("body", string(body)).With("code", resp.StatusCode)
	}

	var user github.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, goerr.Wrap(err, "Failed to parse github user").With("body", string(body)).With("url", url)
	}

	if user.ID == nil {
		return nil, goerr.New("No GitHub user ID").With("user", user)
	}

	str := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	return &model.User{
		UserID:    fmt.Sprintf("%d", *user.ID),
		Login:     str(user.Login),
		Name:      str(user.Name),
		AvatarURL: str(user.AvatarURL),
		URL:       str(user.URL),
	}, nil
}

const pagenationLimit = 100

func (x *GitHubAuthClient) GetInstallations() ([]*github.Installation, error) {
	if x.httpClient == nil {
		return nil, goerr.Wrap(model.ErrNoAuthenticatedClient)
	}

	var result []*github.Installation
	total := 0

	for page := 1; (total == 0 || len(result) < total) && page < pagenationLimit; page++ {
		values := url.Values{}
		values.Set("page", fmt.Sprintf("%d", page))
		values.Set("per_page", "100")

		apiURL := x.apiURL("/user/installations", values)
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

		resp, err := x.httpClient.Do(req)
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, goerr.Wrap(err, "Failed to installations").With("url", apiURL).With("body", string(body)).With("code", resp.StatusCode)
		}

		var listInstallations struct {
			TotalCount    int64                  `json:"total_count"`
			Installations []*github.Installation `json:"installations"`
		}
		if err := json.Unmarshal(body, &listInstallations); err != nil {
			return nil, goerr.Wrap(err, "Failed to parse installations").With("body", string(body)).With("url", apiURL)
		}

		total = int(listInstallations.TotalCount)
		result = append(result, listInstallations.Installations...)
		logger.With("result.length", len(result)).With("total", total).Info("recv installation")
		if total == 0 {
			break
		}
	}

	return result, nil
}

func (x *GitHubAuthClient) GetInstalledRepositories(installID int64) ([]*github.Repository, error) {
	if x.httpClient == nil {
		return nil, goerr.Wrap(model.ErrNoAuthenticatedClient)
	}

	var result []*github.Repository
	total := 0

	for page := 1; (total == 0 || len(result) < total) && page < pagenationLimit; page++ {
		values := url.Values{}
		values.Set("page", fmt.Sprintf("%d", page))
		values.Set("per_page", "100")

		apiURL := x.apiURL(fmt.Sprintf("/user/installations/%d/repositories", installID), values)
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

		resp, err := x.httpClient.Do(req)
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, goerr.Wrap(err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, goerr.Wrap(err, "Failed to get repositories info").With("body", string(body)).With("code", resp.StatusCode)
		}

		var listRepo struct {
			TotalCount   int64                `json:"total_count"`
			Repositories []*github.Repository `json:"repositories"`
		}
		if err := json.Unmarshal(body, &listRepo); err != nil {
			return nil, goerr.Wrap(err, "Failed to parse repositories").With("body", string(body)).With("url", apiURL)
		}

		total = int(listRepo.TotalCount)
		result = append(result, listRepo.Repositories...)
		logger.With("result.length", len(result)).With("total", total).Info("recv repos")
	}

	return result, nil
}
