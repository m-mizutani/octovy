package model

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/pkg/utils"
	"github.com/m-mizutani/zlog"
)

const (
	ContextKeyLogger = "logger"
)

type Context struct {
	base context.Context
	log  *zlog.LogEntity
}

func NewContext() *Context {
	return NewContextWith(context.Background())
}

func NewContextWith(ctx context.Context) *Context {
	newCtx := &Context{
		base: ctx,
	}
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if obj, ok := ginCtx.Get(ContextKeyLogger); ok {
			if log, ok := obj.(*zlog.LogEntity); ok {
				newCtx.log = log
			}
		}
	}

	return newCtx
}

func (x *Context) Deadline() (deadline time.Time, ok bool) {
	return x.base.Deadline()
}
func (x *Context) Done() <-chan struct{}             { return x.base.Done() }
func (x *Context) Err() error                        { return x.base.Err() }
func (x *Context) Value(key interface{}) interface{} { return x.base.Value(key) }

// Logging feature
func (x *Context) Log() *zlog.LogEntity {
	return x.log
}
func (x *Context) With(key string, value interface{}) *zlog.LogEntity {
	if x.log == nil {
		x.log = utils.Logger.Log()
	}
	x.log = x.log.With(key, value)
	return x.log
}
