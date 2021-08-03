package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
)

func NewSQS(region string) (interfaces.SQSClient, error) {
	ssn, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, goerr.Wrap(err).With("region", region)
	}

	return sqs.New(ssn), nil
}

type MockSQS struct {
	Region string
	Input  []*sqs.SendMessageInput
}

func (x *MockSQS) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	x.Input = append(x.Input, input)
	return &sqs.SendMessageOutput{
		MessageId: aws.String(uuid.New().String()),
	}, nil
}

func NewMockSQSSet() (interfaces.NewSQS, *MockSQS) {
	mock := &MockSQS{}
	return func(region string) (interfaces.SQSClient, error) {
		mock.Region = region
		return mock, nil
	}, mock
}
