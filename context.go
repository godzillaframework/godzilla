package godzilla

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

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

func (ctx *context) Context() *fasthttp.RequestCtx {
	return ctx.requestCtx
}

func (ctx *context) SendBytes(value []byte) Context {
	ctx.requestCtx.Response.SetBodyRaw(value)
	return ctx
}

func (ctx *context) SendString(value string) Context {
	ctx.requestCtx.SetBodyString(value)
	return ctx
}

func (ctx *context) SendJSON(in interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	raw, err := json.Marshal(in)

	if err != nil {
		return err
	}

	ctx.requestCtx.Response.Header.SetContentType(MimeApplicationJSON)
	ctx.requestCtx.Response.SetBodyRaw(raw)

	return nil
}

func (ctx *context) Status(status int) Context {
	ctx.requestCtx.Response.SetStatusCode(status)
	return ctx
}

func (ctx *context) Get(key string) string {
	return GetString(ctx.requestCtx.Request.Header.Peek(key))
}

func (ctx *context) Set(key, value string) {
	ctx.requestCtx.Response.Header.Set(key, value)
}

func (ctx *context) Query(key string) string {
	return GetString(ctx.requestCtx.QueryArgs().Peek(key))
}

func (ctx *context) Body() string {
	return GetString(ctx.requestCtx.Request.Body())
}

func (ctx *context) SetLocal(key string, value interface{}) {
	ctx.requestCtx.SetUserValue(key, value)
}

func (ctx *context) GetLocal(key string) interface{} {
	return ctx.requestCtx.UserValue(key)
}

func (ctx *context) ParseBody(out interface{}) error {
	contentType := GetString(ctx.requestCtx.Request.Header.ContentType())
	if strings.HasPrefix(contentType, MimeApplicationJSON) {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		return json.Unmarshal(ctx.requestCtx.Request.Body(), out)
	}

	return fmt.Errorf("content type '%s' is not supported, "+
		"please open a request to support it "+
		"(https://github.com/godzillaframework//godzilla/issues",
		contentType)
}
