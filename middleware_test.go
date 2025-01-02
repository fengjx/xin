package xin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fengjx/go-halo/errs"
)

func TestRecoverer(t *testing.T) {
	// 记录是否触发了 recover 处理
	var recovered bool

	// 创建一个简单的 recover 处理函数
	recoverHandle := func(err any, stack *errs.Stack) {
		fmt.Println("recovered")
		t.Log(err)
		t.Log(stack.StackTrace())
		recovered = true
	}

	// 创建一个会触发 panic 的处理器
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// 使用 recoverer 中间件包装处理器
	handler := recoverer(recoverHandle)(panicHandler)

	// 创建测试请求和响应记录器
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// 执行请求
	handler.ServeHTTP(w, req)

	// 验证是否正确处理了 panic
	if !recovered {
		t.Error("Panic was not recovered")
	}
}
