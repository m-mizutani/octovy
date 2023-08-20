package model

import (
	"context"
	"log/slog"

	"github.com/m-mizutani/octovy/pkg/utils"
)

type Context struct {
	logger *slog.Logger
	context.Context
}

func (x *Context) Logger() *slog.Logger { return x.logger }

func NewContext(options ...Option) *Context {
	ctx := &Context{
		logger:  utils.Logger(),
		Context: context.Background(),
	}

	for _, opt := range options {
		opt(ctx)
	}

	return ctx
}

type Option func(*Context)

func WithLogger(logger *slog.Logger) Option {
	return func(ctx *Context) {
		ctx.logger = logger
	}
}

func WithBase(base context.Context) Option {
	return func(ctx *Context) {
		ctx.Context = base
	}
}
