# errors —— 统一业务错误码

gobase 错误码提供极简的业务错误载体：**两个字段（code + msg）+ 一个翻译查询接口 `Message(code, lang)`**，面向 gRPC / HTTP 统一使用。

## 核心 API

```go
// 结构体——字段不导出，通过 getter 访问，全局变量不可变
// 实现 error 接口，errors.Is 按业务码比较

// 创建错误——msg 自动从 i18n 填充中文
err := errors.New(errors.CodeNotFound)
err := errors.ErrNotFound // 全局变量，直接使用

// 查翻译——不创建错误对象，纯查询
msg := errors.Message("1003", i18n.LangZh) // "资源不存在"
msg := errors.Message("1003", i18n.LangEn) // "Resource not found"
```

## 错误码定义

### 码段划分

所有错误码为 **string 类型**，全部由 gobase 内置（业务**禁止**自定义/注册新码）：

| 码段 | 含义 | 数量 |
|------|------|------|
| `"0"` | 成功 | 1 |
| `"1001"`–`"1999"` | 通用错误 | 10 |
| `"2000"`–`"2999"` | 目录域 | 10 |
| `"3000"`–`"3999"` | 组织架构 | 12 |
| `"4000"`–`"4999"` | 用户 | 14 |
| `"1999"` | 未知错误（兜底） | 1 |

### 通用错误码

| 常量 | 码值 | 含义 |
|------|------|------|
| `CodeOK` | `"0"` | 成功 |
| `CodeInternal` | `"1001"` | 系统内部错误 |
| `CodeInvalidArgument` | `"1002"` | 参数无效 |
| `CodeNotFound` | `"1003"` | 资源不存在 |
| `CodeUnauthenticated` | `"1004"` | 未认证或登录已过期 |
| `CodePermissionDenied` | `"1005"` | 无权限 |
| `CodeTimeout` | `"1006"` | 请求超时 |
| `CodeAlreadyExists` | `"1007"` | 资源已存在 |
| `CodeConflict` | `"1008"` | 依赖冲突 |
| `CodeResourceExhausted` | `"1009"` | 请求过于频繁 |
| `CodeUnknown` | `"1999"` | 未知错误（兜底） |

### 业务域错误码

<details>
<summary>目录域（2000–2999）</summary>

| 常量 | 码值 |
|------|------|
| `CodeDirNotFound` | `"2001"` |
| `CodeDirDomainExists` | `"2002"` |
| `CodeDirAlreadyDisabled` | `"2003"` |
| `CodeDirAlreadyEnabled` | `"2004"` |
| `CodeDirNotEmpty` | `"2005"` |
| `CodeDirDisabled` | `"2006"` |
| `CodeDirCreateFailed` | `"2007"` |
| `CodeDirUpdateFailed` | `"2008"` |
| `CodeDirGetFailed` | `"2009"` |
| `CodeDirDeleteFailed` | `"2010"` |

</details>

<details>
<summary>组织架构（3000–3999）</summary>

| 常量 | 码值 |
|------|------|
| `CodeOrgNotFound` | `"3001"` |
| `CodeOrgCreateFailed` | `"3002"` |
| `CodeOrgParentNotFound` | `"3003"` |
| `CodeOrgCycle` | `"3004"` |
| `CodeOrgHasChildren` | `"3005"` |
| `CodeOrgHasUsers` | `"3006"` |
| `CodeOrgCrossDirectory` | `"3007"` |
| `CodeOrgLevelExceeded` | `"3008"` |
| `CodeOrgUpdateFailed` | `"3009"` |
| `CodeOrgGetFailed` | `"3010"` |
| `CodeOrgDeleteFailed` | `"3011"` |
| `CodeOrgMoveFailed` | `"3012"` |

</details>

<details>
<summary>用户（4000–4999）</summary>

| 常量 | 码值 |
|------|------|
| `CodeUserNotFound` | `"4001"` |
| `CodeUserUsernameExists` | `"4002"` |
| `CodeUserAlreadyDisabled` | `"4003"` |
| `CodeUserAlreadyEnabled` | `"4004"` |
| `CodeUserOrgAlreadySecondary` | `"4005"` |
| `CodeUserOrgIsPrimary` | `"4006"` |
| `CodeUserOrgNotSecondary` | `"4007"` |
| `CodeUserPrimaryOrgNotFound` | `"4008"` |
| `CodeUserCreateFailed` | `"4009"` |
| `CodeUserUpdateFailed` | `"4010"` |
| `CodeUserGetFailed` | `"4011"` |
| `CodeUserDeleteFailed` | `"4012"` |
| `CodeUserEnableFailed` | `"4013"` |
| `CodeUserDisableFailed` | `"4014"` |

</details>

每个码都有对应的全局错误变量（`ErrInternal`、`ErrNotFound`、`ErrUserNotFound` 等），直接使用，不可变。

## BizError

### 结构体

```go
// 字段不导出，通过 GetCode()/GetMessage() 访问
// 全局变量 ErrXxx 不可变，多 goroutine 安全
```

### 方法

| 方法 | 说明 |
|------|------|
| `Error() string` | 实现 `error` 接口，输出 `"code: msg"`，如 `"1001: 系统内部错误"` |
| `GetCode() string` | 返回业务错误码 |
| `GetMessage() string` | 返回错误消息 |
| `Is(target error) bool` | 支持 `errors.Is` 按业务码比较 |

### 使用方式

```go
// 创建——msg 由 i18n 自动填充中文
err := errors.New(errors.CodeNotFound)
err := errors.ErrNotFound         // 全局变量，直接使用
err := errors.ErrUserNotFound     // 每个码都有对应变量

// 错误判断
if errors.Is(err, errors.ErrNotFound) {
    // 资源不存在
}

// 读取 code/msg（通过 getter）
code := err.GetCode()     // "1003"
msg := err.GetMessage()   // "资源不存在"
```

## Message —— 按 code + 语言查翻译

对外提供的翻译查询接口，不创建错误对象：

```go
// 查翻译
msg := errors.Message("1003", i18n.LangZh) // "资源不存在"
msg := errors.Message("1003", i18n.LangEn) // "Resource not found"

// 查找链：指定语言 → 默认语言（zh-CN）→ CodeUnknown 兜底 → 原始 code
msg := errors.Message("99999", i18n.LangZh) // "未知错误"
msg := errors.Message("99999", i18n.LangEn) // "Unknown error"
```

语言常量（定义在 `i18n` 包，全仓库直接使用）：

| 常量 | 值 |
|------|-----|
| `i18n.LangZh` | `"zh-CN"` |
| `i18n.LangEn` | `"en-US"` |

## 翻译（i18n）体系

### 默认翻译

gobase 通过 `go:embed` 内嵌全部 48 个错误码的中英文翻译（`i18n/default/zh-CN.po` + `en-US.po`），开箱即用，无需初始化。

### .po 文件格式

自定义的 key-value 行格式，`key` 为首个空格前内容，`value` 为之后内容：

```
0 成功
1001 系统内部错误
1003 资源不存在
1999 未知错误
4001 用户不存在
```

### 服务覆盖翻译

服务无需修改 gobase，只需在**工作目录**下创建 `i18n/` 目录，按 **code 键** 覆盖即可：

```
your-service/
├── i18n/
│   ├── zh-CN.po
│   └── en-US.po
└── main.go
```

```
# i18n/zh-CN.po —— 只写需要覆盖的码
1001 服务器开小差了，请稍后重试
4001 该用户不存在哦
```

### 查找链

```
Message(code, lang):
  1. i18n.Lookup(lang, code) → 嵌入默认 + 服务文件叠加
  2. i18n.Lookup(lang, CodeUnknown) → 兜底文案
  3. 返回原始 code 字符串
```

## 使用示例

```go
package service

import (
    "errors"
    gobaseErrors "github.com/qkja/gobase/errors"
)

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserResp, error) {
    if req.Username == "" {
        return nil, gobaseErrors.ErrInvalidArgument
    }

    exists, err := s.repo.ExistsByUsername(ctx, req.Username)
    if err != nil {
        return nil, gobaseErrors.ErrInternal
    }
    if exists {
        return nil, gobaseErrors.ErrUserUsernameExists
    }

    return &pb.CreateUserResp{}, nil
}

// 按请求语言返回翻译
func msgByLang(code, acceptLang string) string {
    lang := i18n.LangZh
    if acceptLang == "en" {
        lang = i18n.LangEn
    }
    return gobaseErrors.Message(code, lang)
}
```

## 设计约束

| 规则 | 说明 |
|------|------|
| **只用内置码** | 业务禁止自定义/注册错误码 |
| **消息由 i18n 唯一提供** | 文案来自嵌入 `.po` + 服务覆盖，不硬编码 |
| **纯查询** | `Message()` 不创建对象、无副作用 |

## 目录结构

```
errors/
├── code.go          # 错误码常量 + 便捷构造函数
├── code_test.go     # New / Message / ErrXxx 测试
├── errors.go        # BizError 结构体 + Message 查询
├── errors_test.go   # Error() / Is / GetCode / GetMessage 测试
└── README.md        # 本文档

i18n/
├── i18n.go          # InitI18N / T() / Tf() 全局翻译 API
├── datamap.go       # Lookup() + embed 默认 .po + 服务文件叠加
└── default/
    ├── zh-CN.po     # go:embed 内嵌中文默认翻译（48 条）
    └── en-US.po     # go:embed 内嵌英文默认翻译（48 条）
```
