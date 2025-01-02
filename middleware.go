package xin

import (
	"net/http"

	"github.com/fengjx/go-halo/errs"
)

// recoverer panic 处理中间件
func recoverer(recoverHandle errs.RecoverHandle) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer errs.RecoverFunc(recoverHandle)
			next.ServeHTTP(w, r)
		})
	}
}
