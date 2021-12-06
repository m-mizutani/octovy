package opa

import "context"

type Mock struct {
	MockData func(ctx context.Context, input interface{}, result interface{}) error
}

func (x *Mock) Data(ctx context.Context, input interface{}, result interface{}) error {
	return x.MockData(ctx, input, result)
}
