package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRule(t *testing.T) {
	engine := newServer(t)
	var sev ent.Severity
	var rule ent.Rule

	{ // Create a severity
		w := httptest.NewRecorder()
		req := newRequest("POST", "/api/v1/severity", model.RequestSeverity{
			Label: "critical",
		})
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		bind(t, w.Body, &sev)
	}

	{ // Create a rule
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, newRequest("POST", "/api/v1/rule",
			model.RequestRule{
				Action:     "fail",
				SeverityID: sev.ID,
			}))
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		bind(t, w.Body, &rule)
		assert.NotZero(t, rule.ID)
	}

	{ // Get the created rule
		w := httptest.NewRecorder()

		var rules []*ent.Rule
		engine.ServeHTTP(w, newRequest("GET", "/api/v1/rule", nil))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		bind(t, w.Body, &rules)
		require.Len(t, rules, 1)
		assert.Equal(t, rule.ID, rules[0].ID)
		require.NotNil(t, rules[0].Edges.Severity)
		assert.NotNil(t, "critical", rules[0].Edges.Severity.Label)
	}

	{ // Delete the created rule
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, newRequest("DELETE", fmt.Sprintf("/api/v1/rule/%d", rule.ID), nil))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	}
}

func TestRuleCreateFail(t *testing.T) {
	engine := newServer(t)

	{ // Create a rule with not existing sev ID
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, newRequest("POST", "/api/v1/rule",
			model.RequestRule{
				Action:     "fail",
				SeverityID: 1,
			}))
		assert.NotEqual(t, http.StatusCreated, w.Result().StatusCode)
	}
}

func TestRuleDeleteFail(t *testing.T) {
	engine := newServer(t)

	{ // Create a rule with not existing sev ID
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, newRequest("DELETE", "/api/v1/rule/1", nil))
		assert.NotEqual(t, http.StatusOK, w.Result().StatusCode)
	}
}
