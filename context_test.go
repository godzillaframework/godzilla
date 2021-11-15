package godzilla

import (
	"net/http"
	"testing"
)

func TestNext(t *testing.T) {
	routes := []struct {
		path       string
		middleware handlerFunc
		handler    handlerFunc
	}{
		{path: "/ok", middleware: emptyMiddleware, handler: emptyMiddlewareHandler},
		{path: "/unauthorized", middleware: unAuthorizedHandler, handler: emptyHandler},
	}

	godzilla := setupGodzilla()

	for _, r := range routes {
		wbf.Get(r.path, r.middleware, r.handler)
	}

	startgodzillafr(godzilla)

	testCases := []struct {
		path       string
		statusCode int
	}{
		{path: "/ok", statusCode: StatusOK},
		{path: "/unauthorized", statusCode: StatusUnauthorized},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(MethodGet, tc.path, nil)
		response, err := makeRequest(req, wbf)

		if err != nil {
			t.Fatalf("%s(%s): %s", MethodGet, tc.path, err.Error())
		}

		if response.StatusCode != tc.statusCode {
			t.Fatalf("%s(%s): returned %d expected %d", MethodGet, tc.path, response.StatusCode, tc.statusCode)
		}
	}
}
