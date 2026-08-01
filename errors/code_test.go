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

// TestRegister 业务段位扩展
func TestRegister(t *testing.T) {
	Register("2001", "目录域不存在", "Directory not found", codes.NotFound)
	be := New("2001")
	if be.Message != "目录域不存在" {
		t.Errorf("Register 后 Message = %q, want 目录域不存在", be.Message)
	}
	if be.GRPCCode != codes.NotFound {
		t.Errorf("Register 后 GRPCCode = %v, want NotFound", be.GRPCCode)
	}
	if got := Message("2001", LangEn); got != "Directory not found" {
		t.Errorf("Register 后英文 Message = %q, want Directory not found", got)
	}
}

// TestMessageBilingual 验证 Message 中英双语与回退
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
		if got := Message(tt.code, tt.lang); got != tt.wantZh && got != tt.wantEn {
			t.Errorf("Message(%q,%q) = %q，未命中预期", tt.code, tt.lang, got)
		}
	}
	// 语言一致性：同一 code 下 zh 与 en 应返回各自文案
	if Message(CodeOK, LangZh) != "成功" {
		t.Errorf("CodeOK zh = %q, want 成功", Message(CodeOK, LangZh))
	}
	if Message(CodeOK, LangEn) != "Success" {
		t.Errorf("CodeOK en = %q, want Success", Message(CodeOK, LangEn))
	}
	// 未识别语言按中文处理
	if got := Message(CodeNotFound, "fr-FR"); got != "资源不存在" {
		t.Errorf("未知语言应回退中文, got %q", got)
	}
	// 未注册码兜底为"未知错误"文案（按语言）
	if got := Message("99999", LangZh); got != "未知错误" {
		t.Errorf("未注册码中文 Message = %q, want 未知错误", got)
	}
	if got := Message("99999", LangEn); got != "Unknown error" {
		t.Errorf("未注册码英文 Message = %q, want Unknown error", got)
	}
}

// TestMessageFallback 验证双向回退：中文缺失用英文，英文缺失用中文
func TestMessageFallback(t *testing.T) {
	// 注册一个只有英文、没有中文的业务码
	registry["7101"] = errorMeta{msgEn: "Only English", grpc: codes.NotFound}
	// 中文缺失 → 回退英文
	if got := Message("7101", LangZh); got != "Only English" {
		t.Errorf("中文缺失应回退英文, got %q", got)
	}
	// 英文存在 → 正常返回英文
	if got := Message("7101", LangEn); got != "Only English" {
		t.Errorf("英文 Message = %q, want Only English", got)
	}
	// 清理，避免影响其他测试
	delete(registry, "7101")

	// 注册一个只有中文、没有英文的业务码
	registry["7102"] = errorMeta{msgZh: "只有中文", grpc: codes.NotFound}
	// 英文缺失 → 回退中文
	if got := Message("7102", LangEn); got != "只有中文" {
		t.Errorf("英文缺失应回退中文, got %q", got)
	}
	delete(registry, "7102")
}
