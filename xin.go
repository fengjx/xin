package xin

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fengjx/go-halo/addr"
	"github.com/fengjx/go-halo/errs"
	"github.com/fengjx/go-halo/halo"
)

var Debug = false

func SetDebug(debug bool) {
	Debug = debug
}

// Xin 是核心Web服务器结构体，用于管理HTTP路由和服务器操作
type Xin struct {
	httpServer    *http.Server       // HTTP服务器实例
	router        *Mux               // 路由复用器
	mtx           sync.Mutex         // 用于并发安全的读写锁
	host          string             // 服务器主机地址
	port          int                // 服务器端口
	middlewares   []HTTPMiddleware   // 中间件
	recoverHandle errs.RecoverHandle // panic 处理函数
	started       bool               // 是否已关闭
}

// New 创建一个新的Xin实例
func New() *Xin {
	x := &Xin{}
	mux := NewMux()
	httpServer := &http.Server{
		Handler: mux,
	}
	x.router = mux
	x.httpServer = httpServer
	x.recoverHandle = x.defaultRecoverHandle
	return x
}

func (x *Xin) init() {
	// recover 中间件
	x.router.Use(recoverer(x.recoverHandle))
	// 添加中间件
	x.router.Use(x.middlewares...)
}

// Run 启动HTTP服务器
// sync 是否同步启动
// address 参数格式为 "host:port"，例如 ":8080" 或 "192.168.1.100:8080"
func (x *Xin) Run(address string, sync bool) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	return x.Serve(ln, sync)
}

// Serve 启动HTTP服务器
// sync 是否同步启动
func (x *Xin) Serve(ln net.Listener, sync bool) error {
	x.mtx.Lock()
	if x.started {
		x.mtx.Unlock()
		return nil
	}
	x.started = true
	x.init()
	la := ln.Addr().String()
	host, port, _ := addr.ExtractHostPort(la)
	x.host = host
	x.port, _ = strconv.Atoi(port)

	// 使用 halo 包的优雅关闭功能
	halo.AddShutdownCallback(func() {
		x.Shutdown(60 * time.Second)
	})
	x.mtx.Unlock()
	return x.httpServer.Serve(ln)
}

// Shutdown 优雅地停止服务器
func (x *Xin) Shutdown(timeout time.Duration) error {
	x.mtx.Lock()
	defer x.mtx.Unlock()
	if !x.started {
		return nil
	}
	x.started = false
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 先关闭监听器，停止接收新请求
	if err := x.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	return nil
}

// Recover 设置 panic 处理函数
func (x *Xin) Recover(fn errs.RecoverHandle) *Xin {
	x.recoverHandle = fn
	return x
}

// Use 添加全局中间件
// middlewares 可以添加多个中间件，它们将按照添加顺序依次执行
func (x *Xin) Use(middlewares ...HTTPMiddleware) *Xin {
	x.middlewares = append(x.middlewares, middlewares...)
	return x
}

// Mux 获取路由复用器
func (x *Xin) Mux() *Mux {
	return x.router
}

// Router 注册路由处理函数
// fn 是一个接收路由复用器的函数，用于配置路由规则
func (x *Xin) Router(fn func(r *Mux)) *Xin {
	fn(x.router)
	return x
}

// Handle 注册一个处理特定模式的HTTP处理器
// pattern 格式为 "[METHOD ][HOST]/[PATH]"
// handler 为实现了http.Handler接口的处理器
func (x *Xin) Handle(pattern string, handler http.Handler) *Xin {
	x.router.Handle(pattern+"/", http.StripPrefix(pattern, handler))
	return x
}

// HandleFunc 注册一个处理特定模式的处理函数
// pattern 格式为 "[METHOD ][HOST]/[PATH]"
// hf 为处理HTTP请求的函数
func (x *Xin) HandleFunc(pattern string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(pattern, hf)
	return x
}

// Any alias for HandleFunc
func (x *Xin) Any(pattern string, hf http.HandlerFunc) *Xin {
	return x.HandleFunc(pattern, hf)
}

// POST 注册一个处理POST请求的路由
// relativePath 为相对路径
// hf 为处理HTTP请求的函数
func (x *Xin) POST(relativePath string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(fmt.Sprintf("POST %s", relativePath), hf)
	return x
}

// GET 注册一个处理GET请求的路由
// relativePath 为相对路径
// hf 为处理HTTP请求的函数
func (x *Xin) GET(relativePath string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(fmt.Sprintf("GET %s", relativePath), hf)
	return x
}

// DELETE 绑定 DELETE 请求
func (x *Xin) DELETE(relativePath string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(fmt.Sprintf("DELETE %s", relativePath), hf)
	return x
}

// PATCH 绑定 PATCH 请求
func (x *Xin) PATCH(relativePath string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(fmt.Sprintf("PATCH %s", relativePath), hf)
	return x
}

// PUT 绑定 PUT 请求
func (x *Xin) PUT(relativePath string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(fmt.Sprintf("PUT %s", relativePath), hf)
	return x
}

// OPTIONS 绑定 OPTIONS 请求
func (x *Xin) OPTIONS(relativePath string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(fmt.Sprintf("OPTIONS %s", relativePath), hf)
	return x
}

// HEAD is a shortcut for router.Handle("HEAD", path, handlers).
func (x *Xin) HEAD(relativePath string, hf http.HandlerFunc) *Xin {
	x.router.HandleFunc(fmt.Sprintf("HEAD %s", relativePath), hf)
	return x
}

// Static 注册静态文件服务
// pattern 为URL匹配模式
// root 为静态文件所在的根目录路径
func (x *Xin) Static(pattern string, root string) *Xin {
	x.router.Static(pattern, root)
	return x
}

// StaticFS 注册自定义文件系统的静态文件服务
// pattern 为URL匹配模式
// fs 为实现了fs.FS接口的文件系统
func (x *Xin) StaticFS(pattern string, fs fs.FS) *Xin {
	x.router.StaticFS(pattern, fs)
	return x
}

// HostPort 获取服务器地址和端口
func (x *Xin) HostPort() (host string, port int) {
	return x.host, x.port
}

func (x *Xin) defaultRecoverHandle(err any, stack *errs.Stack) {
	log.Printf("panic: %s %+v\r\n", err, stack)
}
