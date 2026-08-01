package errors

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Domain 错误域标识，用于 ErrorInfo 细节
const Domain = "gobase"

// GRPCStatus 让 BizError 实现 grpc 的 GRPCStatuser 接口。
// 因此 api 层可直接 return nil, errors.ErrNotFound()，无需手动调用 ToGRPCError，
// gRPC 服务端会自动把它序列化为带业务码细节(ErrorInfo)的 status，业务码可跨进程存活。
func (e *BizError) GRPCStatus() *status.Status {
	c, msg := e.GRPCCode, e.message()
	st := status.New(c, msg)
	if c == codes.OK {
		return st
	}
	if st, err := st.WithDetails(e.errorInfo()); err == nil {
		return st
	}
	return status.New(c, msg) // WithDetails 失败兜底（几乎不会发生）
}

// ToGRPCStatus 转为 *status.Status
func (e *BizError) ToGRPCStatus() *status.Status { return e.GRPCStatus() }

// ToGRPCError 转为可直接返回的 gRPC error
func (e *BizError) ToGRPCError() error { return e.GRPCStatus().Err() }

// FromError 从任意 error 提取业务错误，用于 grpc-gateway 错误处理器 / 统一错误拦截器。
//
// 返回：
//   - (be, true)  识别到明确业务码（本进程内 BizError 或跨进程 ErrorInfo 细节）；
//   - (be, false) 未能识别，兜底为 CodeInternal、保留原始消息，便于网关统一回 5xx；
//   - (nil, false) err 为 nil。
func FromError(err error) (*BizError, bool) {
	if err == nil {
		return nil, false
	}
	// 1) 本进程内：BizError 及其 %w 包装链
	var be *BizError
	if errors.As(err, &be) {
		return be, true
	}
	// 2) 跨进程：gRPC status 携带的 ErrorInfo 细节
	st, ok := status.FromError(err)
	if !ok || st.Code() == codes.OK {
		return nil, false
	}
	if be := fromStatus(st); be != nil {
		return be, true
	}
	// 3) 兜底
	return &BizError{Code: CodeInternal, Message: st.Message(), GRPCCode: st.Code()}, false
}

// Is 判断错误是否对应指定业务错误码（基于 FromError 提取，跨进程有效）
func Is(err error, code string) bool {
	be, ok := FromError(err)
	return ok && be.Code == code
}

// fromStatus 从 gRPC status 的 ErrorInfo 细节重建 BizError
func fromStatus(st *status.Status) *BizError {
	for _, d := range st.Details() {
		ei, ok := d.(*errdetails.ErrorInfo)
		if !ok || ei.Reason == "" {
			continue
		}
		be := New(ei.Reason)
		be.GRPCCode = st.Code()
		be.Message = st.Message()
		if m := ei.Metadata; len(m) > 0 {
			args := make(map[string]any, len(m))
			for k, v := range m {
				args[k] = v
			}
			be.Args = args
		}
		return be
	}
	return nil
}

// errorInfo 构造携带业务码的 ErrorInfo 细节
func (e *BizError) errorInfo() *errdetails.ErrorInfo {
	return &errdetails.ErrorInfo{
		Reason:   e.Code,
		Domain:   Domain,
		Metadata: toStringMap(e.Args),
	}
}

// toStringMap 将结构化参数转成 ErrorInfo.Metadata（字符串值）
func toStringMap(args map[string]any) map[string]string {
	if len(args) == 0 {
		return nil
	}
	m := make(map[string]string, len(args))
	for k, v := range args {
		switch t := v.(type) {
		case string:
			m[k] = t
		case error:
			m[k] = t.Error()
		case nil:
			m[k] = ""
		default:
			if b, err := json.Marshal(v); err == nil {
				m[k] = string(b)
			} else {
				m[k] = fmt.Sprint(v)
			}
		}
	}
	return m
}
