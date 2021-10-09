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

	TrivyPath string

	SentryDSN string
	SentryEnv string
}

func (x *Config) CopyWithoutSensitives() *Config {
	copiedConfig := *x
	// Removing sensitive data
	if copiedConfig.GitHubAppPrivateKey != "" {
		copiedConfig.GitHubAppPrivateKey = "[Removed]"
	}
	if copiedConfig.GitHubAppSecret != "" {
		copiedConfig.GitHubAppSecret = "[Removed]"
	}
	if copiedConfig.DBConfig != "" {
		copiedConfig.DBConfig = "[Removed]"
	}
	if copiedConfig.GitHubWebhookSecret != "" {
		copiedConfig.GitHubWebhookSecret = "[Removed]"
	}

	return &copiedConfig
}
