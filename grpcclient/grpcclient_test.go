package grpcclient

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/grpc/metadata"

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
