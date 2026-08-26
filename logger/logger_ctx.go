package logger

import (
	"context"
	"fmt"

	"github.com/qkja/gobase/trace"
)

func formatCtx(ctx context.Context) string {
	if info := trace.GetInfo(ctx); info != nil && info.TraceID != "" {
		return fmt.Sprintf("[TraceId:%s] ", info.TraceID)
	}
	return ""
}

// InfoCtx 带上下文的 Info 日志。自动从 ctx 提取 TraceId 作为前缀。
func InfoCtx(ctx context.Context, format string, v ...any) {
	rootLogger.Infof(formatCtx(ctx)+format, v...)
}

// ErrorCtx 带上下文的 Error 日志。
func ErrorCtx(ctx context.Context, format string, v ...any) {
	rootLogger.Errorf(formatCtx(ctx)+format, v...)
}

// WarnCtx 带上下文的 Warn 日志。
func WarnCtx(ctx context.Context, format string, v ...any) {
	rootLogger.Warnf(formatCtx(ctx)+format, v...)
}

// DebugCtx 带上下文的 Debug 日志。
func DebugCtx(ctx context.Context, format string, v ...any) {
	rootLogger.Debugf(formatCtx(ctx)+format, v...)
}
