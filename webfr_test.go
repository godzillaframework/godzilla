package godzilla

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type fakeConn struct {
	net.Conn
	r bytes.Buffer
	w bytes.Buffer
}

func (c *fakeConn) Close() error {
	return nil
}

func (c *fakeConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

func (c *fakeConn) Write(b []byte) (int, error) {
	return c.w.Write(b)
}

func setupWebfr(settings ...*Settings) *webfr {
	wb := new(webfr)
	wb.registeredRoutes = make([]*Route, 0)

	if len(settings) > 0 {
		wb.settings = settings[0]
	} else {
		wb.settings = &Settings{}
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

	return wb
}

func startWebfr(wb *webfr) {
	wb.setupRouter()
	wb.httpServer = &fasthttp.Server{
		Handler:      wb.router.Handler,
		Logger:       &customLogger{},
		LogAllErrors: false,
	}
}

var emptyHandler = func(ctx Context) {}

var emptyHandlersChain = handlersChain{}

var fakeHandlersChain = handlersChain{emptyHandler}

func makeRequest(request *http.Request, wb *webfr) (*http.Response, error) {
	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		return nil, err
	}

	c := &fakeConn{}
	if _, err = c.r.Write(dumpRequest); err != nil {
		return nil, err
	}

	ch := make(chan error)
	go func() {
		ch <- wb.httpServer.ServeConn(c)
	}()

	if err = <-ch; err != nil {
		return nil, err
	}

	buffer := bufio.NewReader(&c.w)
	resp, err := http.ReadResponse(buffer, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

var handler = func(ctx Context) {}

var errorHandler = func(ctx Context) {
	m := make(map[string]int)
	m["a"] = 0
	ctx.SendString(string(rune(5 / m["a"])))
}

var headerHandler = func(ctx Context) {
	ctx.Set("custom", ctx.Get("my-header"))
}

var queryHandler = func(ctx Context) {
	ctx.SendString(ctx.Query("name"))
}

var bodyHandler = func(ctx Context) {
	ctx.Context().Response.SetBodyString(ctx.Body())
}

var unAuthorizedHandler = func(ctx Context) {
	ctx.Status(StatusUnauthorized)
}

var pingHandler = func(ctx Context) {
	ctx.SendString("pong")
}

var fallbackHandler = func(ctx Context) {
	ctx.Status(StatusNotFound).SendString("custom fallback handler")
}

var emptyMiddleware = func(ctx Context) {
	ctx.SetLocal("test-key", "value")

	ctx.Next()
}

var emptyMiddlewareHandler = func(ctx Context) {
	data, ok := ctx.GetLocal("test-key").(string)
	if !ok || data != "value" {
		panic("test-key value is wrong")
	}
}

func registerRoute(wb Webfr, method, path string, handler func(ctx Context)) {
	switch method {
	case MethodGet:
		wb.Get(path, handler)
	case MethodHead:
		wb.Head(path, handler)
	case MethodPost:
		wb.Post(path, handler)
	case MethodPut:
		wb.Put(path, handler)
	case MethodPatch:
		wb.Patch(path, handler)
	case MethodDelete:
		wb.Delete(path, handler)
	case MethodConnect:
		wb.Connect(path, handler)
	case MethodOptions:
		wb.Options(path, handler)
	case MethodTrace:
		wb.Trace(path, handler)
	}
}

func TestMethods(t *testing.T) {
	routes := []struct {
		method  string
		path    string
		handler func(ctx Context)
	}{
		{method: MethodGet, path: "/order/get", handler: queryHandler},
		{method: MethodPost, path: "/order/add", handler: bodyHandler},
		{method: MethodGet, path: "/books/find", handler: emptyHandler},
		{method: MethodGet, path: "/articles/search", handler: emptyHandler},
		{method: MethodPut, path: "/articles/search", handler: emptyHandler},
		{method: MethodHead, path: "/articles/test", handler: emptyHandler},
		{method: MethodPost, path: "/articles/204", handler: emptyHandler},
		{method: MethodPost, path: "/articles/205", handler: unAuthorizedHandler},
		{method: MethodGet, path: "/ping", handler: pingHandler},
		{method: MethodPut, path: "/posts", handler: emptyHandler},
		{method: MethodPatch, path: "/post/502", handler: emptyHandler},
		{method: MethodDelete, path: "/post/a23011a", handler: emptyHandler},
		{method: MethodConnect, path: "/user/204", handler: headerHandler},
		{method: MethodOptions, path: "/user/204/setting", handler: errorHandler},
		{method: MethodTrace, path: "/users/*", handler: emptyHandler},
	}

	wb := setupWebfr(&Settings{
		CaseInSensitive:        true,
		AutoRecover:            true,
		HandleOPTIONS:          true,
		HandleMethodNotAllowed: true,
	})

	for _, r := range routes {
		registerRoute(wb, r.method, r.path, r.handler)
	}

	startWebfr(wb)

	testCases := []struct {
		method      string
		path        string
		statusCode  int
		requestBody string
		body        string
		headers     map[string]string
	}{
		{method: MethodGet, path: "/order/get?name=art123", statusCode: StatusOK, body: "art123"},
		{method: MethodPost, path: "/order/add", requestBody: "testOrder", statusCode: StatusOK, body: "testOrder"},
		{method: MethodPost, path: "/books/find", statusCode: StatusMethodNotAllowed, body: "Method Not Allowed", headers: map[string]string{"Allow": "GET, OPTIONS"}},
		{method: MethodGet, path: "/articles/search", statusCode: StatusOK},
		{method: MethodGet, path: "/articles/search", statusCode: StatusOK},
		{method: MethodGet, path: "/Articles/search", statusCode: StatusOK},
		{method: MethodOptions, path: "/articles/search", statusCode: StatusOK},
		{method: MethodOptions, path: "*", statusCode: StatusOK},
		{method: MethodOptions, path: "/*", statusCode: StatusOK},
		{method: MethodGet, path: "/articles/searching", statusCode: StatusNotFound, body: "Not Found"},
		{method: MethodHead, path: "/articles/test", statusCode: StatusOK},
		{method: MethodPost, path: "/articles/204", statusCode: StatusOK},
		{method: MethodPost, path: "/articles/205", statusCode: StatusUnauthorized},
		{method: MethodPost, path: "/Articles/205", statusCode: StatusUnauthorized},
		{method: MethodPost, path: "/articles/206", statusCode: StatusNotFound, body: "Not Found"},
		{method: MethodGet, path: "/ping", statusCode: StatusOK, body: "pong"},
		{method: MethodPut, path: "/posts", statusCode: StatusOK},
		{method: MethodPatch, path: "/post/502", statusCode: StatusOK},
		{method: MethodDelete, path: "/post/a23011a", statusCode: StatusOK},
		{method: MethodConnect, path: "/user/204", statusCode: StatusOK, headers: map[string]string{"custom": "testing"}},
		{method: MethodOptions, path: "/user/204/setting", statusCode: StatusInternalServerError, body: "Internal Server Error"},
		{method: MethodTrace, path: "/users/testing", statusCode: StatusOK},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(tc.method, tc.path, strings.NewReader(tc.requestBody))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(tc.requestBody)))
		req.Header.Set("my-header", "testing")

		response, err := makeRequest(req, wb)

		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if response.StatusCode != tc.statusCode {
			t.Fatalf("%s(%s): returned %d expected %d", tc.method, tc.path, response.StatusCode, tc.statusCode)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if string(body) != tc.body {
			t.Fatalf("%s(%s): returned %s expected %s", tc.method, tc.path, body, tc.body)
		}

		for expectedKey, expectedValue := range tc.headers {
			actualValue := response.Header.Get(expectedKey)
			if actualValue != expectedValue {
				t.Errorf(" mismatch for route '%s' parameter '%s' actual '%s', expected '%s'",
					tc.path, expectedKey, actualValue, expectedValue)
			}
		}
	}
}

func TestStatic(t *testing.T) {
	wb := setupWebfr(&Settings{
		CaseInSensitive:        true,
		AutoRecover:            true,
		HandleOPTIONS:          true,
		HandleMethodNotAllowed: true,
	})

	wb.Static("/static/", "./assets/")

	startWebfr(wb)

	testCases := []struct {
		method     string
		path       string
		statusCode int
		body       string
	}{
		{method: MethodGet, path: "/static/webfr.png", statusCode: StatusOK},
	}

	for _, tc := range testCases {

		req, _ := http.NewRequest(tc.method, tc.path, nil)

		response, err := makeRequest(req, wb)

		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if response.StatusCode != tc.statusCode {
			t.Fatalf("%s(%s): returned %d expected %d", tc.method, tc.path, response.StatusCode, tc.statusCode)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if tc.body != "" && string(body) != tc.body {
			t.Fatalf("%s(%s): returned %s expected %s", tc.method, tc.path, body, tc.body)
		}
	}
}

func TestStartWithPrefork(t *testing.T) {
	wb := New(&Settings{
		Prefork: true,
	})

	go func() {
		time.Sleep(1000 * time.Millisecond)
		wb.Stop()
	}()

	wb.Start(":3000")
}

func TestStart(t *testing.T) {
	wb := New()

	go func() {
		time.Sleep(1000 * time.Millisecond)
		wb.Stop()
	}()

	wb.Start(":3010")
}

func TestStartWithTLS(t *testing.T) {
	wb := New(&Settings{
		TLSKeyPath:  "./assets/ssl-cert-snakeoil.key",
		TLSCertPath: "./assets/ssl-cert-snakeoil.crt",
		TLSEnabled:  true,
	})

	errs := make(chan error, 1)

	go func() {
		_, err := tls.DialWithDialer(
			&net.Dialer{
				Timeout: time.Second * 3,
			},
			"tcp",
			"localhost:3050",
			&tls.Config{
				InsecureSkipVerify: true,
			})
		errs <- err
		wb.Stop()
	}()

	wb.Start(":3050")

	err := <-errs
	if err != nil {
		t.Fatalf("StartWithSSL failed to connect with TLS error: %s", err)
	}
}

func TestStartInvalidListener(t *testing.T) {
	wb := New()

	go func() {
		time.Sleep(1000 * time.Millisecond)
		wb.Stop()
	}()

	if err := wb.Start("invalid listener"); err == nil {
		t.Fatalf("invalid listener passed")
	}
}

func TestStop(t *testing.T) {
	wb := New()

	go func() {
		time.Sleep(1000 * time.Millisecond)
		wb.Stop()
	}()

	wb.Start("")
}

func TestNotFound(t *testing.T) {
	wb := setupWebfr()

	wb.Get("/ping", pingHandler)
	wb.NotFound(fallbackHandler)

	startWebfr(wb)

	testCases := []struct {
		method     string
		path       string
		statusCode int
		body       string
	}{
		{method: MethodGet, path: "/ping", statusCode: StatusOK, body: "pong"},
		{method: MethodGet, path: "/error", statusCode: StatusNotFound, body: "custom fallback handler"},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(tc.method, tc.path, nil)
		response, err := makeRequest(req, wb)

		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if response.StatusCode != tc.statusCode {
			t.Fatalf("%s(%s): returned %d expected %d", tc.method, tc.path, response.StatusCode, tc.statusCode)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if string(body) != tc.body {
			t.Fatalf("%s(%s): returned %s expected %s", tc.method, tc.path, body, tc.body)
		}
	}
}

func TestGroupRouting(t *testing.T) {
	wb := setupWebfr()
	routes := []*Route{
		wb.Get("/id", emptyHandler),
		wb.Post("/abc", emptyHandler),
		wb.Post("/abcd", emptyHandler),
	}
	wb.Group("/account", wb.Group("/api", routes))

	startWebfr(wb)

	testCases := []struct {
		method     string
		path       string
		statusCode int
		body       string
	}{
		{method: MethodGet, path: "/account/api/id", statusCode: StatusOK},
		{method: MethodPost, path: "/account/api/abc", statusCode: StatusOK},
		{method: MethodPost, path: "/account/api/abcd", statusCode: StatusOK},
		{method: MethodGet, path: "/id", statusCode: StatusNotFound, body: "Not Found"},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(tc.method, tc.path, nil)
		response, err := makeRequest(req, wb)

		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if response.StatusCode != tc.statusCode {
			t.Fatalf("%s(%s): returned %d expected %d", tc.method, tc.path, response.StatusCode, tc.statusCode)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("%s(%s): %s", tc.method, tc.path, err.Error())
		}

		if string(body) != tc.body {
			t.Fatalf("%s(%s): returned %s expected %s", tc.method, tc.path, body, tc.body)
		}
	}
}

func TestUse(t *testing.T) {
	wb := setupWebfr()

	wb.Get("/ping", pingHandler)

	wb.Use(unAuthorizedHandler)

	startWebfr(wb)

	req, _ := http.NewRequest(MethodGet, "/ping", nil)
	response, err := makeRequest(req, wb)

	if err != nil {
		t.Fatalf("%s(%s): %s", MethodGet, "/ping", err.Error())
	}

	if response.StatusCode != StatusUnauthorized {
		t.Fatalf("%s(%s): returned %d expected %d", MethodGet, "/ping", response.StatusCode, StatusUnauthorized)
	}
}
