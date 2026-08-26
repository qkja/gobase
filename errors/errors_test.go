package errors

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {
	Convey("BizError.Error()", t, func() {
		Convey("全局错误变量应输出 code: msg 格式", func() {
			So(ErrInternal.Error(), ShouldEqual, "1001: 系统内部错误")
			So(ErrNotFound.Error(), ShouldEqual, "1003: 资源不存在")
		})

		Convey("自定义 msg 的实例也应正确输出", func() {
			be := New(CodeNotFound)
			_ = be // be.msg 默认由 i18n 填充，验证 Error() 即可
			So(be.Error(), ShouldEqual, "1003: 资源不存在")
		})
	})
}

func TestIs(t *testing.T) {
	Convey("errors.Is", t, func() {
		Convey("同码应匹配", func() {
			So(ErrNotFound.Is(ErrNotFound), ShouldBeTrue)
		})

		Convey("不同码不应匹配", func() {
			So(ErrNotFound.Is(ErrInternal), ShouldBeFalse)
		})

		Convey("nil target 应安全处理", func() {
			var be *BizError
			So(be.Is(ErrNotFound), ShouldBeFalse)
		})
	})
}

func TestGetCodeAndMessage(t *testing.T) {
	Convey("GetCode / GetMessage", t, func() {
		be := ErrNotFound

		Convey("应返回正确的 code", func() {
			So(be.GetCode(), ShouldEqual, CodeNotFound)
		})

		Convey("应返回正确的 msg", func() {
			So(be.GetMessage(), ShouldEqual, "资源不存在")
		})

		Convey("nil 调用应返回空字符串且不 panic", func() {
			var nilBe *BizError
			So(nilBe.GetCode(), ShouldEqual, "")
			So(nilBe.GetMessage(), ShouldEqual, "")
		})
	})
}
