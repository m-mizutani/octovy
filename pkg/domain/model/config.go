package model

type Config struct {
	FrontendURL string
	WebhookOnly bool

	ServerAddr string
	ServerPort int

	GitHubWebhookSecret string `zlog:"secret"`

	SentryDSN string
	SentryEnv string
}
