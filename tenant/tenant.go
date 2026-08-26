// Package tenant 提供租户上下文信息的存取。
//
// 设计约定（与外部机制的分工）：
//   - 本包只做「结构体 + context 存取」，不做任何 proto/反射注入逻辑；
//   - ctx 中的 *Info 由外部机制（网关入口 / 上游调用方）通过 WithInfo 写入，
//     各服务 logic 层通过 GetInfo 读取，并把 info.TenantID 手动复制进 proto Req。
package tenant

import "context"

// Info 租户上下文信息。
type Info struct {
	TenantID   string         // 租户 ID（写入 proto Req 的 tenant_id）
	Language   string         // 业务语言（可选）
	UILanguage string         // 界面语言（可选）
	Extra      map[string]any // 扩展字段（预留，按需扩展）
}

// ctxKey 是 context 的私有 key 类型，避免与第三方库的 key 冲突。
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
