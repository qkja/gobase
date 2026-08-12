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

// WithInfo 将 info 写入 ctx 并返回新 ctx。ctx 为 nil 时回退 Background。
func WithInfo(ctx context.Context, info *Info) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
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
// 拷贝结构体而非原地修改，避免污染共享的 *Info。
func WithTraceID(ctx context.Context, traceID string) context.Context {
	info := GetInfo(ctx)
	newInfo := &Info{TraceID: traceID}
	if info != nil {
		newInfo.RPCID = info.RPCID
		newInfo.Sampled = info.Sampled
		newInfo.UserID = info.UserID
		newInfo.UserName = info.UserName
	}
	return WithInfo(ctx, newInfo)
}

// TraceIDFromCtx 从 ctx 读取 TraceID，未写入时返回空字符串。
func TraceIDFromCtx(ctx context.Context) string {
	if info := GetInfo(ctx); info != nil {
		return info.TraceID
	}
	return ""
}
