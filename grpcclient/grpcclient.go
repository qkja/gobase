// Package grpcclient 提供服务间 gRPC 调用的统一封装：
// discovery 寻址 + 连接池 + 超时 + tenant/trace metadata 注入 + 业务错误解包。
//
// 前置：调用前需先 discovery.Init(cfgDir)（解析各服务地址）与加载配置。
package grpcclient

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/qkja/gobase/config"
	"github.com/qkja/gobase/constants"
	"github.com/qkja/gobase/discovery"
	"github.com/qkja/gobase/tenant"
	"github.com/qkja/gobase/trace"
)

// connPool 按服务名缓存连接；*grpc.ClientConn 并发安全可复用，内部自动重连。
var connPool sync.Map // svcName -> *grpc.ClientConn

// Dial 按服务名解析地址（discovery.GetAddress）并返回（或复用池化）连接。
// 返回的连接为懒连接（grpc 默认行为），未就绪时在首次 RPC 才暴露错误。
func Dial(ctx context.Context, svcName string) (*grpc.ClientConn, error) {
	if c, ok := connPool.Load(svcName); ok {
		return c.(*grpc.ClientConn), nil
	}
	addr, ok := discovery.GetAddress(svcName)
	if !ok {
		return nil, fmt.Errorf("grpcclient: 服务 %s 未在 discovery 中配置", svcName)
	}
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	actual, loaded := connPool.LoadOrStore(svcName, conn)
	if loaded {
		// 并发首拨时已有连接，释放刚创建的这个，避免泄漏
		_ = conn.Close()
	}
	return actual.(*grpc.ClientConn), nil
}

// Call 通用 gRPC 调用：解析地址→取连接→注入 tenant/trace metadata→套超时→执行 fn。
// fn 收到的 ctx 已注入 outgoing metadata 并带上 grpc.timeout 超时（毫秒，默认 5000）。
// 返回 fn 的原始结果/错误（不做跨进程业务码解包）。
func Call[T any](ctx context.Context, svcName string, fn func(ctx context.Context, conn *grpc.ClientConn) (T, error)) (T, error) {
	conn, err := Dial(ctx, svcName)
	if err != nil {
		var zero T
		return zero, err
	}

	callCtx := injectMetadata(ctx)
	timeout := time.Duration(config.GetValueIntDefault("grpc.timeout", 5000)) * time.Millisecond
	var cancel context.CancelFunc
	if timeout > 0 {
		callCtx, cancel = context.WithTimeout(callCtx, timeout)
		defer cancel()
	}

	return fn(callCtx, conn)
}

// injectMetadata 将租户上下文与 trace 头写入 outgoing gRPC metadata。
// gRPC metadata 键必须小写，故对 constants 头名统一 ToLower。
// 合并 ctx 上已有的 outgoing metadata，不覆盖。
func injectMetadata(ctx context.Context) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	if md == nil {
		md = metadata.New(map[string]string{})
	}
	if info := tenant.GetInfo(ctx); info != nil {
		if info.TenantID != "" {
			md.Append(strings.ToLower(constants.TENANT_HEAD_ID), info.TenantID)
		}
		if info.Language != "" {
			md.Append(strings.ToLower(constants.TENANT_HEAD_LANGUAGE), info.Language)
		}
		if info.UILanguage != "" {
			md.Append(strings.ToLower(constants.TENANT_HEAD_UI_LANGUAGE), info.UILanguage)
		}
	}
	if info := trace.GetInfo(ctx); info != nil {
		pairs := []struct{ k, v string }{
			{constants.TRACE_HEAD_ID, info.TraceID},
			{constants.TRACE_HEAD_RPC_ID, info.RPCID},
			{constants.TRACE_HEAD_SAMPLED, info.Sampled},
			{constants.TRACE_HEAD_USER_ID, info.UserID},
			{constants.TRACE_HEAD_USER_NAME, info.UserName},
			{constants.TRACE_HEAD_REMOTE_APPNAME, config.GetValueStringDefault("application.name", "")},
		}
		for _, p := range pairs {
			if p.v != "" {
				md.Append(strings.ToLower(p.k), p.v)
			}
		}
	}
	return metadata.NewOutgoingContext(ctx, md)
}
