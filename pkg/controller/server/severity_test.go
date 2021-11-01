package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-mizutani/octovy/pkg/controller/server"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoSeverity(t *testing.T) {
	engine := newServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/severity", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	var resp []*ent.Severity
	bind(t, w.Body, &resp)
	assert.Len(t, resp, 0)
}

func TestSeverityCreate(t *testing.T) {
	engine := newServer(t)

	w1 := httptest.NewRecorder()
	body, _ := json.Marshal(server.SeverityRequest{Label: "critical"})
	req1, _ := http.NewRequest("POST", "/api/v1/severity", bytes.NewReader(body))
	engine.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Result().StatusCode)

	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(server.SeverityRequest{Label: "high"})
	req2, _ := http.NewRequest("POST", "/api/v1/severity", bytes.NewReader(body2))
	engine.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Result().StatusCode)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/severity", bytes.NewReader(body))
	engine.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusCreated, w1.Result().StatusCode)
	var resp []*ent.Severity
	bind(t, w3.Body, &resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, "critical", resp[0].Label)
	assert.Equal(t, "high", resp[1].Label)
}

func TestSeverityUpdate(t *testing.T) {
	engine := newServer(t)

	w1 := httptest.NewRecorder()
	body, _ := json.Marshal(server.SeverityRequest{Label: "critical"})
	req1, _ := http.NewRequest("POST", "/api/v1/severity", bytes.NewReader(body))
	engine.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Result().StatusCode)
	var resp1 *ent.Severity
	bind(t, w1.Body, &resp1)

	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(server.SeverityRequest{Label: "high"})
	req2, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/severity/%d", resp1.ID), bytes.NewReader(body2))
	engine.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Result().StatusCode)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/severity", bytes.NewReader(body))
	engine.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Result().StatusCode)
	var resp2 []*ent.Severity
	bind(t, w3.Body, &resp2)
	require.Len(t, resp2, 1)
	assert.Equal(t, "high", resp2[0].Label)
}

func TestSeverityAssign(t *testing.T) {
	engine := newServer(t)

	var sev *ent.Severity
	{
		w1 := httptest.NewRecorder()
		body, _ := json.Marshal(server.SeverityRequest{Label: "critical"})
		req1, _ := http.NewRequest("POST", "/api/v1/severity", bytes.NewReader(body))
		engine.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Result().StatusCode)
		bind(t, w1.Body, &sev)
	}

	{
		w2 := httptest.NewRecorder()
		body2, _ := json.Marshal(ent.Vulnerability{ID: "CVE-2000-1000", Title: "blue"})
		req2, _ := http.NewRequest("POST", "/api/v1/vulnerability", bytes.NewReader(body2))
		engine.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusCreated, w2.Result().StatusCode)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/severity/%d/assign/CVE-2000-1000", sev.ID), nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/vulnerability/CVE-2000-1000", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var vuln *model.RespVulnerability
		bind(t, w.Body, &vuln)
		assert.Equal(t, "blue", vuln.Vulnerability.Title)
		require.NotNil(t, vuln.Vulnerability.Edges.Sev)
		assert.Equal(t, "critical", vuln.Vulnerability.Edges.Sev.Label)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/vulnerability", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var resp server.RespVulnerabilities
		bind(t, w.Body, &resp)
		require.Len(t, resp.Vulnerabilities, 1)
		assert.Equal(t, "blue", resp.Vulnerabilities[0].Title)
		require.NotNil(t, resp.Vulnerabilities[0].Edges.Sev)
		assert.Equal(t, "critical", resp.Vulnerabilities[0].Edges.Sev.Label)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/severity/%d", sev.ID), nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	}

	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/vulnerability/CVE-2000-1000", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var vuln *model.RespVulnerability
		bind(t, w.Body, &vuln)
		assert.Equal(t, "blue", vuln.Vulnerability.Title)
		require.Nil(t, vuln.Vulnerability.Edges.Sev)
	}
}
