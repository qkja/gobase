package errors

// 业务错误码段位划分（按服务分段，每个服务 1000 个码）：
//
//	"0"          成功
//	"1000"-"1999" 公共段（全服务通用；网关自身错误也在本段）
//	"2000"-"2999" identityhubsvr（目录域 + 组织 + 用户 + 用户角色 + 同步）
//	"3000"-"3999" tenantmanagersvr（租户 + 租户管理员 + 租户管理角色）
//	"4000"-"4999" platformsvr（平台账号 + 平台角色）
//	"5000"-"5999" authnexussvr（三域会话 / 登录）
//	"6000"-"6999" auditsvr（审计）
//	"7000"-"7999" 预留段（网关错误使用公共段，本段留空）
//
// 全部错误码由 gobase 内置定义，业务直接使用，禁止自定义/注册新码。
// （对 gobase 本身的扩展属合规路径，各业务服务不得自行注册新码。）
//
// 公共码中 1002/1004/1005/1006/1009 的数值固定 —— 网关据此映射 HTTP 状态
// （400/401/403/504/429），不得改动；其余码可按服务演进重排。

/* 成功 */
const (
	// CodeOK 成功
	CodeOK = "0"
)

/* 公共错误码 1001-1999（全服务通用，数值固定，勿改） */
const (
	// CodeInternal 系统内部错误
	CodeInternal = "1001"
	// CodeInvalidArgument 参数无效
	CodeInvalidArgument = "1002"
	// CodeNotFound 资源不存在
	CodeNotFound = "1003"
	// CodeUnauthenticated 未认证或登录已过期
	CodeUnauthenticated = "1004"
	// CodePermissionDenied 无权限
	CodePermissionDenied = "1005"
	// CodeTimeout 请求超时
	CodeTimeout = "1006"
	// CodeAlreadyExists 资源已存在
	CodeAlreadyExists = "1007"
	// CodeConflict 依赖冲突
	CodeConflict = "1008"
	// CodeResourceExhausted 请求过于频繁
	CodeResourceExhausted = "1009"
)

/* 兜底 */
const (
	// CodeUnknown 未知错误
	CodeUnknown = "1999"
)

/*
identityhubsvr 错误码 2000-2999
子段：2001-2019 目录域；2020-2039 组织；2040-2059 用户；2060-2079 用户角色；2080-2099 同步
*/
const (
	// ---- 目录域（2001-2019）----

	// CodeDirNotFound 目录域不存在
	CodeDirNotFound = "2001"
	// CodeDirNameExists 租户内目录域名称已存在
	CodeDirNameExists = "2002"
	// CodeDirAlreadyDisabled 目录域已是禁用状态
	CodeDirAlreadyDisabled = "2003"
	// CodeDirAlreadyEnabled 目录域已是启用状态
	CodeDirAlreadyEnabled = "2004"
	// CodeDirNotEmpty 目录域非空（仍有组织 / 用户 / 用户角色挂靠），无法删除
	CodeDirNotEmpty = "2005"
	// CodeDirDisabled 目录域已禁用
	CodeDirDisabled = "2006"
	// CodeDirTypeImmutable 目录域类型创建后不可修改
	CodeDirTypeImmutable = "2007"
	// CodeDirLimitExceeded 目录域数量已达上限
	CodeDirLimitExceeded = "2008"

	// ---- 组织（2020-2039）----

	// CodeOrgNotFound 组织不存在
	CodeOrgNotFound = "2020"
	// CodeOrgParentNotFound 父组织不存在
	CodeOrgParentNotFound = "2021"
	// CodeOrgCycle 组织层级存在循环引用
	CodeOrgCycle = "2022"
	// CodeOrgHasChildren 存在子组织，无法删除
	CodeOrgHasChildren = "2023"
	// CodeOrgHasUsers 组织下仍有关联用户，无法删除
	CodeOrgHasUsers = "2024"
	// CodeOrgCrossDirectory 不能跨目录域操作组织
	CodeOrgCrossDirectory = "2025"
	// CodeOrgLevelExceeded 超过组织层级上限
	CodeOrgLevelExceeded = "2026"

	// ---- 用户（2040-2059）----

	// CodeUserNotFound 用户不存在
	CodeUserNotFound = "2040"
	// CodeUserNameExists 目录域内用户名称已存在（名称兼登录账号）
	CodeUserNameExists = "2041"
	// CodeUserAlreadyDisabled 用户已是禁用状态
	CodeUserAlreadyDisabled = "2042"
	// CodeUserAlreadyEnabled 用户已是启用状态
	CodeUserAlreadyEnabled = "2043"
	// CodeUserLocked 用户已锁定
	CodeUserLocked = "2044"
	// CodeUserDisabled 用户已禁用
	CodeUserDisabled = "2045"
	// CodeUserBadCredential 账号或密码错误
	CodeUserBadCredential = "2046"
	// CodeUserExternalSourceReadonly 外部源用户的身份字段不可修改，且不可删除（仅可停用）
	CodeUserExternalSourceReadonly = "2047"

	// ---- 用户角色（2060-2079）----

	// CodeUserRoleNotFound 用户角色不存在
	CodeUserRoleNotFound = "2060"
	// CodeUserRoleNameExists 目录域内用户角色名称已存在
	CodeUserRoleNameExists = "2061"
	// CodeUserRoleAlreadyDisabled 用户角色已是禁用状态
	CodeUserRoleAlreadyDisabled = "2062"
	// CodeUserRoleAlreadyEnabled 用户角色已是启用状态
	CodeUserRoleAlreadyEnabled = "2063"
	// CodeUserRoleMemberCrossDirectory 成员必须与用户角色属于同一目录域
	CodeUserRoleMemberCrossDirectory = "2064"

	// ---- 同步（2080-2099）----

	// CodeSyncUnsupportedDirectory 仅 ad / ldap 目录域支持同步配置，local 目录域拒绝读写
	CodeSyncUnsupportedDirectory = "2080"
)

/*
tenantmanagersvr 错误码 3000-3999
子段：3010-3099 租户；3100-3199 租户管理员；3200-3299 租户管理角色
*/
const (
	// ---- 租户（3010-3099）----

	// CodeTenantNotFound 租户不存在
	CodeTenantNotFound = "3010"
	// CodeTenantAlreadyDisabled 租户已是禁用状态
	CodeTenantAlreadyDisabled = "3011"
	// CodeTenantAlreadyEnabled 租户已是启用状态
	CodeTenantAlreadyEnabled = "3012"
	// CodeTenantDeleted 租户已删除
	CodeTenantDeleted = "3013"
	// CodeTenantDisabled 所属租户已停用，无法登录或操作
	CodeTenantDisabled = "3014"

	// ---- 租户管理员（3100-3199）----

	// CodeTenantAdminNotFound 租户管理员不存在
	CodeTenantAdminNotFound = "3100"
	// CodeTenantAdminNameExists 管理员名称已存在（名称兼登录账号，租户管理员范围内全局唯一）
	CodeTenantAdminNameExists = "3101"
	// CodeTenantAdminAlreadyDisabled 租户管理员已是禁用状态
	CodeTenantAdminAlreadyDisabled = "3102"
	// CodeTenantAdminAlreadyEnabled 租户管理员已是启用状态
	CodeTenantAdminAlreadyEnabled = "3103"
	// CodeTenantAdminLocked 租户管理员已锁定
	CodeTenantAdminLocked = "3104"
	// CodeTenantAdminDisabled 租户管理员已禁用
	CodeTenantAdminDisabled = "3105"
	// CodeTenantAdminBadCredential 账号或密码错误
	CodeTenantAdminBadCredential = "3106"
	// CodeTenantAdminLastSuperAdmin 需至少保留一个启用且绑定超级管理员角色的管理员
	CodeTenantAdminLastSuperAdmin = "3108"

	// ---- 租户管理角色（3200-3299）----

	// CodeTenantRoleNotFound 租户管理角色不存在
	CodeTenantRoleNotFound = "3200"
	// CodeTenantRoleNameExists 租户内角色名称已存在
	CodeTenantRoleNameExists = "3201"
	// CodeTenantRoleBuiltInImmutable 内置角色不可修改、删除或停用
	CodeTenantRoleBuiltInImmutable = "3202"
	// CodeTenantRoleInUse 该角色仍被管理员绑定，请先解绑
	CodeTenantRoleInUse = "3203"
	// CodeTenantRoleInvalidPageCode 页面权限码无效（不存在或作用域不符）
	CodeTenantRoleInvalidPageCode = "3204"
	// CodeTenantRoleScopeRequired 自定义角色至少需选择一个数据范围
	CodeTenantRoleScopeRequired = "3205"
)

/*
platformsvr 错误码 4000-4999
子段：4010-4099 平台账号；4100-4199 平台角色
*/
const (
	// ---- 平台账号（4010-4099）----

	// CodePlatformUserNotFound 平台账号不存在
	CodePlatformUserNotFound = "4010"
	// CodePlatformUserNameExists 账号名称已存在（名称兼登录账号，全局唯一）
	CodePlatformUserNameExists = "4011"
	// CodePlatformUserAlreadyDisabled 平台账号已是禁用状态
	CodePlatformUserAlreadyDisabled = "4012"
	// CodePlatformUserAlreadyEnabled 平台账号已是启用状态
	CodePlatformUserAlreadyEnabled = "4013"
	// CodePlatformUserLocked 平台账号已锁定
	CodePlatformUserLocked = "4014"
	// CodePlatformUserDisabled 平台账号已禁用
	CodePlatformUserDisabled = "4015"
	// CodePlatformUserBadCredential 账号或密码错误
	CodePlatformUserBadCredential = "4016"

	// ---- 平台角色（4100-4199）----

	// CodePlatformRoleNotFound 平台角色不存在
	CodePlatformRoleNotFound = "4100"
	// CodePlatformRoleNameExists 平台角色名称已存在（全局唯一）
	CodePlatformRoleNameExists = "4101"
	// CodePlatformRoleBuiltInImmutable 内置角色不可修改、删除或停用
	CodePlatformRoleBuiltInImmutable = "4102"
	// CodePlatformRoleInUse 该角色仍被账号绑定，请先解绑
	CodePlatformRoleInUse = "4103"
	// CodePlatformRoleInvalidPageCode 页面权限码无效（不存在或作用域不符）
	CodePlatformRoleInvalidPageCode = "4104"
)

/* authnexussvr 错误码 5000-5999（三域会话 / 登录，条件跨域共用，顺序编号） */
const (
	// CodeSessionInvalid 会话不存在或已失效
	CodeSessionInvalid = "5001"
	// CodeRefreshTokenReplayed 刷新凭证无效或已被吊销，请重新登录
	CodeRefreshTokenReplayed = "5002"
	// CodeMustChangePassword 请先修改初始密码
	CodeMustChangePassword = "5003"
)

/* auditsvr 错误码 6000-6999 */
const (
	// CodeAuditLogNotFound 审计记录不存在
	CodeAuditLogNotFound = "6001"
	// CodeAuditTimeRangeTooLarge 查询时间跨度超出上限
	CodeAuditTimeRangeTooLarge = "6002"
	// CodeAuditExportTooLarge 导出数据量超出上限，请收窄条件
	CodeAuditExportTooLarge = "6003"
	// CodeAuditTimeRangeRequired 请指定查询时间区间
	CodeAuditTimeRangeRequired = "6004"
)

// ===== 通用便捷错误实例（全局变量，包初始化时填充，多 goroutine 安全）=====

// 公共段
var (
	ErrInternal          = New(CodeInternal)
	ErrInvalidArgument   = New(CodeInvalidArgument)
	ErrNotFound          = New(CodeNotFound)
	ErrUnauthenticated   = New(CodeUnauthenticated)
	ErrPermissionDenied  = New(CodePermissionDenied)
	ErrTimeout           = New(CodeTimeout)
	ErrAlreadyExists     = New(CodeAlreadyExists)
	ErrConflict          = New(CodeConflict)
	ErrResourceExhausted = New(CodeResourceExhausted)
	ErrUnknown           = New(CodeUnknown)
)

// identityhubsvr —— 目录域
var (
	ErrDirNotFound        = New(CodeDirNotFound)
	ErrDirNameExists      = New(CodeDirNameExists)
	ErrDirAlreadyDisabled = New(CodeDirAlreadyDisabled)
	ErrDirAlreadyEnabled  = New(CodeDirAlreadyEnabled)
	ErrDirNotEmpty        = New(CodeDirNotEmpty)
	ErrDirDisabled        = New(CodeDirDisabled)
	ErrDirTypeImmutable   = New(CodeDirTypeImmutable)
	ErrDirLimitExceeded   = New(CodeDirLimitExceeded)
)

// identityhubsvr —— 组织
var (
	ErrOrgNotFound       = New(CodeOrgNotFound)
	ErrOrgParentNotFound = New(CodeOrgParentNotFound)
	ErrOrgCycle          = New(CodeOrgCycle)
	ErrOrgHasChildren    = New(CodeOrgHasChildren)
	ErrOrgHasUsers       = New(CodeOrgHasUsers)
	ErrOrgCrossDirectory = New(CodeOrgCrossDirectory)
	ErrOrgLevelExceeded  = New(CodeOrgLevelExceeded)
)

// identityhubsvr —— 用户
var (
	ErrUserNotFound               = New(CodeUserNotFound)
	ErrUserNameExists             = New(CodeUserNameExists)
	ErrUserAlreadyDisabled        = New(CodeUserAlreadyDisabled)
	ErrUserAlreadyEnabled         = New(CodeUserAlreadyEnabled)
	ErrUserLocked                 = New(CodeUserLocked)
	ErrUserDisabled               = New(CodeUserDisabled)
	ErrUserBadCredential          = New(CodeUserBadCredential)
	ErrUserExternalSourceReadonly = New(CodeUserExternalSourceReadonly)
)

// identityhubsvr —— 用户角色
var (
	ErrUserRoleNotFound             = New(CodeUserRoleNotFound)
	ErrUserRoleNameExists           = New(CodeUserRoleNameExists)
	ErrUserRoleAlreadyDisabled      = New(CodeUserRoleAlreadyDisabled)
	ErrUserRoleAlreadyEnabled       = New(CodeUserRoleAlreadyEnabled)
	ErrUserRoleMemberCrossDirectory = New(CodeUserRoleMemberCrossDirectory)
)

// identityhubsvr —— 同步
var (
	ErrSyncUnsupportedDirectory = New(CodeSyncUnsupportedDirectory)
)

// tenantmanagersvr —— 租户
var (
	ErrTenantNotFound        = New(CodeTenantNotFound)
	ErrTenantAlreadyDisabled = New(CodeTenantAlreadyDisabled)
	ErrTenantAlreadyEnabled  = New(CodeTenantAlreadyEnabled)
	ErrTenantDeleted         = New(CodeTenantDeleted)
	ErrTenantDisabled        = New(CodeTenantDisabled)
)

// tenantmanagersvr —— 租户管理员
var (
	ErrTenantAdminNotFound        = New(CodeTenantAdminNotFound)
	ErrTenantAdminNameExists      = New(CodeTenantAdminNameExists)
	ErrTenantAdminAlreadyDisabled = New(CodeTenantAdminAlreadyDisabled)
	ErrTenantAdminAlreadyEnabled  = New(CodeTenantAdminAlreadyEnabled)
	ErrTenantAdminLocked          = New(CodeTenantAdminLocked)
	ErrTenantAdminDisabled        = New(CodeTenantAdminDisabled)
	ErrTenantAdminBadCredential   = New(CodeTenantAdminBadCredential)
	ErrTenantAdminLastSuperAdmin  = New(CodeTenantAdminLastSuperAdmin)
)

// tenantmanagersvr —— 租户管理角色
var (
	ErrTenantRoleNotFound         = New(CodeTenantRoleNotFound)
	ErrTenantRoleNameExists       = New(CodeTenantRoleNameExists)
	ErrTenantRoleBuiltInImmutable = New(CodeTenantRoleBuiltInImmutable)
	ErrTenantRoleInUse            = New(CodeTenantRoleInUse)
	ErrTenantRoleInvalidPageCode  = New(CodeTenantRoleInvalidPageCode)
	ErrTenantRoleScopeRequired    = New(CodeTenantRoleScopeRequired)
)

// platformsvr —— 平台账号
var (
	ErrPlatformUserNotFound        = New(CodePlatformUserNotFound)
	ErrPlatformUserNameExists      = New(CodePlatformUserNameExists)
	ErrPlatformUserAlreadyDisabled = New(CodePlatformUserAlreadyDisabled)
	ErrPlatformUserAlreadyEnabled  = New(CodePlatformUserAlreadyEnabled)
	ErrPlatformUserLocked          = New(CodePlatformUserLocked)
	ErrPlatformUserDisabled        = New(CodePlatformUserDisabled)
	ErrPlatformUserBadCredential   = New(CodePlatformUserBadCredential)
)

// platformsvr —— 平台角色
var (
	ErrPlatformRoleNotFound         = New(CodePlatformRoleNotFound)
	ErrPlatformRoleNameExists       = New(CodePlatformRoleNameExists)
	ErrPlatformRoleBuiltInImmutable = New(CodePlatformRoleBuiltInImmutable)
	ErrPlatformRoleInUse            = New(CodePlatformRoleInUse)
	ErrPlatformRoleInvalidPageCode  = New(CodePlatformRoleInvalidPageCode)
)

// authnexussvr —— 三域会话 / 登录
var (
	ErrSessionInvalid       = New(CodeSessionInvalid)
	ErrRefreshTokenReplayed = New(CodeRefreshTokenReplayed)
	ErrMustChangePassword   = New(CodeMustChangePassword)
)

// auditsvr —— 审计
var (
	ErrAuditLogNotFound       = New(CodeAuditLogNotFound)
	ErrAuditTimeRangeTooLarge = New(CodeAuditTimeRangeTooLarge)
	ErrAuditExportTooLarge    = New(CodeAuditExportTooLarge)
	ErrAuditTimeRangeRequired = New(CodeAuditTimeRangeRequired)
)
