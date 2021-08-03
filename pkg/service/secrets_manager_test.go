package service_test

import (
	"encoding/base64"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/aws"
	"github.com/m-mizutani/octovy/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSecrets(t *testing.T) {
	testData := `-----BEGIN TEST PRIVATE KEY-----
MIIEpAIBAAKCAQEAw0lXZsFWxvuazTk9lq6+V2xACoYB0L07GXJJozPhobeu+QNl
d1ep3G0Q4l/96zDQDiTJ6MjS1QmPAfgZ5wfNDWIFMae6W6EgkBnTWg==
-----END TEST PRIVATE KEY-----`

	encodedSecret := base64.StdEncoding.EncodeToString([]byte(testData))
	secretsARN := "arn:aws:secretsmanager:us-east-0:123456789012:secret:testing/blue-jiObOV"
	newSM, mockSM := aws.NewMockSecretsManagerSet()
	mockSM.OutData[secretsARN] = map[string]string{
		"github_app_private_key": encodedSecret,
	}
	cfg := &model.Config{
		SecretsARN: secretsARN,
	}
	svc := service.New(cfg)
	svc.Infra.NewSecretManager = newSM

	secrets, err := svc.GetSecrets()
	require.NoError(t, err)
	assert.Equal(t, secrets.GitHubAppPrivateKey, encodedSecret)
	decodedSecret, err := base64.StdEncoding.DecodeString(secrets.GitHubAppPrivateKey)
	require.NoError(t, err)
	assert.Equal(t, testData, string(decodedSecret))
}
