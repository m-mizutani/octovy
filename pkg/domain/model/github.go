package model

type GitHubRepo struct {
	Owner    string `json:"owner"`
	RepoName string `json:"name"`
}

type GitHubBranch struct {
	GitHubRepo
	Branch string `json:"branch"`
}

type GitHubToken struct {
	UserID                string
	AccessToken           string `json:"access_token"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
}

type GitHubCheckResult struct {
	Conclusion string   `json:"conclusion"`
	Messages   []string `json:"messages"`
}
