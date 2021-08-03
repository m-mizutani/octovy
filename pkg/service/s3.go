package service

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/m-mizutani/goerr"
)

const trivyDBName = "trivy.db.gz"

func trivyDBObjectKey(prefix string) string {
	return fmt.Sprintf("%sdb/%s", prefix, trivyDBName)
}

func (x *Service) UploadTrivyDB(rs io.ReadSeeker) error {
	s3Client, err := x.Infra.NewS3(x.config.S3Region)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:          aws.String(x.config.S3Bucket),
		Key:             aws.String(trivyDBObjectKey(x.config.S3Prefix)),
		Body:            rs,
		ContentType:     aws.String("application/x-gzip"),
		ContentEncoding: aws.String("gzip"),
	}

	if _, err := s3Client.PutObject(input); err != nil {
		return goerr.Wrap(err).With("input", input)
	}

	return nil
}

func (x *Service) downloadTrivyDB(w io.Writer) error {
	s3Client, err := x.Infra.NewS3(x.config.S3Region)
	if err != nil {
		return err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(x.config.S3Bucket),
		Key:    aws.String(trivyDBObjectKey(x.config.S3Prefix)),
	}

	output, err := s3Client.GetObject(input)
	if err != nil {
		return goerr.Wrap(err).With("input", input)
	}

	if _, err := io.Copy(w, output.Body); err != nil {
		return goerr.Wrap(err).With("output", output)
	}

	return nil
}
