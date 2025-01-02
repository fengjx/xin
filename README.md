# xin

Xin 是一个轻量级的 Go Web 框架，专注于简单性和性能。它基于标准库 `net/http` 构建，提供了更便捷的 API 和常用的中间件。

## 特性

- 轻量级，无过多第三方依赖
- 基于标准库 `net/http`，性能优异
- 简单直观的路由系统
- 丰富的中间件支持
- 支持静态文件服务
- 优雅停机

## 安装

```bash
go get github.com/fengjx/xin
```

## 快速开始

```go
package main

import (
	"log"
	"net/http"

	"github.com/fengjx/xin"
	"github.com/fengjx/xin/middleware"
)

func main() {
	app := xin.New()
	app.Use(middleware.Logger)
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		xin.WriteString(w, http.StatusOK, "Hello World!")
	})
	log.Println("Server starting on :8080...")
	app.Run(":8080")
}

```

## 路由

### 基本路由

```go
// GET 请求
app.GET("/users/{id}", handleUsers)

// POST 请求
app.POST("/users/{id}", createUser)

// PUT 请求
app.PUT("/users/{id}", updateUser)

// DELETE 请求
app.DELETE("/users/${id}", deleteUser)
```

### 静态文件服务

```go
// 提供静态文件服务
app.Static("/static", "./static")
```

## 中间件

### 使用内置中间件

```go
// 添加恢复中间件
app.Use(middleware.Recoverer)

// 添加日志中间件
app.Use(middleware.Logger)

// 添加 CORS 中间件
app.Use(middleware.CORS)
```

### 自定义中间件

```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 中间件逻辑
        next.ServeHTTP(w, r)
    })
}

// 使用自定义中间件
app.Use(MyMiddleware)
```

## 优雅关闭

Xin 支持优雅关闭，确保所有请求都被正确处理：

```go
// 优雅关闭超时时间为 30 秒
if err := app.Shutdown(30 * time.Second); err != nil {
    log.Printf("server shutdown error: %v", err)
}
```

## 示例

更多示例可以在 [examples](./examples) 目录中找到：

- [Hello World](./examples/hello/main.go)
- [usage参考示例](./examples/usage/main.go)
- [中间件使用](./examples/middleware/main.go)

## 贡献

欢迎提交 Pull Request 或创建 Issue。

