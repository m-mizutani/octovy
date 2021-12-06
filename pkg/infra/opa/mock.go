package opa

import "context"

type Mock struct {
	MockData func(ctx context.Context, pkg RegoPkg, input interface{}, result interface{}) error
}

func (x *Mock) Data(ctx context.Context, pkg RegoPkg, input interface{}, result interface{}) error {
	return x.MockData(ctx, pkg, input, result)
}
