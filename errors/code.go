package errors

import "google.golang.org/grpc/codes"

// 业务错误码段位划分：
//
//	0         成功（沿袭既有 wire 约定 "0"，见 proto 响应与 server/rsp）
//	1001-1999 通用错误码（本包定义）
//	2000+     业务错误码（由业务模块 init 阶段通过 Register 注册）
//	1999      未知错误（兜底，用于 FromError 未识别时）
//
// 通用段位错误码编号分配：
//
//	1001 系统内部错误 | 1002 参数无效 | 1003 资源不存在 | 1004 未认证
//	1005 无权限 | 1006 请求超时 | 1007 已存在 | 1008 依赖冲突 | 1009 请求过于频繁

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

// errorMeta 错误码元数据：中英文默认消息 + 对应的 gRPC 状态码
type errorMeta struct {
	msgZh string
	msgEn string
	grpc  codes.Code
}

// registry 错误码注册表：code -> 元数据
//
// 用 map 而非 switch，理由：
//  1. FromError 需要「字符串 -> 元数据」的反查，map 天然支持；
//  2. 业务段位可通过 Register 在 init 阶段扩展，switch 必须改本包源码；
//  3. 本表仅 init 阶段填充、之后只读，无并发问题。
var registry = map[string]errorMeta{
	CodeOK:                {msgZh: "成功",        msgEn: "Success",      grpc: codes.OK},
	CodeInternal:          {msgZh: "系统内部错误", msgEn: "Internal error", grpc: codes.Internal},
	CodeInvalidArgument:   {msgZh: "参数无效",     msgEn: "Invalid argument", grpc: codes.InvalidArgument},
	CodeNotFound:          {msgZh: "资源不存在",   msgEn: "Resource not found", grpc: codes.NotFound},
	CodeUnauthenticated:   {msgZh: "未认证或登录已过期", msgEn: "Unauthenticated", grpc: codes.Unauthenticated},
	CodePermissionDenied:  {msgZh: "无权限",       msgEn: "Permission denied", grpc: codes.PermissionDenied},
	CodeTimeout:           {msgZh: "请求超时",     msgEn: "Request timeout", grpc: codes.DeadlineExceeded},
	CodeAlreadyExists:     {msgZh: "资源已存在",   msgEn: "Resource already exists", grpc: codes.AlreadyExists},
	CodeConflict:          {msgZh: "依赖冲突",     msgEn: "Conflict", grpc: codes.FailedPrecondition},
	CodeResourceExhausted: {msgZh: "请求过于频繁", msgEn: "Too many requests", grpc: codes.ResourceExhausted},
	CodeUnknown:           {msgZh: "未知错误",     msgEn: "Unknown error", grpc: codes.Unknown},
	// 目录域（2000-2999）
	CodeDirNotFound:        {msgZh: "目录域不存在",   msgEn: "Directory not found", grpc: codes.NotFound},
	CodeDirDomainExists:    {msgZh: "业务域标识已存在", msgEn: "Domain already exists", grpc: codes.AlreadyExists},
	CodeDirAlreadyDisabled: {msgZh: "目录域已是禁用状态", msgEn: "Directory already disabled", grpc: codes.FailedPrecondition},
	CodeDirAlreadyEnabled:  {msgZh: "目录域已是启用状态", msgEn: "Directory already enabled", grpc: codes.FailedPrecondition},
	CodeDirNotEmpty:        {msgZh: "目录域非空，无法删除", msgEn: "Directory not empty", grpc: codes.FailedPrecondition},
	CodeDirDisabled:        {msgZh: "目录域已禁用",   msgEn: "Directory disabled", grpc: codes.FailedPrecondition},
	CodeDirCreateFailed:    {msgZh: "新增目录域失败", msgEn: "Create directory failed", grpc: codes.Internal},
	CodeDirUpdateFailed:    {msgZh: "编辑目录域失败", msgEn: "Update directory failed", grpc: codes.Internal},
	CodeDirGetFailed:       {msgZh: "查询目录域失败", msgEn: "Get directory failed", grpc: codes.Internal},
	CodeDirDeleteFailed:    {msgZh: "删除目录域失败", msgEn: "Delete directory failed", grpc: codes.Internal},
	// 组织架构（3000-3999）
	CodeOrgNotFound:        {msgZh: "组织不存在",     msgEn: "Organization not found", grpc: codes.NotFound},
	CodeOrgCreateFailed:    {msgZh: "新增组织失败",   msgEn: "Create organization failed", grpc: codes.Internal},
	CodeOrgParentNotFound:  {msgZh: "父组织不存在",   msgEn: "Parent organization not found", grpc: codes.NotFound},
	CodeOrgCycle:           {msgZh: "存在循环引用",   msgEn: "Cyclic reference", grpc: codes.FailedPrecondition},
	CodeOrgHasChildren:     {msgZh: "有子组织，无法删除", msgEn: "Has children, cannot delete", grpc: codes.FailedPrecondition},
	CodeOrgHasUsers:        {msgZh: "有关联用户，无法删除", msgEn: "Has users, cannot delete", grpc: codes.FailedPrecondition},
	CodeOrgCrossDirectory:  {msgZh: "跨目录域操作",   msgEn: "Cross directory operation", grpc: codes.FailedPrecondition},
	CodeOrgLevelExceeded:   {msgZh: "超过层级限制",   msgEn: "Level limit exceeded", grpc: codes.FailedPrecondition},
	CodeOrgUpdateFailed:    {msgZh: "编辑组织失败",   msgEn: "Update organization failed", grpc: codes.Internal},
	CodeOrgGetFailed:       {msgZh: "查询组织失败",   msgEn: "Get organization failed", grpc: codes.Internal},
	CodeOrgDeleteFailed:    {msgZh: "删除组织失败",   msgEn: "Delete organization failed", grpc: codes.Internal},
	CodeOrgMoveFailed:      {msgZh: "移动组织失败",   msgEn: "Move organization failed", grpc: codes.Internal},
	// 用户（4000-4999）
	CodeUserNotFound:            {msgZh: "用户不存在",     msgEn: "User not found", grpc: codes.NotFound},
	CodeUserUsernameExists:      {msgZh: "用户名已存在",   msgEn: "Username already exists", grpc: codes.AlreadyExists},
	CodeUserAlreadyDisabled:     {msgZh: "用户已是禁用状态", msgEn: "User already disabled", grpc: codes.FailedPrecondition},
	CodeUserAlreadyEnabled:      {msgZh: "用户已是启用状态", msgEn: "User already enabled", grpc: codes.FailedPrecondition},
	CodeUserOrgAlreadySecondary: {msgZh: "组织已是副组织",   msgEn: "Organization already secondary", grpc: codes.FailedPrecondition},
	CodeUserOrgIsPrimary:        {msgZh: "组织是主组织",     msgEn: "Organization is primary", grpc: codes.FailedPrecondition},
	CodeUserOrgNotSecondary:     {msgZh: "组织不是副组织",   msgEn: "Organization not secondary", grpc: codes.FailedPrecondition},
	CodeUserPrimaryOrgNotFound:  {msgZh: "主组织不存在",   msgEn: "Primary organization not found", grpc: codes.NotFound},
	CodeUserCreateFailed:        {msgZh: "新增用户失败",   msgEn: "Create user failed", grpc: codes.Internal},
	CodeUserUpdateFailed:        {msgZh: "编辑用户失败",   msgEn: "Update user failed", grpc: codes.Internal},
	CodeUserGetFailed:           {msgZh: "查询用户失败",   msgEn: "Get user failed", grpc: codes.Internal},
	CodeUserDeleteFailed:        {msgZh: "删除用户失败",   msgEn: "Delete user failed", grpc: codes.Internal},
	CodeUserEnableFailed:        {msgZh: "启用用户失败",   msgEn: "Enable user failed", grpc: codes.Internal},
	CodeUserDisableFailed:       {msgZh: "禁用用户失败",   msgEn: "Disable user failed", grpc: codes.Internal},
}

// Register 注册自定义业务错误码（业务段位 2000+ 用）。
// 必须在业务包的 init() 阶段调用；禁止在请求处理路径并发调用。
func Register(code, msgZh, msgEn string, grpcCode codes.Code) {
	registry[code] = errorMeta{msgZh: msgZh, msgEn: msgEn, grpc: grpcCode}
}

// 目录域便捷构造函数
func ErrDirNotFound() *BizError          { return New(CodeDirNotFound) }
func ErrDirDomainExists() *BizError      { return New(CodeDirDomainExists) }
func ErrDirAlreadyDisabled() *BizError   { return New(CodeDirAlreadyDisabled) }
func ErrDirAlreadyEnabled() *BizError    { return New(CodeDirAlreadyEnabled) }
func ErrDirNotEmpty() *BizError          { return New(CodeDirNotEmpty) }
func ErrDirDisabled() *BizError          { return New(CodeDirDisabled) }
func ErrDirCreateFailed() *BizError      { return New(CodeDirCreateFailed) }
func ErrDirUpdateFailed() *BizError      { return New(CodeDirUpdateFailed) }
func ErrDirGetFailed() *BizError         { return New(CodeDirGetFailed) }
func ErrDirDeleteFailed() *BizError      { return New(CodeDirDeleteFailed) }

// 组织架构便捷构造函数
func ErrOrgNotFound() *BizError       { return New(CodeOrgNotFound) }
func ErrOrgCreateFailed() *BizError   { return New(CodeOrgCreateFailed) }
func ErrOrgParentNotFound() *BizError { return New(CodeOrgParentNotFound) }
func ErrOrgCycle() *BizError          { return New(CodeOrgCycle) }
func ErrOrgHasChildren() *BizError    { return New(CodeOrgHasChildren) }
func ErrOrgHasUsers() *BizError       { return New(CodeOrgHasUsers) }
func ErrOrgCrossDirectory() *BizError { return New(CodeOrgCrossDirectory) }
func ErrOrgLevelExceeded() *BizError  { return New(CodeOrgLevelExceeded) }
func ErrOrgUpdateFailed() *BizError   { return New(CodeOrgUpdateFailed) }
func ErrOrgGetFailed() *BizError      { return New(CodeOrgGetFailed) }
func ErrOrgDeleteFailed() *BizError   { return New(CodeOrgDeleteFailed) }
func ErrOrgMoveFailed() *BizError     { return New(CodeOrgMoveFailed) }

// 用户便捷构造函数
func ErrUserNotFound() *BizError            { return New(CodeUserNotFound) }
func ErrUserUsernameExists() *BizError      { return New(CodeUserUsernameExists) }
func ErrUserAlreadyDisabled() *BizError     { return New(CodeUserAlreadyDisabled) }
func ErrUserAlreadyEnabled() *BizError      { return New(CodeUserAlreadyEnabled) }
func ErrUserOrgAlreadySecondary() *BizError { return New(CodeUserOrgAlreadySecondary) }
func ErrUserOrgIsPrimary() *BizError        { return New(CodeUserOrgIsPrimary) }
func ErrUserOrgNotSecondary() *BizError     { return New(CodeUserOrgNotSecondary) }
func ErrUserPrimaryOrgNotFound() *BizError  { return New(CodeUserPrimaryOrgNotFound) }
func ErrUserCreateFailed() *BizError        { return New(CodeUserCreateFailed) }
func ErrUserUpdateFailed() *BizError        { return New(CodeUserUpdateFailed) }
func ErrUserGetFailed() *BizError           { return New(CodeUserGetFailed) }
func ErrUserDeleteFailed() *BizError        { return New(CodeUserDeleteFailed) }
func ErrUserEnableFailed() *BizError        { return New(CodeUserEnableFailed) }
func ErrUserDisableFailed() *BizError       { return New(CodeUserDisableFailed) }
