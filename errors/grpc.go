package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCStatus 让 BizError 实现 grpc 的 GRPCStatuser 接口。
// 因此 api 层可直接 return nil, errors.ErrNotFound()，gRPC 服务端会自动把它序列化为带消息的 status。
// 跨进程业务码还原（ErrorInfo / FromError）已移除：业务只在本服务内使用 gobase 错误码。
func (e *BizError) GRPCStatus() *status.Status {
	return status.New(e.GRPCCode, e.message())
}

// ToGRPCStatus 转为 *status.Status
func (e *BizError) ToGRPCStatus() *status.Status { return e.GRPCStatus() }

// ToGRPCError 转为可直接返回的 gRPC error
func (e *BizError) ToGRPCError() error { return e.GRPCStatus().Err() }

// 保证 codes 包被引用（GRPCCode 类型），避免误删。
var _ codes.Code
