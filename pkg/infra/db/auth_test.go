package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	t.Run("create and lookup user", func(t *testing.T) {
		client := newTestTable(t)
		user := &model.User{
			UserID:    "123",
			Login:     "mizutani",
			Name:      "mizutani",
			AvatarURL: "https://avatars.githubusercontent.com/u/605953?s=60&v=4",
			URL:       "https://github.com/m-mizutani",
		}

		require.NoError(t, client.PutUser(user))

		t.Run("Found created user", func(t *testing.T) {
			u1, err := client.GetUser("123")
			require.NoError(t, err)
			assert.Equal(t, user, u1)
		})

		t.Run("Not found not created user", func(t *testing.T) {
			u1, err := client.GetUser("999")
			require.NoError(t, err)
			assert.Nil(t, u1)
		})
	})
}
