package model

import (
	"strings"

	"github.com/Netflix/go-env"
)

type Metadata struct {
	FrontendURL  string `env:"FRONTEND_URL"`
	GitHubWebURL string `env:"GITHUB_WEB_URL"`
	HomepageURL  string `env:"HOMEPAGE_URL"`
}

type Config struct {
	AwsRegion            string `env:"AWS_REGION"`
	TableName            string `env:"TABLE_NAME"`
	SecretsARN           string `env:"SECRETS_ARN"`
	ScanRequestQueue     string `env:"SCAN_REQUEST_QUEUE"`
	FeedbackRequestQueue string `env:"FEEDBACK_REQUEST_QUEUE"`
	GitHubEndpoint       string `env:"GITHUB_ENDPOINT"`

	Metadata

	RulePullReqCommentTriggers string `env:"RULE_PR_COMMENT_TRIGGERS"`
	RuleFailCheckIfVuln        string `env:"RULE_FAIL_CHECK_IF_VULN"`

	S3Region string `env:"S3_REGION"`
	S3Bucket string `env:"S3_BUCKET"`
	S3Prefix string `env:"S3_PREFIX"`

	TrivyDBPath string
}

func NewConfig() *Config {
	var config Config
	if _, err := env.UnmarshalFromEnviron(&config); err != nil {
		panic("Failed UnmarshalFromEnviron to Config: " + err.Error())
	}
	return &config
}

// FrontendBaseURL returns frontend URL trimmed last slash
func (x *Config) FrontendBaseURL() string {
	return strings.TrimSuffix(x.FrontendURL, "/")
}

func (x *Config) ShouldCommentPR(event string) bool {
	for _, trigger := range strings.Split(x.RulePullReqCommentTriggers, "|") {
		if event == trigger {
			return true
		}
	}
	return false
}

func (x *Config) ShouldFailIfVuln() bool {
	return x.RuleFailCheckIfVuln != ""
}
