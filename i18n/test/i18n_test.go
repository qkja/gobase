package test

import (
	"testing"

	. "github.com/qkja/gobase/i18n"
)

func TestI18N(t *testing.T) {
	_ = InitI18N("zh-CN")
	_ = LoadI18NLanguage("en-US")

	t.Logf(T("msgid"))
	t.Logf(T("msgstr"))
	t.Logf(T("msgfmt"))
	t.Logf(T("msgcn"))
	t.Logf(Tf("msgf2", "rarnu", 2333))

}

// TestLookup 验证 Lookup：嵌入默认 + 服务文件叠加 + 默认语言回退。
// CWD 为 i18n/test，服务文件 i18n/zh-CN.po / en-US.po 提供 msgid 等键。
func TestLookup(t *testing.T) {
	// 嵌入默认：框架错误码
	if got, ok := Lookup("zh-CN", "1001"); !ok || got != "系统内部错误" {
		t.Errorf(`Lookup("zh-CN","1001") = %q,%v want "系统内部错误",true`, got, ok)
	}
	if got, ok := Lookup("en-US", "1001"); !ok || got != "Internal error" {
		t.Errorf(`Lookup("en-US","1001") = %q,%v want "Internal error",true`, got, ok)
	}
	// 服务文件叠加：zh-CN.po 提供 msgid -> "hello"
	if got, ok := Lookup("zh-CN", "msgid"); !ok || got != "hello" {
		t.Errorf(`Lookup("zh-CN","msgid") = %q,%v want "hello",true`, got, ok)
	}
	// en-US 服务文件缺 msgcn -> 回退默认语言(zh-CN)
	if got, ok := Lookup("en-US", "msgcn"); !ok || got != "这是默认的" {
		t.Errorf(`Lookup("en-US","msgcn") = %q,%v want "这是默认的",true`, got, ok)
	}
	// 缺失键
	if _, ok := Lookup("zh-CN", "no-such-key"); ok {
		t.Error(`Lookup("zh-CN","no-such-key") 应返回 not found`)
	}
}
