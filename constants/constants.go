package constants

const (
	BeanNameGormPre  = "redis_"
	BeanNameXormPre  = "xorm_"
	BeanNameRedisPre = "redis_"
	BeanNameEtcdPre  = "etcd_"
)

const (
	TRACE_HEAD_ID             = "t-head-traceId"
	TRACE_HEAD_RPC_ID         = "t-head-rpcId"
	TRACE_HEAD_SAMPLED        = "t-head-sampled"
	TRACE_HEAD_USER_ID        = "t-head-userId"
	TRACE_HEAD_USER_NAME      = "t-head-userName"
	TRACE_HEAD_REMOTE_IP      = "t-head-remoteIp"
	TRACE_HEAD_REMOTE_APPNAME = "t-head-remoteAppName"
	TRACE_HEAD_ORIGNAL_URL    = "t-head-orignal-url"
)

const (
	SAVE_HEAD        = "head"
	SAVE_REMOTE_ADDR = "remote_addr"
)

// 租户上下文请求头：网关入口（或上游调用方）通过这三个头下发租户信息，
// 各服务经 middleware 读取后写入 ctx（tenant.WithInfo），再注入 proto Req。
const (
	TENANT_HEAD_ID          = "t-head-tenantId"
	TENANT_HEAD_LANGUAGE    = "t-head-tenantLanguage"
	TENANT_HEAD_UI_LANGUAGE = "t-head-tenantUILanguage"
)
