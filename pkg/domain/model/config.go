package model

type Config struct {
	FrontendURL string

	DisableFrontend      bool
	DisableWebhookGitHub bool
	DisableWebhookTrivy  bool

	ServerAddr string
	ServerPort int

	GitHubWebhookSecret string `zlog:"secret"`

	SentryDSN string
	SentryEnv string
}
