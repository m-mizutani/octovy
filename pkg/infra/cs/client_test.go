package cs_test

import (
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/infra/cs"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func TestCloudStorage(t *testing.T) {
	bucket := utils.LoadEnv(t, "TEST_CLOUD_STORAGE_BUCKET")

	t.Run("Put and Get", func(t *testing.T) {
		client, err := cs.New(context.Background(), bucket)
		gt.NoError(t, err)

		key := "test-key/" + uuid.NewString() + ".txt"
		r := strings.NewReader("blue")

		gt.NoError(t, client.Put(context.Background(), key, ioutil.NopCloser(r)))

		r2, err := client.Get(context.Background(), key)
		gt.NoError(t, err)
		defer r2.Close()

		data, err := io.ReadAll(r2)
		gt.NoError(t, err)
		gt.Equal(t, "blue", string(data))
	})
}
