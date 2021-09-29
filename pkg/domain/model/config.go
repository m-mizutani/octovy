package model

type Config struct {
	DBType      string
	DBConfig    string
	FrontendURL string
	WebhookOnly bool

	ServerAddr string
	ServerPort int

	GitHubAppID         int64
	GitHubAppPrivateKey string
	GitHubAppClientID   string
	GitHubAppSecret     string
	GitHubWebhookSecret string

	TrivyDBPath string

	SentryDSN string
	SentryEnv string
}
