package aws_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/m-mizutani/octovy/pkg/infra/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3Mock(t *testing.T) {
	bucketName := "blue"
	key := "path/to/doom"
	_, mock := aws.NewMockS3()
	_, err := mock.PutObject(&s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   strings.NewReader("boom!"),
	})
	require.NoError(t, err)

	out, err := mock.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	require.NoError(t, err)
	raw, err := ioutil.ReadAll(out.Body)
	require.NoError(t, err)
	assert.Equal(t, "boom!", string(raw))
}
