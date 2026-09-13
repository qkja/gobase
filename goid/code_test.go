package goid

import (
	"strings"
	"testing"
)

// crockfordAlphabet ULID 使用的最小化混淆 Base32 字母表（不含 I/L/O/U）
const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func isCrockford(s string) bool {
	for i := 0; i < len(s); i++ {
		if !strings.ContainsRune(crockfordAlphabet, rune(s[i])) {
			return false
		}
	}
	return true
}

// TestGenerateCodeFormat 前缀拼接、总长与 ULID 字符集
func TestGenerateCodeFormat(t *testing.T) {
	prefix := "AUD"
	code := GenerateCode(prefix)

	if !strings.HasPrefix(code, prefix) {
		t.Fatalf("GenerateCode(%q) = %q, want prefix %q", prefix, code, prefix)
	}
	body := strings.TrimPrefix(code, prefix)
	if len(body) != 26 {
		t.Fatalf("ULID body length = %d, want 26 (code=%q)", len(body), code)
	}
	if !isCrockford(body) {
		t.Fatalf("ULID body %q contains non-Crockford chars", body)
	}
}

// TestGenerateCodeEmptyPrefix 空前缀时仅返回 26 位 ULID
func TestGenerateCodeEmptyPrefix(t *testing.T) {
	code := GenerateCode("")
	if len(code) != 26 {
		t.Fatalf("GenerateCode(\"\") length = %d, want 26", len(code))
	}
	if !isCrockford(code) {
		t.Fatalf("code %q contains non-Crockford chars", code)
	}
}

// TestGenerateCodeOrdered 连续生成的 code 严格递增（时间戳为基数，字典序即时间序）
func TestGenerateCodeOrdered(t *testing.T) {
	const n = 5000
	prev := ""
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		code := GenerateCode("AUD")
		if _, dup := seen[code]; dup {
			t.Fatalf("duplicated code at index %d: %q", i, code)
		}
		seen[code] = struct{}{}
		if code <= prev {
			t.Fatalf("code not increasing: prev=%q, cur=%q (index=%d)", prev, code, i)
		}
		prev = code
	}
}
