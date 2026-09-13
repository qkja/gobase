package errors

import (
	"context"
	"testing"

	"github.com/qkja/gobase/platform"
	"github.com/qkja/gobase/tenant"
)

func TestGetTenantMsg(t *testing.T) {
	const wantZh = "资源不存在"
	const wantEn = "Resource not found"

	// ctx 无租户信息 → 默认英文
	if got := GetTenantMsg(context.Background(), CodeNotFound); got != wantEn {
		t.Fatalf("no-ctx: got %q, want %q", got, wantEn)
	}

	// 中文租户（UILanguage zh-CN）→ 中文
	ctxZh := tenant.WithInfo(context.Background(), &tenant.Info{UILanguage: "zh-CN"})
	if got := GetTenantMsg(ctxZh, CodeNotFound); got != wantZh {
		t.Fatalf("zh: got %q, want %q", got, wantZh)
	}

	// 英文租户（UILanguage en_US）→ 英文
	ctx := tenant.WithInfo(context.Background(), &tenant.Info{UILanguage: "en_US"})
	if got := GetTenantMsg(ctx, CodeNotFound); got != wantEn {
		t.Fatalf("en: got %q, want %q", got, wantEn)
	}

	// 方法版本 err.GetTenantMsg(ctx)（取错误自身 code）
	if got := New(CodeNotFound).GetTenantMsg(ctx); got != wantEn {
		t.Fatalf("method: got %q, want %q", got, wantEn)
	}

	// nil 接收者安全
	var nilErr *BizError
	if got := nilErr.GetTenantMsg(ctx); got != "" {
		t.Fatalf("nil-receiver: got %q, want empty", got)
	}
}

func TestGetPlatformMsg(t *testing.T) {
	const wantZh = "资源不存在"
	const wantEn = "Resource not found"

	// ctx 无平台管理员信息 → 默认英文
	if got := GetPlatformMsg(context.Background(), CodeNotFound); got != wantEn {
		t.Fatalf("no-ctx: got %q, want %q", got, wantEn)
	}

	// 中文平台管理员（UILanguage zh-CN）→ 中文
	ctxZh := platform.WithInfo(context.Background(), &platform.Info{UILanguage: "zh-CN"})
	if got := GetPlatformMsg(ctxZh, CodeNotFound); got != wantZh {
		t.Fatalf("zh: got %q, want %q", got, wantZh)
	}

	// 英文平台管理员 → 英文
	ctx := platform.WithInfo(context.Background(), &platform.Info{UILanguage: "en_US"})
	if got := GetPlatformMsg(ctx, CodeNotFound); got != wantEn {
		t.Fatalf("en: got %q, want %q", got, wantEn)
	}

	// 方法版本 err.GetPlatformMsg(ctx)
	if got := New(CodeNotFound).GetPlatformMsg(ctx); got != wantEn {
		t.Fatalf("method: got %q, want %q", got, wantEn)
	}
}
