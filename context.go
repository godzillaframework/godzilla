package godzilla

import "github.com/valyala/fasthttp"

const (
	MimeApplicationJSON = "application/json"
)

type Context interface {
	Next()
	Context() *fasthttp.RequestCtx
	Param(key string) string
	Query(key string) string
	SendBytes(value []byte) Context
	SendString(value string) Context
	SendJSON(in interface{}) error
	Status(status int) Context
	Set(key string, value string)
	Get(key string) string
	SetLocal(key string, value interface{})
	GetLocal(key string) interface{}
	Body() string
	ParseBody(out interface{}) error
}

type handlerFunc func(ctx Context)

type handlersChain []handlerFunc

type context struct {
	requestCtx  *fasthttp.RequestCtx
	paramValues map[string]string
	handlers    handlersChain
	index       int
}

func (ctx *context) Next() {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *context) Param(key string) string {
	return ctx.paramValues[key]
}

func (ctx *context) Context() *fasthttp.Request {
	return ctx.requestCtx
}

func (ctx *context) SendBytes(value []byte) Context {
	ctx.requestCtx.Response.SetBodyRaw(value)
	return ctx
}
