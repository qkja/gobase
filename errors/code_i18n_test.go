package errors

import (
	"os"
	"regexp"
	"testing"

	"github.com/qkja/gobase/i18n"
)

// TestAllCodesHaveI18n 遍历 code.go 中全部错误码，断言 zh-CN / en-US 两个 .po 均有非兜底文案。
//
// 背景（docs P14/T8）：gobase 无任何机制保证错误码有对应 i18n 条目，
// 漏加 .po 会静默降级为「未知错误」且现有测试不失败。此测试把「人工维护清单必然漏项」提前到 CI。
func TestAllCodesHaveI18n(t *testing.T) {
	src, err := os.ReadFile("code.go")
	if err != nil {
		t.Fatalf("read code.go: %v", err)
	}

	re := regexp.MustCompile(`Code\w+ = "(\d+)"`)
	seen := map[string]bool{}
	var codes []string
	for _, m := range re.FindAllSubmatch(src, -1) {
		code := string(m[1])
		if !seen[code] {
			seen[code] = true
			codes = append(codes, code)
		}
	}
	if len(codes) == 0 {
		t.Fatal("no codes extracted from code.go")
	}

	fallbackZh, _ := i18n.Lookup(i18n.LangZh, CodeUnknown)
	fallbackEn, _ := i18n.Lookup(i18n.LangEn, CodeUnknown)

	for _, code := range codes {
		// CodeUnknown(1999) 本身即兜底码，其文案就是兜底文案，跳过
		if code == CodeUnknown {
			continue
		}
		zh, okZh := i18n.Lookup(i18n.LangZh, code)
		if !okZh || zh == "" || zh == fallbackZh || zh == code {
			t.Errorf("code %s: zh-CN 缺翻译（当前 %q），会在运行时静默降级为「未知错误」", code, zh)
		}
		en, okEn := i18n.Lookup(i18n.LangEn, code)
		if !okEn || en == "" || en == fallbackEn || en == code {
			t.Errorf("code %s: en-US 缺翻译（当前 %q）", code, en)
		}
	}
}
