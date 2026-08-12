package logger

import (
	"context"
	"fmt"
	"slices"

	"github.com/qkja/gobase/config"
	"github.com/qkja/gobase/tenant"
	"github.com/sirupsen/logrus"
)

// 租户级 debug：排障时为特定租户开启 debug 日志，不影响其他租户。
var (
	// tenantDebugLogger 恒为 DebugLevel，绕开 rootLogger 的 info 级别过滤。
	tenantDebugLogger *logrus.Logger
	// debugTenantIDs 开启 debug 的租户列表（logger.debug_tenant_ids）。
	debugTenantIDs []string
)

// reloadDebugTenants 从配置重读 logger.debug_tenant_ids（热更新入口）。
func reloadDebugTenants() {
	arr := config.GetValueArray("logger.debug_tenant_ids")
	debugTenantIDs = make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok && s != "" {
			debugTenantIDs = append(debugTenantIDs, s)
		}
	}
}

// debugTenantEnabled 判断指定租户是否开启 debug：
//   - 全局级别为 debug/trace 时，所有租户开启；
//   - 否则仅 debug_tenant_ids 列表中的租户开启。
func debugTenantEnabled(tenantID string) bool {
	if rootLogger.IsLevelEnabled(logrus.DebugLevel) {
		return true
	}
	if tenantID == "" {
		return false
	}
	return slices.Contains(debugTenantIDs, tenantID)
}

// DebugTenant 租户级 debug 日志。
// 仅当「全局 debug 开启」或「ctx 中租户在 debug_tenant_ids 列表」时输出，
// 自动带 [TenantId:xxx] [TraceId:xxx] 前缀。
//
// 用法：在业务逻辑中需要细粒度排查的地方调用。
//
//	logger.DebugTenant(ctx, "processing order: %v", order)
func DebugTenant(ctx context.Context, format string, v ...any) {
	info := tenant.GetInfo(ctx)
	tenantID := ""
	if info != nil {
		tenantID = info.TenantID
	}
	if !debugTenantEnabled(tenantID) {
		return
	}

	msg := fmt.Sprintf(format, v...)
	prefix := formatCtx(ctx) // [TraceId:xxx] 前缀（无 trace 则空）
	if tenantID != "" {
		tenantDebugLogger.Debugf("[TenantId:%s] %s%s", tenantID, prefix, msg)
	} else {
		tenantDebugLogger.Debugf("%s%s", prefix, msg)
	}
}
