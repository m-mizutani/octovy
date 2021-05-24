package aws

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/backend/pkg/domain/interfaces"
)

func NewS3(region string) (interfaces.S3Client, error) {
	ssn, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, goerr.Wrap(err).With("region", region)
	}

	return s3.New(ssn), nil
}

type MockS3 struct {
	Region   string
	Objects  map[string]map[string][]byte
	GetInput []*s3.GetObjectInput
}

func (x *MockS3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	x.GetInput = append(x.GetInput, input)

	bucket, ok := x.Objects[*input.Bucket]
	if !ok {
		return nil, goerr.New(s3.ErrCodeNoSuchBucket)
	}
	object, ok := bucket[*input.Key]
	if !ok {
		return nil, goerr.New(s3.ErrCodeNoSuchKey)
	}

	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader(object)),
	}, nil
}

func (x *MockS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	bucket, ok := x.Objects[*input.Bucket]
	if !ok {
		bucket = make(map[string][]byte)
		x.Objects[*input.Bucket] = bucket
	}
	raw, err := ioutil.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	bucket[*input.Key] = raw

	return &s3.PutObjectOutput{}, nil
}

func NewMockS3() (interfaces.NewS3, *MockS3) {
	mock := &MockS3{
		Objects: make(map[string]map[string][]byte),
	}
	return func(region string) (interfaces.S3Client, error) {
		mock.Region = region
		return mock, nil
	}, mock
}
