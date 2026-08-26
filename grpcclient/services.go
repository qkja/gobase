package grpcclient

// 后端服务名常量（discovery.json 的 key）。
// 集中维护：新增后端服务时在此加一个常量，供 grpcclient.Call(ctx, svcName, ...) 按服务名解析地址。
const (
	IdentityhubSvr   = "identityhubsvr"
	TenantManagerSvr = "tenantmanagersvr"
)
