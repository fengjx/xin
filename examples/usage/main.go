package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fengjx/xin"
	"github.com/fengjx/xin/pprof"
)

// 中间件示例：记录请求日志
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 创建 Xin 实例
	app := xin.New()

	// 添加全局中间件
	app.Use(loggerMiddleware)

	// 注册路由处理函数
	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Xin!")
	})

	// 处理 POST 请求
	app.POST("/api/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Create user endpoint")
	})

	// 处理 POST 请求
	app.GET("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "get user endpoint %s", r.PathValue("id"))
	})

	mux := app.Mux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service is healthy")
	})

	mux.Sub("/api/v1", func(sub *xin.Mux) {
		sub.HandleFunc("GET /foo", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "foo v1")
		})

		sub.HandleFunc("GET /bar", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "bar v1")
		})
	})

	// 提供静态文件服务
	app.Static("/static", "./public")

	// 开启 pprof，使用basic认证，用户名和密码为foo/bar
	app.Handle(pprof.DefaultPrefix, pprof.Profiler(map[string]string{
		"foo": "bar",
	}))

	// 启动服务器
	log.Println("Server starting on :8080...")
	app.Run(":8080")
}