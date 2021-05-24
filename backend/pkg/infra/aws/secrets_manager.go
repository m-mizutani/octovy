package aws

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/m-mizutani/golambda"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
)

func NewSecretsManager(region string) (interfaces.SecretsManagerClient, error) {
	ssn, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, golambda.WrapError(err).With("region", region)
	}

	return secretsmanager.New(ssn), nil
}

type MockSecretsManager struct {
	InputLog []*secretsmanager.GetSecretValueInput
	OutData  map[string]interface{}
}

func (x *MockSecretsManager) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	v, ok := x.OutData[*input.SecretId]
	if !ok {
		return nil, golambda.NewError("Secret not found")
	}

	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(string(raw)),
	}, nil
}

func NewMockSecretsManagerSet() (interfaces.NewSecretManager, *MockSecretsManager) {
	client := &MockSecretsManager{
		OutData: make(map[string]interface{}),
	}
	return func(region string) (interfaces.SecretsManagerClient, error) {
		return client, nil
	}, client
}
