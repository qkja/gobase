// Package trace 提供请求链路追踪上下文信息的存取。
//
// 设计约定：
//   - 本包只做「结构体 + context 存取」，不做任何 HTTP/gRPC/反射逻辑；
//   - ctx 中的 *Info 由网关入口中间件通过 WithInfo 写入，
//     各服务通过 GetInfo 读取，日志/grpcclient 自动提取并注入。
package trace

import "context"

// Info 链路追踪上下文信息。
type Info struct {
	TraceID  string
	RPCID    string
	Sampled  string
	UserID   string
	UserName string
}

type ctxKey struct{}

// WithInfo 将 info 写入 ctx 并返回新 ctx。
func WithInfo(ctx context.Context, info *Info) context.Context {
	return context.WithValue(ctx, ctxKey{}, info)
}

// GetInfo 从 ctx 读取 *Info；未写入或 ctx 为 nil 时返回 nil。
func GetInfo(ctx context.Context) *Info {
	if ctx == nil {
		return nil
	}
	info, _ := ctx.Value(ctxKey{}).(*Info)
	return info
}

// WithTraceID 将 traceID 写入 ctx（便捷方法，保留已有 Info 的其他字段）。
func WithTraceID(ctx context.Context, traceID string) context.Context {
	info := GetInfo(ctx)
	if info == nil {
		info = &Info{}
	}
	info.TraceID = traceID
	return WithInfo(ctx, info)
}

// TraceIDFromCtx 从 ctx 读取 TraceID，未写入时返回空字符串。
func TraceIDFromCtx(ctx context.Context) string {
	if info := GetInfo(ctx); info != nil {
		return info.TraceID
	}
	return ""
}
