package service

import "github.com/Netflix/go-env"

type Config struct {
	AwsRegion        string `env:"AWS_REGION"`
	TableName        string `env:"TABLE_NAME"`
	SecretsARN       string `env:"SECRETS_ARN"`
	ScanRequestQueue string `env:"SCAN_REQUEST_QUEUE"`
	GitHubEndpoint   string `env:"GITHUB_ENDPOINT"`
}

func NewConfig() *Config {
	var config Config
	if _, err := env.UnmarshalFromEnviron(&config); err != nil {
		panic("Failed UnmarshalFromEnviron to Config: " + err.Error())
	}
	return &config
}
