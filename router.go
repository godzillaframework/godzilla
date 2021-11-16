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

func (r *router) handle(method, path string, handlers handlersChain) {
	if path == "" {
		panic("path is empty")
	} else if method == "" {
		panic("method is empty")
	} else if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	} else if len(handlers) == 0 {
		panic("no handlers provided with path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = createRootNode()
		r.trees[method] = root
	}

	root.addRoute(path, handlers)
}

func (r *router) allowed(reqMethod, path string, ctx *context) string {
	var allow string

	pathLen := len(path)

	if (pathLen == 1 && path[0] == '*') || (pathLen > 1 && path[1] == '*') {
		for method := range r.trees {
			if method == MethodOptions {
				continue
			}

			if allow != "" {
				allow += ", " + method
			} else {
				allow = method
			}
		}
		return allow
	}

	for method, tree := range r.trees {
		if method == reqMethod || method == MethodOptions {
			continue
		}

		handlers := tree.matchRoute(path, ctx)
		if handlers != nil {
			if allow != "" {
				allow += ", " + method
			} else {
				allow = method
			}
		}
	}

	if len(allow) > 0 {
		allow += ", " + MethodOptions
	}
	return allow
}
