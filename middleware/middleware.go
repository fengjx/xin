package middleware

// Ported from chi's middleware, source:
// https://github.com/go-chi/chi/tree/master/middleware

import (
	"net/http"

	"github.com/fengjx/xin"
)

// New will create a new middleware handler from a http.Handler.
func New(h http.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "http/middleware context value " + k.name
}

// NewRouter 创建一个新的路由器实例
func NewRouter() *xin.Mux {
	return xin.NewMux()
}
