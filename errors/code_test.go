package errors

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/qkja/gobase/i18n"
)

func TestGlobalVars(t *testing.T) {
	Convey("全局错误变量", t, func() {
		Convey("通用错误应预填充正确的 code 和中文 msg", func() {
			So(ErrInternal.GetCode(), ShouldEqual, CodeInternal)
			So(ErrInternal.GetMessage(), ShouldEqual, "系统内部错误")
			So(ErrNotFound.GetCode(), ShouldEqual, CodeNotFound)
			So(ErrNotFound.GetMessage(), ShouldEqual, "资源不存在")
			So(ErrUnknown.GetCode(), ShouldEqual, CodeUnknown)
			So(ErrUnknown.GetMessage(), ShouldEqual, "未知错误")
		})

		Convey("目录域错误应预填充正确的 code 和中文 msg", func() {
			So(ErrDirNotFound.GetCode(), ShouldEqual, CodeDirNotFound)
			So(ErrDirNotFound.GetMessage(), ShouldEqual, "目录域不存在")
			So(ErrDirCreateFailed.GetCode(), ShouldEqual, CodeDirCreateFailed)
			So(ErrDirCreateFailed.GetMessage(), ShouldEqual, "新增目录域失败")
		})

		Convey("组织错误应预填充正确的 code 和中文 msg", func() {
			So(ErrOrgNotFound.GetCode(), ShouldEqual, CodeOrgNotFound)
			So(ErrOrgNotFound.GetMessage(), ShouldEqual, "组织不存在")
			So(ErrOrgCycle.GetCode(), ShouldEqual, CodeOrgCycle)
			So(ErrOrgCycle.GetMessage(), ShouldEqual, "存在循环引用")
		})

		Convey("用户错误应预填充正确的 code 和中文 msg", func() {
			So(ErrUserNotFound.GetCode(), ShouldEqual, CodeUserNotFound)
			So(ErrUserNotFound.GetMessage(), ShouldEqual, "用户不存在")
			So(ErrUserUsernameExists.GetCode(), ShouldEqual, CodeUserUsernameExists)
			So(ErrUserUsernameExists.GetMessage(), ShouldEqual, "用户名已存在")
		})
	})
}

func TestNew(t *testing.T) {
	Convey("New(code)", t, func() {
		Convey("应调用 i18n.Lookup 获取中文 msg", func() {
			patches := ApplyFunc(i18n.Lookup, func(lang, code string) (string, bool) {
				if code == "1003" {
					return "资源不存在", true
				}
				return "", false
			})
			defer patches.Reset()

			be := New(CodeNotFound)
			So(be.GetCode(), ShouldEqual, CodeNotFound)
			So(be.GetMessage(), ShouldEqual, "资源不存在")
		})

		Convey("未注册的 code 应兜底为 CodeUnknown 文案", func() {
			patches := ApplyFunc(i18n.Lookup, func(lang, code string) (string, bool) {
				if code == CodeUnknown {
					return "未知错误", true
				}
				return "", false
			})
			defer patches.Reset()

			be := New("9999")
			So(be.GetMessage(), ShouldEqual, "未知错误")
		})
	})
}

func TestMessage(t *testing.T) {
	Convey("Message(code, lang)", t, func() {
		Convey("中文查询应返回中文翻译", func() {
			patches := ApplyFunc(i18n.Lookup, func(lang, code string) (string, bool) {
				if lang == i18n.LangZh && code == CodeInternal {
					return "系统内部错误", true
				}
				return "", false
			})
			defer patches.Reset()

			So(Message(CodeInternal, i18n.LangZh), ShouldEqual, "系统内部错误")
		})

		Convey("英文查询应返回英文翻译", func() {
			patches := ApplyFunc(i18n.Lookup, func(lang, code string) (string, bool) {
				if lang == i18n.LangEn && code == CodeInternal {
					return "Internal error", true
				}
				return "", false
			})
			defer patches.Reset()

			So(Message(CodeInternal, i18n.LangEn), ShouldEqual, "Internal error")
		})

		Convey("未知语言应回退到 CodeUnknown 兜底文案", func() {
			patches := ApplyFunc(i18n.Lookup, func(lang, code string) (string, bool) {
				if code == CodeUnknown {
					return "Unknown error", true
				}
				return "", false
			})
			defer patches.Reset()

			So(Message(CodeNotFound, "fr-FR"), ShouldEqual, "Unknown error")
		})

		Convey("所有 lookup 都失败应返回原始 code", func() {
			patches := ApplyFunc(i18n.Lookup, func(lang, code string) (string, bool) {
				return "", false
			})
			defer patches.Reset()

			So(Message("9999", i18n.LangZh), ShouldEqual, "9999")
		})

		Convey("patch 恢复后真实调用应正常工作", func() {
			So(Message(CodeNotFound, i18n.LangZh), ShouldEqual, "资源不存在")
			So(Message(CodeNotFound, i18n.LangEn), ShouldEqual, "Resource not found")
		})
	})
}
