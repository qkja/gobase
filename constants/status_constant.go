package constants

/* 业务 code 前缀（统一在 gobase 定义，各服务引用，避免字符串漂移） */
const (
	AuditLogCodePrefix = "aud"
)

/* 通用状态：enable/disable（平台账号、租户管理员等） */
const (
	StatusEnable  = "enable"
	StatusDisable = "disable"
)

/* 租户生命周期状态 */
const (
	TenantStatusTrial     = "trial"
	TenantStatusPaid      = "paid"
	TenantStatusSuspended = "suspended"
	TenantStatusDeleted   = "deleted"
)

/* 平台权限点（P12：代码固定，用户不可增删；网关 RequirePermission 与前端 use_permission 对齐） */
const (
	PermTenantRead           = "tenant:read"
	PermTenantWrite          = "tenant:write"
	PermPlatformAccountRead  = "platform_account:read"
	PermPlatformAccountWrite = "platform_account:write"
	PermRoleRead             = "role:read"
	PermRoleWrite            = "role:write"
	PermAuditRead            = "audit:read"
)

/* 内置角色 key（P12：仅播种与迁移用标识；角色本体为数据库记录，可自定义增删） */
const (
	BuiltInRoleSuperAdmin = "super_admin"
	BuiltInRoleOps        = "ops"
	BuiltInRoleReadonly   = "readonly"
)

/* 审计 result / source / target_type（docs 操作审计日志需求说明 §4.1） */
const (
	AuditResultSuccess    = "success"
	AuditResultBizFailure = "biz_failure"
	AuditResultDenied     = "denied"
	AuditResultError      = "error"

	AuditSourceGateway      = "gateway"
	AuditSourceServiceEvent = "service_event"
	AuditSourceCLI          = "cli"
	AuditSourceSystem       = "system"

	AuditTargetTenant          = "tenant"
	AuditTargetPlatformAccount = "platform_account"
	AuditTargetRole            = "role"
	AuditTargetSession         = "session"
	AuditTargetAudit           = "audit"
	AuditTargetSystem          = "system"
)

/* 审计动作标识（集中定义，避免各服务拼字符串漂移） */
const (
	AuditActionTenantCreate        = "tenant.create"
	AuditActionTenantUpdate        = "tenant.update"
	AuditActionTenantDisable       = "tenant.disable"
	AuditActionTenantEnable        = "tenant.enable"
	AuditActionTenantDelete        = "tenant.delete"
	AuditActionTenantResetAdminPwd = "tenant.reset_admin_password"

	AuditActionAccountCreate         = "platform_account.create"
	AuditActionAccountUpdate         = "platform_account.update"
	AuditActionAccountStatusChange   = "platform_account.status_change"
	AuditActionAccountRolesChange    = "platform_account.roles_change"
	AuditActionAccountResetPassword  = "platform_account.reset_password"
	AuditActionAccountChangePassword = "platform_account.change_password"
	AuditActionAccountUnlock         = "platform_account.unlock"
	AuditActionAccountDelete         = "platform_account.delete"
	AuditActionAccountLocked         = "platform_account.locked"
	AuditActionAccountCliReset       = "platform_account.cli_reset"
	AuditActionAccountSeeded         = "platform_account.seeded"

	AuditActionRoleCreate = "role.create"
	AuditActionRoleUpdate = "role.update"
	AuditActionRoleDelete = "role.delete"

	AuditActionAuthLogin            = "auth.login"
	AuditActionAuthLoginFailed      = "auth.login_failed"
	AuditActionAuthLogout           = "auth.logout"
	AuditActionAuthTokenReplay      = "auth.token_replay"
	AuditActionSessionRevoke        = "session.revoke"
	AuditActionSessionRevokeCascade = "session.revoke_cascade"

	AuditActionAuditQuery  = "audit.query"
	AuditActionAuditExport = "audit.export"
)
