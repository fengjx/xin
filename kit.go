package xin

import (
	"net"
	"net/http"
	"strings"

	"github.com/fengjx/go-halo/json"
)

var (
	trueClientIP  = http.CanonicalHeaderKey("True-Client-IP")
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

// Map is a map of string to any.
type Map map[string]any

// GetRealIP 从 request 获取真实客户端ip
func GetRealIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

// Write 写入响应内容
func Write(w http.ResponseWriter, code int, contentType string, message any) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(message)
}

// WriteString 写入响应
func WriteString(w http.ResponseWriter, code int, message any) error {
	return Write(w, code, "text/plain", message)
}

// WriteJSON 写入JSON响应
func WriteJSON(w http.ResponseWriter, code int, data any) error {
	jsonData, err := json.ToJson(data)
	if err != nil {
		return err
	}
	return Write(w, code, "application/json", jsonData)
}

// WriteNoContent 只返回响应码，不返回内容
func WriteNoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}
