package model

import "github.com/Netflix/go-env"

type Config struct {
	AwsRegion            string `env:"AWS_REGION"`
	TableName            string `env:"TABLE_NAME"`
	SecretsARN           string `env:"SECRETS_ARN"`
	ScanRequestQueue     string `env:"SCAN_REQUEST_QUEUE"`
	FeedbackRequestQueue string `env:"FEEDBACK_REQUEST_QUEUE"`
	GitHubEndpoint       string `env:"GITHUB_ENDPOINT"`
	FrontendURL          string `env:"FRONTEND_URL"`

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
