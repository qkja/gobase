// Package platform 提供平台管理员（平台运营控制台）上下文信息的存取。
//
// 与 tenant 包对偶：租户侧请求由网关经 tenant.WithInfo 写入租户上下文；
// 平台管理员操作则写入本包（platform.WithInfo），二者是两个独立的 ctx 信息体。
// 字段语义与 tenant.Info 一致（Language / UILanguage 供 errors.GetPlatformMsg 做消息本地化）。
//
// 设计约定：本包只做「结构体 + context 存取」，不做任何 proto/反射注入逻辑；
// ctx 中的 *Info 由外部机制（网关入口 / 上游调用方）通过 WithInfo 写入。
package platform

import "context"

// Info 平台管理员上下文信息。
type Info struct {
	AccountCode string         // 平台账号 code（写入下游 proto Req 的 operator_code 等账号标识字段）
	Language    string         // 业务语言（可选）
	UILanguage  string         // 界面语言（可选）
	Extra       map[string]any // 扩展字段（预留，按需扩展）
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
