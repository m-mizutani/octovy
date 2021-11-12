package rule

import (
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
	return &checkRule{
		compiler: compiler,
	}, nil
}

type CheckFactory func(query string) (Check, error)
type Check interface {
	Result(ctx *model.Context, inv *model.PackageInventory) (*model.GitHubCheckResult, error)
}

type checkRule struct {
	compiler *ast.Compiler
}

func (x *checkRule) Result(ctx *model.Context, inv *model.PackageInventory) (*model.GitHubCheckResult, error) {
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
		return nil, goerr.Wrap(model.ErrInvalidRuleResult, "only 1 result is acceptable").With("rego.ResultSet", rs)
	}

	response, ok := rs[0].Bindings["response"]
	if !ok {
		return nil, goerr.Wrap(model.ErrInvalidRuleResult, "'response' is empty")
	}
	respMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, goerr.Wrap(model.ErrInvalidRuleResult, "'response' type is invalid")
	}

	obj, ok := respMap["result"]
	if !ok {
		return nil, goerr.Wrap(model.ErrInvalidRuleResult, "'result' field is not found").With("response", respMap)
	}
	result, ok := obj.(string)
	if !ok {
		return nil, goerr.Wrap(model.ErrInvalidRuleResult, "'result' field must be string")
	}

	var msg string
	if obj, ok := respMap["msg"]; ok {
		if m, ok := obj.(string); ok {
			msg = m
		} else {
			ctx.Log().With("msg", obj).Warn("Check rule result has 'msg' field, but not string type")
		}
	}

	switch result {
	case "action_required", "cancelled", "failure", "neutral", "success", "skipped", "stale", "timed_out":
		return &model.GitHubCheckResult{
			Conclusion: result,
			Message:    msg,
		}, nil

	default:
		return nil, goerr.Wrap(model.ErrInvalidRuleResult, "Unsupported GitHub check conclusion").With("result", result)
	}
}
