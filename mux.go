package xin

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/fengjx/go-halo/errs"
)

// HTTPMiddleware http.Handler 请求中间件
type HTTPMiddleware func(http.Handler) http.Handler

// MiddlewareFunc http.HandlerFunc 请求中间件
type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

// Mux http 路由
type Mux struct {
	*http.ServeMux
	middlewares []HTTPMiddleware
	handler     http.Handler
}

// NewMux 创建一个新的 HTTP 路由复用器
func NewMux() *Mux {
	mux := http.NewServeMux()
	router := &Mux{
		ServeMux: mux,
	}
	router.then(mux)
	return router
}

// Use 注册中间件
func (mux *Mux) Use(middlewares ...HTTPMiddleware) *Mux {
	mux.middlewares = append(mux.middlewares, middlewares...)
	mux.then(mux.ServeMux)
	return mux
}

// Group 注册路由组
func (mux *Mux) Group(prefix string) *Mux {
	group := NewMux()
	// 确保不以 / 结尾
	prefix = strings.TrimSuffix(prefix, "/")
	mux.Handle(prefix+"/", http.StripPrefix(prefix, group))
	return group
}

func (mux *Mux) then(h http.Handler) {
	mux.handler = HandlerChain(h, mux.middlewares...)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errs.Recover()
	mux.handler.ServeHTTP(w, r)
}

// Handle 注册HTTP处理器 参考 http.ServeMux.Handle
// [METHOD ][HOST]/[PATH]
func (mux *Mux) Handle(pattern string, handler http.Handler) *Mux {
	mux.ServeMux.Handle(pattern, handler)
	return mux
}

// HandleFunc 注册HTTP处理函数 参考 http.ServeMux.HandleFunc
// [METHOD][HOST]/[PATH]
func (mux *Mux) HandleFunc(pattern string, hf http.HandlerFunc) *Mux {
	mux.ServeMux.HandleFunc(pattern, hf)
	return mux
}

// Any alias for HandleFunc
func (mux *Mux) Any(pattern string, hf http.HandlerFunc) *Mux {
	return mux.HandleFunc(pattern, hf)
}

// POST 绑定 POST 请求
func (mux *Mux) POST(relativePath string, hf http.HandlerFunc) *Mux {
	mux.HandleFunc(fmt.Sprintf("POST %s", relativePath), hf)
	return mux
}

// GET 绑定 GET 请求
func (mux *Mux) GET(relativePath string, hf http.HandlerFunc) *Mux {
	mux.HandleFunc(fmt.Sprintf("GET %s", relativePath), hf)
	return mux
}

// DELETE 绑定 DELETE 请求
func (mux *Mux) DELETE(relativePath string, hf http.HandlerFunc) *Mux {
	mux.HandleFunc(fmt.Sprintf("DELETE %s", relativePath), hf)
	return mux
}

// PATCH 绑定 PATCH 请求
func (mux *Mux) PATCH(relativePath string, hf http.HandlerFunc) *Mux {
	mux.HandleFunc(fmt.Sprintf("PATCH %s", relativePath), hf)
	return mux
}

// PUT 绑定 PUT 请求
func (mux *Mux) PUT(relativePath string, hf http.HandlerFunc) *Mux {
	mux.HandleFunc(fmt.Sprintf("PUT %s", relativePath), hf)
	return mux
}

// OPTIONS 绑定 OPTIONS 请求
func (mux *Mux) OPTIONS(relativePath string, hf http.HandlerFunc) *Mux {
	mux.HandleFunc(fmt.Sprintf("OPTIONS %s", relativePath), hf)
	return mux
}

// HEAD is a shortcut for router.Handle("HEAD", path, handlers).
func (mux *Mux) HEAD(relativePath string, hf http.HandlerFunc) *Mux {
	mux.HandleFunc(fmt.Sprintf("HEAD %s", relativePath), hf)
	return mux
}

// Static 注册静态文件服务
// 默认不显示文件目录
func (mux *Mux) Static(pattern string, root string) *Mux {
	return mux.StaticFS(pattern, Dir(root, false))
}

// StaticFS 注册静态文件服务，自定义文件系统
// fs 可以使用 luchen.Dir() 创建
func (mux *Mux) StaticFS(pattern string, fs fs.FS) *Mux {
	prefix := pattern
	// 处理 [METHOD /path] 格式
	arr := strings.Fields(pattern)
	if len(arr) > 1 {
		prefix = arr[1]
	}
	mux.ServeMux.Handle(pattern, FileHandler(prefix, fs))
	return mux
}

// HandlerChain 使用中间件包装 handler
func HandlerChain(h http.Handler, middlewares ...HTTPMiddleware) http.Handler {
	size := len(middlewares)
	for i := range middlewares {
		h = middlewares[size-1-i](h)
	}
	return h
}

// WrapMiddleware wraps `func(http.Handler) http.Handler` into `xin.MiddlewareFunc`
func WrapMiddleware(m HTTPMiddleware) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			m(next).ServeHTTP(w, r)
		}
	}
}

// WrapHandler wraps `http.Handler` into `http.HandlerFunc`.
func WrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
