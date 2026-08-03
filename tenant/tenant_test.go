package tenant

import (
	"context"
	"testing"
)

func TestWithInfo_GetInfo_RoundTrip(t *testing.T) {
	info := &Info{TenantID: "t1", Language: "zh-CN", UILanguage: "zh-CN"}
	ctx := WithInfo(context.Background(), info)
	got := GetInfo(ctx)
	if got == nil || got.TenantID != "t1" || got.Language != "zh-CN" {
		t.Fatalf("GetInfo = %+v, want TenantID=t1 Language=zh-CN", got)
	}
}

func TestGetInfo_Missing_ReturnsNil(t *testing.T) {
	if got := GetInfo(context.Background()); got != nil {
		t.Fatalf("GetInfo(empty) = %+v, want nil", got)
	}
}

// 派生 ctx（子 ctx 又塞入其他值）仍应能读到租户信息。
func TestGetInfo_ChildContext_Inherits(t *testing.T) {
	ctx := WithInfo(context.Background(), &Info{TenantID: "t2"})
	child := context.WithValue(ctx, "some-key", "v")
	if got := GetInfo(child); got == nil || got.TenantID != "t2" {
		t.Fatalf("GetInfo(child) = %+v, want t2", got)
	}
}
