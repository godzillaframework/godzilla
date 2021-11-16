package godzilla

import (
	"sync"

	"github.com/valyala/fasthttp"
)

var (
	defaultContentType = []byte("text/plain; charset=utf-8")
)

type router struct {
	trees    map[string]*node
	cache    map[string]*matchResult
	cacheLen int
	mutex    sync.RWMutex
	notFound handlersChain
	settings *Settings
	pool     sync.Pool
}

type matchResult struct {
	handlers handlersChain
	params   map[string]string
}

func (r *router) acquireCtx(fctx *fasthttp.RequestCtx) *context {
	ctx := r.pool.Get().(*context)

	ctx.index = 0
	ctx.paramValues = make(map[string]string)
	ctx.requestCtx = fctx

	return ctx
}

func (r *router) releaseCtx(ctx *context) {
	ctx.handlers = nil
	ctx.paramValues = nil
	ctx.requestCtx = nil
	r.pool.Put(ctx)
}
