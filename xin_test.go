package xin_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/fengjx/xin"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func waitForServer(url string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("server didn't respond in %s", timeout)
		case <-time.After(100 * time.Millisecond):
			_, err := http.Get(url)
			if err == nil {
				return nil
			}
		}
	}
}

func TestXin(t *testing.T) {
	// 获取空闲端口
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	addr := fmt.Sprintf(":%d", port)
	url := fmt.Sprintf("http://localhost%s", addr)

	// 创建服务器
	app := xin.New()
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "foo")
	})

	// 启动服务器
	errCh := make(chan error, 1)
	go func() {
		t.Logf("Server starting on %s...", addr)
		if err := app.Run(addr, true); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// 等待服务器启动
	if err := waitForServer(url, 5*time.Second); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}

	t.Run("Basic GET request", func(t *testing.T) {
		// 发送请求
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// 检查状态码
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		// 验证响应内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		if string(body) != "foo" {
			t.Errorf("Expected 'foo', got '%s'", string(body))
		}
	})

	t.Run("Concurrent requests", func(t *testing.T) {
		const concurrentRequests = 10
		errors := make(chan error, concurrentRequests)

		for i := 0; i < concurrentRequests; i++ {
			go func() {
				resp, err := http.Get(url)
				if err != nil {
					errors <- err
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					errors <- fmt.Errorf("expected status OK, got %v", resp.Status)
					return
				}
				errors <- nil
			}()
		}

		for i := 0; i < concurrentRequests; i++ {
			if err := <-errors; err != nil {
				t.Errorf("Concurrent request failed: %v", err)
			}
		}
	})

	// 优雅关闭服务器
	if err := app.Shutdown(10 * time.Second); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}

	// 检查服务器是否有错误
	if err := <-errCh; err != nil {
		t.Errorf("Server error: %v", err)
	}
}
