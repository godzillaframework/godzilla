package godzilla

import "sync"

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
