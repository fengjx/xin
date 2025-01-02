package pprof

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/fengjx/xin"
	"github.com/fengjx/xin/middleware"
)

const (
	// DefaultPrefix 默认的 pprof 路由前缀
	DefaultPrefix = "/debug/pprof"
)

// Profiler pprof 路由
// creds: basic 认证的用户名和密码，支持多组
func Profiler(creds map[string]string) http.Handler {
	r := xin.NewMux()
	r.HandleFunc("/", pprof.Index)
	r.HandleFunc("/cmdline", pprof.Cmdline)
	r.HandleFunc("/profile", pprof.Profile)
	r.HandleFunc("/symbol", pprof.Symbol)
	r.HandleFunc("/trace", pprof.Trace)
	r.Handle("/vars", expvar.Handler())
	r.Handle("/allocs", pprof.Handler("allocs"))
	r.Handle("/block", pprof.Handler("block"))
	r.Handle("/goroutine", pprof.Handler("goroutine"))
	r.Handle("/heap", pprof.Handler("heap"))
	r.Handle("/mutex", pprof.Handler("mutex"))
	r.Handle("/threadcreate", pprof.Handler("threadcreate"))
	if len(creds) > 0 {
		r.Use(middleware.BasicAuth("pprof", creds))
	}
	return r
}
