package errors

import (
	"testing"

	"google.golang.org/grpc/codes"
)

// TestRegistryContainsAllConsts 遍历全部常量码断言 registry 有对应项
func TestRegistryContainsAllConsts(t *testing.T) {
	codesToCheck := []string{
		CodeOK, CodeInternal, CodeInvalidArgument, CodeNotFound,
		CodeUnauthenticated, CodePermissionDenied, CodeTimeout,
		CodeAlreadyExists, CodeConflict, CodeResourceExhausted, CodeUnknown,
	}
	for _, c := range codesToCheck {
		if _, ok := registry[c]; !ok {
			t.Errorf("registry 缺少错误码 %q", c)
		}
	}
}

// TestNew 验证默认消息与 gRPC 映射
func TestNew(t *testing.T) {
	tests := []struct {
		code     string
		wantMsg  string
		wantCode codes.Code
	}{
		{CodeOK, "成功", codes.OK},
		{CodeInternal, "系统内部错误", codes.Internal},
		{CodeInvalidArgument, "参数无效", codes.InvalidArgument},
		{CodeNotFound, "资源不存在", codes.NotFound},
		{CodeTimeout, "请求超时", codes.DeadlineExceeded},
		{CodeUnknown, "未知错误", codes.Unknown},
	}
	for _, tt := range tests {
		be := New(tt.code)
		if be.Message != tt.wantMsg {
			t.Errorf("New(%q).Message = %q, want %q", tt.code, be.Message, tt.wantMsg)
		}
		if be.GRPCCode != tt.wantCode {
			t.Errorf("New(%q).GRPCCode = %v, want %v", tt.code, be.GRPCCode, tt.wantCode)
		}
	}
}

// TestNewUnregistered 未注册码降级为 Unknown
func TestNewUnregistered(t *testing.T) {
	be := New("9999")
	if be.GRPCCode != codes.Unknown {
		t.Errorf("未注册码 GRPCCode = %v, want Unknown", be.GRPCCode)
	}
	if be.Code != "9999" {
		t.Errorf("Code = %q, want 9999", be.Code)
	}
}

// TestMessageBilingual 验证 lookupMessage 中英双语与回退
func TestMessageBilingual(t *testing.T) {
	tests := []struct {
		code   string
		lang   string
		wantZh string
		wantEn string
	}{
		{CodeOK, LangZh, "成功", "Success"},
		{CodeInternal, LangZh, "系统内部错误", "Internal error"},
		{CodeNotFound, LangEn, "资源不存在", "Resource not found"},
		{CodeInvalidArgument, LangEn, "参数无效", "Invalid argument"},
	}
	for _, tt := range tests {
		if got := lookupMessage(tt.code, tt.lang); got != tt.wantZh && got != tt.wantEn {
			t.Errorf("lookupMessage(%q,%q) = %q，未命中预期", tt.code, tt.lang, got)
		}
	}
	// 语言一致性：同一 code 下 zh 与 en 应返回各自文案
	if lookupMessage(CodeOK, LangZh) != "成功" {
		t.Errorf("CodeOK zh = %q, want 成功", lookupMessage(CodeOK, LangZh))
	}
	if lookupMessage(CodeOK, LangEn) != "Success" {
		t.Errorf("CodeOK en = %q, want Success", lookupMessage(CodeOK, LangEn))
	}
	// 未识别语言按中文处理
	if got := lookupMessage(CodeNotFound, "fr-FR"); got != "资源不存在" {
		t.Errorf("未知语言应回退中文, got %q", got)
	}
	// 未注册码兜底为"未知错误"文案（按语言）
	if got := lookupMessage("99999", LangZh); got != "未知错误" {
		t.Errorf("未注册码中文 Message = %q, want 未知错误", got)
	}
	if got := lookupMessage("99999", LangEn); got != "Unknown error" {
		t.Errorf("未注册码英文 Message = %q, want Unknown error", got)
	}
}

// TestMessageFallback 验证 i18n 回退语义：未知语言→默认中文；未配置 i18n 的码→兜底"未知错误"
func TestMessageFallback(t *testing.T) {
	// 未知语言 → 回退默认语言(zh-CN)
	if got := lookupMessage(CodeInternal, "fr-FR"); got != "系统内部错误" {
		t.Errorf("未知语言应回退默认中文, got %q", got)
	}
	// 未配置 i18n 的码 → 按语言兜底"未知错误"
	if got := lookupMessage("5999", LangZh); got != "未知错误" {
		t.Errorf("未配置码中文 = %q, want 未知错误", got)
	}
	if got := lookupMessage("5999", LangEn); got != "Unknown error" {
		t.Errorf("未配置码英文 = %q, want Unknown error", got)
	}
}
