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
