package grpcclient

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/qkja/gobase/errors"
	"github.com/qkja/gobase/store"
	"github.com/qkja/gobase/tenant"
)

// TestInjectMetadata 验证 tenant 与 store 的 trace 值被写入 outgoing metadata（键名小写）。
func TestInjectMetadata(t *testing.T) {
	ctx := tenant.WithInfo(context.Background(), &tenant.Info{
		TenantID: "T1",
		Language: "zh-CN",
	})
	store.Put("t-head-traceId", "trace-123")
	store.Put("t-head-userId", "u9")
	defer store.Clean()

	out := injectMetadata(ctx)
	md, _ := metadata.FromOutgoingContext(out)

	// 键名应小写（gRPC 规范）：t-head-tenantId -> t-head-tenantid
	if v := md.Get(strings.ToLower("t-head-tenantId")); len(v) == 0 || v[0] != "T1" {
		t.Errorf("tenantId metadata = %v, want [T1]", v)
	}
	if v := md.Get(strings.ToLower("t-head-traceId")); len(v) == 0 || v[0] != "trace-123" {
		t.Errorf("traceId metadata = %v, want [trace-123]", v)
	}
	if v := md.Get(strings.ToLower("t-head-userId")); len(v) == 0 || v[0] != "u9" {
		t.Errorf("userId metadata = %v, want [u9]", v)
	}
}

// TestFromGRPCError 验证 gRPC 错误可解包回 gobase 业务错误（ErrorInfo 细节）。
func TestFromGRPCError(t *testing.T) {
	// ErrNotFound().ToGRPCError() 序列化为带 ErrorInfo(Reason=1003) 的 status
	grpcErr := errors.ErrNotFound().ToGRPCError()

	be, ok := FromGRPCError(grpcErr)
	if !ok {
		t.Fatal("FromGRPCError 应识别业务错误")
	}
	if be.Code != errors.CodeNotFound {
		t.Errorf("解包 Code = %q, want %q", be.Code, errors.CodeNotFound)
	}

	// 非业务错误（普通错误）不识别
	if _, ok := FromGRPCError(context.Canceled); ok {
		t.Error("普通错误不应被识别为业务错误")
	}
}
