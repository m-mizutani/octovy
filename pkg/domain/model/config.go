package model

type Config struct {
	DBType      string
	DBConfig    string
	FrontendURL string

	ServerAddr string
	ServerPort int

	GitHubAppID             int64
	GitHubAppPrivateKeyPath string
	GitHubAppClientID       string
	GitHubAppSecret         string
}
