package errors

import (
	"fmt"

	"github.com/qkja/gobase/i18n"
	"google.golang.org/grpc/codes"
)

// BizError 业务错误
//
// 面向 gRPC 服务的统一业务错误载体：业务错误码(string) + 用户消息 + gRPC 状态码映射。
// 实现了 error 接口；实现 GRPCStatus() 后，api 层可直接 return nil, errors.ErrXxx()，
// gRPC 服务端会自动把它序列化为带业务码细节的 status（见 grpc.go）。
type BizError struct {
	Code     string         // 业务错误码
	Message  string         // 错误消息（面向用户/前端）
	GRPCCode codes.Code     // 映射的 gRPC 状态码
	Args     map[string]any // 附加结构化参数（用于日志 / ErrorInfo 透传）
	cause    error          // 原始错误（可选，支持 Unwrap 包装链）
}

// New 创建业务错误，gRPC 映射取自注册表，默认消息（中文）由 i18n 提供。
// code 未注册时降级为 codes.Unknown，便于业务段位先行使用。
func New(code string) *BizError {
	m, ok := registry[code]
	if !ok {
		return &BizError{Code: code, GRPCCode: codes.Unknown}
	}
	return &BizError{Code: code, Message: lookupMessage(code, LangZh), GRPCCode: m.grpc}
}

// Newf 创建业务错误并格式化消息
func Newf(code, format string, args ...any) *BizError {
	be := New(code)
	be.Message = fmt.Sprintf(format, args...)
	return be
}

// NewWithArgs 创建业务错误并附加结构化参数
func NewWithArgs(code string, args map[string]any) *BizError {
	be := New(code)
	be.Args = args
	return be
}

// 便捷构造函数（通用段位，薄封装 New）
func ErrInternal() *BizError          { return New(CodeInternal) }
func ErrInvalidArgument() *BizError   { return New(CodeInvalidArgument) }
func ErrNotFound() *BizError          { return New(CodeNotFound) }
func ErrUnauthenticated() *BizError   { return New(CodeUnauthenticated) }
func ErrPermissionDenied() *BizError  { return New(CodePermissionDenied) }
func ErrTimeout() *BizError           { return New(CodeTimeout) }
func ErrAlreadyExists() *BizError     { return New(CodeAlreadyExists) }
func ErrConflict() *BizError          { return New(CodeConflict) }
func ErrResourceExhausted() *BizError { return New(CodeResourceExhausted) }
func ErrUnknown() *BizError           { return New(CodeUnknown) }

// Error 实现 error 接口，输出 "code: msg"，便于日志检索
func (e *BizError) Error() string {
	if e == nil {
		return ""
	}
	if m := e.message(); m != "" {
		return e.Code + ": " + m
	}
	return e.Code
}

// GetCode 返回业务错误码
func (e *BizError) GetCode() string {
	if e == nil {
		return ""
	}
	return e.Code
}

// GetMessage 返回业务错误消息（自定义消息优先，否则注册表默认中文）
func (e *BizError) GetMessage() string {
	if e == nil {
		return ""
	}
	return e.message()
}

// Unwrap 返回原始错误，支持 errors.Is / errors.As 遍历包装链
func (e *BizError) Unwrap() error { return e.cause }

// Is 支持 errors.Is(err, ErrXxx()) 按业务错误码比较
func (e *BizError) Is(target error) bool {
	t, ok := target.(*BizError)
	return ok && t != nil && t.Code == e.Code
}

// WithMessage 返回带自定义消息的错误副本（不可变式 builder）
func (e *BizError) WithMessage(msg string) *BizError {
	c := *e
	c.Message = msg
	return &c
}

// WithMessagef 返回带格式化消息的错误副本
func (e *BizError) WithMessagef(format string, args ...any) *BizError {
	return e.WithMessage(fmt.Sprintf(format, args...))
}

// WithArgs 返回附加结构化参数的错误副本
func (e *BizError) WithArgs(args map[string]any) *BizError {
	c := *e
	c.Args = args
	return &c
}

// WithCause 返回携带原始错误的副本。
// 推荐用它承载底层错误，而不是用 fmt.Errorf("%w") 包装后返回给 gRPC——
// grpc v1.41.0 的 status.FromError 只做顶层类型断言，%w 包装会让它识别不到业务码。
func (e *BizError) WithCause(err error) *BizError {
	c := *e
	c.cause = err
	return &c
}

// 语言标记
const (
	// LangZh 简体中文
	LangZh = "zh-CN"
	// LangEn 英文
	LangEn = "en-US"
)

// message 返回最终消息：优先自定义 Message，其次 i18n 默认（中文）
func (e *BizError) message() string {
	if e.Message != "" {
		return e.Message
	}
	return lookupMessage(e.Code, LangZh)
}

// lookupMessage 内部消息查找：由 i18n 唯一提供（嵌入默认 + 服务 i18n/<lang>.po 按 code 键覆盖）。
// lang 取 LangZh / LangEn；未配置的 code 兜底为"未知错误"（CodeUnknown）。
// 业务不直接调用——通过 errors.New/ErrXxx 创建时消息已随 BizError 带上（err.GetMessage()）。
func lookupMessage(code, lang string) string {
	if s, ok := i18n.Lookup(lang, code); ok {
		return s
	}
	if s, ok := i18n.Lookup(lang, CodeUnknown); ok {
		return s
	}
	return code
}
