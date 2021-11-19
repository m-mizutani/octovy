package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoLabel(t *testing.T) {
	ctx := model.NewContext()
	dbClient := db.NewMock(t)
	repo, err := dbClient.CreateRepo(ctx, &ent.Repository{
		Owner: "blue",
		Name:  "five",
	})
	require.NoError(t, err)
	engine := newServerWithDB(t, dbClient)

	{
		w := httptest.NewRecorder()
		req := newRequest("GET", "/api/v1/repo-label", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp []*ent.RepoLabel
		bind(t, w.Body, &resp)
		assert.Len(t, resp, 0)
	}

	{
		w := httptest.NewRecorder()
		req := newRequest("POST", "/api/v1/repo-label", model.RequestRepoLabel{
			Name: "magic",
		})
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	}

	var label *ent.RepoLabel
	{
		w := httptest.NewRecorder()
		req := newRequest("GET", "/api/v1/repo-label", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp []*ent.RepoLabel
		bind(t, w.Body, &resp)
		require.Len(t, resp, 1)
		assert.Equal(t, "magic", resp[0].Name)
		label = resp[0]
	}

	{
		w := httptest.NewRecorder()
		req := newRequest("POST", fmt.Sprintf("/api/v1/repo-label/%d/assign/%d", label.ID, repo.ID), nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	}

	{
		w := httptest.NewRecorder()
		req := newRequest("GET", "/api/v1/repository", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp []*ent.Repository
		bind(t, w.Body, &resp)
		require.Len(t, resp, 1)
		require.Len(t, resp[0].Edges.Labels, 1)
		assert.Equal(t, "magic", resp[0].Edges.Labels[0].Name)
	}

	{
		w := httptest.NewRecorder()
		req := newRequest("DELETE", fmt.Sprintf("/api/v1/repo-label/%d/assign/%d", label.ID, repo.ID), nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	}

	{
		w := httptest.NewRecorder()
		req := newRequest("GET", "/api/v1/repository", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp []*ent.Repository
		bind(t, w.Body, &resp)
		require.Len(t, resp, 1)
		require.Len(t, resp[0].Edges.Labels, 0)
	}
}
