package errors

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
)

// TestErrorFormat 验证 Error() 输出格式
func TestErrorFormat(t *testing.T) {
	be := ErrInternal().WithMessage("磁盘写入失败")
	if got := be.Error(); got != "1001: 磁盘写入失败" {
		t.Errorf("Error() = %q, want %q", got, "1001: 磁盘写入失败")
	}
	if got := ErrNotFound().Error(); got != "1003: 资源不存在" {
		t.Errorf("Error() = %q, want %q", got, "1003: 资源不存在")
	}
}

// TestNewf 验证格式化消息
func TestNewf(t *testing.T) {
	be := Newf(CodeNotFound, "用户 %s 不存在", "qkj")
	if be.Message != "用户 qkj 不存在" {
		t.Errorf("Message = %q, want %q", be.Message, "用户 qkj 不存在")
	}
}

// TestIsAndUnwrap 验证 errors.Is / errors.As 包装链
func TestIsAndUnwrap(t *testing.T) {
	inner := ErrNotFound()
	wrapped := fmt.Errorf("wrap: %w", inner)

	if !errors.Is(wrapped, ErrNotFound()) {
		t.Error("errors.Is 应识别同码 BizError")
	}
	if errors.Is(wrapped, ErrInternal()) {
		t.Error("errors.Is 不应匹配不同码")
	}

	var be *BizError
	if !errors.As(wrapped, &be) {
		t.Fatal("errors.As 应提取到 *BizError")
	}
	if be.Code != CodeNotFound {
		t.Errorf("提取到的 Code = %q, want %q", be.Code, CodeNotFound)
	}
}

// TestWithCause 验证 WithCause 承载底层错误
func TestWithCause(t *testing.T) {
	root := fmt.Errorf("db connection refused")
	be := ErrInternal().WithCause(root)

	if !errors.Is(be, root) {
		t.Error("errors.Is 应穿透 WithCause 的 cause")
	}
	// Unwrap 直接返回 cause
	if unwrapped := errors.Unwrap(be); unwrapped != root {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, root)
	}
	// 但 Error() 输出仍是业务码格式
	if got := be.Error(); got != "1001: 系统内部错误" {
		t.Errorf("Error() = %q, want %q", got, "1001: 系统内部错误")
	}
}

// TestBuilder 验证不可变式 builder
func TestBuilder(t *testing.T) {
	base := ErrInvalidArgument()
	withMsg := base.WithMessage("用户名不能为空")
	if withMsg.Message != "用户名不能为空" {
		t.Errorf("WithMessage 未生效: %q", withMsg.Message)
	}
	// 原对象不受影响
	if base.Message != "参数无效" {
		t.Errorf("builder 应不可变, base.Message = %q", base.Message)
	}
	_ = codes.InvalidArgument // 引用 codes 避免未使用
}
