package godzilla

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/prefork"
)

const (
	Version = "1.0.0"
	Name    = "WebFramework"
	banner  = `
 
  WEB FRAMEWORK v%s
Listening on %s`
)

const (
	defaultCacheSize = 1000

	defaultConcurrency = 256 * 1024

	defaultMaxRequestBodySize = 4 * 1024 * 1024

	defaultMaxRouteParams = 1024

	defaultMaxRequestURLLength = 2048
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

const (
	StatusContinue           = 100
	StatusSwitchingProtocols = 101
	StatusProcessing         = 102

	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent            = 204
	StatusResetContent         = 205
	StatusPartialContent       = 206
	StatusMultiStatus          = 207
	StatusAlreadyReported      = 208
	StatusIMUsed               = 226

	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	_                       = 306
	StatusTemporaryRedirect = 307
	StatusPermanentRedirect = 308

	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
	StatusUnprocessableEntity          = 422
	StatusLocked                       = 423
	StatusFailedDependency             = 424
	StatusUpgradeRequired              = 426
	StatusPreconditionRequired         = 428
	StatusTooManyRequests              = 429
	StatusRequestHeaderFieldsTooLarge  = 431
	StatusUnavailableForLegalReasons   = 451

	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusVariantAlsoNegotiates         = 506
	StatusInsufficientStorage           = 507
	StatusLoopDetected                  = 508
	StatusNotExtended                   = 510
	StatusNetworkAuthenticationRequired = 511
)

type Webfr interface {
	Start(address string) error
	Stop() error
	Get(path string, handlers ...handlerFunc) *Route
	Head(path string, handlers ...handlerFunc) *Route
	Post(path string, handlers ...handlerFunc) *Route
	Put(path string, handlers ...handlerFunc) *Route
	Patch(path string, handlers ...handlerFunc) *Route
	Delete(path string, handlers ...handlerFunc) *Route
	Connect(path string, handlers ...handlerFunc) *Route
	Options(path string, handlers ...handlerFunc) *Route
	Trace(path string, handlers ...handlerFunc) *Route
	Group(prefix string, routes []*Route) []*Route
	Static(prefix, root string)
	NotFound(handlers ...handlerFunc)
	Use(middlewares ...handlerFunc)
}

type webfr struct {
	httpServer       *fasthttp.Server
	router           *router
	registeredRoutes []*Route
	address          string
	middlewares      handlersChain
	settings         *Settings
}

type Settings struct {
	CaseInSensitive bool

	CacheSize int

	HandleMethodNotAllowed bool

	HandleOPTIONS bool

	AutoRecover bool

	ServerName string

	MaxRequestBodySize int

	MaxRouteParams int

	MaxRequestURLLength int

	Concurrency int

	Prefork bool

	DisableCaching bool

	DisableStartupMessage bool

	DisableKeepalive bool

	DisableDefaultDate bool

	DisableDefaultContentType bool

	DisableHeaderNormalizing bool

	ReadTimeout time.Duration

	WriteTimeout time.Duration

	IdleTimeout time.Duration

	TLSEnabled bool

	TLSCertPath string

	TLSKeyPath string
}

type Route struct {
	Method   string
	Path     string
	Handlers handlersChain
}

func New(settings ...*Settings) Webfr {
	wb := new(webfr)
	wb.registeredRoutes = make([]*Route, 0)

	if len(settings) > 0 {
		wb.settings = settings[0]
	} else {
		wb.settings = &Settings{}
	}

	if wb.settings.CacheSize <= 0 {
		wb.settings.CacheSize = defaultCacheSize
	}

	if wb.settings.MaxRequestBodySize <= 0 {
		wb.settings.MaxRequestBodySize = defaultMaxRequestBodySize
	}

	if wb.settings.MaxRouteParams <= 0 || wb.settings.MaxRouteParams > defaultMaxRouteParams {
		wb.settings.MaxRouteParams = defaultMaxRouteParams
	}

	if wb.settings.MaxRequestURLLength <= 0 || wb.settings.MaxRequestURLLength > defaultMaxRequestURLLength {
		wb.settings.MaxRequestURLLength = defaultMaxRequestURLLength
	}

	if wb.settings.Concurrency <= 0 {
		wb.settings.Concurrency = defaultConcurrency
	}

	wb.router = &router{
		settings: wb.settings,
		cache:    make(map[string]*matchResult),
		pool: sync.Pool{
			New: func() interface{} {
				return new(context)
			},
		},
	}

	wb.httpServer = wb.newHTTPServer()

	return wb
}

func (wb *webfr) Start(address string) error {
	wb.setupRouter()

	if wb.settings.Prefork {
		if !wb.settings.DisableStartupMessage {
			printStartupMessage(address)
		}

		pf := prefork.New(wb.httpServer)
		pf.Reuseport = true
		pf.Network = "tcp4"

		if wb.settings.TLSEnabled {
			return pf.ListenAndServeTLS(address, wb.settings.TLSCertPath, wb.settings.TLSKeyPath)
		}
		return pf.ListenAndServe(address)
	}

	ln, err := net.Listen("tcp4", address)
	if err != nil {
		return err
	}
	wb.address = address

	if !wb.settings.DisableStartupMessage {
		printStartupMessage(address)
	}

	if wb.settings.TLSEnabled {
		return wb.httpServer.ServeTLS(ln, wb.settings.TLSCertPath, wb.settings.TLSKeyPath)
	}
	return wb.httpServer.Serve(ln)
}

type customLogger struct{}

func (dl *customLogger) Printf(format string, args ...interface{}) {
}

func (wb *webfr) newHTTPServer() *fasthttp.Server {
	return &fasthttp.Server{
		Handler:                       wb.router.Handler,
		Logger:                        &customLogger{},
		LogAllErrors:                  false,
		Name:                          wb.settings.ServerName,
		Concurrency:                   wb.settings.Concurrency,
		NoDefaultDate:                 wb.settings.DisableDefaultDate,
		NoDefaultContentType:          wb.settings.DisableDefaultContentType,
		DisableHeaderNamesNormalizing: wb.settings.DisableHeaderNormalizing,
		DisableKeepalive:              wb.settings.DisableKeepalive,
		NoDefaultServerHeader:         wb.settings.ServerName == "",
		ReadTimeout:                   wb.settings.ReadTimeout,
		WriteTimeout:                  wb.settings.WriteTimeout,
		IdleTimeout:                   wb.settings.IdleTimeout,
	}
}

func (wb *webfr) registerRoute(method, path string, handlers handlersChain) *Route {
	if wb.settings.CaseInSensitive {
		path = strings.ToLower(path)
	}

	route := &Route{
		Path:     path,
		Method:   method,
		Handlers: handlers,
	}

	wb.registeredRoutes = append(wb.registeredRoutes, route)
	return route
}

func (wb *webfr) setupRouter() {
	for _, route := range wb.registeredRoutes {
		wb.router.handle(route.Method, route.Path, append(wb.middlewares, route.Handlers...))
	}

	wb.registeredRoutes = nil
	wb.middlewares = nil
}

func (wb *webfr) Stop() error {
	err := wb.httpServer.Shutdown()

	if err == nil && wb.address != "" {
		log.Printf("%s stopped listening on %s", Name, wb.address)
		return nil
	}

	return err
}

func (wb *webfr) Get(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodGet, path, handlers)
}

func (wb *webfr) Head(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodHead, path, handlers)
}

func (wb *webfr) Post(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodPost, path, handlers)
}

func (wb *webfr) Put(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodPut, path, handlers)
}

func (wb *webfr) Patch(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodPatch, path, handlers)
}

func (wb *webfr) Delete(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodDelete, path, handlers)
}

func (wb *webfr) Connect(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodConnect, path, handlers)
}

func (wb *webfr) Options(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodOptions, path, handlers)
}

func (wb *webfr) Trace(path string, handlers ...handlerFunc) *Route {
	return wb.registerRoute(MethodTrace, path, handlers)
}

func (wb *webfr) Group(prefix string, routes []*Route) []*Route {
	for _, route := range routes {
		route.Path = prefix + route.Path
	}
	return routes
}

func (wb *webfr) Static(prefix, root string) {
	if wb.settings.CaseInSensitive {
		prefix = strings.ToLower(prefix)
	}

	if len(root) > 1 && root[len(root)-1] == '/' {
		root = root[:len(root)-1]
	}

	if len(prefix) > 1 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	fs := &fasthttp.FS{
		Root:       root,
		IndexNames: []string{"index.html"},
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			path := ctx.Path()

			if len(path) >= len(prefix) {
				path = path[len(prefix):]
			}

			if len(path) > 0 && path[0] != '/' {
				path = append([]byte("/"), path...)
			} else if len(path) == 0 {
				path = []byte("/")
			}
			return path
		},
	}

	fileHandler := fs.NewRequestHandler()
	handler := func(ctx Context) {
		fctx := ctx.Context()

		fileHandler(fctx)

		status := fctx.Response.StatusCode()
		if status != StatusNotFound && status != StatusForbidden {
			return
		}

		if wb.router.notFound != nil {
			wb.router.notFound[0](ctx)
			return
		}

		fctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound),
			fasthttp.StatusNotFound)
	}

	wb.Get(prefix, handler)

	if len(prefix) > 1 && prefix[len(prefix)-1] != '*' {
		wb.Get(prefix+"/*", handler)
	}
}

func (wb *webfr) NotFound(handlers ...handlerFunc) {
	wb.router.SetNotFound(handlers)
}

func (wb *webfr) Use(middlewares ...handlerFunc) {
	wb.middlewares = append(wb.middlewares, middlewares...)
}

func printStartupMessage(addr string) {
	if prefork.IsChild() {
		log.Printf("Started child proc #%v\n", os.Getpid())
	} else {
		log.Printf(banner, Version, addr)
	}
}
