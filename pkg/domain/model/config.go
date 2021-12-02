package model

type Config struct {
	DBType      string
	DBConfig    string `zlog:"secret"`
	FrontendURL string
	WebhookOnly bool

	ServerAddr string
	ServerPort int

	GitHubAppID         int64
	GitHubAppPrivateKey string `zlog:"secret"`
	GitHubAppClientID   string
	GitHubAppSecret     string `zlog:"secret"`
	GitHubWebhookSecret string `zlog:"secret"`

	CheckPolicyData string

	OPAServerURL    string
	OPAUseGoogleIAP bool

	TrivyPath string

	SentryDSN string
	SentryEnv string
}
