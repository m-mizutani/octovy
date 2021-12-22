package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newServer(t *testing.T) *gin.Engine {
	uc := usecase.NewTest(t)
	engine := server.New(uc, server.DisableAuth())
	return engine
}

func newServerWithDB(t *testing.T, client *db.Client, options ...server.Option) *gin.Engine {
	uc := usecase.NewTest(t, usecase.OptInjectDB(client))
	options = append(options, server.DisableAuth())
	engine := server.New(uc, options...)
	return engine
}

func bind(t *testing.T, body *bytes.Buffer, v interface{}) {
	var resp server.BaseResponse
	require.NoError(t, json.Unmarshal(body.Bytes(), &resp))

	raw, err := json.Marshal(resp.Data)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, v))
}

func newRequest(method, url string, data interface{}) *http.Request {
	var body io.Reader
	if data != nil {
		raw, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return req
}

func TestDisable(t *testing.T) {
	t.Run("test disable webhook-github", func(t *testing.T) {
		mock := db.NewMock(t)
		s := newServerWithDB(t, mock, server.DisableWebhookGitHub())

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "http://localhost/webhook/github", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		}

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "http://localhost/webhook/trivy", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusNotFound, w.Result().StatusCode)
		}

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost/api/v1/repository", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusNotFound, w.Result().StatusCode)
		}
	})

	t.Run("test disable webhook-trivy", func(t *testing.T) {
		mock := db.NewMock(t)
		s := newServerWithDB(t, mock, server.DisableWebhookTrivy())

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "http://localhost/webhook/github", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusNotFound, w.Result().StatusCode)
		}

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "http://localhost/webhook/trivy", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		}

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost/api/v1/repository", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusNotFound, w.Result().StatusCode)
		}
	})

	t.Run("test disable frontend", func(t *testing.T) {
		mock := db.NewMock(t)
		s := newServerWithDB(t, mock, server.DisableFrontend())

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "http://localhost/webhook/github", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusNotFound, w.Result().StatusCode)
		}

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "http://localhost/webhook/trivy", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusNotFound, w.Result().StatusCode)
		}

		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost/api/v1/repository", nil)
			require.NoError(t, err)
			s.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		}
	})
}
