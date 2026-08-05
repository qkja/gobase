package errors

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestGRPCStatus 验证 grpc 自动识别 BizError，序列化为带消息的 status（无 ErrorInfo 细节）
func TestGRPCStatus(t *testing.T) {
	be := ErrInvalidArgument().WithArgs(map[string]any{"field": "name"})
	st := status.Convert(be)

	if st.Code() != codes.InvalidArgument {
		t.Errorf("st.Code() = %v, want InvalidArgument", st.Code())
	}
	if st.Message() != "参数无效" {
		t.Errorf("st.Message() = %q, want 参数无效", st.Message())
	}
	// 跨进程业务码还原已移除：status 不应携带 ErrorInfo 细节
	if len(st.Details()) != 0 {
		t.Errorf("status 不应携带 ErrorInfo 细节，got %d details", len(st.Details()))
	}
}

// TestToGRPCError 验证 ToGRPCError 返回可识别的 gRPC error
func TestToGRPCError(t *testing.T) {
	err := ErrNotFound().ToGRPCError()
	if status.Code(err) != codes.NotFound {
		t.Errorf("status.Code(err) = %v, want NotFound", status.Code(err))
	}
}
