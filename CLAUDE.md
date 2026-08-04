# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 概述

isc-gobase（`github.com/qkja/gobase`）是一套借鉴 Spring Boot 思想的 Go 工具框架，沉淀自 Java 转 Go 的实践经验。它是一个**库**，不是应用：使用方在自己的 `main()` 中引入各包并驱动它，本仓库没有可运行的二进制。

**语言约定：** 代码注释、README、配置示例和测试均使用**中文**，新代码需保持该风格。`application.yaml`/TOML 配置键与 API 约定均对标 Spring Boot。

## 构建与测试

```shell
# 构建框架核心包（jvm 和 system/process 在裸 Windows 机器上无法构建）
go build ./config/... ./server/... ./logger/... ./errors/... ./isc/...

# 全量测试（需 bash，在 Git Bash / WSL 下执行）
sh go_test.sh          # 或：bash go_test.sh

# 单个测试包
go test ./config/test
go test ./logger

# 单个用例
go test ./errors -run TestRegistryContainsAllConsts
```

- CI（`.github/workflows/go.yml`）在 Go 1.18 下执行 `go mod tidy` 并运行同样 13 个测试包；`jvm` 包有 CGo 依赖，CI 会安装 JDK 17 并软链 `jni.h`/`libjvm.so`。
- **不要在裸机器上直接 `go build ./...` / `go test ./...`** —— `jvm` 需要 JNI 头文件，`system/process` 的 Windows 架构（386/amd64）构建标签在本地会失败，CI 是显式安装了 JNI 的。
- 规范测试目录见 `go_test.sh`：`config`、`isc`、`validate`、`compress`、`cron`、`encoding`、`file`、`goid`、`i18n`、`coder`、`time`、`listener`、`bean` 在 `<pkg>/test/` 下。注意 `server`、`logger`、`errors`、`cache`、`tenant`、`discovery`、`db` 改用**包内** `_test.go`。
- 多数测试从工作目录加载 `application.yaml`（每个 `<pkg>/test/` 自带配置），因此要以测试目录作为 CWD 运行。

## 核心架构

各包构成分层栈，几乎所有包都依赖 `isc`。

- **`isc/`** —— 基础工具包（34 个文件）。类型转换（`isc.ToString/ToInt/ToBool/...`）、反射数据绑定（`isc.DataToObject`）、JSON/YAML/TOML/properties 互转（`YamlToProperties`、`PropertiesToMap`、`JsonToYaml`、`TomlToMap`），以及各包广泛使用的 `ISCString`/`ISCList` 包装类型。**优先用它们而非标准库**，比如 `isc.ObjectToJson`、`isc.ToString` 随处可见。
- **`config/`** —— 仿 Spring 的配置加载 + 热加载。这是整套框架的基石。
- **`server/`** —— 基于 gin 的 HTTP 服务封装。
- **`goid/` + `store/`** —— Go 版 ThreadLocal 与单请求上下文。
- **`logger/`** —— 基于 logrus 的分组日志。
- **`errors/`** —— 统一的 gRPC 业务错误码。
- **`listener/`** —— 事件发布订阅（服务生命周期 + 配置变更事件）。
- **`bean/`** —— 运行时对象注册表，用于线上调试（通过 HTTP 端点做属性查看/修改、方法调用）。

### 配置系统（`config/`，最值得理解的部分）

- 从 CWD 加载 `application.{yaml,yml,properties,json,toml}`；支持 Spring profile，通过 `profiles.active`（环境变量优先于文件）加载 `application-{profile}.ext`。还会追加 `./config/application-default.yml`（可用环境变量 `config.additional-location` 覆盖）。
- **环境变量契约：** 外部环境变量随 `base.` 前缀移除同步改名：`base.profiles.active` → `profiles.active`、`base.config.additional-location` → `config.additional-location`。
- **双存储**：`ApplicationProperty` 内部有 `ValueMap`（扁平点号键，如 `"server.port"`）和 `ValueDeepMap`（嵌套 map）。所有 loader 都归一化写入两者；`SetValue` 通过 properties↔yaml 往返转换保持二者同步。
- **两种取值方式：** `config.GetValueString/Int/Bool/...Default(key)`（实时、点号键）与 `config.GetValueObject("", &BaseCfg)`（结构体绑定、快照、非实时）。README 明确推荐实时 getter。`config/base_properties.go` 中的 `BaseCfg`/`RedisCfg`/`EtcdCfg`/`EmqxCfg` 在加载时一次性填充。
- **热加载：** `config/watch.go` 每秒轮询 `./config/application.toml`（size+ModTime 快速跳过、sha256 确认——因为 k8s ConfigMap 的 symlink 原子替换会让 inotify 失效）。整文件重载时只发布一个**合成** `appconfig.reload` 事件（Key=`"appconfig.reload"`），不产生逐 key 事件；HTTP `config/update` 接口的逐 key 变更则按 key 发布普通 `ConfigChangeEvent`。需要支持热加载的新特性必须处理 `appconfig.reload`（参考 `logger.ConfigChangeListener`）。
- 配置变更经 `listener.PublishEvent(listener.ConfigChangeEvent{Key, Value})` 广播；用 `listener.AddListener(listener.EventOfConfigChange, fn)` 订阅。

### Web 服务（`server/`）

- `server.Run()` 启动 gin；`server.Get/Post/Put/Delete/All/GetPost(path, handler)` 注册路由。`server.init()` 会调用 `config.LoadConfig()`。
- 路由自动加前缀：`/{api.prefix}/{api-module}/{path}`（见 `getPathAppendApiModel`）。`api-module` 是顶层配置键；路径本身以 `api` 开头时按原样使用。
- 内置中间件链：CORS、`gin.Recovery()`+`ErrHandler()`、`RequestSaveHandler()`（把请求头写入 `store`，请求结束时 `store.Clean()`）、`rsp.ResponseHandler()`。可选端点由配置开关控制：健康检查（`endpoint.health.enable`）、配置查看/修改（`endpoint.config.enable`）、bean 调试（`endpoint.bean.enable`）、pprof（`server.gin.pprof.enable`）、swagger（`swagger.enable`），均挂在 `/{prefix}/{module}/...` 下。
- **基于请求头的 API 版本化：** `server.GetWith(path, header, versionName, handler)` 可为同一路径注册多个 handler，按请求头取值匹配分发（`api_version.go` 中的 `ApiPath`/`ApiVersion`）。
- 响应使用 `server/rsp`：`SuccessOfStandard`/`FailedOfStandard` 产出 `{"code":0,"data":...,"message":"success"}`，其中 `code` 为 **int**，并带 `isc-biz-code`/`isc-biz-message` 响应头。另有泛型 `DataResponse[T]`、`PagedResponse[T]`，以及分页用 `req.PageRequest[T]`。

### 线程局部与请求上下文（`goid/` + `store/`）

- `goid/` 实现基于 goroutine ID 的 Java 式 ThreadLocal。**关键陷阱：** 要把局部存储数据带过协程，必须用 `goid.Go(func(){...})` 创建协程，不能用 `go func(){...}` —— 原生 `go` 会丢数据。
- `store/` 用 `goid.LocalStorage` + concurrent-map 封装请求级 KV。`RequestSaveHandler` 从 trace/tenant 请求头（`constants.TRACE_HEAD_*`、`TENANT_HEAD_*`）填充它，并在每次请求后 `store.Clean()`。logger 从这里读取 `traceId`/`userId`。请求级数据请用 `store.Put/Get`，不要另建 `goid` 存储。

### 日志（`logger/`）

- 基于 logrus 的**分组**日志：`logger.Group("name")` 返回带独立级别的分组 logger（`logger.group.<name>.level`，缺省回退 `logger.level`）。包级 `logger.Info/Warn/Error/...` 走 `root` 分组。每个分组通过 file-rotatelogs + lfshook 按级别写文件到 `logger.home`（默认 `./logs/`）。
- 日志行通过自定义 `StandardFormatter` 内嵌 `store` 中的 `traceId`/`userId`，链路上下文自动带进日志。
- `logger.DebugWithTenant(tenantID, ...)` 仅在全局级别为 debug/trace、或该租户在 `logger.debug_tenant_ids` 中时才输出——用于生产环境的按租户调试。
- 级别热更新在 `ConfigChangeListener` 中处理；收到 `appconfig.reload` 时重读配置并应用到已创建的分组（不会新建分组）。

### 业务错误（`errors/` 与 `server/rsp` 是两层，别混淆）

- **`errors/`** 面向 gRPC 业务错误。`BizError`（实现 `error` + `GRPCStatus()`）可直接在 gRPC handler 里 `return nil, errors.ErrXxx()`；业务码通过 `errdetails.ErrorInfo` 跨进程传递。`errors.FromError(err)` 能从任意 error 中提取业务码（进程内 `BizError`、`%w` 包装链、或跨进程 `ErrorInfo`）。
- **码段划分：** `0` 成功，`1001–1999` 通用（本包定义），`2000+` 业务段由业务模块在 `init` 阶段 `errors.Register(code, msgZh, msgEn, grpcCode)` 注册。注册表是 `map[string]errorMeta`（而非 switch），支撑码→（中英消息、grpc 码）反查和 `Message(code, lang)` 双语回退。
- **不要用 `fmt.Errorf("%w")` 包装要返回给 gRPC 的 `BizError`** —— grpc v1.41.0 的 `status.FromError` 只做顶层类型断言，`%w` 包装会识别不到业务码；改用 `errors.ErrXxx().WithCause(err)`（见 `errors/grpc.go` 注释）。
- **`server/rsp`** 是另一套 HTTP 响应层，用 **int** 码和 `{"code":0,...}` 信封，与上面无关。

### 基础设施包（db / discovery / tenant）

这三个是较新的包（当前 `main` 上未提交的工作，与 errors/logger/config 改动一同演进），面向 wire 装配的 gRPC 服务：

- **`db/`** —— 为 `google/wire` 提供 MongoDB provider：`db.MongoSet = wire.NewSet(NewClient, NewDatabase)`，从 `database.mongodb`（uri/database/connect_timeout）读取配置。`EnsureIndexesFromFile` 按 JSON 清单同步集合索引（补齐缺失、删掉清单外的索引、保留 `_id_`），支持唯一索引与配合假删除的部分索引过滤条件。
- **`discovery/`** —— 基于 `discovery.json` 的服务发现。两层：服务自身 `<cfgDir>/discovery.json`（高优先级）、内置 `discovery.default.json`（兜底，go:embed）。启动时调 `discovery.Init(cfgDir)`（把 `{namespace}` 占位符替换为 `POD_NAMESPACE` 或 `WithNamespace` 指定值），之后用 `discovery.GetAddress(svcName)`。未配置返回 `("", false)`。
- **`tenant/`** —— 仅通过 `context.Context` 传递租户上下文。`tenant.WithInfo(ctx, &tenant.Info{...})` 写入，`tenant.GetInfo(ctx)` 读取。设计约定：本包**不做**任何 proto/反射注入；由外部网关中间件写入 info，logic 层读取后手动把 `TenantID` 复制进 proto 请求。请求头名在 `constants.TENANT_HEAD_*`。

## 约定与陷阱

- **`isc` 优先：** 任何字符串/类型/map/JSON 转换先查 `isc` 再看标准库。跨包代码依赖 `isc.ToString`/`ObjectToJson`/`DataToObject`/`YamlToProperties`。
- **框架配置键都是顶层键**（`server.enable`、`logger.level`、`endpoint.*`）；`api-module` 是顶层键，控制 URL 模块段。
- **中文注释与消息是常态**；注册表中的默认错误消息保持中英双语（zh + en）。
- **请求作用域结束必须 `store.Clean()`**；server 中间件会做，但其它填充 `store` 的地方也要记得。
- 配置结构体绑定（`GetValueObject`）是快照——需要响应热加载的场景不要用它，改用实时 getter 或监听配置变更事件。
