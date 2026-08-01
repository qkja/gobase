package errors

import (
	"fmt"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestGRPCStatus 验证 grpc 自动识别 BizError + 携带 ErrorInfo 细节
func TestGRPCStatus(t *testing.T) {
	be := ErrInvalidArgument().WithArgs(map[string]any{"field": "name"})
	st := status.Convert(be)

	if st.Code() != codes.InvalidArgument {
		t.Errorf("st.Code() = %v, want InvalidArgument", st.Code())
	}
	if st.Message() != "参数无效" {
		t.Errorf("st.Message() = %q, want 参数无效", st.Message())
	}
	// 细节里应能找到业务码
	found := false
	for _, d := range st.Details() {
		ei, ok := d.(*errdetails.ErrorInfo)
		if ok && ei.Reason == CodeInvalidArgument {
			found = true
			if ei.Metadata["field"] != "name" {
				t.Errorf("Metadata[field] = %q, want name", ei.Metadata["field"])
			}
		}
	}
	if !found {
		t.Error("status 细节中未找到 ErrorInfo 业务码")
	}
}

// TestToGRPCError 验证 ToGRPCError 返回可识别的 gRPC error
func TestToGRPCError(t *testing.T) {
	err := ErrNotFound().ToGRPCError()
	if status.Code(err) != codes.NotFound {
		t.Errorf("status.Code(err) = %v, want NotFound", status.Code(err))
	}
}

// TestFromErrorRoundTrip 核心：BizError → 上线 → 还原
func TestFromErrorRoundTrip(t *testing.T) {
	orig := ErrInvalidArgument().
		WithMessage("用户 qkj 不存在").
		WithArgs(map[string]any{"userId": "123", "retryable": true})

	// 模拟跨进程：序列化为 gRPC error（发上线）
	wire := orig.ToGRPCError()

	// 客户端/网关侧还原
	be, ok := FromError(wire)
	if !ok {
		t.Fatal("FromError 未识别业务错误")
	}
	if be.Code != CodeInvalidArgument {
		t.Errorf("Code = %q, want %q", be.Code, CodeInvalidArgument)
	}
	if be.Message != "用户 qkj 不存在" {
		t.Errorf("Message = %q, want %q", be.Message, "用户 qkj 不存在")
	}
	if be.GRPCCode != codes.InvalidArgument {
		t.Errorf("GRPCCode = %v, want InvalidArgument", be.GRPCCode)
	}
	if be.Args["userId"] != "123" {
		t.Errorf("Args[userId] = %v, want 123", be.Args["userId"])
	}
	if be.Args["retryable"] != "true" {
		t.Errorf("Args[retryable] = %v, want \"true\"（bool 序列化为字符串）", be.Args["retryable"])
	}
}

// TestFromErrorDirectAndWrapped 本进程内直接 / 包装识别
func TestFromErrorDirectAndWrapped(t *testing.T) {
	if be, ok := FromError(ErrNotFound()); !ok || be.Code != CodeNotFound {
		t.Error("FromError 直接识别 BizError 失败")
	}
	wrapped := fmt.Errorf("inner: %w", ErrTimeout())
	if be, ok := FromError(wrapped); !ok || be.Code != CodeTimeout {
		t.Error("FromError 未识别 %w 包装的 BizError")
	}
}

// TestFromErrorUnknownFallback 未识别错误兜底为 CodeInternal
func TestFromErrorUnknownFallback(t *testing.T) {
	plain := status.Error(codes.NotFound, "other service not found")
	be, ok := FromError(plain)
	if ok {
		t.Error("普通 status 错误不应识别为业务错误")
	}
	if be == nil || be.Code != CodeInternal {
		t.Errorf("兜底 Code = %v, want CodeInternal", be)
	}
	if be.Message != "other service not found" {
		t.Errorf("兜底应保留原始消息: %q", be.Message)
	}
}

// TestFromErrorNil nil 返回 (nil, false)
func TestFromErrorNil(t *testing.T) {
	be, ok := FromError(nil)
	if be != nil || ok {
		t.Error("FromError(nil) 应返回 (nil, false)")
	}
}

// TestIs 包级 Is 按业务码判断
func TestIs(t *testing.T) {
	if !Is(ErrInternal(), CodeInternal) {
		t.Error("Is 应匹配 CodeInternal")
	}
	if Is(ErrInternal(), CodeNotFound) {
		t.Error("Is 不应匹配 CodeNotFound")
	}
}
