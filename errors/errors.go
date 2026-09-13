package errors

import (
	"context"
	"strings"

	"github.com/qkja/gobase/i18n"
	"github.com/qkja/gobase/platform"
	"github.com/qkja/gobase/tenant"
)

// BizError 业务错误——仅包含错误码和消息两个字段（不导出，通过 getter 访问）。
// 实现了 error 接口，Is 按业务码比较。
// 全局变量 ErrXxx 不可变，多 goroutine 安全。
type BizError struct {
	code string
	msg  string
}

// New 创建业务错误，消息由 i18n 自动填充（默认中文）。
// code 未在 i18n 中配置时，兜底为 CodeUnknown 对应的文案。
func New(code string) *BizError {
	return &BizError{code: code, msg: Message(code, i18n.LangZh)}
}

// Error 实现 error 接口，输出 "code: msg" 格式，便于日志检索。
func (e *BizError) Error() string {
	if e == nil {
		return ""
	}
	if e.msg != "" {
		return e.code + ": " + e.msg
	}
	return e.code
}

// GetCode 返回业务错误码。
func (e *BizError) GetCode() string {
	if e == nil {
		return ""
	}
	return e.code
}

// GetMessage 返回业务错误消息。
func (e *BizError) GetMessage() string {
	if e == nil {
		return ""
	}
	return e.msg
}

// Is 支持 errors.Is(err, ErrXxx) 按业务错误码比较。
func (e *BizError) Is(target error) bool {
	if e == nil {
		return false
	}
	t, ok := target.(*BizError)
	return ok && t != nil && t.code == e.code
}

// Message 根据错误码和语言从 i18n 查询翻译。
// 语言参数使用 i18n.LangZh / i18n.LangEn。
// 查找链：指定语言 → 默认语言（zh-CN）→ CodeUnknown 兜底文案 → 原始 code。
func Message(code, lang string) string {
	if s, ok := i18n.Lookup(lang, code); ok && s != "" {
		return s
	}
	if s, ok := i18n.Lookup(lang, CodeUnknown); ok && s != "" {
		return s
	}
	return code
}

// GetTenantMsg 按错误码 + ctx 中租户上下文（tenant.Info）的当前页面语言（中/英）生成文案。
// 语言取自 ctx 中 tenant.Info（由网关经 tenant.WithInfo 写入）：
//   - UILanguage：租户控制台（界面）语言，即"当前页面"语言，优先；
//   - Language：租户（业务）语言，UILanguage 为空时回退。
//
// 识别为 zh 前缀时用中文（i18n.LangZh）；无租户信息或其它语言一律默认英文（i18n.LangEn）。
func GetTenantMsg(ctx context.Context, code string) string {
	if info := tenant.GetInfo(ctx); info != nil {
		ui := info.UILanguage
		if ui == "" {
			ui = info.Language
		}
		if strings.HasPrefix(ui, "zh") {
			return Message(code, i18n.LangZh)
		}
	}
	return Message(code, i18n.LangEn)
}

// GetPlatformMsg 按错误码 + ctx 中平台管理员上下文（platform.Info）的当前页面语言（中/英）生成文案。
// 语言取自 ctx 中 platform.Info（由网关经 platform.WithInfo 写入），字段语义同 GetTenantMsg。
//
// 识别为 zh 前缀时用中文（i18n.LangZh）；无平台管理员信息或其它语言一律默认英文（i18n.LangEn）。
func GetPlatformMsg(ctx context.Context, code string) string {
	if info := platform.GetInfo(ctx); info != nil {
		ui := info.UILanguage
		if ui == "" {
			ui = info.Language
		}
		if strings.HasPrefix(ui, "zh") {
			return Message(code, i18n.LangZh)
		}
	}
	return Message(code, i18n.LangEn)
}

// GetTenantMsg 返回按租户上下文当前页面语言本地化的错误消息（取错误自身携带的 code）。
// 用法：err.GetTenantMsg(ctx)。
func (e *BizError) GetTenantMsg(ctx context.Context) string {
	if e == nil {
		return ""
	}
	return GetTenantMsg(ctx, e.code)
}

// GetPlatformMsg 返回按平台管理员上下文当前页面语言本地化的错误消息（取错误自身携带的 code）。
// 用法：err.GetPlatformMsg(ctx)。
func (e *BizError) GetPlatformMsg(ctx context.Context) string {
	if e == nil {
		return ""
	}
	return GetPlatformMsg(ctx, e.code)
}
