package policy

import (
	"encoding/json"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func NewCheck(query string) (Check, error) {
	compiler, err := ast.CompileModules(map[string]string{
		"check.rego": query,
	})
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	return &checkPolicy{
		compiler: compiler,
	}, nil
}

type CheckFactory func(query string) (Check, error)
type Check interface {
	Result(ctx *model.Context, inv *model.ScanReport) (*model.GitHubCheckResult, error)
}

type checkPolicy struct {
	compiler *ast.Compiler
}

func (x *checkPolicy) Result(ctx *model.Context, inv *model.ScanReport) (*model.GitHubCheckResult, error) {
	policy := rego.New(
		rego.Query(`response = data.octovy.check`),
		rego.Compiler(x.compiler),
		rego.Input(inv),
	)
	rs, err := policy.Eval(ctx)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	if len(rs) != 1 {
		return nil, goerr.Wrap(model.ErrInvalidPolicyResult, "only 1 result is acceptable").With("rego.ResultSet", rs)
	}

	response, ok := rs[0].Bindings["response"]
	if !ok {
		return nil, goerr.Wrap(model.ErrInvalidPolicyResult, "'response' is empty")
	}

	raw, err := json.Marshal(response)
	if err != nil {
		return nil, model.ErrInvalidPolicyResult.Wrap(err)
	}
	var result model.GitHubCheckResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, model.ErrInvalidPolicyResult.Wrap(err)
	}

	switch result.Conclusion {
	case "action_required", "cancelled", "failure", "neutral", "success", "skipped", "stale", "timed_out":
		return &result, nil

	default:
		return nil, goerr.Wrap(model.ErrInvalidPolicyResult, "Unsupported GitHub check conclusion").With("resultSet", rs)
	}
}
