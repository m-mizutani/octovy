package model

type Config struct {
	DBType      string
	DBConfig    string
	FrontendURL string

	GitHubAppID    int64
	GitHubClientID int64

	ServerAddr string
	ServerPort int
}
