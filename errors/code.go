package errors

// 业务错误码段位划分：
//
//	"0"         成功
//	"1001"-"1999" 通用错误码
//	"2000"-"2999" 目录域
//	"3000"-"3999" 组织架构
//	"4000"-"4999" 用户
//	"5000"-"5099" 租户
//	"5100"-"5199" 平台账号与角色
//	"5200"-"5299" 审计
//	"1999"      未知错误（兜底）
//
// 全部错误码由 gobase 内置定义，业务直接使用，禁止自定义/注册新码。 （对 gobase 本身的扩展属合规路径，各业务服务不得自行注册新码。）

/* 成功 */
const (
	// CodeOK 成功
	CodeOK = "0"
)

/* 通用错误码 1001-1999 */
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

/* 目录域错误码 2000-2999 */
const (
	// CodeDirNotFound 目录域不存在
	CodeDirNotFound = "2001"
	// CodeDirDomainExists domain 已存在
	CodeDirDomainExists = "2002"
	// CodeDirAlreadyDisabled 已是禁用状态
	CodeDirAlreadyDisabled = "2003"
	// CodeDirAlreadyEnabled 已是启用状态
	CodeDirAlreadyEnabled = "2004"
	// CodeDirNotEmpty 目录域非空（有待清理的资源）
	CodeDirNotEmpty = "2005"
	// CodeDirDisabled 目录域已禁用
	CodeDirDisabled = "2006"
	// CodeDirCreateFailed 新增目录域失败
	CodeDirCreateFailed = "2007"
	// CodeDirUpdateFailed 编辑目录域失败
	CodeDirUpdateFailed = "2008"
	// CodeDirGetFailed 查询目录域失败
	CodeDirGetFailed = "2009"
	// CodeDirDeleteFailed 删除目录域失败
	CodeDirDeleteFailed = "2010"
)

/* 组织架构错误码 3000-3999 */
const (
	// CodeOrgNotFound 组织不存在
	CodeOrgNotFound = "3001"
	// CodeOrgCreateFailed 新增组织失败
	CodeOrgCreateFailed = "3002"
	// CodeOrgParentNotFound 父组织不存在
	CodeOrgParentNotFound = "3003"
	// CodeOrgCycle 循环引用
	CodeOrgCycle = "3004"
	// CodeOrgHasChildren 有子组织，无法删除
	CodeOrgHasChildren = "3005"
	// CodeOrgHasUsers 有关联用户，无法删除
	CodeOrgHasUsers = "3006"
	// CodeOrgCrossDirectory 跨目录域操作
	CodeOrgCrossDirectory = "3007"
	// CodeOrgLevelExceeded 超过层级限制
	CodeOrgLevelExceeded = "3008"
	// CodeOrgUpdateFailed 编辑组织失败
	CodeOrgUpdateFailed = "3009"
	// CodeOrgGetFailed 查询组织失败
	CodeOrgGetFailed = "3010"
	// CodeOrgDeleteFailed 删除组织失败
	CodeOrgDeleteFailed = "3011"
	// CodeOrgMoveFailed 移动组织失败
	CodeOrgMoveFailed = "3012"
)

/* 用户错误码 4000-4999 */
const (
	// CodeUserNotFound 用户不存在
	CodeUserNotFound = "4001"
	// CodeUserUsernameExists 用户名已存在
	CodeUserUsernameExists = "4002"
	// CodeUserAlreadyDisabled 已是禁用状态
	CodeUserAlreadyDisabled = "4003"
	// CodeUserAlreadyEnabled 已是启用状态
	CodeUserAlreadyEnabled = "4004"
	// CodeUserOrgAlreadySecondary 组织已是副组织
	CodeUserOrgAlreadySecondary = "4005"
	// CodeUserOrgIsPrimary 组织是主组织
	CodeUserOrgIsPrimary = "4006"
	// CodeUserOrgNotSecondary 组织不是副组织
	CodeUserOrgNotSecondary = "4007"
	// CodeUserPrimaryOrgNotFound 主组织不存在
	CodeUserPrimaryOrgNotFound = "4008"
	// CodeUserCreateFailed 新增用户失败
	CodeUserCreateFailed = "4009"
	// CodeUserUpdateFailed 编辑用户失败
	CodeUserUpdateFailed = "4010"
	// CodeUserGetFailed 查询用户失败
	CodeUserGetFailed = "4011"
	// CodeUserDeleteFailed 删除用户失败
	CodeUserDeleteFailed = "4012"
	// CodeUserEnableFailed 启用用户失败
	CodeUserEnableFailed = "4013"
	// CodeUserDisableFailed 禁用用户失败
	CodeUserDisableFailed = "4014"
)

/* 租户错误码 5000-5099 */
const (
	// CodeTenantNotFound 租户不存在
	CodeTenantNotFound = "5001"
	// CodeTenantDomainExists 租户 domain 已存在
	CodeTenantDomainExists = "5002"
	// CodeTenantPhoneExists 手机号已被使用
	CodeTenantPhoneExists = "5003"
	// CodeTenantEmailExists 邮箱已被使用
	CodeTenantEmailExists = "5004"
	// CodeTenantAlreadyDisabled 租户已是禁用状态
	CodeTenantAlreadyDisabled = "5005"
	// CodeTenantAlreadyEnabled 租户已是启用状态
	CodeTenantAlreadyEnabled = "5006"
	// CodeTenantDeleted 租户已删除
	CodeTenantDeleted = "5007"
	// CodeTenantDomainImmutable 租户 domain 不可修改
	CodeTenantDomainImmutable = "5008"
	// CodeTenantCreateFailed 创建租户失败
	CodeTenantCreateFailed = "5009"
	// CodeTenantUpdateFailed 编辑租户失败
	CodeTenantUpdateFailed = "5010"
	// CodeTenantGetFailed 查询租户失败
	CodeTenantGetFailed = "5011"
	// CodeTenantDeleteFailed 删除租户失败
	CodeTenantDeleteFailed = "5012"
)

/* 平台账号与角色错误码 5100-5199 */
const (
	// CodePlatformAccountNotFound 平台账号不存在
	CodePlatformAccountNotFound = "5101"
	// CodePlatformUsernameExists 账号名已被使用
	CodePlatformUsernameExists = "5102"
	// CodePlatformPhoneExists 手机号已被使用
	CodePlatformPhoneExists = "5103"
	// CodePlatformEmailExists 邮箱已被使用
	CodePlatformEmailExists = "5104"
	// CodePlatformAccountLocked 账号已锁定
	CodePlatformAccountLocked = "5105"
	// CodePlatformAccountDisabled 账号已禁用
	CodePlatformAccountDisabled = "5106"
	// CodePlatformLastAdmin 系统需至少保留一个拥有该权限的启用账号
	CodePlatformLastAdmin = "5107"
	// CodePlatformCannotOperateSelf 不能对自己执行此操作
	CodePlatformCannotOperateSelf = "5108"
	// CodePlatformMustChangePassword 请先修改初始密码
	CodePlatformMustChangePassword = "5109"
	// CodePlatformRoleNotFound 角色不存在
	CodePlatformRoleNotFound = "5110"
	// CodePlatformRoleKeyExists 角色标识已被使用
	CodePlatformRoleKeyExists = "5111"
	// CodePlatformRoleNameExists 角色名称已被使用
	CodePlatformRoleNameExists = "5112"
	// CodePlatformRoleBuiltInImmutable 内置角色不可修改或删除
	CodePlatformRoleBuiltInImmutable = "5113"
	// CodePlatformRoleInUse 该角色仍被账号使用，请先解绑
	CodePlatformRoleInUse = "5114"
	// CodePlatformInvalidPermission 权限点无效
	CodePlatformInvalidPermission = "5115"
	// CodePlatformRoleRequired 账号至少需绑定一个角色
	CodePlatformRoleRequired = "5116"
	// CodePlatformSessionNotFound 会话不存在或已失效
	CodePlatformSessionNotFound = "5117"
	// CodePlatformRefreshTokenReplayed 凭证异常，请重新登录
	CodePlatformRefreshTokenReplayed = "5118"
)

/* 审计错误码 5200-5299 */
const (
	// CodeAuditLogNotFound 审计记录不存在
	CodeAuditLogNotFound = "5201"
	// CodeAuditTimeRangeTooLarge 查询时间跨度超出上限
	CodeAuditTimeRangeTooLarge = "5202"
	// CodeAuditExportTooLarge 导出数据量超出上限，请收窄条件
	CodeAuditExportTooLarge = "5203"
	// CodeAuditTimeRangeRequired 请指定查询时间区间
	CodeAuditTimeRangeRequired = "5204"
)

// 通用便捷错误实例（全局变量，包初始化时填充，线程安全）
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

// 目录域便捷错误实例
var (
	ErrDirNotFound        = New(CodeDirNotFound)
	ErrDirDomainExists    = New(CodeDirDomainExists)
	ErrDirAlreadyDisabled = New(CodeDirAlreadyDisabled)
	ErrDirAlreadyEnabled  = New(CodeDirAlreadyEnabled)
	ErrDirNotEmpty        = New(CodeDirNotEmpty)
	ErrDirDisabled        = New(CodeDirDisabled)
	ErrDirCreateFailed    = New(CodeDirCreateFailed)
	ErrDirUpdateFailed    = New(CodeDirUpdateFailed)
	ErrDirGetFailed       = New(CodeDirGetFailed)
	ErrDirDeleteFailed    = New(CodeDirDeleteFailed)
)

// 组织架构便捷错误实例
var (
	ErrOrgNotFound       = New(CodeOrgNotFound)
	ErrOrgCreateFailed   = New(CodeOrgCreateFailed)
	ErrOrgParentNotFound = New(CodeOrgParentNotFound)
	ErrOrgCycle          = New(CodeOrgCycle)
	ErrOrgHasChildren    = New(CodeOrgHasChildren)
	ErrOrgHasUsers       = New(CodeOrgHasUsers)
	ErrOrgCrossDirectory = New(CodeOrgCrossDirectory)
	ErrOrgLevelExceeded  = New(CodeOrgLevelExceeded)
	ErrOrgUpdateFailed   = New(CodeOrgUpdateFailed)
	ErrOrgGetFailed      = New(CodeOrgGetFailed)
	ErrOrgDeleteFailed   = New(CodeOrgDeleteFailed)
	ErrOrgMoveFailed     = New(CodeOrgMoveFailed)
)

// 用户便捷错误实例
var (
	ErrUserNotFound            = New(CodeUserNotFound)
	ErrUserUsernameExists      = New(CodeUserUsernameExists)
	ErrUserAlreadyDisabled     = New(CodeUserAlreadyDisabled)
	ErrUserAlreadyEnabled      = New(CodeUserAlreadyEnabled)
	ErrUserOrgAlreadySecondary = New(CodeUserOrgAlreadySecondary)
	ErrUserOrgIsPrimary        = New(CodeUserOrgIsPrimary)
	ErrUserOrgNotSecondary     = New(CodeUserOrgNotSecondary)
	ErrUserPrimaryOrgNotFound  = New(CodeUserPrimaryOrgNotFound)
	ErrUserCreateFailed        = New(CodeUserCreateFailed)
	ErrUserUpdateFailed        = New(CodeUserUpdateFailed)
	ErrUserGetFailed           = New(CodeUserGetFailed)
	ErrUserDeleteFailed        = New(CodeUserDeleteFailed)
	ErrUserEnableFailed        = New(CodeUserEnableFailed)
	ErrUserDisableFailed       = New(CodeUserDisableFailed)
)

// 租户便捷错误实例
var (
	ErrTenantNotFound        = New(CodeTenantNotFound)
	ErrTenantDomainExists    = New(CodeTenantDomainExists)
	ErrTenantPhoneExists     = New(CodeTenantPhoneExists)
	ErrTenantEmailExists     = New(CodeTenantEmailExists)
	ErrTenantAlreadyDisabled = New(CodeTenantAlreadyDisabled)
	ErrTenantAlreadyEnabled  = New(CodeTenantAlreadyEnabled)
	ErrTenantDeleted         = New(CodeTenantDeleted)
	ErrTenantDomainImmutable = New(CodeTenantDomainImmutable)
	ErrTenantCreateFailed    = New(CodeTenantCreateFailed)
	ErrTenantUpdateFailed    = New(CodeTenantUpdateFailed)
	ErrTenantGetFailed       = New(CodeTenantGetFailed)
	ErrTenantDeleteFailed    = New(CodeTenantDeleteFailed)
)

// 平台账号与角色便捷错误实例
var (
	ErrPlatformAccountNotFound      = New(CodePlatformAccountNotFound)
	ErrPlatformUsernameExists       = New(CodePlatformUsernameExists)
	ErrPlatformPhoneExists          = New(CodePlatformPhoneExists)
	ErrPlatformEmailExists          = New(CodePlatformEmailExists)
	ErrPlatformAccountLocked        = New(CodePlatformAccountLocked)
	ErrPlatformAccountDisabled      = New(CodePlatformAccountDisabled)
	ErrPlatformLastAdmin            = New(CodePlatformLastAdmin)
	ErrPlatformCannotOperateSelf    = New(CodePlatformCannotOperateSelf)
	ErrPlatformMustChangePassword   = New(CodePlatformMustChangePassword)
	ErrPlatformRoleNotFound         = New(CodePlatformRoleNotFound)
	ErrPlatformRoleKeyExists        = New(CodePlatformRoleKeyExists)
	ErrPlatformRoleNameExists       = New(CodePlatformRoleNameExists)
	ErrPlatformRoleBuiltInImmutable = New(CodePlatformRoleBuiltInImmutable)
	ErrPlatformRoleInUse            = New(CodePlatformRoleInUse)
	ErrPlatformInvalidPermission    = New(CodePlatformInvalidPermission)
	ErrPlatformRoleRequired         = New(CodePlatformRoleRequired)
	ErrPlatformSessionNotFound      = New(CodePlatformSessionNotFound)
	ErrPlatformRefreshTokenReplayed = New(CodePlatformRefreshTokenReplayed)
)

// 审计便捷错误实例
var (
	ErrAuditLogNotFound       = New(CodeAuditLogNotFound)
	ErrAuditTimeRangeTooLarge = New(CodeAuditTimeRangeTooLarge)
	ErrAuditExportTooLarge    = New(CodeAuditExportTooLarge)
	ErrAuditTimeRangeRequired = New(CodeAuditTimeRangeRequired)
)
