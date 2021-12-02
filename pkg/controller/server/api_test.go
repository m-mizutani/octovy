package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/require"
)

func newServer(t *testing.T) *gin.Engine {
	uc := usecase.NewTest(t)
	engine := server.New(uc, &server.Option{DisableAuth: true})
	return engine
}

func newServerWithDB(t *testing.T, client *db.Client) *gin.Engine {
	uc := usecase.NewTest(t, usecase.OptInjectDB(client))
	engine := server.New(uc, &server.Option{DisableAuth: true})
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
