package service

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/model"
)

func (x *Service) sendSQSMessage(msg interface{}, url string) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return goerr.Wrap(model.ErrInvalidValue, err.Error())
	}

	client, err := x.Infra.NewSQS(x.config.AwsRegion)
	if err != nil {
		return goerr.Wrap(model.ErrSystem, err.Error())
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    &url,
		MessageBody: aws.String(string(raw)),
	}
	logger.With("input", input).Debug("Sending SQS")
	output, err := client.SendMessage(input)
	if err != nil {
		return goerr.Wrap(err).With("input", input)
	}
	logger.With("output", output).Debug("Sent SQS")

	return nil
}

func (x *Service) SendScanRequest(req *model.ScanRepositoryRequest) error {
	if req == nil {
		return goerr.New("req is not set")
	}
	if err := req.IsValid(); err != nil {
		return err
	}

	return x.sendSQSMessage(req, x.config.ScanRequestQueue)
}

func (x *Service) SendFeedbackRequest(req *model.FeedbackRequest) error {
	if req == nil {
		return goerr.New("req is not set")
	}
	if err := req.IsValid(); err != nil {
		return err
	}

	return x.sendSQSMessage(req, x.config.FeedbackRequestQueue)
}
