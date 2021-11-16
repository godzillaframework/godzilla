package godzilla

import "testing"

func catchPanic(f func()) (recv interface{}) {
	defer func() {
		recv = recover()
	}()

	f()

	return
}

type testRouter struct {
	path     string
	conflict bool
}

func TestAddRoute(t *testing.T) {
	tree := createRootNode()

	routes := []testRoute{
		{"/cmd/:tool/:sub", false},
		{"/cmd/vet", true},
		{"/src/*", false},
		{"/src/*", true},
		{"/src/test", true},
		{"/src/:test", true},
		{"/src/", false},
		{"/src1/", false},
		{"/src1/*", false},
		{"/search/:query", false},
		{"/search/invalid", true},
		{"/user_:name", false},
		{"/user_x", false},
		{"/id:id", false},
		{"/id/:id", false},
		{"/id/:value", true},
		{"/id/:id/settings", false},
		{"/id/:id/:type", true},
		{"/*", true},
		{"books/*/get", true},
		{"/file/test", false},
		{"/file/test", true},
		{"/file/:test", true},
		{"/orders/:id/settings/:id", true},
		{"/accounts/*/settings", true},
		{"/results/*", false},
		{"/results/*/view", true},
	}
	for _, route := range routes {
		recv := catchPanic(func() {
			tree.addRoute(route.path, emptyHandlersChain)
		})

		if route.conflict {
			if recv == nil {
				t.Errorf("no panic for conflicting route '%s'", route.path)
			}
		} else if recv != nil {
			t.Errorf("unexpected panic for route '%s': %v", route.path, recv)
		}
	}
}

type testRequest []struct {
	path   string
	match  bool
	params map[string]string
}
