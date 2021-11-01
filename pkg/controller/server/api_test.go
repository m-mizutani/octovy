package server_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/usecase"
	"github.com/stretchr/testify/require"
)

func newServer(t *testing.T) *gin.Engine {
	uc := usecase.NewTest(t)
	require.NoError(t, uc.Init())
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
