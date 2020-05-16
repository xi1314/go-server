// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package router

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

type Context struct {
	context iris.Context
	err     error
}

func (c *Context) Param(key string) string {
	return c.context.Params().Get(key)
}

func (c *Context) Header(key string, value string) {
	c.context.Header(key, value)
}

func (c *Context) StatusCode(code int) {
	c.context.StatusCode(code)
}

func (c *Context) Request() *http.Request {
	return c.context.Request()
}

func (c *Context) Writer() http.ResponseWriter {
	return c.context.ResponseWriter()
}

func (c *Context) Application() context.Application {
	return c.context.Application()
}

func (c *Context) GetStatusCode() int {
	return c.context.GetStatusCode()
}

func (c *Context) JSON(err error, data interface{}, meta *schema.Meta) {
	res := schema.Response{}
	if err != nil {
		res.Message = err.Error()

		if t, ok := err.(exception.Error); ok {
			res.Status = t.Code()
		} else {
			res.Status = exception.Unknown.Code()
		}
		res.Data = nil
		res.Meta = nil
	} else {
		res.Data = data
		res.Status = schema.StatusSuccess
		res.Meta = meta
	}

	_, _ = c.context.JSON(res)
}

func (c *Context) SetContext(key string, value interface{}) {
	c.context.Values().Set(key, value)
}

func (c *Context) GetContext(key string) interface{} {
	return c.context.Values().Get(key)
}

func Handler(handler func(c Context)) iris.Handler {
	return func(c iris.Context) {
		handler(Context{context: c})
	}
}