package godzilla

const (
	Version = "1.0.0"
	Name    = "Gdozilla"

	banner = `

	GODZILLA
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
