package trace

import (
	"context"
	"testing"
)

func TestWithInfoGetInfoRoundTrip(t *testing.T) {
	ctx := context.Background()
	ctx = WithInfo(ctx, &Info{TraceID: "t1", UserID: "u1"})

	info := GetInfo(ctx)
	if info == nil {
		t.Fatal("GetInfo 返回 nil")
	}
	if info.TraceID != "t1" || info.UserID != "u1" {
		t.Errorf("GetInfo = %+v, want traceID=t1 userID=u1", info)
	}
}

func TestGetInfoNilCtx(t *testing.T) {
	if info := GetInfo(nil); info != nil {
		t.Errorf("GetInfo(nil) = %v, want nil", info)
	}
}

func TestWithInfoNilCtx(t *testing.T) {
	// 不应 panic，回退 Background
	ctx := WithInfo(nil, &Info{TraceID: "t1"})
	if info := GetInfo(ctx); info == nil || info.TraceID != "t1" {
		t.Errorf("WithInfo(nil) 未正确写入, info=%v", info)
	}
}

func TestWithTraceIDCopiesNotMutates(t *testing.T) {
	shared := &Info{TraceID: "orig", UserID: "u1"}
	ctx := WithInfo(context.Background(), shared)

	// WithTraceID 应拷贝，不污染 shared
	ctx2 := WithTraceID(ctx, "new")
	if shared.TraceID != "orig" {
		t.Errorf("shared.TraceID 被修改 = %q, want orig", shared.TraceID)
	}
	if info := GetInfo(ctx2); info.TraceID != "new" || info.UserID != "u1" {
		t.Errorf("WithTraceID 未保留其他字段: %+v", info)
	}
}

func TestTraceIDFromCtx(t *testing.T) {
	if got := TraceIDFromCtx(context.Background()); got != "" {
		t.Errorf("空 ctx TraceID = %q, want \"\"", got)
	}
	ctx := WithInfo(context.Background(), &Info{TraceID: "abc"})
	if got := TraceIDFromCtx(ctx); got != "abc" {
		t.Errorf("TraceID = %q, want abc", got)
	}
}
