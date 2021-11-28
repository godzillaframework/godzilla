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
	version = "2.0.0"
	name    = "Godzilla Framework"

	banner = `
Server Running on %s`
)

const (
	// defaultCacheSize is number of entries that can cache hold
	defaultCacheSize = 1000

	// defaultConcurrency is the maximum number of concurrent connections
	defaultConcurrency = 256 * 1024

	// defaultMaxRequestBodySize is the maximum request body size the server
	defaultMaxRequestBodySize = 4 * 1024 * 1024

	// defaultMaxRouteParams is the maximum number of routes params
	defaultMaxRouteParams = 1024

	// defaultMaxRequestURLLength is the maximum request url length
	defaultMaxRequestURLLength = 2048
)

// HTTP methods were copied from net/http.
const (
	MethodGet     = "GET"     // RFC 7231, 4.3.1
	MethodHead    = "HEAD"    // RFC 7231, 4.3.2
	MethodPost    = "POST"    // RFC 7231, 4.3.3
	MethodPut     = "PUT"     // RFC 7231, 4.3.4
	MethodPatch   = "PATCH"   // RFC 5789
	MethodDelete  = "DELETE"  // RFC 7231, 4.3.5
	MethodConnect = "CONNECT" // RFC 7231, 4.3.6
	MethodOptions = "OPTIONS" // RFC 7231, 4.3.7
	MethodTrace   = "TRACE"   // RFC 7231, 4.3.8
)

// HTTP status codes were copied from net/http.
const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1

	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusPartialContent       = 206 // RFC 7233, 4.1
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
	StatusFound             = 302 // RFC 7231, 6.4.3
	StatusSeeOther          = 303 // RFC 7231, 6.4.4
	StatusNotModified       = 304 // RFC 7232, 4.1
	StatusUseProxy          = 305 // RFC 7231, 6.4.5
	_                       = 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
	StatusTeapot                       = 418 // RFC 7168, 2.3.3
	StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

type Godzilla interface {
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

type godzilla struct {
	httpServer       *fasthttp.Server
	router           *router
	registeredRoutes []*Route
	address          string // server address
	middlewares      handlersChain
	settings         *Settings
}

// Settings struct holds server settings
type Settings struct {
	// Enable case sensitive routing
	CaseInSensitive bool // default false

	// Maximum size of LRU cache that will be used in routing if it's enabled
	CacheSize int // default 1000

	// Enables answering with HTTP status code 405 if request does not match
	// with any route, but there are another methods are allowed for that route
	// otherwise answer with Not Found handlers or status code 404.
	HandleMethodNotAllowed bool // default false

	// Enables automatic replies to OPTIONS requests if there are no handlers
	// registered for that route
	HandleOPTIONS bool // default false

	// Enables automatic recovering from panic while executing handlers by
	// answering with HTTP status code 500 and logging error message without
	// stopping service
	AutoRecover bool // default false

	// ServerName for sending in response headers
	ServerName string // default ""

	// Maximum request body size
	MaxRequestBodySize int // default 4 * 1024 * 1024

	// Maximum number of route params count
	MaxRouteParams int // default 1024

	// Max request url length
	MaxRequestURLLength int // default 2048

	// Maximum number of concurrent connections
	Concurrency int // default 256 * 1024

	// This will spawn multiple Go processes listening on the same port
	// Default: false
	Prefork bool

	// LRU caching used to speed up routing
	DisableCaching bool // default false

	DisableStartupMessage bool // default false

	// Disable keep-alive connections, the server will close incoming connections after sending the first response to client
	DisableKeepalive bool // default false

	// When set to true causes the default date header to be excluded from the response
	DisableDefaultDate bool // default false

	// When set to true, causes the default Content-Type header to be excluded from the Response
	DisableDefaultContentType bool // default false

	// By default all header names are normalized: conteNT-tYPE -> Content-Type
	DisableHeaderNormalizing bool // default false

	// The amount of time allowed to read the full request including body
	ReadTimeout time.Duration // default unlimited

	// The maximum duration before timing out writes of the response
	WriteTimeout time.Duration // default unlimited

	// The maximum amount of time to wait for the next request when keep-alive is enabled
	IdleTimeout time.Duration // default unlimited

	// Enable TLS or not
	TLSEnabled bool // default false

	// The path of the TLS certificate
	TLSCertPath string // default ""

	// The path of the TLS key
	TLSKeyPath string // default ""
}

// Route struct which holds each route info
type Route struct {
	Method   string
	Path     string
	Handlers handlersChain
}

func New(settings ...*Settings) Godzilla {
	gz := new(godzilla)
	gz.registeredRoutes = make([]*Route, 0)

	if len(settings) > 0 {
		gz.settings = settings[0]
	} else {
		gz.settings = &Settings{}
	}

	// set default settings for settings that don't have values set
	if gz.settings.CacheSize <= 0 {
		gz.settings.CacheSize = defaultCacheSize
	}

	if gz.settings.MaxRequestBodySize <= 0 {
		gz.settings.MaxRequestBodySize = defaultMaxRequestBodySize
	}

	if gz.settings.MaxRouteParams <= 0 || gz.settings.MaxRouteParams > defaultMaxRouteParams {
		gz.settings.MaxRouteParams = defaultMaxRouteParams
	}

	if gz.settings.MaxRequestURLLength <= 0 || gz.settings.MaxRequestURLLength > defaultMaxRequestURLLength {
		gz.settings.MaxRequestURLLength = defaultMaxRequestURLLength
	}

	if gz.settings.Concurrency <= 0 {
		gz.settings.Concurrency = defaultConcurrency
	}

	// Initialize router
	gz.router = &router{
		settings: gz.settings,
		cache:    make(map[string]*matchResult),
		pool: sync.Pool{
			New: func() interface{} {
				return new(context)
			},
		},
	}

	gz.httpServer = gz.newHTTPServer()

	return gz
}

// Start handling requests
func (gz *godzilla) Start(address string) error {
	// Setup router
	gz.setupRouter()

	if gz.settings.Prefork {
		if !gz.settings.DisableStartupMessage {
			printStartupMessage(address)
		}

		pf := prefork.New(gz.httpServer)
		pf.Reuseport = true
		pf.Network = "tcp4"

		if gz.settings.TLSEnabled {
			return pf.ListenAndServeTLS(address, gz.settings.TLSCertPath, gz.settings.TLSKeyPath)
		}
		return pf.ListenAndServe(address)
	}

	ln, err := net.Listen("tcp4", address)
	if err != nil {
		return err
	}
	gz.address = address

	if !gz.settings.DisableStartupMessage {
		printStartupMessage(address)
	}

	if gz.settings.TLSEnabled {
		return gz.httpServer.ServeTLS(ln, gz.settings.TLSCertPath, gz.settings.TLSKeyPath)
	}
	return gz.httpServer.Serve(ln)
}

// customLogger Customized logger used to filter logging messages
type customLogger struct{}

func (dl *customLogger) Printf(format string, args ...interface{}) {
}

// newHTTPServer returns a new instance of fasthttp server
func (gz *godzilla) newHTTPServer() *fasthttp.Server {
	return &fasthttp.Server{
		Handler:                       gz.router.Handler,
		Logger:                        &customLogger{},
		LogAllErrors:                  false,
		Name:                          gz.settings.ServerName,
		Concurrency:                   gz.settings.Concurrency,
		NoDefaultDate:                 gz.settings.DisableDefaultDate,
		NoDefaultContentType:          gz.settings.DisableDefaultContentType,
		DisableHeaderNamesNormalizing: gz.settings.DisableHeaderNormalizing,
		DisableKeepalive:              gz.settings.DisableKeepalive,
		NoDefaultServerHeader:         gz.settings.ServerName == "",
		ReadTimeout:                   gz.settings.ReadTimeout,
		WriteTimeout:                  gz.settings.WriteTimeout,
		IdleTimeout:                   gz.settings.IdleTimeout,
	}
}

// registerRoute registers handlers with method and path
func (gz *godzilla) registerRoute(method, path string, handlers handlersChain) *Route {
	if gz.settings.CaseInSensitive {
		path = strings.ToLower(path)
	}

	route := &Route{
		Path:     path,
		Method:   method,
		Handlers: handlers,
	}

	// Add route to registered routes
	gz.registeredRoutes = append(gz.registeredRoutes, route)
	return route
}

// setupRouter initializes router with registered routes
func (gz *godzilla) setupRouter() {
	for _, route := range gz.registeredRoutes {
		gz.router.handle(route.Method, route.Path, append(gz.middlewares, route.Handlers...))
	}

	// Frees intermediate stores after initializing router
	gz.registeredRoutes = nil
	gz.middlewares = nil
}

// Stop serving
func (gz *godzilla) Stop() error {
	err := gz.httpServer.Shutdown()

	// check if shutdown was ok and server had valid address
	if err == nil && gz.address != "" {
		log.Printf("%s stopped listening on %s", name, gz.address)
		return nil
	}

	return err
}

// Get registers an http relevant method
func (gz *godzilla) Get(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodGet, path, handlers)
}

// Head registers an http relevant method
func (gz *godzilla) Head(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodHead, path, handlers)
}

// Post registers an http relevant method
func (gz *godzilla) Post(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodPost, path, handlers)
}

// Put registers an http relevant method
func (gz *godzilla) Put(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodPut, path, handlers)
}

// Patch registers an http relevant method
func (gz *godzilla) Patch(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodPatch, path, handlers)
}

// Delete registers an http relevant method
func (gz *godzilla) Delete(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodDelete, path, handlers)
}

// Connect registers an http relevant method
func (gz *godzilla) Connect(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodConnect, path, handlers)
}

// Options registers an http relevant method
func (gz *godzilla) Options(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodOptions, path, handlers)
}

// Trace registers an http relevant method
func (gz *godzilla) Trace(path string, handlers ...handlerFunc) *Route {
	return gz.registerRoute(MethodTrace, path, handlers)
}

// Group appends a prefix to registered routes
func (gz *godzilla) Group(prefix string, routes []*Route) []*Route {
	for _, route := range routes {
		route.Path = prefix + route.Path
	}
	return routes
}

// Static serves files in root directory under specific prefix
func (gz *godzilla) Static(prefix, root string) {
	if gz.settings.CaseInSensitive {
		prefix = strings.ToLower(prefix)
	}

	// remove trailing slash
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

		// Pass to custom not found handlers if there are
		if gz.router.notFound != nil {
			gz.router.notFound[0](ctx)
			return
		}

		// Default Not Found response
		fctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound),
			fasthttp.StatusNotFound)
	}

	// TODO: Improve
	gz.Get(prefix, handler)

	if len(prefix) > 1 && prefix[len(prefix)-1] != '*' {
		gz.Get(prefix+"/*", handler)
	}
}

// NotFound registers an http handlers that will be called when no other routes
// match with request
func (gz *godzilla) NotFound(handlers ...handlerFunc) {
	gz.router.SetNotFound(handlers)
}

func (gz *godzilla) Use(middlewares ...handlerFunc) {
	gz.middlewares = append(gz.middlewares, middlewares...)
}

func printStartupMessage(addr string) {
	if prefork.IsChild() {
		log.Printf("Started child proc #%v\n", os.Getpid())
	} else {
		log.Println(banner, addr, version)
	}
}
