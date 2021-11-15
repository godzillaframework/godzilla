package godzilla

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	http.Pusher

	Status() int
	Written() bool
	Size() int
	Before(BeforeFunc)
}

type BeforeFunc func(ResponseWriter)
