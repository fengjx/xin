# xin

xin 是一个轻量级的 Go Web 框架，专注于简单性和性能。它基于标准库 `net/http` 构建，提供了更便捷的 API 和常用的中间件。

## 特性

- 轻量级，无过多第三方依赖
- 基于标准库 `net/http`，性能优异
- 简单直观的路由系统
- 丰富的中间件支持
- 支持静态文件服务

## 版本要求

- Go 版本 >= 1.18

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
	app.Run(":8080", true)
}

```

## 路由系统

### HTTP 方法支持

```go
// GET 请求
app.GET("/users", handleUsers)

// POST 请求
app.POST("/users", createUser)

// PUT 请求
app.PUT("/users/:id", updateUser)

// DELETE 请求
app.DELETE("/users/:id", deleteUser)

// PATCH 请求
app.PATCH("/users/:id", patchUser)

// OPTIONS 请求
app.OPTIONS("/users", optionsUser)

// HEAD 请求
app.HEAD("/users", headUser)

// 绑定任意 method
app.Any("/any", anyHandler)
```

### 子路由

```go
app.Sub("/api", func(r *xin.Mux) {
	r.GET("/users", handleUsers)
	r.POST("/users", createUser)
})

g := app.Group("/api/v1")
g.HandleFunc("GET /foo", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "foo v1")
})
```

### 静态文件服务

```go
// 静态页面，支持 index.html 自动查找
app.Static("/static", "./static")

// 自定义文件系统
app.StaticFS("/assets", myCustomFS)
```


## 请求参数处理

### 参数绑定
```go
// 绑定 URL 参数和表单数据到结构体
type UserForm struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    var form UserForm
    if err := xin.ShouldBind(r, &form); err != nil {
        // 处理错误
        return
    }
}

// 绑定 JSON 数据到结构体
func jsonHandler(w http.ResponseWriter, r *http.Request) {
    var user UserForm
    if err := xin.ShouldBindJSON(r, &user); err != nil {
        // 处理错误
        return
    }
}
```

### 获取请求参数
```go
// 获取查询参数
name := xin.GetQuery(r, "name")                    // 如果不存在返回空字符串
page := xin.GetQueryDefault(r, "page", "1")        // 如果不存在返回默认值

// 获取表单参数
email := xin.GetForm(r, "email")                   // 如果不存在返回空字符串
role := xin.GetFormDefault(r, "role", "user")      // 如果不存在返回默认值

// 获取请求头
token := xin.GetHeader(r, "Authorization")         // 如果不存在返回空字符串
lang := xin.GetHeaderDefault(r, "Accept-Language", "en-US") // 如果不存在返回默认值

// 获取 Cookie
sessionID, err := xin.GetCookie(r, "session_id")   // 如果不存在返回错误
userID := xin.GetCookieDefault(r, "user_id", "")   // 如果不存在返回默认值
```

## 中间件

### 内置中间件

```go
// 日志中间件
app.Use(middleware.Logger)

// CORS 中间件
app.Use(middleware.CORS)

// 压缩中间件
app.Use(middleware.Compress(5))

// 请求超时中间件
app.Use(middleware.Timeout(30 * time.Second))

// 请求 ID 中间件
app.Use(middleware.RequestID)
```

### 自定义中间件

```go
func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 前置处理
		next.ServeHTTP(w, r)
		// 后置处理
	})
}

app.Use(MyMiddleware)
```

## 错误处理

```go
// 在 context 中设置错误
ctx = xin.WithError(ctx, err)

// 从 context 中获取错误
if err := xin.CtxError(ctx); err != nil {
	// 处理错误
}

// 在请求中设置错误
r = xin.WithErrRequest(r, err)

// 从请求中获取错误
if err := xin.CtxRequestErr(r); err != nil {
	// 处理错误
}
```

## 示例

更多示例可以在 [examples](./examples) 目录中找到：

- [Hello World](./examples/hello/main.go) - 基本的 HTTP 服务器
- [完整示例](./examples/usage/main.go) - 展示框架的主要功能
- [中间件使用](./examples/middleware/main.go) - 中间件的使用方法

## 性能

- 轻量级设计，基于 `net/http` 标准库
- 支持中间件链和子路由模式

## 许可

本项目采用 MIT 许可证，详见 [LICENSE](./LICENSE) 文件。
