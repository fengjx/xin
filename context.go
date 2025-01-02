package xin

import (
	"context"
	"net/http"
)

// errorKey 是用于在 context 中存储错误的键类型
type (
	errorKey struct{}
)

// WithError 将错误添加到 context 中
// ctx 是原始的 context
// err 是要存储的错误
// 返回包含错误的新 context
func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey{}, err)
}

// CtxError 从 context 中获取错误
// ctx 是包含错误的 context
// 返回存储在 context 中的错误，如果没有错误则返回 nil
func CtxError(ctx context.Context) error {
	err := ctx.Value(errorKey{})
	if err == nil {
		return nil
	}
	return err.(error)
}

// WithErrRequest 将错误添加到 http.Request 的 context 中
// r 是原始的 http.Request
// err 是要存储的错误
// 返回包含错误 context 的新 http.Request
func WithErrRequest(r *http.Request, err error) *http.Request {
	ctx := WithError(r.Context(), err)
	return r.WithContext(ctx)
}

// CtxRequestErr 从 http.Request 的 context 中获取错误
// r 是包含错误的 http.Request
// 返回存储在 request context 中的错误，如果没有错误则返回 nil
func CtxRequestErr(r *http.Request) error {
	return CtxError(r.Context())
}
